package entities

import (
	"time"

	"gorm.io/gorm"
)

type Address struct {
	gorm.Model
	ZipCode        *string `gorm:"size:16;index"`
	State          *string `gorm:"size:100;index"`
	City           *string `gorm:"size:100;index"`
	Neighborhood   *string `gorm:"size:100"` // Important in Brazil
	Street         *string `gorm:"size:255"`
	Number         *string `gorm:"size:16"`
	Complement     *string `gorm:"size:255"`
	ReferencePoint *string `gorm:"size:255"` // Common in Brazil
	AddressType    *string `gorm:"size:50"`  // residential, commercial, rural

	// Geocoding
	Latitude  *float64
	Longitude *float64

	// Validation
	ValidatedByPost bool `gorm:"default:false"`
	RiskScore       *int `gorm:"index"` // Address fraud risk

	IsCurrent          bool `gorm:"default:true;index"`
	IsCorrespondence   bool `gorm:"default:false"`
	MovedInDate        *time.Time
	MovedOutDate       *time.Time
	VerificationStatus string `gorm:"size:50;default:'unverified'"` // verified, unverified, disputed
}
