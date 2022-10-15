package dgclient

import (
	"errors"
	"fmt"
)

type validRemoveRoleParams struct {
	Name string
}

func removeRoleHelp() string {
	return `Usage: !role remove <role name>
	
	Examples:
	!role remove valorant`
}

// !role remove <0: role name>
func validateParamsForRemove(params RoleCommandParams) (validRemoveRoleParams, error) {
	removeRoleParams := validRemoveRoleParams{}
	RoleName := params.Rest[0]

	if !validateUserHasRole(params.Session, params.Message.GuildID, params.Message.Author.ID, params.Client.roleAddRoleID) {
		return removeRoleParams, errors.New("you do not have permission to remove roles")
	}

	if !params.Client.db.RoleIsNameTaken(RoleName) {
		return removeRoleParams, errors.New("no role with that name exists")
	}

	return validRemoveRoleParams{
		Name: RoleName,
	}, nil
}

func handleRemoveAction(params RoleCommandParams) {
	removeRoleParams, err := validateParamsForRemove(params)

	if err != nil {
		params.Session.ChannelMessageSendReply(params.Message.ChannelID, fmt.Sprintf("⚠ Error: %s\n%s", err, removeRoleHelp()), params.Message.Reference())
		return
	}

	id := params.Client.db.RoleGetIdByName(removeRoleParams.Name)
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
