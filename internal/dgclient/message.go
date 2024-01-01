package dgclient

import (
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (client *DiscordGoClient) GetOnMessageHandler() func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if isRoleCommand(m.Content) {
			action, rest := splitRoleCommand(m.Content)
			roleCommandHandler(RoleCommandParams{
				Session: s,
				Message: m,
				Client:  client,
				Action:  action,
				Rest:    rest,
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

func splitRoleCommand(command string) (string, []string) {
	// Split and slice out the !role part of the command
	split := strings.Split(command, " ")[1:]

	if len(split) == 0 || (split[0] == "help" && len(split) == 1) {
		return Actions.Help, []string{Actions.Help}
	}

	return strings.ToLower(split[0]), split[1:]
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

	switch params.Action {
	case Actions.Add:
		handleAddAction(params)
	case Actions.Update:
		handleUpdateAction(params)
	case Actions.Remove:
		handleRemoveAction(params)
	case Actions.Help:
		handleHelpAction(params)
	default:
		handleHelpAction(params)
	}
}
