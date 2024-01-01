package dgclient

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
	"github.com/zaptross/reactroles/internal/pgdb"
)

type validHelpParams struct {
	Action string
}

func helpRoleHelp() string {
	return fmt.Sprintf(`Usage: /role help <action>
The available actions are: %s

Examples:
/role help add`,
		strings.Join(Actions.All(), ", "))
}

func validateParamsForHelp(params RoleCommandParams) (validHelpParams, error) {
	helpParams := validHelpParams{}

	restAction := ""
	if len(params.Rest) > 0 {
		restAction = params.Rest[0]
	}

	action := strings.ToLower(restAction)

	if action == "" || !validAction(action) {
		return helpParams, errors.New("invalid action")
	}

	return validHelpParams{
		Action: action,
	}, nil
}

func validAction(action string) bool {
	for _, a := range Actions.All() {
		if a == action {
			return true
		}
	}

	return false
}

// /role help <0: action>
func handleHelpAction(params RoleCommandParams) {
	helpParams, err := validateParamsForHelp(params)

	if err != nil {
		params.Reply(fmt.Sprintf("⚠ Error: %s\n%s", "No such action exists, available actions are:", strings.Join(Actions.All(), ", ")))
	}

	switch helpParams.Action {
	case Actions.Add:
		params.Reply(addRoleHelp())
	case Actions.Update:
		params.Reply(updateRoleHelp())
	case Actions.Remove:
		params.Reply(removeRoleHelp())
	default:
		params.Reply(helpRoleHelp())
	}
}

func helpSlashCommand() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        Actions.Help,
		Description: "Get help for a role command.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "command",
				Description: "The command to get help for.",
				Choices: lo.Map(Actions.All(), func(action string, _ int) *discordgo.ApplicationCommandOptionChoice {
					return &discordgo.ApplicationCommandOptionChoice{Name: action, Value: action}
				}),
			},
		},
	}
}

func handleHelpSlashCommand(client *DiscordGoClient, s *discordgo.Session, i *discordgo.InteractionCreate, server *pgdb.ServerConfiguration) {
	sc := i.ApplicationCommandData().Options[0]
	command := ""
	if len(sc.Options) > 0 {
		command = sc.Options[0].StringValue()
	}

	params := RoleCommandParams{
		Server:      server,
		Session:     s,
		Interaction: i,
		Client:      client,
		Action:      Actions.Help,
		Rest:        []string{command},
	}

	handleHelpAction(params)
}
