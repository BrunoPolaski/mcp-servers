package entities

import (
	"time"

	valueobjects "github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities/value_objects"
)

type Session struct {
	UUID         string                `gorm:"primaryKey;size:36"`
	ApiKey       string                `gorm:"size:36"`
	UserID       uint                  `gorm:"not null"`
	UserType     valueobjects.UserType `gorm:"size:20;not null"`
	CreatedAt    time.Time
	LastActivity time.Time
}
