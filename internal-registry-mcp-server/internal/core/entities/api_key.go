package entities

import (
	"time"
)

type ApiKey struct {
	UUID      string `gorm:"primaryKey;size:36"`
	Slug      string `gorm:"size:50;uniqueIndex;not null"`
	CreatedAt time.Time
}
