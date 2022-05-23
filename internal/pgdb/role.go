package pgdb

import "time"

type Role struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `gorm:"index"`
	ID        string    `gorm:"primarykey"`
	Emoji     string
}

func (db *ReactRolesDatabase) RoleAdd(id string, emoji string) {
	db.DB.Create(&Role{ID: id, Emoji: emoji})
}

func (db *ReactRolesDatabase) RoleRemove(id string) {
	db.DB.Delete(&Role{ID: id})
}

func (db *ReactRolesDatabase) RoleIsEmojiTaken(emoji string) bool {
	var role Role
	db.DB.Where("emoji = ?", emoji).First(&role)
	return role.ID != ""
}
