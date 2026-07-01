package entities

import (
	"time"

	"gorm.io/gorm"
)

type PaymentHistory struct {
	gorm.Model
	PersonID        uint  `gorm:"not null;index"`
	CreditAccountID *uint `gorm:"index"`
	DebtID          *uint `gorm:"index"`

	PaymentDate time.Time `gorm:"not null;index"`
	DueDate     *time.Time
	Amount      float64 `gorm:"not null"`
	AmountDue   *float64
	Status      *string `gorm:"size:50;index"` // on_time, late, missed, partial
	DaysLate    int     `gorm:"default:0;index"`
}
