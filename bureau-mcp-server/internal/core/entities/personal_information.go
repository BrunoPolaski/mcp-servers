package entities

import (
	"time"

	valueobjects "github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities/value_objects"
	"gorm.io/gorm"
)

type PersonalInformation struct {
	gorm.Model

	FullName      string    `gorm:"size:255;not null;index"`
	MotherName    string    `gorm:"size:255;index"`
	BirthDate     time.Time `gorm:"not null;index"`
	Gender        *string   `gorm:"size:20"`
	Nationality   string    `gorm:"size:100;not null;default:'Brazilian'"`
	MaritalStatus *string   `gorm:"size:50"`

	Document    valueobjects.Document `gorm:"size:11;uniqueIndex;not null"`
	RG          *string               `gorm:"size:20;index"`
	RGIssuer    *string               `gorm:"size:50"`
	RGIssueDate *time.Time
	VoterID     *string `gorm:"size:20;index"`
	WorkCard    *string `gorm:"size:20"`

	PrimaryPhone     valueobjects.PhoneNumber  `gorm:"size:15;not null;index"`
	SecondaryPhone   *valueobjects.PhoneNumber `gorm:"size:15"`
	Email            string                    `gorm:"size:255;index"`
	AlternativeEmail *string                   `gorm:"size:255"`

	Addresses []PersonAddress `gorm:"foreignKey:PersonalInformationID"`

	ProfilePhotoID *uint
	ProfilePhoto   *File
	Documents      []PersonDocument `gorm:"foreignKey:PersonalInformationID"`

	DocumentValidated    bool   `gorm:"default:false"`
	EmailVerified        bool   `gorm:"default:false"`
	PhoneVerified        bool   `gorm:"default:false"`
	BiometricValidated   bool   `gorm:"default:false"`
	ReceitaFederalStatus string `gorm:"size:50"`
}
