package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

var BotToken string

func LoadConfig() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    BotToken = os.Getenv("BOT_TOKEN")
    if BotToken == "" {
        log.Fatalf("BOT_TOKEN is not set in the environment variables")
    }
}
