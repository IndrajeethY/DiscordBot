package handlers

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {
	RegisterCommand("purge", purge, "Deletes messages from the command message to the current message.")
	RegisterCommand("del", delete, "Deletes the command message.")
}

func purge(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	channelID := m.ChannelID

	var startMessageID string

	if m.MessageReference != nil {
		startMessageID = m.MessageReference.MessageID
	} else {
		startMessageID = m.ID
	}
	messages, err := s.ChannelMessages(channelID, 100, m.ID, startMessageID, "")
	if err != nil {
		return err
	}
	var messageIDs []string
	for _, message := range messages {
		messageIDs = append(messageIDs, message.ID)
	}
	if len(messageIDs) > 1 {
		err = s.ChannelMessagesBulkDelete(channelID, messageIDs)
	} else if len(messageIDs) == 1 {
		err = s.ChannelMessageDelete(channelID, messageIDs[0])
	}
	deL, _ := s.ChannelMessageSend(channelID, fmt.Sprintf("Deleted %d messages.", len(messageIDs)))
	time.Sleep(3 * time.Second)
	s.ChannelMessageDelete(channelID, deL.ID)
	return err
}

func delete(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	channelID := m.ChannelID
	var messageIDs []string
	if m.MessageReference != nil {
		messageIDs = append(messageIDs, m.MessageReference.MessageID)
	}
	messageIDs = append(messageIDs, m.ID)
	if len(messageIDs) > 1 {
		err := s.ChannelMessagesBulkDelete(channelID, messageIDs)
		return err
	} else if len(messageIDs) == 1 {
		err := s.ChannelMessageDelete(channelID, messageIDs[0])
		return err
	}
	return nil
}
