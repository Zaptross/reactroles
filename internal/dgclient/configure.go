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
				Description: "The role which gives permission to remove roles. (Make a new role, or use the same role as add.)",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "update",
				Description: "The role which gives permission to update roles. (Make a new role, or use the same role as add.)",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "channel-creation",
				Description: "Is channel creation enabled?",
				Required:    true,
			}, {
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "create-channel",
				Description: "The role which gives permission to create channels. (Make a new role for this.)",
				Required:    true,
			}, {
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "remove-channel",
				Description: "The role which gives permission to remove channels. (Make a new role, or use the same as create.)",
				Required:    true,
			},
			{
				Type:         discordgo.ApplicationCommandOptionChannel,
				Name:         "category",
				Description:  "The category to create channels in.",
				Required:     true,
				ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildCategory},
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "cascade-delete",
				Description: "Should channels be deleted when the role is removed? (If false, channels must be deleted manually.)",
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
	channelCreate := sc.Options[4].BoolValue()
	channelCreateRole := sc.Options[5].RoleValue(s, i.GuildID)
	channelRemoveRole := sc.Options[6].RoleValue(s, i.GuildID)
	channelCategory := sc.Options[7].ChannelValue(s)
	cascadeDelete := sc.Options[8].BoolValue()

	err := validateServerConfiguration(client, i.GuildID, channel, addRole, removeRole, updateRole, channelCreate, channelCreateRole, channelRemoveRole, cascadeDelete, i.Member)

	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: fmt.Sprintf("Error configuring server: %s", err.Error()),
		})
		return
	}

	serverConfig := client.db.ServerConfigurationGet(i.GuildID)
	oldConfig := serverConfig.Clone()

	if serverConfig.GuildID != "" {
		serverConfig = client.db.ServerConfigurationUpdate(i.GuildID, addRole.ID, removeRole.ID, updateRole.ID, channel.ID, channelCreate, channelCreateRole.ID, channelRemoveRole.ID, channelCategory.ID, cascadeDelete)
	} else {
		client.db.ServerConfigurationCreate(i.GuildID, addRole.ID, removeRole.ID, updateRole.ID, channel.ID, channelCreate, channelCreateRole.ID, channelRemoveRole.ID, channelCategory.ID, cascadeDelete)
	}

	client.updateRoleSelectorMessage(i.GuildID)

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: fmt.Sprintf(":white_check_mark: React Roles has %s.", updatedOrConfigured(client.Session, serverConfig, oldConfig)),
	})
}

func validateServerConfiguration(client *DiscordGoClient, guildId string, channel *discordgo.Channel, addRole *discordgo.Role, removeRole *discordgo.Role, updateRole *discordgo.Role, channelCreate bool, channelCreateRole *discordgo.Role, channelRemoveRole *discordgo.Role, cascadeDelete bool, member *discordgo.Member) error {
	if member.Permissions&discordgo.PermissionManageWebhooks != discordgo.PermissionManageWebhooks {
		return errors.New("you must have the Manage Webhooks permission to configure ReactRoles")
	}

	if channel.Type != discordgo.ChannelTypeGuildText {
		return errors.New("role channel must be a text channel")
	}

	if addRole.ID == "" || updateRole.ID == "" || removeRole.ID == "" || channelCreateRole.ID == "" || channelRemoveRole.ID == "" {
		return errors.New("add-role, remove-role, and update-role are required")
	}

	if channelCreate && channelCreateRole.ID == "" {
		return errors.New("create-role is required if channel-creation is enabled")
	}

	if cascadeDelete && channelRemoveRole.ID == "" {
		return errors.New("remove-role is required if cascade-delete is enabled")
	}

	managedRoles := lo.Map(client.db.RoleGetAll(guildId), func(i pgdb.Role, _ int) string { return i.ID })
	for _, role := range []*discordgo.Role{addRole, removeRole, updateRole, channelCreateRole, channelRemoveRole} {
		if lo.Contains(managedRoles, role.ID) {
			return fmt.Errorf("%s is managed by ReactRoles\n Permission roles may not be managed by ReactRoles", role.Mention())
		}
	}

	return nil
}

func updatedOrConfigured(session *discordgo.Session, serverConfig *pgdb.ServerConfiguration, oldConfig *pgdb.ServerConfiguration) string {
	if oldConfig.GuildID != "" {
		return oldConfig.Diff(serverConfig)
	}

	return "been configured for this server"
}
