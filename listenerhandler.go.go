package listener

import (
    "context"
    "log"
    "math/big"

    "pra-exchange/service"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
)

type EventHandler struct {
    tradingService *service.TradingService
}

func NewEventHandler(tradingService *service.TradingService) *EventHandler {
    return &EventHandler{
        tradingService: tradingService,
    }
}

// OnTradeCreated : وقتی معامله جدید ایجاد شد
func (h *EventHandler) OnTradeCreated(vLog types.Log) {
    log.Printf("🆕 New Trade Created | TxHash: %s", vLog.TxHash.Hex())
    // می‌توانی اطلاعات بیشتری از vLog.Data استخراج کنی
    // و در دیتابیس ذخیره کنی یا اعلان بفرستی
}

// OnPaymentConfirmed : وقتی خریدار پرداخت را تأیید کرد (آفلاین)
func (h *EventHandler) OnPaymentConfirmed(vLog types.Log) {
    log.Printf("💳 Payment Confirmed | TxHash: %s", vLog.TxHash.Hex())
    // در اینجا می‌توانی عملیات خودکار مثل تسویه بانکی را شروع کنی
}

// OnTradeCompleted : وقتی معامله کامل شد و توکن به خریدار رسید
func (h *EventHandler) OnTradeCompleted(vLog types.Log) {
    log.Printf("✅ Trade Completed | TxHash: %s", vLog.TxHash.Hex())
    // در اینجا می‌توانی وضعیت معامله را در دیتابیس به روز کنی
}

// OnTradeCancelled : وقتی معامله لغو شد
func (h *EventHandler) OnTradeCancelled(vLog types.Log) {
    log.Printf("❌ Trade Cancelled | TxHash: %s", vLog.TxHash.Hex())
}

// OnTokenTransfer : وقتی توکن PRA جابجا شد
func (h *EventHandler) OnTokenTransfer(from, to common.Address, amount *big.Int) {
    log.Printf("🔄 Token Transfer: %s -> %s | Amount: %s", from.Hex(), to.Hex(), amount.String())
    // می‌توانی موجودی کاربران را به‌روز کنی
}