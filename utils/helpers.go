package utils

import (
	"strings"
	"unicode"

	"github.com/bwmarrin/discordgo"
)

func IsBotMessage(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	return m.Author.ID == s.State.User.ID
}

func SplitArgs(input string) []string {
	return strings.FieldsFunc(input, func(r rune) bool {
		return unicode.IsSpace(r)
	})
}
