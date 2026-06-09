package database

import (
    "fmt"
    "log"

    "pra-exchange/config"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
        config.AppConfig.DBHost,
        config.AppConfig.DBUser,
        config.AppConfig.DBPassword,
        config.AppConfig.DBName,
        config.AppConfig.DBPort,
    )

    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Auto migrate
    err = DB.AutoMigrate(&Trade{}, &TransferEvent{}, &Notification{})
    if err != nil {
        log.Fatal("Failed to migrate database:", err)
    }

    log.Println("✅ Database connected and migrated")
}