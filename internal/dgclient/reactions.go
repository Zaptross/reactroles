package dgclient

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
	"github.com/zaptross/reactroles/internal/pgdb"
	"github.com/zaptross/reactroles/internal/utils"
)

func (client *DiscordGoClient) GetOnReactionAddHandler() func(*discordgo.Session, *discordgo.MessageReactionAdd) {
	return func(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
		// ignore the bot's reactions
		if m.UserID == s.State.User.ID {
			return
		}

		selectors := client.db.SelectorGetAll(m.GuildID)
		selector, reactingToSelector := isReactingToSelector(selectors, m.MessageID)
		if !reactingToSelector {
			return
		}

		// if bell to the preamble selector, add the notify role to the user
		if selector.RoleID == "" && m.Emoji.Name == utils.EMOJI_BELL {
			config := client.db.ServerConfigurationGet(m.GuildID)
			roleErr := s.GuildMemberRoleAdd(m.GuildID, m.UserID, config.NotifyRoleID)

			if roleErr != nil {
				log.Println(roleErr.Error())
			}

			return
		}

		// If the emoji is not used for a role, ignore it
		if !client.db.RoleIsEmojiTaken(m.Emoji.Name, m.GuildID) {
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
		// ignore the bot's reactions
		if m.UserID == s.State.User.ID {
			return
		}

		selectors := client.db.SelectorGetAll(m.GuildID)
		selector, reactingToSelector := isReactingToSelector(selectors, m.MessageID)
		if !reactingToSelector {
			return
		}

		// if bell to the preamble selector, remove the notify role from the user
		if selector.RoleID == "" && m.Emoji.Name == utils.EMOJI_BELL {
			config := client.db.ServerConfigurationGet(m.GuildID)
			roleErr := s.GuildMemberRoleRemove(m.GuildID, m.UserID, config.NotifyRoleID)

			if roleErr != nil {
				log.Println(roleErr.Error())
			}
			return
		}

		// If the emoji is not used for a role, ignore it
		if !client.db.RoleIsEmojiTaken(m.Emoji.Name, m.GuildID) {
			return
		}

		roleErr := s.GuildMemberRoleRemove(m.GuildID, m.UserID, client.db.RoleGetIdByEmoji(m.Emoji.Name))

		if roleErr != nil {
			log.Println(roleErr.Error())
		}
	}
}

func isReactingToSelector(selectors []pgdb.Selector, messageID string) (*pgdb.Selector, bool) {
	for _, selector := range selectors {
		if selector.ID == messageID {
			return &selector, true
		}
	}

	return nil, false
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
