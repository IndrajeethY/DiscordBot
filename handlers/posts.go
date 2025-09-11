package handlers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	customBtnRe = regexp.MustCompile(`::btn\[(?P<label>[^\]]+)\]\{(?P<url>https?://[^\s}]+)\}`)
)

func init() {
	RegisterCommand("post", postCmd, "Repost the replied message to a channel: reply + post [channelID]. If no channelID provided, reposts to same channel.")
}

func postCmd(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if m.MessageReference == nil || m.MessageReference.MessageID == "" {
		_, _ = s.ChannelMessageSendReply(m.ChannelID, "Reply to a message with `post [channelID]` to repost it.", m.Reference())
		return nil
	}

	targetChannel := m.ChannelID
	if len(args) >= 1 && args[0] != "" {
		targetChannel = args[0]
	}

	refChan := m.MessageReference.ChannelID
	refMsgID := m.MessageReference.MessageID
	refMsg, err := s.ChannelMessage(refChan, refMsgID)
	if err != nil {
		_, _ = s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("failed to fetch referenced message: %v", err), m.Reference())
		return err
	}

	content := stripButtonSyntax(refMsg.Content)
	components := parseButtonsFromText(refMsg.Content)

	msgSend := &discordgo.MessageSend{
		Content:    content,
		Components: components,
		Flags:      discordgo.MessageFlagsSuppressEmbeds,
	}

	_, err = s.ChannelMessageSendComplex(targetChannel, msgSend)
	if err != nil {
		_, _ = s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("failed to post message to channel %s: %v", targetChannel, err), m.Reference())
		return err
	}

	_, _ = s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("Reposted message to <#%s>.", targetChannel), m.Reference())
	return nil
}

func parseButtonsFromText(text string) []discordgo.MessageComponent {
	btns := make([]discordgo.Button, 0, 5)
	customMatches := customBtnRe.FindAllStringSubmatch(text, -1)
	for _, mat := range customMatches {
		if len(mat) < 3 {
			continue
		}
		label := strings.TrimSpace(mat[1])
		url := strings.TrimSpace(mat[2])
		if label == "" || url == "" {
			continue
		}
		btns = append(btns, discordgo.Button{
			Label: label,
			Style: discordgo.LinkButton,
			URL:   url,
		})
		if len(btns) >= 5 {
			break
		}
	}

	if len(btns) == 0 {
		return nil
	}

	row := discordgo.ActionsRow{
		Components: make([]discordgo.MessageComponent, 0, len(btns)),
	}
	for _, b := range btns {
		row.Components = append(row.Components, b)
	}
	return []discordgo.MessageComponent{row}
}

func stripButtonSyntax(text string) string {
	return customBtnRe.ReplaceAllString(text, "")
}
