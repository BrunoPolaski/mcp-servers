package entities

import (
	"time"

	"gorm.io/gorm"
)

// BankAccountProfile aggregates the banking-behavior view a person exposes
// through Open Finance data sharing.
type BankAccountProfile struct {
	gorm.Model
	PersonID    uint      `gorm:"not null;uniqueIndex:idx_person_bank_profile_date"`
	ProfileDate time.Time `gorm:"not null;uniqueIndex:idx_person_bank_profile_date"`

	// Banking Behavior (from Open Finance)
	BankingRelationships int  // Number of active bank relationships
	AccountAgeAverage    *int // Average age in months
	HasCheckingAccount   bool
	HasSavingsAccount    bool
	HasInvestmentAccount bool
	InvestmentsValue     *float64 // Total invested balance visible via Open Finance
}
