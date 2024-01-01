package dgclient

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/reactroles/internal/pgdb"
)

const (
	ROLES_PER_SELECTOR = 20
)

func (client *DiscordGoClient) updateAllRoleSelectorMessages() {
	servers := client.db.GetAllServerConfigurations()
	for i, server := range servers {
		log.Printf("[dgclient] Updating role selector message for server %d/%d...\n", i+1, len(servers))
		client.updateRoleSelectorMessage(server.GuildID)
	}
}

func (client *DiscordGoClient) updateRoleSelectorMessage(guildId string) {
	server := client.db.ServerConfigurationGet(guildId)
	roles := client.db.RoleGetAll(guildId)

	roleLines := []string{
		"**Role Selector**",
		"To join a role, react with the corresponding emoji to this message.",
		"To leave a role, remove the reaction from this message.",
		"If the emoji is missing, you may have to add and/or remove that emoji again.",
		"",
		"**Roles**",
	}

	ver, err := getVersionMessageIfPossible()

	if err == nil {
		roleLines = append([]string{ver}, roleLines...)
	}

	selectors := lookupMessagesForSelectors(client, client.db.SelectorGetAll(guildId))

	if len(selectors) == 0 {
		message, err := client.Session.ChannelMessageSend(server.SelectorChannelID, "Setting up role assignment message...")

		if err != nil {
			log.Fatal(err)
		}

		client.db.SelectorCreate(message, guildId)
		selectors = append(selectors, message)

		log.Printf("[dgclient] Role selector 0 created: %s\n", message.ID)
	}

	if len(roles) == 0 {
		roleLines = append(roleLines, "No roles")

		_, err := client.Session.ChannelMessageEdit(server.SelectorChannelID, selectors[0].ID, strings.Join(roleLines, "\n"))
		if err != nil {
			log.Println(err.Error())
		}
		return
	}

	requiredSelectors := int(math.Ceil(float64(len(roles)) / ROLES_PER_SELECTOR))

	if len(selectors) < requiredSelectors {
		for i := len(selectors); i < requiredSelectors; i++ {
			message, err := client.Session.ChannelMessageSend(server.SelectorChannelID, "Setting up role assignment message...")

			if err != nil {
				log.Fatal(err)
			}

			client.db.SelectorCreate(message, guildId)
			selectors = append(selectors, message)

			log.Printf("[dgclient] Role selector %d created: %s\n", i, message.ID)
		}
	}

	if len(selectors) > requiredSelectors {
		for i := len(selectors); i > requiredSelectors; i-- {
			toDeleteId := selectors[i-1].ID

			err = client.Session.ChannelMessageDelete(server.SelectorChannelID, toDeleteId)

			if err != nil {
				log.Println(err.Error())
			}

			client.db.SelectorDelete(selectors[i-1], guildId)
			selectors = selectors[:i-1]

			log.Printf("[dgclient] Role selector %d deleted: %s\n", i, toDeleteId)
		}
	}

	for i, selector := range selectors {
		for j := i * ROLES_PER_SELECTOR; j < (i+1)*ROLES_PER_SELECTOR && j < len(roles); j++ {
			roleLines = append(roleLines, fmt.Sprintf("%s %s", roles[j].Emoji, roles[j].Name))
		}

		_, err := client.Session.ChannelMessageEdit(selector.ChannelID, selector.ID, strings.Join(roleLines, "\n"))

		if err != nil {
			log.Println(err.Error())
		}

		roleLines = []string{}
	}
}

func getVersionMessageIfPossible() (string, error) {
	dat, err := os.ReadFile("/etc/program-version")

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("V %s", string(dat)), nil
}

func findSelectorForRole(selectors []*discordgo.Message, role pgdb.Role) (*discordgo.Message, error) {
	for _, selector := range selectors {
		if strings.Contains(selector.Content, role.Name) {
			return selector, nil
		}
	}

	return nil, errors.New("no selector found for role")
}
