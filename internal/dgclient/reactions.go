package dgclient

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func (client *DiscordGoClient) GetOnReactionAddHandler() func(*discordgo.Session, *discordgo.MessageReactionAdd) {
	return func(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
		if m.UserID == s.State.User.ID || m.MessageID != client.roleMessage.ID || !client.db.RoleIsEmojiTaken(m.Emoji.Name) {
			return
		}

		roleErr := s.GuildMemberRoleAdd(m.GuildID, m.UserID, client.db.RoleGetIdByEmoji(m.Emoji.Name))

		if roleErr != nil {
			log.Println(roleErr.Error())
		}
	}
}

func (client *DiscordGoClient) GetOnReactionRemoveHandler() func(*discordgo.Session, *discordgo.MessageReactionRemove) {
	return func(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
		if m.UserID == s.State.User.ID || m.MessageID != client.roleMessage.ID || !client.db.RoleIsEmojiTaken(m.Emoji.Name) {
			return
		}

		roleErr := s.GuildMemberRoleRemove(m.GuildID, m.UserID, client.db.RoleGetIdByEmoji(m.Emoji.Name))

		if roleErr != nil {
			log.Println(roleErr.Error())
		}
	}
}
