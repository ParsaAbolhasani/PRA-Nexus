package database

import (
    "time"
)

type Trade struct {
    ID          uint      `gorm:"primaryKey"`
    TradeID     uint64    `gorm:"uniqueIndex"`
    Buyer       string    `gorm:"index"`
    Seller      string    `gorm:"index"`
    PRAAmount   float64
    IRRAmount   float64
    Status      string    // pending, paid, completed, cancelled
    TxHash      string    `gorm:"index"`
    CreatedAt   time.Time
    PaidAt      *time.Time
    CompletedAt *time.Time
}

type TransferEvent struct {
    ID        uint      `gorm:"primaryKey"`
    From      string    `gorm:"index"`
    To        string    `gorm:"index"`
    Amount    float64
    TxHash    string    `gorm:"uniqueIndex"`
    BlockNum  uint64
    CreatedAt time.Time
}

type Notification struct {
    ID        uint      `gorm:"primaryKey"`
    Type      string    // trade_created, payment_confirmed, trade_completed, transfer
    Title     string
    Message   string
    TradeID   *uint64
    TxHash    string
    SentAt    time.Time
}