package listener

import (
    "log"
    "math/big"
    "time"

    "pra-exchange/database"
    "pra-exchange/service"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
)

type EventHandler struct {
    tradingService     *service.TradingService
    notificationService *service.NotificationService
}

func NewEventHandler(tradingService *service.TradingService, notificationService *service.NotificationService) *EventHandler {
    return &EventHandler{
        tradingService:     tradingService,
        notificationService: notificationService,
    }
}

func (h *EventHandler) OnTradeCreated(vLog types.Log) {
    log.Printf("🆕 New Trade Created | TxHash: %s", vLog.TxHash.Hex())
    
    // استخراج اطلاعات از vLog.Data (بسته به ساختار رویداد)
    // اینجا ساده شده، در واقع باید ABI را unpack کنی
    
    tradeID := uint64(1) // نمونه، از vLog استخراج کن
    
    // ذخیره در دیتابیس
    trade := database.Trade{
        TradeID:   tradeID,
        Status:    "pending",
        TxHash:    vLog.TxHash.Hex(),
        CreatedAt: time.Now(),
    }
    database.DB.Create(&trade)
    
    // اعلان
    h.notificationService.SaveNotification(
        "trade_created",
        "🆕 معامله جدید ایجاد شد",
        "یک معامله جدید با شناسه "+string(rune(tradeID))+" ایجاد شد.",
        &tradeID,
        vLog.TxHash.Hex(),
    )
}

func (h *EventHandler) OnPaymentConfirmed(vLog types.Log) {
    log.Printf("💳 Payment Confirmed | TxHash: %s", vLog.TxHash.Hex())
    
    var tradeID uint64
    // استخراج tradeID از vLog
    
    // به‌روزرسانی دیتابیس
    database.DB.Model(&database.Trade{}).
        Where("trade_id = ?", tradeID).
        Updates(map[string]interface{}{
            "status": "paid",
            "paid_at": time.Now(),
        })
    
    h.notificationService.SaveNotification(
        "payment_confirmed",
        "💳 پرداخت تأیید شد",
        "پرداخت برای معامله تأیید شد.",
        &tradeID,
        vLog.TxHash.Hex(),
    )
    
    // شروع تسویه خودکار (در صورت نیاز)
    go h.tradingService.ProcessPayout(tradeID)
}

func (h *EventHandler) OnTradeCompleted(vLog types.Log) {
    log.Printf("✅ Trade Completed | TxHash: %s", vLog.TxHash.Hex())
    
    var tradeID uint64
    
    database.DB.Model(&database.Trade{}).
        Where("trade_id = ?", tradeID).
        Updates(map[string]interface{}{
            "status": "completed",
            "completed_at": time.Now(),
        })
    
    h.notificationService.SaveNotification(
        "trade_completed",
        "✅ معامله تکمیل شد",
        "معامله با موفقیت انجام شد.",
        &tradeID,
        vLog.TxHash.Hex(),
    )
}

func (h *EventHandler) OnTradeCancelled(vLog types.Log) {
    log.Printf("❌ Trade Cancelled | TxHash: %s", vLog.TxHash.Hex())
    
    var tradeID uint64
    
    database.DB.Model(&database.Trade{}).
        Where("trade_id = ?", tradeID).
        Update("status", "cancelled")
    
    h.notificationService.SaveNotification(
        "trade_cancelled",
        "❌ معامله لغو شد",
        "معامله لغو شد.",
        &tradeID,
        vLog.TxHash.Hex(),
    )
}

func (h *EventHandler) OnTokenTransfer(from, to common.Address, amount *big.Int) {
    amountFloat := new(big.Float).SetInt(amount)
    divisor := new(big.Float).SetFloat64(1e18)
    amountFloat.Quo(amountFloat, divisor)
    
    log.Printf("🔄 Token Transfer: %s -> %s | Amount: %.6f", from.Hex(), to.Hex(), amountFloat)
    
    // ذخیره در دیتابیس
    transfer := database.TransferEvent{
        From:      from.Hex(),
        To:        to.Hex(),
        Amount:    amountFloat,
        CreatedAt: time.Now(),
    }
    database.DB.Create(&transfer)
    
    h.notificationService.SaveNotification(
        "transfer",
        "🔄 انتقال توکن",
        amountFloat.String()+" PRA از "+from.Hex()+" به "+to.Hex()+" منتقل شد.",
        nil,
        "",
    )
}