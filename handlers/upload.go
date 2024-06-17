package handlers

import (
	"DiscordBot/utils"
	"github.com/bwmarrin/discordgo"
)

func init() {
	RegisterCommand("upload", uploadFile, "Uploads a specified file.")
}

func uploadFile(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	filePath := "Pottu Thakku.mp3"
	return utils.UploadFile(s, m.ChannelID, filePath)
}
