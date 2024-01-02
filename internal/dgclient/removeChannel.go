package dgclient

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
	"github.com/zaptross/reactroles/internal/pgdb"
)

func removeChannelSlashCommand() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Name:        Actions.RemoveChannel,
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Description: "Remove a channel for a given role.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "role",
				Description: "The role to remove a channel for.",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "channel",
				Description: "The channel to remove. This channel must be linked to the role.",
				Required:    true,
			},
		},
	}
}

func handleRemoveChannelSlashCommand(client *DiscordGoClient, s *discordgo.Session, i *discordgo.InteractionCreate, server *pgdb.ServerConfiguration) {
	sc := i.ApplicationCommandData().Options[0]
	role := sc.Options[0].RoleValue(s, i.GuildID)
	channel := sc.Options[1].ChannelValue(s)

	channelType := ""
	if channel.Type == discordgo.ChannelTypeGuildText {
		channelType = "text"
	} else if channel.Type == discordgo.ChannelTypeGuildVoice {
		channelType = "voice"
	}

	err := validateRemoveChannelCommand(client, i.GuildID, channel, role, server, i, channelType)
	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: err.Error(),
		})
		return
	}

	_, err = s.ChannelDelete(channel.ID)
	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: "Error: failed to delete channel",
		})
		println(err.Error())
		return
	}

	client.db.RoleChannelRemove(role.ID, i.GuildID, channelType)

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: fmt.Sprintf("Channel %s removed for role %s", channel.Name, role.Name),
	})
}

func validateRemoveChannelCommand(client *DiscordGoClient, guildID string, channel *discordgo.Channel, role *discordgo.Role, server *pgdb.ServerConfiguration, i *discordgo.InteractionCreate, channelType string) error {
	if !lo.Contains(i.Member.Roles, server.ChannelRemoveRoleID) {
		return fmt.Errorf("you do not have permission to remove channels")
	}

	if channelType == "" {
		return errors.New("channel type is not valid")
	}

	if channel.GuildID != guildID {
		return fmt.Errorf("channel %s is not in this server", channel.Name)
	}

	dbRole := client.db.RoleGetById(role.ID, guildID)

	if dbRole.TextChannelID != channel.ID && dbRole.VoiceChannelID != channel.ID {
		return fmt.Errorf("channel %s is not linked to role %s", channel.Mention(), role.Mention())
	}

	return nil
}
