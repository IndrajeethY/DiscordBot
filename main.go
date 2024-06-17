package main

import (
	"log"
	"os"
	"os/signal"

	"DiscordBot/config"
	"DiscordBot/handlers"

	"github.com/bwmarrin/discordgo"
)

func main() {
	config.LoadConfig()
	dg, err := discordgo.New("Bot " + config.BotToken)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
		return
	}
	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) { log.Println("Bot is up!") })
	dg.AddHandler(handlers.MessageCreate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages
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
