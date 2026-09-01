package entities

import (
	"time"

	"gorm.io/gorm"
)

type PersonDocument struct {
	gorm.Model
	PersonalInformationID uint `gorm:"not null;index"`
	FileID                uint `gorm:"not null"`
	File                  *File
	DocumentType          string `gorm:"size:100;not null"` // rg, cpf, proof_residence, income_proof, etc.
	IsVerified            bool   `gorm:"default:false"`
	VerifiedAt            *time.Time
	VerifiedBy            *string `gorm:"size:255"`
	ExpirationDate        *time.Time
}
