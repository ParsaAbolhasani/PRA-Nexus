package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ======================== تنظیمات شبکه ========================
var (
	RPC_URL          = "https://sepolia.infura.io/v3/YOUR_PROJECT_ID"
	CONTRACT_ADDRESS = common.HexToAddress("0x...") // آدرس توکن PRA
	PRIVATE_KEY      = "your_private_key_here"
	TOKEN_DECIMALS   = 18
)

// ======================== Converter Logic ========================
type Converter struct {
	RatePRAYPerIRR float64 // هر 1 IRR چند PRA
}

func NewConverter(rate float64) *Converter {
	return &Converter{RatePRAYPerIRR: rate}
}

// IRR to PRA
func (c *Converter) IRRToPRA(irrAmount float64) float64 {
	return irrAmount * c.RatePRAYPerIRR
}

// PRA to IRR
func (c *Converter) PRAToIRR(praAmount float64) float64 {
	return praAmount / c.RatePRAYPerIRR
}

// ======================== Wallet Manager (بلاکچینی) ========================
type WalletManager struct {
	Client      *ethclient.Client
	PrivateKey  *ecdsa.PrivateKey
	FromAddress common.Address
}

func NewWalletManager(rpcURL, privKeyHex string) (*WalletManager, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.HexToECDSA(privKeyHex)
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("invalid public key")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	return &WalletManager{
		Client:      client,
		PrivateKey:  privateKey,
		FromAddress: fromAddress,
	}, nil
}

// GetPRABalance : دریافت موجودی واقعی توکن PRA از بلاکچین
func (wm *WalletManager) GetPRABalance(address common.Address) (*big.Float, error) {
	// در اینجا باید ABI توکن را بخوانی
	return big.NewFloat(1000), nil // placeholder
}

// TransferPRA : انتقال توکن PRA به آدرس دیگر
func (wm *WalletManager) TransferPRA(to common.Address, amount *big.Int) (string, error) {
	nonce, err := wm.Client.PendingNonceAt(context.Background(), wm.FromAddress)
	if err != nil {
		return "", err
	}

	gasPrice, err := wm.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}

	tx := types.NewTransaction(nonce, CONTRACT_ADDRESS, big.NewInt(0), 210000, gasPrice, nil)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(11155111)), wm.PrivateKey)
	if err != nil {
		return "", err
	}

	err = wm.Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}

// ======================== Payment Gateway ========================
type PaymentGateway struct{}

func (pg *PaymentGateway) ProcessPayment(amountIRR float64) bool {
	fmt.Printf("\n💳 اتصال به درگاه پرداخت... مبلغ: %.0f IRR\n", amountIRR)
	time.Sleep(1 * time.Second)
	fmt.Println("✅ پرداخت موفق")
	return true
}

type PayoutSystem struct{}

func (ps *PayoutSystem) ProcessPayout(amountIRR float64, bankAccount string) bool {
	fmt.Printf("\n🏦 واریز به حساب بانکی %s مبلغ: %.0f IRR\n", bankAccount, amountIRR)
	time.Sleep(1 * time.Second)
	fmt.Println("✅ تسویه انجام شد")
	return true
}

// ======================== فرآیند خرید و فروش PRA ========================
type TradingService struct {
	Converter        *Converter
	Wallet           *WalletManager
	PaymentGW        *PaymentGateway
	PayoutGW         *PayoutSystem
	UserRialBalance  float64
}

func (ts *TradingService) BuyPRA(userAddress common.Address, amountIRR float64) {
	fmt.Printf("\n🟢 درخواست خرید PRA با مبلغ %.0f IRR\n", amountIRR)
	praAmount := ts.Converter.IRRToPRA(amountIRR)
	fmt.Printf("💰 نرخ: هر 1 IRR = %.4f PRA | دریافت PRA: %.4f\n", ts.Converter.RatePRAYPerIRR, praAmount)

	if ts.UserRialBalance < amountIRR {
		needed := amountIRR - ts.UserRialBalance
		if !ts.PaymentGW.ProcessPayment(needed) {
			fmt.Println("❌ خطا در پرداخت")
			return
		}
		ts.UserRialBalance += needed
	}

	ts.UserRialBalance -= amountIRR

	txHash, err := ts.Wallet.TransferPRA(userAddress, big.NewInt(int64(praAmount*1e18)))
	if err != nil {
		fmt.Printf("❌ خطا در انتقال بلاکچین: %v\n", err)
		return
	}

	fmt.Printf("✅ خرید موفق! %f PRA به آدرس %s ارسال شد.\n", praAmount, userAddress.Hex())
	fmt.Printf("🔗 تراکنش: %s\n", txHash)
}

func (ts *TradingService) SellPRA(userAddress common.Address, praAmount float64, bankAccount string) {
	fmt.Printf("\n🔴 درخواست فروش %.4f PRA\n", praAmount)
	irrAmount := ts.Converter.PRAToIRR(praAmount)
	fmt.Printf("💰 مبلغ دریافتی: %.0f IRR\n", irrAmount)

	balance, err := ts.Wallet.GetPRABalance(userAddress)
	if err != nil {
		fmt.Printf("❌ خطا در خواندن موجودی: %v\n", err)
		return
	}
	balanceFloat, _ := balance.Float64()
	if balanceFloat < praAmount {
		fmt.Println("❌ موجودی PRA کافی نیست")
		return
	}

	_, err = ts.Wallet.TransferPRA(ts.Wallet.FromAddress, big.NewInt(int64(praAmount*1e18)))
	if err != nil {
		fmt.Printf("❌ خطا در انتقال PRA به صرافی: %v\n", err)
		return
	}
	fmt.Println("✅ توکن‌ها به خزانه صرافی منتقل شد.")

	if !ts.PayoutGW.ProcessPayout(irrAmount, bankAccount) {
		fmt.Println("❌ خطا در تسویه ریالی")
		return
	}

	fmt.Printf("✅ فروش موفق! مبلغ %.0f IRR به حساب %s واریز شد.\n", irrAmount, bankAccount)
}

// ======================== Main ========================
func main() {
	converter := NewConverter(0.00002) // هر 1 IRR = 0.00002 PRA

	wm, err := NewWalletManager(RPC_URL, PRIVATE_KEY)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✅ متصل به بلاکچین | آدرس صرافی: %s\n", wm.FromAddress.Hex())

	ts := &TradingService{
		Converter:        converter,
		Wallet:           wm,
		PaymentGW:        &PaymentGateway{},
		PayoutGW:         &PayoutSystem{},
		UserRialBalance:  500000,
	}

	userAddress := common.HexToAddress("0xAb8483F64d9C6d1EcF9b849Ae677dD3315835cb2")
	bankAccount := "IR1234567890"

	// خرید PRA
	ts.BuyPRA(userAddress, 100000)

	// فروش PRA
	ts.SellPRA(userAddress, 2.5, bankAccount)
}