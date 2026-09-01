package entities

import (
	"time"

	"gorm.io/gorm"
)

// RecurringTransaction captures recurring incomes and fixed expenses
// identified from the shared transaction history.
type RecurringTransaction struct {
	gorm.Model
	PersonID        uint   `gorm:"not null;index"`
	TransactionType string `gorm:"size:20;not null;index"` // income, expense
	Category        string `gorm:"size:100;not null"`      // salary, rent, utility, subscription
	Description     string `gorm:"size:255"`

	Amount    float64 `gorm:"not null"`
	Frequency string  `gorm:"size:50;not null"` // monthly, weekly, biweekly

	Counterparty *string `gorm:"size:255"`

	FirstDetectedDate  time.Time `gorm:"not null"`
	LastOccurrenceDate time.Time `gorm:"not null;index"`
	IsActive           bool      `gorm:"default:true;index"`
}
