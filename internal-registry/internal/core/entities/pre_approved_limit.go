package entities

import (
	"time"

	"gorm.io/gorm"
)

// PreApprovedLimit represents a pre-approved credit limit granted to the
// customer based on internal credit policies.
type PreApprovedLimit struct {
	gorm.Model
	PersonID       uint    `gorm:"not null;index"`
	ProductType    string  `gorm:"size:50;not null;index"` // credit_card, personal_loan, overdraft
	ApprovedAmount float64 `gorm:"not null"`
	InterestRate   *float64
	CalculatedDate time.Time `gorm:"not null"`
	ValidUntil     time.Time `gorm:"not null;index"`
	PolicyVersion  string    `gorm:"size:50"`
	IsActive       bool      `gorm:"default:true;index"`
}
