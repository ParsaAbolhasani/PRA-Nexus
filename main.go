package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "pra-exchange/config"
    "pra-exchange/listener"
    "pra-exchange/service"
)

func main() {
    // بارگذاری تنظیمات
    config.LoadConfig()

    // سرویس اصلی معاملات (در صورت نیاز)
    tradingService := service.NewTradingService()

    // هندلر رویدادها
    eventHandler := listener.NewEventHandler(tradingService)

    // راه‌اندازی شنونده WebSocket
    eventListener, err := listener.NewEventListener(eventHandler)
    if err != nil {
        log.Fatal("Failed to create event listener:", err)
    }

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // شروع شنود (در یک گوروتین)
    go eventListener.StartListening(ctx)

    log.Println("🚀 PRA Exchange Backend started with WebSocket listener")

    // اضافه کردن یک تایمر ساده برای نمایش اینکه برنامه زنده است
    go func() {
        ticker := time.NewTicker(30 * time.Second)
        defer ticker.Stop()
        for {
            select {
            case <-ticker.C:
                log.Println("💓 Backend is alive and listening to blockchain events...")
            case <-ctx.Done():
                return
            }
        }
    }()

    // منتظر سیگنال خروج
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("Shutting down gracefully...")
    cancel()
    time.Sleep(2 * time.Second)
}
