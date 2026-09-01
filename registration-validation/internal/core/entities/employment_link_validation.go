package entities

import (
	"time"

	"gorm.io/gorm"
)

// EmploymentLinkValidation represents an employment relationship validated
// against eSocial.
type EmploymentLinkValidation struct {
	gorm.Model
	PersonID       uint      `gorm:"not null;index"`
	ValidationDate time.Time `gorm:"not null;index"`

	EmployerName     string `gorm:"size:255;not null"`
	EmployerDocument string `gorm:"size:14"`                // CNPJ
	EmploymentType   string `gorm:"size:50"`                // CLT, estatutario, temporary
	Status           string `gorm:"size:50;not null;index"` // active, terminated

	StartDate time.Time `gorm:"not null"`
	EndDate   *time.Time

	Source   string `gorm:"size:100;default:'eSocial'"`
	Verified bool   `gorm:"default:false;index"`
}
