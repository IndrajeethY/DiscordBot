package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func init() {
	RegisterCommand("help", help, "Displays help information for available commands.")
}

func help(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	helpMessage := "Available commands:\n"
	for name, help := range HelpMessages {
		helpMessage += fmt.Sprintf("• %s: %s\n", name, help)
	}
	helpMessage += "Use one of the following prefixes to invoke a command: !, ?, /"
	_, err := s.ChannelMessageSendReply(m.ChannelID, helpMessage, m.Reference())
	return err
}
