package entities

import (
	"time"

	"gorm.io/gorm"
)

// IncomeDeclaration represents an income declared by the customer during
// onboarding/relationship with the institution, optionally backed by a proof
// document.
type IncomeDeclaration struct {
	gorm.Model
	PersonID        uint      `gorm:"not null;index"`
	DeclarationDate time.Time `gorm:"not null;index"`
	IncomeType      string    `gorm:"size:100;not null"` // salary, business, rental, investment
	MonthlyAmount   float64   `gorm:"not null"`
	YearlyAmount    *float64
	Source          *string `gorm:"size:255"`
	Verified        bool    `gorm:"default:false"`
	VerifiedBy      *string `gorm:"size:100"`
	ProofFileID     *uint
	ProofFile       *File
}
