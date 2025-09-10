package handlers

import (
	"DiscordBot/utils"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var devLoadOnce sync.Once

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if utils.IsBotMessage(s, m) {
		return
	}
	if !utils.IsDevUser(m.Author.ID) {
		return
	}

	parts := utils.SplitArgs(m.Content)
	if len(parts) == 0 {
		return
	}

	if len(parts[0]) < 2 {
		return
	}

	prefix := rune(parts[0][0])
	if commands, exists := CommandMap[prefix]; exists {
		command := strings.ToLower(parts[0][1:])
		if metadata, exists := commands[command]; exists {
			if err := metadata.Handler(s, m, parts[1:]); err != nil {
				_, _ = s.ChannelMessageSendReply(m.ChannelID, "Error: "+err.Error(), m.Reference())
			}
		}
	}
}
