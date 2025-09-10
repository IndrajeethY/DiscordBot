package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var BotToken string
var DevUserIDs []string

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	BotToken = os.Getenv("BOT_TOKEN")
	if BotToken == "" {
		log.Fatalf("BOT_TOKEN is not set in the environment variables")
	}

	devs := os.Getenv("DEV_USERS")
	if devs == "" {
		DevUserIDs = []string{}
		return
	}

	parts := strings.Split(devs, ",")
	DevUserIDs = make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			DevUserIDs = append(DevUserIDs, p)
		}
	}
}
