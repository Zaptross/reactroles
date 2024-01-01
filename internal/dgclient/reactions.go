package dgclient

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
	"github.com/zaptross/reactroles/internal/pgdb"
)

func (client *DiscordGoClient) GetOnReactionAddHandler() func(*discordgo.Session, *discordgo.MessageReactionAdd) {
	return func(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
		selectors := lookupMessagesForSelectors(client, client.db.SelectorGetAll(m.GuildID))
		if m.UserID == s.State.User.ID || !isReactingToSelector(selectors, m.MessageID) || !client.db.RoleIsEmojiTaken(m.Emoji.Name, m.GuildID) {
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
		selectors := lookupMessagesForSelectors(client, client.db.SelectorGetAll(m.GuildID))
		if m.UserID == s.State.User.ID || !isReactingToSelector(selectors, m.MessageID) || !client.db.RoleIsEmojiTaken(m.Emoji.Name, m.GuildID) {
			return
		}

		roleErr := s.GuildMemberRoleRemove(m.GuildID, m.UserID, client.db.RoleGetIdByEmoji(m.Emoji.Name))

		if roleErr != nil {
			log.Println(roleErr.Error())
		}
	}
}

func isReactingToSelector(selectors []*discordgo.Message, messageID string) bool {
	for _, selector := range selectors {
		if selector.ID == messageID {
			return true
		}
	}

	return false
}

func lookupMessagesForSelectors(client *DiscordGoClient, selectors []pgdb.Selector) []*discordgo.Message {
	return lo.Map(
		selectors,
		func(selector pgdb.Selector, _ int) *discordgo.Message {
			return lookupMessageForSelector(client, selector)
		},
	)
}

func lookupMessageForSelector(client *DiscordGoClient, selector pgdb.Selector) *discordgo.Message {
	selectorMessage, err := client.Session.ChannelMessage(selector.ChannelID, selector.ID)

	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return selectorMessage
}
