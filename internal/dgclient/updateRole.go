package dgclient

import (
	"errors"
	"fmt"
)

type validUpdateParams struct {
	RoleName             string
	RoleField            string
	RoleFieldValueString string
	RoleFieldValueInt    int
}

func updateRoleHelp() string {
	return `Usage: !role update <role> <field> <value>
!role update valorant color #ff0000
!role update valorant name vlorarronat
!role update valorant emoji 🎮`
}

func validateParamsForUpdate(params RoleCommandParams) (validUpdateParams, error) {
	updateParams := validUpdateParams{}

	if len(params.Rest) < 3 {
		return updateParams, errors.New("not enough arguments")
	}

	roleName := params.Rest[0]
	roleField := params.Rest[1]
	roleFieldValueString := params.Rest[2]
	roleFieldValueInt := 0

	if roleField == "name" && isRoleTaken(params, roleName) != nil {
		return updateParams, errors.New("role name already taken")
	}

	if roleField == "emoji" && isEmojiTaken(params, roleName) != nil {
		return updateParams, errors.New("emoji already taken")
	}

	switch roleField {
	case "color":
		roleFieldValueInt = parseColorToInt(roleFieldValueString)
	}

	return validUpdateParams{
		RoleName:             roleName,
		RoleField:            roleField,
		RoleFieldValueString: roleFieldValueString,
		RoleFieldValueInt:    roleFieldValueInt,
	}, nil
}

// !role update <0: role> <1: field> <2: value>
func handleUpdateAction(params RoleCommandParams) {
	updateParams, err := validateParamsForUpdate(params)

	if err != nil {
		params.Session.ChannelMessageSendReply(params.Message.ChannelID, fmt.Sprintf("⚠ Error: %s\n%s", err.Error(), updateRoleHelp()), params.Message.Reference())
		return
	}

	role := params.Client.db.GetRoleByName(updateParams.RoleName)

	if role.Name == "" {
		params.Session.ChannelMessageSendReply(params.Message.ChannelID, fmt.Sprintf("⚠ Error: %s %s %s", "Role", updateParams.RoleName, "not found"), params.Message.Reference())
		return
	}

	switch updateParams.RoleField {
	case "color":
		_, editErr := params.Session.GuildRoleEdit(params.Message.GuildID, role.ID, role.Name, updateParams.RoleFieldValueInt, false, 0, true)
		if editErr != nil {
			params.Session.ChannelMessageSendReply(params.Message.ChannelID, "Error updating role", params.Message.Reference())
			println(editErr.Error())
			return
		}
	case "name":
		params.Client.db.RoleUpdate(role.ID, role.Emoji, updateParams.RoleFieldValueString)

		guildRoles, grErr := params.Session.GuildRoles(params.Message.GuildID)
		if grErr != nil {
			params.Session.ChannelMessageSendReply(params.Message.ChannelID, "Error getting Discord roles", params.Message.Reference())
		}

		for _, guildRole := range guildRoles {
			if guildRole.ID == role.ID {
				_, editErr := params.Session.GuildRoleEdit(params.Message.GuildID, role.ID, updateParams.RoleFieldValueString, guildRole.Color, false, 0, true)
				if editErr != nil {
					params.Session.ChannelMessageSendReply(params.Message.ChannelID, "Error updating role", params.Message.Reference())
					println(editErr.Error())
					return
				}
			}
		}
	case "emoji":
		params.Client.db.RoleUpdate(role.ID, updateParams.RoleFieldValueString, role.Name)
	}

	params.Session.ChannelMessageSendReply(params.Message.ChannelID, fmt.Sprintf("✅ Successfully updated role %s's %s to %s", role.Name, updateParams.RoleField, updateParams.RoleFieldValueString), params.Message.Reference())
}
