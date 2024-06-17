package main

import (
	"log"
	"os"
	"os/signal"

	"DiscordBot/config"
	"DiscordBot/database"
	"DiscordBot/handlers"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/bwmarrin/discordgo"
)

var (
	db *gorm.DB
)

func main() {
	// Connect to SQLite database
	err := database.Init("bot.db")
	if err != nil {
		log.Println("Failed to connect database:", err)
		return
	}
	defer database.Close()
	config.LoadConfig()
	dg, err := discordgo.New("Bot " + config.BotToken)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
		return
	}
	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) { log.Println("Bot is up!") })
	dg.AddHandler(handlers.MessageCreate)
	err = dg.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
		return
	}
	defer dg.Close()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutting down")
}
