package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"slices"
)

func init() {
	RegisterCommand("rules", rulesCmd, "Send rules buttons: reply+rules <memberRoleID> OR rules <threadID> <memberRoleID>")
	RegisterComponentPrefix("rules_accept:", rulesAcceptHandler)
	RegisterComponent("rules_deny", rulesDenyHandler)
}

func rulesCmd(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if m.MessageReference != nil && m.MessageReference.MessageID != "" {
		var targetChan string
		var memberRoleID string
		if len(args) == 0 {
			_, _ = s.ChannelMessageSendReply(m.ChannelID, "Usage when replying: `rules <memberRoleID>` or `rules <targetChannelID> <memberRoleID>`", m.Reference())
			return nil
		}
		if len(args) == 1 {
			targetChan = m.ChannelID
			memberRoleID = args[0]
		} else {
			targetChan = args[0]
			memberRoleID = args[1]
		}
		refChan := m.MessageReference.ChannelID
		refMsgID := m.MessageReference.MessageID
		refMsg, err := s.ChannelMessage(refChan, refMsgID)
		if err != nil {
			_, _ = s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("Failed to fetch referenced message: `%v`", err), m.Reference())
			return err
		}
		content := refMsg.Content
		components := buildRuleButtons(memberRoleID)
		_, err = s.ChannelMessageSendComplex(targetChan, &discordgo.MessageSend{
			Content:    content,
			Embeds:     refMsg.Embeds,
			Components: components,
		})
		if err != nil {
			_, _ = s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("Failed to repost referenced message with buttons to channel `%s`: `%v`", targetChan, err), m.Reference())
			return err
		}
		_, _ = s.ChannelMessageSendReply(m.ChannelID, "Reposted message with rules buttons.", m.Reference())
		return nil
	}

	if len(args) >= 2 {
		threadID := args[0]
		memberRoleID := args[1]
		content := "**Server Rules**\n\n1. Be respectful.\n2. No spamming.\n3. Follow staff instructions.\n\nPlease respond using the buttons below."
		components := buildRuleButtons(memberRoleID)
		_, err := s.ChannelMessageSendComplex(threadID, &discordgo.MessageSend{
			Content:    content,
			Components: components,
		})
		if err != nil {
			_, _ = s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("Failed to send rules to channel `%s`: `%v`", threadID, err), m.Reference())
			return err
		}
		_, _ = s.ChannelMessageSendReply(m.ChannelID, "Rules posted successfully.", m.Reference())
		return nil
	}

	_, _ = s.ChannelMessageSendReply(m.ChannelID, "Reply to a message and use `rules <memberRoleID>` or use `rules <threadID> <memberRoleID>`", m.Reference())
	return nil
}

func buildRuleButtons(memberRoleID string) []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "I ᴀɢʀᴇᴇ ᴛᴏ ᴛʜᴇ ʀᴜʟᴇs",
					Style:    discordgo.PrimaryButton,
					CustomID: "rules_accept:" + memberRoleID,
				},
				discordgo.Button{
					Label:    "I ᴅᴏɴ'ᴛ ᴀɢʀᴇᴇ - ᴋɪᴄᴋ ᴍᴇ",
					Style:    discordgo.DangerButton,
					CustomID: "rules_deny",
				},
			},
		},
	}
}

func rulesAcceptHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	parts := strings.SplitN(i.MessageComponentData().CustomID, ":", 2)
	if len(parts) != 2 {
		_ = respondEphemeral(s, i, "Internal error: invalid accept button")
		return
	}
	roleID := parts[1]
	guildID := i.GuildID
	if guildID == "" || i.Member == nil || i.Member.User == nil {
		_ = respondEphemeral(s, i, "This action must be performed within a server.")
		return
	}
	userID := i.Member.User.ID

	roles, err := s.GuildRoles(guildID)
	if err != nil {
		_ = respondEphemeral(s, i, fmt.Sprintf("Failed to fetch guild roles: `%v`", err))
		return
	}
	var roleName string
	found := false
	for _, r := range roles {
		if r.ID == roleID {
			roleName = r.Name
			found = true
			break
		}
	}
	if !found {
		_ = respondEphemeral(s, i, fmt.Sprintf("Role not found in this server. Please verify the role ID: `%s`.", roleID))
		return
	}

	member, err := s.GuildMember(guildID, userID)
	if err == nil {
		if slices.Contains(member.Roles, roleID) {
			_ = respondEphemeral(s, i, "You are already a member.")
			return
		}
	}

	if err := s.GuildMemberRoleAdd(guildID, userID, roleID); err != nil {
		if discordErr, ok := err.(*discordgo.RESTError); ok && discordErr.Message != nil {
			_ = respondEphemeral(s, i, fmt.Sprintf("Failed to add role `%s`: `%s`", roleID, discordErr.Message.Message))
			return
		}
		_ = respondEphemeral(s, i, fmt.Sprintf("Failed to add role `%s`: `%v`", roleID, err))
		return
	}
	if roleName == "" {
		roleName = roleID
	}
	_ = respondEphemeral(s, i, fmt.Sprintf("Role `%s` assigned — welcome!", roleName))
}

func rulesDenyHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guildID := i.GuildID
	if guildID == "" || i.Member == nil || i.Member.User == nil {
		_ = respondEphemeral(s, i, "This action must be performed within a server.")
		return
	}
	userID := i.Member.User.ID
	_ = respondEphemeral(s, i, "You chose not to agree. You will be removed from the server.")
	go func() {
		time.Sleep(500 * time.Millisecond)
		_ = s.GuildMemberDelete(guildID, userID)
	}()
}

func respondEphemeral(s *discordgo.Session, i *discordgo.InteractionCreate, content string) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
