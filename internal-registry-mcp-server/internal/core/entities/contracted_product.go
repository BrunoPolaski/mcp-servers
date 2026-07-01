package entities

import (
	"time"

	"gorm.io/gorm"
)

// ContractedProduct represents a product the customer has contracted with the
// institution.
type ContractedProduct struct {
	gorm.Model
	PersonID       uint      `gorm:"not null;index"`
	ProductType    string    `gorm:"size:50;not null;index"` // checking_account, credit_card, loan, insurance, investment
	ProductName    string    `gorm:"size:255;not null"`
	ContractNumber string    `gorm:"size:100;index"`
	ContractedDate time.Time `gorm:"not null"`
	Status         string    `gorm:"size:50;not null;index"` // active, closed, suspended

	Balance      *float64
	MonthlyValue *float64
}
