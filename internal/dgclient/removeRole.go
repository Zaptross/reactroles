package dgclient

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type validRemoveRoleParams struct {
	Name string
}

func removeRoleHelp() string {
	return `Usage: /role remove @<role>
	
	Examples:
	/role remove @valorant`
}

// !role remove <0: role name>
func validateParamsForRemove(params RoleCommandParams) (validRemoveRoleParams, error) {
	removeRoleParams := validRemoveRoleParams{}
	RoleName := params.Rest[0]

	if !validateUserHasRole(params.Session, params.GuildID(), params.AuthorID(), params.Client.roleAddRoleID) {
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
		params.Reply(fmt.Sprintf("⚠ Error: %s\n%s", err, removeRoleHelp()))
		return
	}

	id := params.Client.db.RoleGetIdByName(removeRoleParams.Name)
	role := params.Client.db.RoleGetById(id)

	deleteErr := params.Session.GuildRoleDelete(params.GuildID(), id)
	if deleteErr != nil {
		params.Reply("Error deleting role")
		println(deleteErr.Error())
		return
	}

	selectorForRole, err := findSelectorForRole(params.Client.selectors, role)

	if err != nil {
		params.Reply("Error finding selector for role")
		println(err.Error())
		return
	}

	reactRemoveErr := params.Session.MessageReactionsRemoveEmoji(params.Client.RoleChannel, selectorForRole.ID, role.Emoji)
	if reactRemoveErr != nil {
		params.Reply("Error removing reaction")
		println(reactRemoveErr.Error())
		return
	}

	params.Client.db.RoleRemove(id)
	params.Reply(fmt.Sprintf("Removed role %s %s", removeRoleParams.Name, role.Emoji))
}

func removeRoleSlashCommand() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Name:        Actions.Remove,
		Description: "Remove a role",
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "role",
				Description: "The role to remove",
				Required:    true,
			},
		},
	}
}

func handleRemoveRoleSlashCommand(client *DiscordGoClient, s *discordgo.Session, i *discordgo.InteractionCreate) {
	sc := i.ApplicationCommandData().Options[0]
	role := sc.Options[0].RoleValue(s, i.GuildID)

	params := RoleCommandParams{
		Session:     s,
		Interaction: i,
		Rest:        []string{role.Name},
		Client:      client,
		Action:      Actions.Remove,
	}

	handleRemoveAction(params)
}
