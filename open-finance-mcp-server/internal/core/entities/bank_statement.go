package entities

import (
	"time"

	"gorm.io/gorm"
)

// BankStatement represents a bank statement shared via Open Finance,
// typically covering the last 90 days of a given account.
type BankStatement struct {
	gorm.Model
	PersonID            uint   `gorm:"not null;index"`
	Institution         string `gorm:"size:255;not null"`
	InstitutionDocument string `gorm:"size:14"`          // CNPJ
	AccountType         string `gorm:"size:50;not null"` // checking, savings, payment

	PeriodStart time.Time `gorm:"not null;index"`
	PeriodEnd   time.Time `gorm:"not null;index"`

	OpeningBalance float64 `gorm:"not null"`
	ClosingBalance float64 `gorm:"not null"`
	TotalCredits   float64 `gorm:"not null"`
	TotalDebits    float64 `gorm:"not null"`

	TransactionCount int    `gorm:"default:0"`
	Currency         string `gorm:"size:3;default:'BRL'"`
}
