package handlers

import "github.com/bwmarrin/discordgo"

type CommandHandler func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error

type CommandMetadata struct {
	Handler CommandHandler
	Help    string
}

var CommandMap = make(map[rune]map[string]CommandMetadata)

var HelpMessages = make(map[string]string)

func RegisterCommand(name string, handler CommandHandler, help string) {
	prefixes := []rune{'!', '?', '/'}
	for _, prefix := range prefixes {
		if _, exists := CommandMap[prefix]; !exists {
			CommandMap[prefix] = make(map[string]CommandMetadata)
		}
		CommandMap[prefix][name] = CommandMetadata{
			Handler: handler,
			Help:    HelpMessages[name],
		}
	}
	if help != "" {
		HelpMessages[name] = help
	}
}
