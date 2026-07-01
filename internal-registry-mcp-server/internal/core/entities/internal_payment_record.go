package entities

import (
	"time"

	"gorm.io/gorm"
)

// InternalPaymentRecord holds the internal payment/delinquency history of the
// customer for products contracted with the institution.
type InternalPaymentRecord struct {
	gorm.Model
	PersonID            uint  `gorm:"not null;index"`
	ContractedProductID *uint `gorm:"index"`

	ReferenceMonth time.Time `gorm:"not null;index"`
	DueDate        time.Time `gorm:"not null"`
	PaymentDate    *time.Time

	AmountDue  float64 `gorm:"not null"`
	AmountPaid float64 `gorm:"not null"`
	Status     string  `gorm:"size:50;not null;index"` // on_time, late, missed, partial
	DaysLate   int     `gorm:"default:0;index"`
}
