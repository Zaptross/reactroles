package dgclient

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
	"github.com/zaptross/reactroles/internal/pgdb"
)

func linkChannelSlashCommand() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Name:        Actions.LinkChannel,
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Description: "Link a channel to a role.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "role",
				Description: "The role to link.",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "channel",
				Description: "The channel to link.",
				Required:    true,
			},
		},
	}
}

func handleLinkChannelSlashCommand(client *DiscordGoClient, s *discordgo.Session, i *discordgo.InteractionCreate, server *pgdb.ServerConfiguration) {
	sc := i.ApplicationCommandData().Options[0]
	role := sc.Options[0].RoleValue(s, i.GuildID)
	channel := sc.Options[1].ChannelValue(s)

	channelType := ""
	if channel.Type == discordgo.ChannelTypeGuildText {
		channelType = "text"
	} else if channel.Type == discordgo.ChannelTypeGuildVoice {
		channelType = "voice"
	}

	err := validateLinkChannelCommand(client, i.GuildID, channel, role, server, i, channelType)
	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: fmt.Sprintf("Error: %s", err.Error()),
		})
		return
	}

	requiredOverrides := []*discordgo.PermissionOverwrite{}

	if !channelHasRolePermissionConfiguration(channel, role.ID, discordgo.PermissionViewChannel) {
		requiredOverrides = append(requiredOverrides, &discordgo.PermissionOverwrite{
			ID:    role.ID,
			Type:  discordgo.PermissionOverwriteTypeRole,
			Allow: discordgo.PermissionViewChannel,
		})
	}

	if !channelHasRolePermissionConfiguration(channel, role.ID, discordgo.PermissionViewChannel) {
		requiredOverrides = append(requiredOverrides, &discordgo.PermissionOverwrite{
			ID:    s.State.User.ID,
			Type:  discordgo.PermissionOverwriteTypeMember,
			Allow: discordgo.PermissionViewChannel,
		})
	}

	if !channelHasRolePermissionConfiguration(channel, i.GuildID, discordgo.PermissionViewChannel) {
		requiredOverrides = append(requiredOverrides, &discordgo.PermissionOverwrite{
			ID:   i.GuildID,
			Type: discordgo.PermissionOverwriteTypeRole,
			Deny: discordgo.PermissionViewChannel,
		})
	}

	if len(requiredOverrides) > 0 {
		_, err = s.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
			PermissionOverwrites: append(channel.PermissionOverwrites, requiredOverrides...),
		})

		if err != nil {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: "Error: failed to update permissions for channel",
			})
			println(err.Error())
			return
		}
	}

	client.db.RoleLinkChannel(channel.ID, role.ID, i.GuildID, channelType)

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: fmt.Sprintf("Channel %s linked to role %s.", channel.Mention(), role.Mention()),
	})
}

func validateLinkChannelCommand(client *DiscordGoClient, guildId string, channel *discordgo.Channel, role *discordgo.Role, server *pgdb.ServerConfiguration, i *discordgo.InteractionCreate, channelType string) error {
	if !(lo.Contains(i.Member.Roles, server.ChannelCreateRoleID) && lo.Contains(i.Member.Roles, server.ChannelRemoveRoleID)) {
		return errors.New("you must have the channel create and remove roles to link a channel")
	}

	dbRole := client.db.RoleGetById(role.ID, guildId)

	if channelType != "text" && channelType != "voice" {
		return errors.New("channel type must be text or voice")
	}

	if channelType == "voice" && dbRole.VoiceChannelID != "" {
		return errors.New("role is already linked to a voice channel")
	}

	if channelType == "text" && dbRole.TextChannelID != "" {
		return errors.New("role is already linked to a text channel")
	}

	if channel.ParentID != server.ChannelCategoryID {
		return errors.New("channel must be in the configured category")
	}

	return nil
}

func channelHasRolePermissionConfiguration(channel *discordgo.Channel, roleId string, permission int64) bool {
	return lo.ContainsBy(channel.PermissionOverwrites, func(po *discordgo.PermissionOverwrite) bool {
		return po.ID == roleId && po.Allow&permission == permission
	})
}
