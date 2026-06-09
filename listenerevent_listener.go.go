package listener

import (
    "context"
    "log"
    "math/big"
    "strings"

    "pra-exchange/config"

    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/ethclient"
)

type EventListener struct {
    wsClient   *ethclient.Client
    httpClient *ethclient.Client
    escrowAddr common.Address
    tokenAddr  common.Address
    handlers   *EventHandler
}

func NewEventListener(handlers *EventHandler) (*EventListener, error) {
    wsClient, err := ethclient.Dial(config.AppConfig.RpcWsURL)
    if err != nil {
        return nil, err
    }

    httpClient, err := ethclient.Dial(config.AppConfig.RpcHttpURL)
    if err != nil {
        return nil, err
    }

    return &EventListener{
        wsClient:   wsClient,
        httpClient: httpClient,
        escrowAddr: common.HexToAddress(config.AppConfig.PraEscrowAddr),
        tokenAddr:  common.HexToAddress(config.AppConfig.PraTokenAddr),
        handlers:   handlers,
    }, nil
}

// StartListening : شروع شنود رویدادها (با گوروتین)
func (el *EventListener) StartListening(ctx context.Context) {
    // فیلتر رویدادهای Escrow
    escrowQuery := ethereum.FilterQuery{
        Addresses: []common.Address{el.escrowAddr},
    }

    // فیلتر رویدادهای توکن (Transfer, Approval)
    tokenQuery := ethereum.FilterQuery{
        Addresses: []common.Address{el.tokenAddr},
    }

    // کانال دریافت لاگ‌ها
    logs := make(chan types.Log)

    // اشتراک WebSocket برای Escrow
    escrowSub, err := el.wsClient.SubscribeFilterLogs(ctx, escrowQuery, logs)
    if err != nil {
        log.Fatal("Escrow subscribe error:", err)
    }

    // اشتراک WebSocket برای Token
    tokenSub, err := el.wsClient.SubscribeFilterLogs(ctx, tokenQuery, logs)
    if err != nil {
        log.Fatal("Token subscribe error:", err)
    }

    log.Println("✅ WebSocket connected. Listening for events...")

    // پردازش رویدادها در یک لوپ جداگانه (گوروتین)
    go func() {
        for {
            select {
            case err := <-escrowSub.Err():
                log.Println("Escrow subscription error:", err)
                // تلاش برای reconnect
                return
            case err := <-tokenSub.Err():
                log.Println("Token subscription error:", err)
                return
            case vLog := <-logs:
                el.processLog(vLog)
            case <-ctx.Done():
                log.Println("Stopping listener...")
                return
            }
        }
    }()
}

// processLog : پردازش هر لاگ بر اساس آدرس قرارداد و signature
func (el *EventListener) processLog(vLog types.Log) {
    // تشخیص قرارداد
    if vLog.Address == el.escrowAddr {
        el.processEscrowEvent(vLog)
    } else if vLog.Address == el.tokenAddr {
        el.processTokenEvent(vLog)
    }
}

// processEscrowEvent : پردازش رویدادهای Escrow
func (el *EventListener) processEscrowEvent(vLog types.Log) {
    // Topic0 نشان‌دهنده نام رویداد است
    if len(vLog.Topics) == 0 {
        return
    }

    eventSignature := vLog.Topics[0].Hex()

    switch eventSignature {
    case "0x"...: // TradeCreated
        el.handlers.OnTradeCreated(vLog)
    case "0x"...: // PaymentConfirmed
        el.handlers.OnPaymentConfirmed(vLog)
    case "0x"...: // TradeCompleted
        el.handlers.OnTradeCompleted(vLog)
    case "0x"...: // TradeCancelled
        el.handlers.OnTradeCancelled(vLog)
    default:
        log.Printf("Unknown escrow event: %s", eventSignature)
    }
}

// processTokenEvent : پردازش رویدادهای توکن PRA
func (el *EventListener) processTokenEvent(vLog types.Log) {
    if len(vLog.Topics) < 3 {
        return
    }

    eventSignature := vLog.Topics[0].Hex()
    // signature Transfer(address,address,uint256)
    transferSig := "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

    if eventSignature == transferSig {
        from := common.HexToAddress(vLog.Topics[1].Hex())
        to := common.HexToAddress(vLog.Topics[2].Hex())
        amount := new(big.Int).SetBytes(vLog.Data)

        el.handlers.OnTokenTransfer(from, to, amount)
    }
}