package dgclient

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/reactroles/internal/pgdb"
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

	if !validateUserHasRole(params.Session, params.GuildID(), params.AuthorID(), params.Server.RoleRemoveRoleID) {
		return removeRoleParams, errors.New("you do not have permission to remove roles")
	}

	if !params.Client.db.RoleIsNameTaken(RoleName, params.GuildID()) {
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
	role := params.Client.db.RoleGetById(id, params.GuildID())

	deleteErr := params.Session.GuildRoleDelete(params.GuildID(), id)
	if deleteErr != nil {
		params.Reply("Error deleting role")
		println(deleteErr.Error())
		return
	}

	selectors := params.Client.db.SelectorGetAll(params.GuildID())
	selectorForRole, err := findSelectorForRole(lookupMessagesForSelectors(params.Client, selectors), role)

	if err != nil {
		params.Reply("Error finding selector for role")
		println(err.Error())
		return
	}

	reactRemoveErr := params.Session.MessageReactionsRemoveEmoji(params.Server.SelectorChannelID, selectorForRole.ID, role.Emoji)
	if reactRemoveErr != nil {
		println(reactRemoveErr.Error())
	}

	params.Client.db.RoleRemove(id, params.GuildID())
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

func handleRemoveRoleSlashCommand(client *DiscordGoClient, s *discordgo.Session, i *discordgo.InteractionCreate, server *pgdb.ServerConfiguration) {
	sc := i.ApplicationCommandData().Options[0]
	role := sc.Options[0].RoleValue(s, i.GuildID)

	params := RoleCommandParams{
		Server:      server,
		Session:     s,
		Interaction: i,
		Rest:        []string{role.Name},
		Client:      client,
		Action:      Actions.Remove,
	}

	handleRemoveAction(params)
}
