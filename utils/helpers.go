package utils

import (
	"log"
	"strings"
	"sync"
	"unicode"

	"github.com/bwmarrin/discordgo"
)

var (
	devMu    sync.RWMutex
	devUsers = make(map[string]bool)
)

func IsBotMessage(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	return m.Author.ID == s.State.User.ID
}

func SplitArgs(input string) []string {
	return strings.FieldsFunc(input, func(r rune) bool {
		return unicode.IsSpace(r)
	})
}

func RegisterDevUsers(ids []string) {
	devMu.Lock()
	defer devMu.Unlock()
	devUsers = make(map[string]bool, len(ids))
	for _, id := range ids {
		if id != "" {
			devUsers[id] = true
		}
	}
	log.Printf("Registered %d dev users", len(devUsers))
}

func IsDevUser(id string) bool {
	devMu.RLock()
	defer devMu.RUnlock()
	if len(devUsers) == 0 {
		return false
	}
	return devUsers[id]
}
