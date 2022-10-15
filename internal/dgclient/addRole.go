package dgclient

import (
	"errors"
	"fmt"
)

type validAddRoleParams struct {
	Name  string
	Emoji string
	Color int
}

func addRoleHelp() string {
	return `Usage: !role add <role name> <emoji> [colour hex]
Note: [] is optional

Examples:
!role add valorant :gun: #d34454
!role add valorant :gun:`
}

// !role add <0: role name> <1: emoji> [2: color]
func validateParamsForAdd(params RoleCommandParams) (validAddRoleParams, error) {
	var err error
	addRoleParams := validAddRoleParams{}
	RoleName := params.Rest[0]
	RoleEmoji := params.Rest[1]
	RoleColor := params.Rest[1]

	if !validateUserHasRole(params.Session, params.Message.GuildID, params.Message.Author.ID, params.Client.roleAddRoleID) {
		return addRoleParams, errors.New("you do not have permission to add roles")
	}

	if RoleEmoji == "" {
		return addRoleParams, errors.New("no emoji specified")
	}

	if !isValidEmoji(RoleEmoji) {
		return addRoleParams, errors.New("invalid emoji specified")
	}

	err = isRoleTaken(params, RoleName)
	if err != nil {
		return addRoleParams, err
	}

	err = isEmojiTaken(params, RoleEmoji)

	if err != nil {
		return addRoleParams, err
	}

	return validAddRoleParams{
		Name:  RoleName,
		Emoji: RoleEmoji,
		Color: parseColorToInt(RoleColor),
	}, nil
}

func handleAddAction(params RoleCommandParams) {
	addRoleParams, err := validateParamsForAdd(params)

	if err != nil {
		params.Session.ChannelMessageSendReply(params.Message.ChannelID, fmt.Sprintf("⚠ Error: %s\n%s", err, addRoleHelp()), params.Message.Reference())
		return
	}

	role, roleCreateErr := params.Session.GuildRoleCreate(params.Message.GuildID)

	if roleCreateErr != nil {
		params.Session.ChannelMessageSendReply(params.Message.ChannelID, "Error creating role", params.Message.Reference())
		println(roleCreateErr.Error())
		return
	}

	_, editErr := params.Session.GuildRoleEdit(params.Message.GuildID, role.ID, addRoleParams.Name, addRoleParams.Color, false, 0, true)
	if editErr != nil {
		params.Session.ChannelMessageSendReply(params.Message.ChannelID, "Error editing role", params.Message.Reference())
		println(editErr.Error())
		return
	}

	params.Client.db.RoleAdd(role.ID, addRoleParams.Emoji, addRoleParams.Name)

	reactAddErr := params.Session.MessageReactionAdd(params.Client.roleMessage.ChannelID, params.Client.roleMessage.ID, addRoleParams.Emoji)
	if reactAddErr != nil {
		params.Session.ChannelMessageSendReply(params.Message.ChannelID, "Error adding reaction", params.Message.Reference())
		println(reactAddErr.Error())
		return
	}
}
