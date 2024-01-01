package dgclient

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/reactroles/internal/pgdb"
)

type validUpdateParams struct {
	RoleName             string
	RoleField            string
	RoleFieldValueString string
	RoleFieldValueInt    int
}

func updateRoleHelp() string {
	return `Usage: /role update <role> <field> <value>
/role update valorant color #ff0000
/role update valorant name vlorarronat
/role update valorant emoji 🎮`
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

	if roleField == "name" && isRoleTaken(params, roleFieldValueString) != nil {
		return updateParams, errors.New("role name already taken")
	}

	if roleField == "emoji" && isEmojiTaken(params, roleFieldValueString) != nil {
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
		params.Reply(fmt.Sprintf("⚠ Error: %s\n%s", err.Error(), updateRoleHelp()))
		return
	}

	role := params.Client.db.GetRoleByName(updateParams.RoleName)

	if role.Name == "" {
		params.Reply(fmt.Sprintf("⚠ Error: %s %s %s", "Role", updateParams.RoleName, "not found"))
		return
	}

	switch updateParams.RoleField {
	case "color":
		_, editErr := params.Session.GuildRoleEdit(params.GuildID(), role.ID, role.Name, updateParams.RoleFieldValueInt, false, 0, true)
		if editErr != nil {
			params.Reply("Error updating role")
			println(editErr.Error())
			return
		}
	case "name":
		params.Client.db.RoleUpdate(role.ID, role.Emoji, updateParams.RoleFieldValueString, params.GuildID())

		guildRoles, grErr := params.Session.GuildRoles(params.GuildID())
		if grErr != nil {
			params.Reply("Error getting Discord roles")
		}

		for _, guildRole := range guildRoles {
			if guildRole.ID == role.ID {
				_, editErr := params.Session.GuildRoleEdit(params.GuildID(), role.ID, updateParams.RoleFieldValueString, guildRole.Color, false, 0, true)
				if editErr != nil {
					params.Reply("Error updating role")
					println(editErr.Error())
					return
				}
			}
		}
	case "emoji":
		selectorForRole, err := findSelectorForRole(lookupMessagesForSelectors(params.Client, params.Client.db.SelectorGetAll(params.GuildID())), role)

		if err != nil {
			params.Reply("Error finding selector for role")
			println(err.Error())
			return
		}

		reactRemoveErr := params.Session.MessageReactionsRemoveEmoji(params.Server.SelectorChannelID, selectorForRole.ID, role.Emoji)
		if reactRemoveErr != nil {
			params.Reply("Error removing reaction")
			println(reactRemoveErr.Error())
			return
		}

		reactAddErr := params.Session.MessageReactionAdd(params.Server.SelectorChannelID, selectorForRole.ID, updateParams.RoleFieldValueString)
		if reactAddErr != nil {
			params.Reply("Error adding reaction")
			println(reactAddErr.Error())
			return
		}

		params.Client.db.RoleUpdate(role.ID, updateParams.RoleFieldValueString, role.Name, params.GuildID())
	}

	params.Reply(fmt.Sprintf("Successfully updated role %s's %s to %s", role.Name, updateParams.RoleField, updateParams.RoleFieldValueString))
}

func updateRoleSlashCommand() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Name:        Actions.Update,
		Description: "Update any part of a role (name, emoji, color)",
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "role",
				Description: "The role to update.",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "field",
				Description: "The field to update.",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{Name: "name", Value: "name"},
					{Name: "emoji", Value: "emoji"},
					{Name: "color", Value: "color"},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "value",
				Description: "The value to update the field to.",
				Required:    true,
			},
		},
	}
}

func handleUpdateRoleSlashCommand(client *DiscordGoClient, s *discordgo.Session, i *discordgo.InteractionCreate, server *pgdb.ServerConfiguration) {
	sc := i.ApplicationCommandData().Options[0]
	role := sc.Options[0].RoleValue(s, i.GuildID)
	field := sc.Options[1].StringValue()
	value := sc.Options[2].StringValue()

	params := RoleCommandParams{
		Server:      server,
		Session:     s,
		Interaction: i,
		Client:      client,
		Action:      Actions.Update,
		Rest:        []string{role.Name, field, value},
	}

	handleUpdateAction(params)
}
