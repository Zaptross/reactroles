package dgclient

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/reactroles/internal/pgdb"
)

type validAddRoleParams struct {
	Name  string
	Emoji string
	Color int
}

func addRoleHelp() string {
	return `Usage: /role add <role name> <emoji> [colour hex]
Note: [] is optional

Examples:
/role add valorant :gun: #d34454
/role add valorant :gun:`
}

// !role add <0: role name> <1: emoji> [2: color]
func validateParamsForAdd(params RoleCommandParams) (validAddRoleParams, error) {
	var err error
	addRoleParams := validAddRoleParams{}
	RoleName := params.Rest[0]
	RoleEmoji := params.Rest[1]
	RoleColor := params.Rest[2]

	if !validateUserHasRole(params.Session, params.GuildID(), params.AuthorID(), params.Server.RoleAddRoleID) {
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
		params.Reply(fmt.Sprintf("⚠ Error: %s\n%s", err, addRoleHelp()))
		return
	}

	role, roleCreateErr := params.Session.GuildRoleCreate(params.GuildID())

	if roleCreateErr != nil {
		params.Reply("Error creating role")
		println(roleCreateErr.Error())
		return
	}

	_, editErr := params.Session.GuildRoleEdit(params.GuildID(), role.ID, addRoleParams.Name, addRoleParams.Color, false, 0, true)
	if editErr != nil {
		params.Reply("Error editing role")
		println(editErr.Error())
		return
	}

	params.Client.db.RoleAdd(role.ID, addRoleParams.Emoji, addRoleParams.Name, params.GuildID())
	params.Reply(fmt.Sprintf("Role %s %s added", role.Mention(), addRoleParams.Emoji))

	rolesCount := params.Client.db.RoleGetCount(params.GuildID())

	if rolesCount%ROLES_PER_SELECTOR > 0 {
		// update selectors early if we need a new one
		params.Client.updateRoleSelectorMessage(params.GuildID())
	}

	selectors := params.Client.db.SelectorGetAll(params.GuildID())
	lastSelectorId := selectors[len(selectors)-1].ID
	reactAddErr := params.Session.MessageReactionAdd(params.Server.SelectorChannelID, lastSelectorId, addRoleParams.Emoji)
	if reactAddErr != nil {
		params.Reply("Error adding reaction")
		println(reactAddErr.Error())
		return
	}
}

func addRoleSlashCommand() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Name:        Actions.Add,
		Description: "Create a new role associated with an emoji.",
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "The name of the new role.",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "emoji",
				Description: "The emoji to associate with the new role.",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "color",
				Description: "The color of the new role.",
				Required:    false,
			},
		},
	}
}

func handleAddRoleSlashCommand(client *DiscordGoClient, s *discordgo.Session, i *discordgo.InteractionCreate, server *pgdb.ServerConfiguration) {
	sc := i.ApplicationCommandData().Options[0]
	roleName := sc.Options[0].StringValue()
	roleEmoji := sc.Options[1].StringValue()
	roleColor := ""
	if len(sc.Options) >= 3 {
		roleColor = sc.Options[2].StringValue()
	}

	params := RoleCommandParams{
		Server:      server,
		Session:     s,
		Interaction: i,
		Rest:        []string{roleName, roleEmoji, roleColor},
		Client:      client,
		Action:      Actions.Add,
	}

	handleAddAction(params)
}
