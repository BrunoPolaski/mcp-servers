package entities

import (
	"time"

	"gorm.io/gorm"
)

type PersonRelationship struct {
	gorm.Model
	PersonID           uint   `gorm:"not null;index"`
	RelatedPersonID    uint   `gorm:"not null;index"`
	RelationType       string `gorm:"size:100;not null"` // spouse, business_partner, guarantor, co_signer, family
	StartDate          *time.Time
	EndDate            *time.Time
	IsActive           bool   `gorm:"default:true"`
	VerificationStatus string `gorm:"size:50;default:'unverified'"`
}
