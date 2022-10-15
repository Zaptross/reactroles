package dgclient

import (
	"errors"
	"fmt"
	"strings"
)

type validHelpParams struct {
	Action string
}

func helpRoleHelp() string {
	return `Usage: !role help <action>
!role help add
!role help remove`
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

// !role help <0: action>
func handleHelpAction(params RoleCommandParams) {
	helpParams, err := validateParamsForHelp(params)

	if err != nil {
		params.Session.ChannelMessageSendReply(params.Message.ChannelID, fmt.Sprintf("⚠ Error: %s\n%s", "No such action exists, available actions are:", strings.Join(Actions.All(), ", ")), params.Message.Reference())
	}

	switch helpParams.Action {
	case Actions.Add:
		params.Session.ChannelMessageSendReply(params.Message.ChannelID, addRoleHelp(), params.Message.Reference())
	case Actions.Remove:
		params.Session.ChannelMessageSendReply(params.Message.ChannelID, removeRoleHelp(), params.Message.Reference())
	default:
		params.Session.ChannelMessageSendReply(params.Message.ChannelID, helpRoleHelp(), params.Message.Reference())
	}
}
