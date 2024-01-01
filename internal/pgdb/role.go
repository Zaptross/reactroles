package pgdb

import "time"

type Role struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `gorm:"index"`
	ID        string    `gorm:"primarykey"`
	GuildID   string
	Name      string
	Emoji     string
}

func (db *ReactRolesDatabase) RoleGetIdByEmoji(emoji string) string {
	var role Role
	db.DB.Where("emoji = ?", emoji).First(&role)
	return role.ID
}

func (db *ReactRolesDatabase) RoleGetIdByName(name string) string {
	var role Role
	db.DB.Where("name = ?", name).First(&role)
	return role.ID
}

func (db *ReactRolesDatabase) GetRoleByName(name string) Role {
	var role Role
	db.DB.Where("name = ?", name).First(&role)
	return role
}

func (db *ReactRolesDatabase) RoleAdd(id string, emoji string, name string, guildId string) {
	db.DB.Create(&Role{ID: id, Emoji: emoji, Name: name, GuildID: guildId})
}

func (db *ReactRolesDatabase) RoleUpdate(id string, emoji string, name string, guildId string) {
	db.DB.Model(&Role{}).Where("id = ? AND guild_id = ?", id, guildId).Updates(Role{Emoji: emoji, Name: name})
}

func (db *ReactRolesDatabase) RoleRemove(id string, guildId string) {
	db.DB.Delete(&Role{ID: id, GuildID: guildId})
}

func (db *ReactRolesDatabase) RoleGetById(id string, guildId string) Role {
	var role Role
	db.DB.Where("id = ? AND guild_id = ?", id, guildId).First(&role)
	return role
}

func (db *ReactRolesDatabase) RoleGetAll(guildId string) []Role {
	var roles []Role
	db.DB.Where("guild_id = ?", guildId).Find(&roles)
	return roles
}

func (db *ReactRolesDatabase) RoleIsEmojiTaken(emoji string, guildId string) bool {
	var role Role
	db.DB.Where("emoji = ? AND guild_id = ?", emoji, guildId).First(&role)
	return role.ID != ""
}

func (db *ReactRolesDatabase) RoleIsNameTaken(name string, guildId string) bool {
	var role Role
	db.DB.Where("name = ? AND guild_id = ?", name, guildId).First(&role)
	return role.ID != ""
}

func (db *ReactRolesDatabase) RoleGetCount(guildId string) int {
	var count int64
	db.DB.Model(&Role{}).Where("guild_id = ?", guildId).Count(&count)
	return int(count)
}
