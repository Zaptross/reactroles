package pgdb

import "time"

type Role struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `gorm:"index"`
	ID        string    `gorm:"primarykey"`
	Emoji     string
}
