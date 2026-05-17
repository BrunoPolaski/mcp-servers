package entities

import (
	"time"

	"gorm.io/gorm"
)

type Debt struct {
	gorm.Model
	PersonID         uint   `gorm:"not null;index"`
	DebtType         string `gorm:"size:100;not null;index"` // credit_card, loan, utility, medical, tax
	Creditor         string `gorm:"size:255;not null;index"`
	CreditorDocument string `gorm:"size:14"`

	OriginalAmount float64 `gorm:"not null"`
	CurrentAmount  float64 `gorm:"not null;index"`
	InterestRate   *float64
	Fees           *float64

	OriginDate time.Time `gorm:"not null"`
	DueDate    time.Time `gorm:"not null;index"`
	Status     string    `gorm:"size:50;not null;index"` // active, settled, defaulted, in_collection

	// Collection Info
	InCollection     bool `gorm:"default:false;index"`
	CollectionDate   *time.Time
	CollectionAgency *string `gorm:"size:255"`

	SettlementAmount *float64
	SettlementDate   *time.Time
}
