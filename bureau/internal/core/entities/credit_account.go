package entities

import (
	"time"

	"gorm.io/gorm"
)

type CreditAccount struct {
	gorm.Model
	PersonID         uint    `gorm:"not null;index"`
	AccountType      string  `gorm:"size:50;not null;index"` // credit_card, loan, mortgage, auto_loan, overdraft
	Creditor         *string `gorm:"size:255"`
	CreditorDocument *string `gorm:"size:14"` // CNPJ
	AccountNumber    *string `gorm:"size:100"`

	// Account Details
	OpenedDate *time.Time
	ClosedDate *time.Time
	Status     *string `gorm:"size:50;index"` // active, closed, charged_off, defaulted

	// Credit Limits & Balances
	CreditLimit     *float64 `gorm:"index"`
	CurrentBalance  *float64 `gorm:"index"`
	AvailableCredit *float64

	// Loan Specific
	OriginalAmount    *float64
	RemainingAmount   *float64
	InterestRate      *float64
	MonthlyPayment    *float64
	PaymentDueDay     *int
	NumberOfPayments  *int
	RemainingPayments *int

	// Performance
	PaymentStatus   *string `gorm:"size:50;index"` // current, late, default
	DaysLate        int     `gorm:"default:0;index"`
	HighestDaysLate int     `gorm:"default:0"`
	TimesLate30Days int     `gorm:"default:0"`
	TimesLate60Days int     `gorm:"default:0"`
	TimesLate90Days int     `gorm:"default:0"`

	LastReportedDate *time.Time
}
