package entities

import (
	"time"

	"gorm.io/gorm"
)

type IncomeDeclaration struct {
	gorm.Model
	PersonID        uint      `gorm:"not null;index"`
	DeclarationDate time.Time `gorm:"not null;index"`
	IncomeType      string    `gorm:"size:100;not null"` // salary, business, rental, investment
	MonthlyAmount   float64   `gorm:"not null"`
	YearlyAmount    *float64
	Source          string  `gorm:"size:255"`
	Verified        bool    `gorm:"default:false"`
	VerifiedBy      *string `gorm:"size:100"`
	ProofFileID     *uint
	ProofFile       *File
}
