package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var rolesCache = map[string]string{}

func main() {
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		panic(err)
	}
	println("Bot is now running.")

	roleAssignMessageId := os.Getenv("DISCORD_ROLES_MESSAGE_ID")
	roleAssignChannelId := os.Getenv("DISCORD_ROLES_CHANNEL_ID")

	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		println(fmt.Sprintf("%s: %s", m.Author.Username, m.Content))

		if m.Content == "!ping" {
			reactErr := s.MessageReactionAdd(m.ChannelID, m.ID, "🏓")

			if reactErr != nil {
				println(reactErr.Error())
			}
		}

		if strings.Contains(m.Content, "!role") {
			action, roleName, emojiName := splitRoleCommand(m.Content)
			defer (func() {
				if s.ChannelMessageDelete(m.ChannelID, m.ID) != nil {
					fmt.Printf("Failed to delete message %s:%s\n", m.ChannelID, m.ID)
				}
			})()

			if roleName == "" || emojiName == "" || (action != "add" && action != "remove") || (action == "add" && rolesCache[emojiName] != "") || (action == "remove" && rolesCache[emojiName] == "") {
				s.ChannelMessageSendReply(m.ChannelID, "Usage: !role <add/remove> <role name> <emoji>", m.Reference())
				return
			}

			if action == "add" {
				role, roleCreateErr := s.GuildRoleCreate(m.GuildID)

				if roleCreateErr != nil {
					s.ChannelMessageSendReply(m.ChannelID, "Error creating role", m.Reference())
					println(roleCreateErr.Error())
					return
				}

				_, editErr := s.GuildRoleEdit(m.GuildID, role.ID, roleName, 0, false, 0, true)
				if editErr != nil {
					s.ChannelMessageSendReply(m.ChannelID, "Error editing role", m.Reference())
					println(editErr.Error())
					return
				}
				rolesCache[emojiName] = role.ID

				reactAddErr := s.MessageReactionAdd(roleAssignChannelId, roleAssignMessageId, emojiName)
				if reactAddErr != nil {
					s.ChannelMessageSendReply(m.ChannelID, "Error adding reaction", m.Reference())
					println(reactAddErr.Error())
					return
				}
			}

			if action == "remove" {
				if rolesCache[emojiName] == "" {
					s.ChannelMessageSendReply(m.ChannelID, "No role found for that emoji", m.Reference())
					return
				}

				deleteErr := s.GuildRoleDelete(m.GuildID, rolesCache[emojiName])
				if deleteErr != nil {
					s.ChannelMessageSendReply(m.ChannelID, "Error deleting role", m.Reference())
					println(deleteErr.Error())
					return
				}

				reactRemoveErr := s.MessageReactionsRemoveEmoji(roleAssignChannelId, roleAssignMessageId, emojiName)
				if reactRemoveErr != nil {
					s.ChannelMessageSendReply(m.ChannelID, "Error removing reaction", m.Reference())
					println(reactRemoveErr.Error())
					return
				}

				delete(rolesCache, emojiName)
			}
		}
	})

	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
		if m.UserID == s.State.User.ID || m.MessageID != roleAssignMessageId || rolesCache[m.Emoji.Name] == "" {
			return
		}

		println(fmt.Sprintf("adding %s to %s %s", m.UserID, rolesCache[m.Emoji.Name], m.Emoji.Name))

		roleErr := s.GuildMemberRoleAdd(m.GuildID, m.UserID, rolesCache[m.Emoji.Name])

		if roleErr != nil {
			println(roleErr.Error())
		}
	})

	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
		if m.UserID == s.State.User.ID || m.MessageID != roleAssignMessageId {
			return
		}

		println(fmt.Sprintf("removing %s from %s %s", m.UserID, rolesCache[m.Emoji.Name], m.Emoji.Name))

		roleErr := s.GuildMemberRoleRemove(m.GuildID, m.UserID, rolesCache[m.Emoji.Name])

		if roleErr != nil {
			println(roleErr.Error())
		}
	})

	strconv.ParseInt(os.Getenv("DISCORD_BOT_INTENT"), 10, 64)

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessageReactions | discordgo.IntentsGuildMessages)

	dg.Open()

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Open()
}

func splitRoleCommand(command string) (string, string, string) {
	split := strings.Split(command, " ")

	if len(split) == 1 {
		return "", "", ""
	}

	return split[1], split[2], split[3]
}
