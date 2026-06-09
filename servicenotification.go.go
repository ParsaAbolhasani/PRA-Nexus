package service

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"

    "pra-exchange/config"
    "pra-exchange/database"
)

type NotificationService struct{}

func NewNotificationService() *NotificationService {
    return &NotificationService{}
}

// SendTelegramMessage : ارسال پیام به تلگرام
func (ns *NotificationService) SendTelegramMessage(message string) error {
    if config.AppConfig.TelegramBotToken == "" || config.AppConfig.TelegramChatID == "" {
        log.Println("⚠️ Telegram not configured, skipping notification")
        return nil
    }

    url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.AppConfig.TelegramBotToken)
    
    payload := map[string]string{
        "chat_id": config.AppConfig.TelegramChatID,
        "text":    message,
        "parse_mode": "HTML",
    }
    
    jsonPayload, _ := json.Marshal(payload)
    
    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    log.Printf("📨 Telegram notification sent: %s", message[:min(50, len(message))])
    return nil
}

// SaveNotification : ذخیره در دیتابیس
func (ns *NotificationService) SaveNotification(notifType, title, message string, tradeID *uint64, txHash string) {
    notification := database.Notification{
        Type:    notifType,
        Title:   title,
        Message: message,
        TradeID: tradeID,
        TxHash:  txHash,
        SentAt:  time.Now(),
    }
    
    if err := database.DB.Create(&notification).Error; err != nil {
        log.Printf("Failed to save notification: %v", err)
    }
    
    // ارسال به تلگرام
    go ns.SendTelegramMessage(fmt.Sprintf("<b>%s</b>\n%s", title, message))
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}