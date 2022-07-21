package dgclient

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type RoleCommandParams struct {
	Session   *discordgo.Session
	Message   *discordgo.MessageCreate
	Client    *DiscordGoClient
	Action    string
	RoleName  string
	EmojiName string
	Color     int
}

func (client *DiscordGoClient) GetOnMessageHandler() func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if isRoleCommand(m.Content) {
			action, roleName, emojiName, color := splitRoleCommand(m.Content)
			roleCommandHandler(RoleCommandParams{
				Session:   s,
				Message:   m,
				Client:    client,
				Action:    action,
				RoleName:  roleName,
				EmojiName: emojiName,
				Color:     color,
			})
		}
	}
}

func validateUserHasRole(s *discordgo.Session, guildID string, authorID string, role string) bool {

	member, gmErr := s.GuildMember(guildID, authorID)

	if gmErr != nil {
		log.Println(gmErr.Error())
		return false
	}

	for _, r := range member.Roles {
		if r == role {
			return true
		}
	}

	return false
}

func splitRoleCommand(command string) (string, string, string, int) {
	split := strings.Split(command, " ")[1:]

	for len(split) < 4 {
		split = append(split, "")
	}

	return strings.ToLower(split[0]), strings.ToLower(split[1]), split[2], parseColorToInt(split[3])
}

func parseColorToInt(col string) int {
	if col == "" {
		return 0
	}

	if col[0] == '#' {
		col = col[1:]
	}

	i, pErr := strconv.ParseInt(col, 16, 32)

	if pErr == nil {
		return int(i)
	}

	return 0
}

func isRoleCommand(command string) bool {
	return strings.Contains(command, "!role")
}

func roleCommandHandler(params RoleCommandParams) {
	defer params.Client.updateRoleSelectorMessage()

	validateErr := validateRoleCommand(params)

	if validateErr == nil {
		switch params.Action {
		case "add":
			handleAddAction(params)
		case "remove":
			handleRemoveAction(params)
		}

		if params.Session.ChannelMessageDelete(params.Message.ChannelID, params.Message.ID) != nil {
			log.Printf("Failed to delete message %s:%s\n", params.Message.ChannelID, params.Message.ID)
		}

	} else {
		params.Session.ChannelMessageSendReply(params.Message.ChannelID, fmt.Sprintf("⚠ Error: %s\nUsage: !role <add/remove> <role name> <emoji> [colour hex]\nNote: [] is optional\n\nExamples:\n!role add valorant :gun: #d34454\n!role add valorant :gun:\n!role remove valorant", validateErr), params.Message.Reference())
	}
}

func validateRoleCommand(params RoleCommandParams) error {
	if params.Action != "add" && params.Action != "remove" {
		return errors.New("invalid action")
	}

	if params.RoleName == "" {
		return errors.New("no role name specified")
	}

	if params.Action == "add" {
		if !validateUserHasRole(params.Session, params.Message.GuildID, params.Message.Author.ID, params.Client.roleAddRoleID) {
			return errors.New("you do not have permission to add roles")
		}

		if params.EmojiName == "" {
			return errors.New("no emoji specified")
		}

		if params.Client.db.RoleIsNameTaken(params.RoleName) {
			return errors.New("role name is already taken")
		}

		if params.Client.db.RoleIsEmojiTaken(params.EmojiName) {
			return errors.New("emoji is already taken")
		}
	}

	if params.Action == "remove" {
		if !validateUserHasRole(params.Session, params.Message.GuildID, params.Message.Author.ID, params.Client.roleRemoveRoleID) {
			return errors.New("you do not have permission to remove roles")
		}

		if !params.Client.db.RoleIsNameTaken(params.RoleName) {
			return errors.New("role name does not exist")
		}
	}

	return nil
}

func handleAddAction(params RoleCommandParams) {
	role, roleCreateErr := params.Session.GuildRoleCreate(params.Message.GuildID)

	if roleCreateErr != nil {
		params.Session.ChannelMessageSendReply(params.Message.ChannelID, "Error creating role", params.Message.Reference())
		println(roleCreateErr.Error())
		return
	}

	_, editErr := params.Session.GuildRoleEdit(params.Message.GuildID, role.ID, params.RoleName, params.Color, false, 0, true)
	if editErr != nil {
		params.Session.ChannelMessageSendReply(params.Message.ChannelID, "Error editing role", params.Message.Reference())
		println(editErr.Error())
		return
	}

	params.Client.db.RoleAdd(role.ID, params.EmojiName, params.RoleName)

	reactAddErr := params.Session.MessageReactionAdd(params.Client.roleMessage.ChannelID, params.Client.roleMessage.ID, params.EmojiName)
	if reactAddErr != nil {
		params.Session.ChannelMessageSendReply(params.Message.ChannelID, "Error adding reaction", params.Message.Reference())
		println(reactAddErr.Error())
		return
	}
}

func handleRemoveAction(params RoleCommandParams) {
	id := params.Client.db.RoleGetIdByName(params.RoleName)
	role := params.Client.db.RoleGetById(id)

	deleteErr := params.Session.GuildRoleDelete(params.Message.GuildID, id)
	if deleteErr != nil {
		params.Session.ChannelMessageSendReply(params.Message.ChannelID, "Error deleting role", params.Message.Reference())
		println(deleteErr.Error())
		return
	}

	reactRemoveErr := params.Session.MessageReactionsRemoveEmoji(params.Client.roleMessage.ChannelID, params.Client.roleMessage.ID, role.Emoji)
	if reactRemoveErr != nil {
		params.Session.ChannelMessageSendReply(params.Message.ChannelID, "Error removing reaction", params.Message.Reference())
		println(reactRemoveErr.Error())
		return
	}

	params.Client.db.RoleRemove(id)
}
