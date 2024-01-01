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

func (db *ReactRolesDatabase) SelectorCreate(message *discordgo.Message, guildId string) *Selector {
	selector := &Selector{
		ID:        message.ID,
		ChannelID: message.ChannelID,
		GuildID:   guildId,
	}

	db.DB.Create(selector)

	return selector
}

func (db *ReactRolesDatabase) SelectorDelete(message *discordgo.Message, guildId string) {
	db.DB.Delete(&Selector{}, "id = ? AND guild_id = ?", message.ID, guildId)
}

func (db *ReactRolesDatabase) SelectorGetAll(guildId string) []Selector {
	var selectors []Selector
	db.DB.Where("guild_id = ?", guildId).Find(&selectors).Order("created_at ASC")
	return selectors
}
