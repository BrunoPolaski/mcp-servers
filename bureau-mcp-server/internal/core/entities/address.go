package entities

import (
	"time"

	"gorm.io/gorm"
)

type Address struct {
	gorm.Model
	ZipCode        string  `gorm:"size:16;not null;index"`
	State          string  `gorm:"size:100;not null;index"`
	City           string  `gorm:"size:100;not null;index"`
	Neighborhood   string  `gorm:"size:100"` // Important in Brazil
	Street         string  `gorm:"size:255;not null"`
	Number         string  `gorm:"size:16;not null"`
	Complement     *string `gorm:"size:255"`
	ReferencePoint *string `gorm:"size:255"` // Common in Brazil
	AddressType    string  `gorm:"size:50"`  // residential, commercial, rural

	// Geocoding
	Latitude  *float64
	Longitude *float64

	// Validation
	ValidatedByPost bool `gorm:"default:false"`
	RiskScore       *int `gorm:"index"` // Address fraud risk

	IsCurrent          bool      `gorm:"default:true;index"`
	IsCorrespondence   bool      `gorm:"default:false"`
	MovedInDate        time.Time `gorm:"not null"`
	MovedOutDate       *time.Time
	VerificationStatus string `gorm:"size:50;default:'unverified'"` // verified, unverified, disputed
}
