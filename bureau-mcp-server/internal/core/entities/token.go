package entities

import "time"

type Token struct {
	UUID      string `gorm:"primaryKey;size:36"`
	ApiKey    string `gorm:"size:36"`
	CreatedAt time.Time
}
