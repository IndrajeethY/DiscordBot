package handlers

import (
	"github.com/bwmarrin/discordgo"
)

func init() {
	RegisterCommand("start", start, "Starts the bot and responds with a message.")
}

func start(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	_, err := s.ChannelMessageSendReply(m.ChannelID, "Bot is running!", m.Reference())
	return err
}
