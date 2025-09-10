package handlers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type ComponentHandler func(s *discordgo.Session, i *discordgo.InteractionCreate)

var componentHandlers = make(map[string]ComponentHandler)
var componentPrefixHandlers = make(map[string]ComponentHandler)

func RegisterComponent(id string, h ComponentHandler) {
	componentHandlers[id] = h
}

func RegisterComponentPrefix(prefix string, h ComponentHandler) {
	componentPrefixHandlers[prefix] = h
}

func RegisterComponentHandlers(session *discordgo.Session) {
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionMessageComponent {
			return
		}
		customID := i.MessageComponentData().CustomID
		if customID == "" {
			return
		}
		if h, ok := componentHandlers[customID]; ok {
			h(s, i)
			return
		}
		for p, h := range componentPrefixHandlers {
			if strings.HasPrefix(customID, p) {
				h(s, i)
				return
			}
		}
	})
}
