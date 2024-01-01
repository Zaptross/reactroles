package pgdb

import "time"

type ServerConfiguration struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	GuildID   string `gorm:"primarykey"`

	// These are the role IDs that the bot will use to determine if a user has
	// permission to perform certain actions.
	RoleAddRoleID    string
	RoleRemoveRoleID string
	RoleUpdateRoleID string

	// The channel ID where the bot will listen for role reactions and send role
	// selector messages.
	SelectorChannelID string
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

func (db *ReactRolesDatabase) ServerConfigurationCreate(guildId string, addRole string, removeRole string, updateRole string, selectorChannel string) *ServerConfiguration {
	config := &ServerConfiguration{
		GuildID:           guildId,
		RoleAddRoleID:     addRole,
		RoleRemoveRoleID:  removeRole,
		RoleUpdateRoleID:  updateRole,
		SelectorChannelID: selectorChannel,
	}

	db.DB.Create(config)

	return config
}

func (db *ReactRolesDatabase) ServerConfigurationUpdate(guildId string, addRole string, removeRole string, updateRole string, selectorChannel string) {
	db.DB.Model(&ServerConfiguration{}).Where("guild_id = ?", guildId).Updates(ServerConfiguration{
		RoleAddRoleID:     addRole,
		RoleRemoveRoleID:  removeRole,
		RoleUpdateRoleID:  updateRole,
		SelectorChannelID: selectorChannel,
	})
}
