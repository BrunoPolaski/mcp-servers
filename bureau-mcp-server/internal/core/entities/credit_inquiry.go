package entities

import (
	"time"

	"gorm.io/gorm"
)

type CreditInquiry struct {
	gorm.Model
	PersonID         uint      `gorm:"not null;index"`
	InquiryDate      time.Time `gorm:"not null;index"`
	InquiryType      *string   `gorm:"size:50"` // hard, soft
	Creditor         *string   `gorm:"size:255"`
	CreditorDocument *string   `gorm:"size:14"`
	Purpose          *string   `gorm:"size:100"` // credit_card, loan, mortgage, rental, employment
	Amount           *float64  // Amount requested
	Result           *string   `gorm:"size:50"` // approved, denied, pending
}
