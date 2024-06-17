package handlers

import (
	"github.com/bwmarrin/discordgo"
)

func init() {
	RegisterCommand("ping", ping, "Responds with 'Pong!'.")
}

func ping(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	_, err := s.ChannelMessageSendReply(m.ChannelID, "Pong!", m.Reference())
	return err
}
