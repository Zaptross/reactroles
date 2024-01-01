package pgdb

import "time"

type ServerConfiguration struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	GuildID   string `gorm:"primarykey"`

	// Permissions
	// These are the role IDs that the bot will use to determine if a user has
	// permission to perform certain actions.
	RoleAddRoleID    string
	RoleRemoveRoleID string
	RoleUpdateRoleID string

	//// Selectors
	// The channel ID where the bot will listen for role reactions and send role
	SelectorChannelID string

	//// Channel Creation
	ChannelCreation      bool
	ChannelCategoryID    string
	ChannelCascadeDelete bool
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
	channelCategoryID string,
	channelCascadeDelete bool,
) {
	db.DB.Model(&ServerConfiguration{}).Where("guild_id = ?", guildId).Updates(ServerConfiguration{
		RoleAddRoleID:        addRole,
		RoleRemoveRoleID:     removeRole,
		RoleUpdateRoleID:     updateRole,
		SelectorChannelID:    selectorChannel,
		ChannelCreation:      channelCreation,
		ChannelCategoryID:    channelCategoryID,
		ChannelCascadeDelete: channelCascadeDelete,
	})
}
