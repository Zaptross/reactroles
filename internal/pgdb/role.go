package pgdb

import "time"

type Role struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `gorm:"index"`
	ID        string    `gorm:"primarykey"`
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

func (db *ReactRolesDatabase) RoleAdd(id string, emoji string, name string) {
	db.DB.Create(&Role{ID: id, Emoji: emoji, Name: name})
}

func (db *ReactRolesDatabase) RoleRemove(id string) {
	db.DB.Delete(&Role{ID: id})
}

func (db *ReactRolesDatabase) RoleGetById(id string) Role {
	var role Role
	db.DB.Where("id = ?", id).First(&role)
	return role
}

func (db *ReactRolesDatabase) RoleGetAll() []Role {
	var roles []Role
	db.DB.Find(&roles)
	return roles
}

func (db *ReactRolesDatabase) RoleIsEmojiTaken(emoji string) bool {
	var role Role
	db.DB.Where("emoji = ?", emoji).First(&role)
	return role.ID != ""
}

func (db *ReactRolesDatabase) RoleIsNameTaken(name string) bool {
	var role Role
	db.DB.Where("name = ?", name).First(&role)
	return role.ID != ""
}
