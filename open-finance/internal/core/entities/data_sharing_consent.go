package entities

import (
	"time"

	"gorm.io/gorm"
)

// DataSharingConsent represents the Open Finance consent that authorizes the
// institution to access the person's shared data.
type DataSharingConsent struct {
	gorm.Model
	PersonID    uint   `gorm:"not null;index"`
	ConsentID   string `gorm:"size:100;uniqueIndex;not null"` // External Open Finance consent id
	Institution string `gorm:"size:255;not null"`
	Status      string `gorm:"size:50;not null;index"` // granted, revoked, expired, awaiting
	Scope       string `gorm:"type:json"`              // JSON array of permissions

	GrantedAt time.Time `gorm:"not null"`
	ExpiresAt *time.Time
	RevokedAt *time.Time
}
