package pgdb

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type Selector struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `gorm:"index"`
	ID        string    `gorm:"primarykey"`
	ChannelID string
	GuildID   string
}

func (db *ReactRolesDatabase) SelectorCreate(message *discordgo.Message) *Selector {
	selector := &Selector{
		ID:        message.ID,
		ChannelID: message.ChannelID,
		GuildID:   message.GuildID,
	}

	db.DB.Create(selector)

	return selector
}

func (db *ReactRolesDatabase) SelectorDelete(message *discordgo.Message) {
	db.DB.Delete(&Selector{}, "id = ?", message.ID)
}

func (db *ReactRolesDatabase) SelectorGetAll() []Selector {
	var selectors []Selector
	db.DB.Find(&selectors).Order("created_at ASC")
	return selectors
}
