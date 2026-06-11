package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    RpcWsURL       string
    RpcHttpURL     string
    PraTokenAddr   string
    PraEscrowAddr  string
    PrivateKey     string
}

var AppConfig *Config

func LoadConfig() {
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found, using system env")
    }

    AppConfig = &Config{
        RpcWsURL:      os.Getenv("RPC_WS_URL"),
        RpcHttpURL:    os.Getenv("RPC_HTTP_URL"),
        PraTokenAddr:  os.Getenv("PRA_TOKEN_ADDRESS"),
        PraEscrowAddr: os.Getenv("PRA_ESCROW_ADDRESS"),
        PrivateKey:    os.Getenv("PRIVATE_KEY"),
    }
}
