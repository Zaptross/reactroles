package pgdb

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type ServerConfiguration struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	GuildID   string `gorm:"primarykey"`

	// Whether Channel Creation is enabled.
	ChannelCreation bool
	// The channel ID where the bot will listen for role reactions and send role
	SelectorChannelID string
	// The category ID where the bot will create channels for roles.
	ChannelCategoryID string
	// Whether to delete channels when a role is removed.
	ChannelCascadeDelete bool
	// Notify role ID is the role ID that will be notified when a new role is added.
	NotifyRoleID string

	/// Permissions
	/// These are the role IDs that the bot will use to determine if a user has
	/// permission to perform certain actions.

	// RoleAddRoleID is the role ID that allows a user to add roles to themselves.
	RoleAddRoleID string
	// RoleRemoveRoleID is the role ID that allows a user to remove roles from themselves.
	RoleRemoveRoleID string
	// RoleUpdateRoleID is the role ID that allows a user to update roles for themselves.
	RoleUpdateRoleID string
	// The role ID that allows a user to create channels for roles.
	ChannelCreateRoleID string
	// The role ID that allows a user to remove channels for roles.
	ChannelRemoveRoleID string
}

func (db *ReactRolesDatabase) GetAllServerConfigurations() []ServerConfiguration {
	var configs []ServerConfiguration
	db.DB.Find(&configs)
	return configs
}

func (db *ReactRolesDatabase) ServerConfigurationGet(guildId string) *ServerConfiguration {
	var config ServerConfiguration
	db.DB.Where("guild_id = ?", guildId).First(&config)
	return &config
}

func (db *ReactRolesDatabase) ServerConfigurationCreate(
	guildId string,
	addRole string,
	removeRole string,
	updateRole string,
	selectorChannel string,
	channelCreation bool,
	channelCreateRoleID string,
	channelRemoveRoleID string,
	notifyRoleID string,
	channelCategoryID string,
	channelCascadeDelete bool,
) *ServerConfiguration {
	config := &ServerConfiguration{
		GuildID:              guildId,
		RoleAddRoleID:        addRole,
		RoleRemoveRoleID:     removeRole,
		RoleUpdateRoleID:     updateRole,
		SelectorChannelID:    selectorChannel,
		ChannelCreation:      channelCreation,
		ChannelCreateRoleID:  channelCreateRoleID,
		ChannelRemoveRoleID:  channelRemoveRoleID,
		NotifyRoleID:         notifyRoleID,
		ChannelCategoryID:    channelCategoryID,
		ChannelCascadeDelete: channelCascadeDelete,
	}

	db.DB.Create(config)

	return config
}

func (db *ReactRolesDatabase) ServerConfigurationUpdate(guildId string,
	addRole string,
	removeRole string,
	updateRole string,
	selectorChannel string,
	channelCreation bool,
	channelCreateRoleID string,
	channelRemoveRoleID string,
	notifyRoleID string,
	channelCategoryID string,
	channelCascadeDelete bool,
) *ServerConfiguration {
	config := &ServerConfiguration{GuildID: guildId}
	db.DB.Model(config).Where("guild_id = ?", guildId).Updates(ServerConfiguration{
		RoleAddRoleID:        addRole,
		RoleRemoveRoleID:     removeRole,
		RoleUpdateRoleID:     updateRole,
		SelectorChannelID:    selectorChannel,
		ChannelCreation:      channelCreation,
		ChannelCreateRoleID:  channelCreateRoleID,
		ChannelRemoveRoleID:  channelRemoveRoleID,
		NotifyRoleID:         notifyRoleID,
		ChannelCategoryID:    channelCategoryID,
		ChannelCascadeDelete: channelCascadeDelete,
	})

	return config
}

func (sc *ServerConfiguration) Diff(other *ServerConfiguration) string {
	out := []string{}

	for _, field := range reflect.VisibleFields(reflect.TypeOf(*sc)) {
		if field.Type == reflect.TypeOf(time.Time{}) {
			continue
		}

		thisField := reflect.ValueOf(*sc).FieldByName(field.Name).String()
		otherField := reflect.ValueOf(*other).FieldByName(field.Name).String()

		if field.Type == reflect.TypeOf(true) {
			thisField = strconv.FormatBool(reflect.ValueOf(*sc).FieldByName(field.Name).Bool())
			otherField = strconv.FormatBool(reflect.ValueOf(*other).FieldByName(field.Name).Bool())
		}

		if thisField != otherField {
			out = append(out, fmt.Sprintf("%s: %s -> %s", field.Name, thisField, otherField))
		}
	}

	if len(out) == 0 {
		return "updated configuration: no changes"
	}

	return fmt.Sprintf("updated configuration:\n%s", strings.Join(out, "\n"))
}

func (sc *ServerConfiguration) Clone() *ServerConfiguration {
	return &ServerConfiguration{
		GuildID:              sc.GuildID,
		RoleAddRoleID:        sc.RoleAddRoleID,
		RoleRemoveRoleID:     sc.RoleRemoveRoleID,
		RoleUpdateRoleID:     sc.RoleUpdateRoleID,
		SelectorChannelID:    sc.SelectorChannelID,
		ChannelCreation:      sc.ChannelCreation,
		ChannelCreateRoleID:  sc.ChannelCreateRoleID,
		ChannelRemoveRoleID:  sc.ChannelRemoveRoleID,
		NotifyRoleID:         sc.NotifyRoleID,
		ChannelCategoryID:    sc.ChannelCategoryID,
		ChannelCascadeDelete: sc.ChannelCascadeDelete,
	}
}
