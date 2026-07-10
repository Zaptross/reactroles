package dgclient

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
	"github.com/zaptross/reactroles/internal/pgdb"
	"github.com/zaptross/reactroles/internal/utils"
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
	roles := client.db.RolesGetAllOrderByAge(guildId)
	selectors := client.db.SelectorGetAll(guildId)
	s, c := utils.GetVersionRaw()

	preambleSections := []string{
		"# ReactRoles",
		fmt.Sprintf("%s (%s)", s, c),
		"",
		"## Joining Roles",
		"To join a role, react with the corresponding emoji to the message for that role.",
		"To leave a role, remove that reaction.",
		"",
		"## Roles",
	}

	roleCommands := []string{
		"## Creating Roles",
		fmt.Sprintf("To add roles, %s users can use the `/role add` command.", roleMention(server.RoleAddRoleID)),
		fmt.Sprintf("To update roles, %s users can use the `/role update` command.", roleMention(server.RoleUpdateRoleID)),
		fmt.Sprintf("To remove roles %s users can use the `/role remove` command.", roleMention(server.RoleRemoveRoleID)),
		"",
	}

	channelCommands := []string{}
	if server.ChannelCreation {
		channelCommands = []string{
			"## Creating Channels for Roles",
			fmt.Sprintf("To create a text or voice channel for that role, %s users can use the `/role create-channel` command.", roleMention(server.ChannelCreateRoleID)),
			fmt.Sprintf("To remove text or voice channels for roles %s users can use the `/role remove-channel` command.", roleMention(server.ChannelRemoveRoleID)),
			"",
		}
	}

	notifySection := []string{
		"## Notifications",
		fmt.Sprintf("When a new role is added, %s users will be notified.", roleMention(server.NotifyRoleID)),
		fmt.Sprintf("If you'd like to be notified when a new role is added, react with %s to this message.", utils.EMOJI_BELL),
		"",
	}

	preambleSections = append(preambleSections, roleCommands...)
	preambleSections = append(preambleSections, channelCommands...)
	preambleSections = append(preambleSections, notifySection...)

	if len(selectors) == 0 {
		message, err := client.Session.ChannelMessageSend(server.SelectorChannelID, "Setting up role assignment message...")

		if err != nil {
			log.Fatal(err)
		}

		selector := client.db.SelectorCreate(message, guildId, "")
		selectors = append(selectors, *selector)

		log.Printf("[dgclient] Role selector 0 created: %s\n", message.ID)
	}

	preambleSelector, ok := lo.Find(selectors, func(selector pgdb.Selector) bool {
		return selector.RoleID == ""
	})

	if !ok {
		// shouldn't be possible
		log.Println("[dgclient] No preamble selector found")
		return
	}

	preamble := strings.Join(preambleSections, "\n")
	if len(roles) == 0 {
		preamble = preamble + "\n\nNo roles."
	}
	preambleMessage, err := client.Session.ChannelMessageEdit(server.SelectorChannelID, preambleSelector.ID, preamble)
	if err != nil {
		log.Println(err.Error())
	}

	if len(preambleMessage.Reactions) == 0 {
		err = client.Session.MessageReactionAdd(preambleSelector.ChannelID, preambleSelector.ID, utils.EMOJI_BELL)
		if err != nil {
			log.Println(err.Error())
		}
	}

	// check if any selectors need to be deleted
	for _, selector := range selectors {
		if selector.RoleID == "" {
			continue // preamble selector
		}

		// if no role exists for this selector, delete it
		if !lo.ContainsBy(roles, func(role pgdb.Role) bool {
			return role.ID == selector.RoleID
		}) {
			err = client.Session.ChannelMessageDelete(server.SelectorChannelID, selector.ID)

			if err != nil {
				log.Println(err.Error())
			}

			client.db.SelectorDelete(selector.GuildID, selector.ID)
		}
	}

	for _, role := range roles {
		selector, ok := lo.Find(selectors, func(selector pgdb.Selector) bool {
			return selector.RoleID == role.ID
		})

		// if no selector exists for this role, create one
		if !ok {
			message, err := client.Session.ChannelMessageSend(server.SelectorChannelID, formatRoleMessage(role, server.NotifyRoleID))

			if err != nil {
				log.Fatal(err)
			}

			client.db.SelectorCreate(message, guildId, role.ID)
			log.Printf("[dgclient] Role selector created for role %s: %s\n", role.Name, message.ID)

			err = client.Session.MessageReactionAdd(message.ChannelID, message.ID, role.Emoji)
			if err != nil {
				log.Println(err.Error())
			}
		} else {
			// if selector exists, update it
			_, err := client.Session.ChannelMessageEdit(server.SelectorChannelID, selector.ID, formatRoleMessage(role, server.NotifyRoleID))

			if err != nil {
				log.Println(err.Error())
			}
		}
	}
}

func findSelectorForRole(selectors []*discordgo.Message, role pgdb.Role) (*discordgo.Message, error) {
	for _, selector := range selectors {
		if strings.Contains(selector.Content, role.Name) {
			return selector, nil
		}
	}

	return nil, errors.New("no selector found for role")
}

func formatRoleMessage(role pgdb.Role, notifyID string) string {
	channelsAndNotify := []string{}
	if role.TextChannelID != "" {
		channelsAndNotify = append(channelsAndNotify, fmt.Sprintf("%s <#%s>", utils.EMOJI_KEYBOARD, role.TextChannelID))
	}
	if role.VoiceChannelID != "" {
		channelsAndNotify = append(channelsAndNotify, fmt.Sprintf("%s <#%s>", utils.EMOJI_STUDIO_MICROPHONE, role.VoiceChannelID))
	}
	channelsAndNotify = append(channelsAndNotify, roleMention(notifyID))
	return fmt.Sprintf("%s %s %s", role.Emoji, role.Name, strings.Join(channelsAndNotify, " "))
}

func roleMention(roleID string) string {
	return fmt.Sprintf("<@&%s>", roleID)
}
