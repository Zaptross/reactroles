package dgclient

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func (client *DiscordGoClient) updateRoleSelectorMessage() {
	roles := client.db.RoleGetAll()

	roleLines := []string{
		"**Role Selector**",
		"To join a role, react with the corresponding emoji to this message.",
		"To leave a role, remove the reaction from this message.\n",
		"**Roles**",
	}

	ver, vErr := getVersionMessageIfPossible()

	if vErr == nil {
		roleLines = append([]string{ver}, roleLines...)
	}

	if len(roles) == 0 {
		roleLines = append(roleLines, "No roles")
	} else {
		for _, role := range roles {
			roleLines = append(roleLines, fmt.Sprintf("%s %s", role.Emoji, role.Name))

			client.Session.MessageReactionAdd(client.roleMessage.ChannelID, client.roleMessage.ID, role.Emoji)
		}
	}

	_, err := client.Session.ChannelMessageEdit(client.roleMessage.ChannelID, client.roleMessage.ID, strings.Join(roleLines, "\n"))

	if err != nil {
		log.Println(err.Error())
	}
}

func getVersionMessageIfPossible() (string, error) {
	dat, err := os.ReadFile("/etc/program-version")

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("V %s", string(dat)), nil
}
