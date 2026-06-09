package config

import (
    "log"
    "os"
    "strconv"

    "github.com/joho/godotenv"
)

type Config struct {
    // Blockchain
    RpcWsURL      string
    RpcHttpURL    string
    PraTokenAddr  string
    PraEscrowAddr string
    PrivateKey    string

    // Database
    DBHost     string
    DBPort     int
    DBUser     string
    DBPassword string
    DBName     string

    // Telegram
    TelegramBotToken string
    TelegramChatID   string

    // API
    APIPort int
}

var AppConfig *Config

func LoadConfig() {
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found, using system env")
    }

    dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
    apiPort, _ := strconv.Atoi(getEnv("API_PORT", "8080"))

    AppConfig = &Config{
        RpcWsURL:      getEnv("RPC_WS_URL", ""),
        RpcHttpURL:    getEnv("RPC_HTTP_URL", ""),
        PraTokenAddr:  getEnv("PRA_TOKEN_ADDRESS", ""),
        PraEscrowAddr: getEnv("PRA_ESCROW_ADDRESS", ""),
        PrivateKey:    getEnv("PRIVATE_KEY", ""),

        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     dbPort,
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", ""),
        DBName:     getEnv("DB_NAME", "pra_exchange"),

        TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
        TelegramChatID:   getEnv("TELEGRAM_CHAT_ID", ""),

        APIPort: apiPort,
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}