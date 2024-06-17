package utils

import (
	"github.com/bwmarrin/discordgo"
	"os"
)

// UploadFile uploads a file to a specified Discord channel.
func UploadFile(s *discordgo.Session, channelID, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = s.ChannelFileSend(channelID, file.Name(), file)
	return err
}
