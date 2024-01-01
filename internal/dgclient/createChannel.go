package dgclient

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/reactroles/internal/pgdb"
)

func createChannelSlashCommand() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Name:        Actions.CreateChannel,
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Description: "Create a new text or voice channel for a given role.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "role",
				Description: "The role to create a channel for.",
				Required:    true,
			},
			{
				Type:         discordgo.ApplicationCommandOptionChannel,
				Name:         "category",
				Description:  "The category to create the channel in.",
				Required:     true,
				ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildCategory},
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "The name of the channel to create.",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "type",
				Description: "The type of channel to create.",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{Name: "text", Value: "text"},
					{Name: "voice", Value: "voice"},
				},
			},
		},
	}
}

func handleCreateChannelSlashCommand(client *DiscordGoClient, s *discordgo.Session, i *discordgo.InteractionCreate, server *pgdb.ServerConfiguration) {
	sc := i.ApplicationCommandData().Options[0]
	role := sc.Options[0].RoleValue(s, i.GuildID)
	category := sc.Options[1].ChannelValue(s)
	name := sc.Options[2].StringValue()
	channelType := sc.Options[3].StringValue()

	err := validateCreateChannelCommand(client, i.GuildID, role, category, name, channelType, server, i)
	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: fmt.Sprintf("Error: %s", err.Error()),
		})
		return
	}

	createdChannel, err := client.Session.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
		Name:     name,
		Type:     channelTypeFromString(channelType),
		ParentID: category.ID,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				ID:    role.ID,
				Allow: discordgo.PermissionViewChannel,
			},
			{
				ID:   i.GuildID,
				Deny: discordgo.PermissionViewChannel,
			},
		},
	})

	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: fmt.Sprintf("Error creating channel: %s", err.Error()),
		})
		return
	}

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: fmt.Sprintf("Channel %s created for role %s in category %s", createdChannel.Mention(), role.Mention(), category.Mention()),
	})

	if channelType == "text" {
		client.db.RoleUpdateTextChannel(role.ID, createdChannel.ID, i.GuildID)
	} else {
		client.db.RoleUpdateVoiceChannel(role.ID, createdChannel.ID, i.GuildID)
	}
}

func validateCreateChannelCommand(client *DiscordGoClient, guildId string, role *discordgo.Role, category *discordgo.Channel, name string, channelType string, server *pgdb.ServerConfiguration, i *discordgo.InteractionCreate) error {

	if role == nil {
		return fmt.Errorf("no such role exists")
	}

	if category == nil {
		return fmt.Errorf("no such category exists")
	}

	allChannels, err := client.Session.GuildChannels(guildId)

	if err != nil {
		return err
	}

	for _, channel := range allChannels {
		if channel.Name == name {
			return fmt.Errorf("a channel named %s already exists", name)
		}
	}

	if channelType != "text" && channelType != "voice" {
		return fmt.Errorf("invalid channel type: %s", channelType)
	}

	if client.db.RoleGetById(role.ID, guildId).TextChannelID != "" && channelType == "text" {
		return fmt.Errorf("%s already has a text channel", role.Name)
	}

	if client.db.RoleGetById(role.ID, guildId).VoiceChannelID != "" && channelType == "voice" {
		return fmt.Errorf("%s already has a voice channel", role.Name)
	}

	return nil
}

func channelTypeFromString(channelType string) discordgo.ChannelType {
	if channelType == "text" {
		return discordgo.ChannelTypeGuildText
	}

	return discordgo.ChannelTypeGuildVoice
}
