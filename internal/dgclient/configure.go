package dgclient

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
	"github.com/zaptross/reactroles/internal/pgdb"
)

func configureServerSlashCommand() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Name:        ActionConfigure,
		Description: "Configure React Roles for this server.",
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "channel",
				Description: "Role selection channel. (It is recommended to create a new channel for this purpose.)",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "add",
				Description: "The role which gives permission to add roles. (Make a new role for this.)",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "remove",
				Description: "The role which gives permission to remove roles. (Make a new role, or use the same role as the add.)",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "update",
				Description: "The role which gives permission to update roles. (Make a new role, or use the same role as the add.)",
				Required:    true,
			},
		},
	}
}

func handleConfigureSlashCommand(client *DiscordGoClient, s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: ":gear: Working...",
			Flags:   uint64(discordgo.MessageFlagsEphemeral),
		},
	})

	sc := i.ApplicationCommandData().Options[0]
	channel := sc.Options[0].ChannelValue(s)
	addRole := sc.Options[1].RoleValue(s, i.GuildID)
	removeRole := sc.Options[2].RoleValue(s, i.GuildID)
	updateRole := sc.Options[3].RoleValue(s, i.GuildID)

	err := validateServerConfiguration(channel, addRole, removeRole, updateRole)

	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: fmt.Sprintf("Error configuring server: %s", err.Error()),
		})
	}

	server := client.db.ServerConfigurationGet(i.GuildID)
	oldRoles := []string{}

	if server.GuildID != "" {
		client.db.ServerConfigurationUpdate(i.GuildID, addRole.ID, removeRole.ID, updateRole.ID, channel.ID)
		oldRoles = []string{server.SelectorChannelID, server.RoleAddRoleID, server.RoleRemoveRoleID, server.RoleUpdateRoleID}
	} else {
		client.db.ServerConfigurationCreate(i.GuildID, addRole.ID, removeRole.ID, updateRole.ID, channel.ID)
	}

	client.updateRoleSelectorMessage(i.GuildID)

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: fmt.Sprintf(":white_check_mark: React Roles has been %s.", updatedOrConfigured(client.Session, server, channel, []*discordgo.Role{addRole, removeRole, updateRole}, oldRoles, channel.GuildID)),
	})
}

func validateServerConfiguration(channel *discordgo.Channel, addRole *discordgo.Role, removeRole *discordgo.Role, updateRole *discordgo.Role) error {
	if channel.Type != discordgo.ChannelTypeGuildText {
		return errors.New("role channel must be a text channel")
	}

	if addRole.ID == "" || updateRole.ID == "" || removeRole.ID == "" {
		return errors.New("add-role, remove-role, and update-role are required")
	}

	return nil
}

func updatedOrConfigured(session *discordgo.Session, server *pgdb.ServerConfiguration, channel *discordgo.Channel, roles []*discordgo.Role, oldRoles []string, guildId string) string {
	if server.GuildID != "" {
		oldChannel, err := session.Channel(server.SelectorChannelID)

		oldChannelMention := server.SelectorChannelID
		if err == nil {
			oldChannelMention = oldChannel.Mention()
		}

		return fmt.Sprintf("updated from %s, %s, %s, and %s to %s, %s, %s, and %s",
			oldChannelMention,
			getRoleById(session, oldRoles[0], guildId).Name,
			getRoleById(session, oldRoles[1], guildId).Name,
			getRoleById(session, oldRoles[2], guildId).Name,
			channel.Mention(),
			roles[0].Mention(),
			roles[1].Mention(),
			roles[2].Mention(),
		)
	}

	return "configured for this server"
}

func getRoleById(s *discordgo.Session, roleId string, guildId string) *discordgo.Role {
	guild, success := lo.Find(s.State.Guilds, func(guild *discordgo.Guild) bool {
		return guild.ID == guildId
	})

	if !success {
		return nil
	}

	for _, role := range guild.Roles {
		if role.ID == roleId {
			return role
		}
	}

	return nil
}
