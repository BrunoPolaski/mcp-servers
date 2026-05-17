package entities

import (
	"time"

	"gorm.io/gorm"
)

type FinancialProfile struct {
	gorm.Model
	PersonID    uint      `gorm:"not null;uniqueIndex:idx_person_profile_date"`
	ProfileDate time.Time `gorm:"not null;uniqueIndex:idx_person_profile_date"`

	// Income
	DeclaredMonthlyIncome  *float64 `gorm:"index"`
	EstimatedMonthlyIncome *float64 `gorm:"index"`
	IncomeSource           *string  `gorm:"size:100"`

	// Assets
	TotalAssets      *float64
	RealEstateValue  *float64
	VehiclesValue    *float64
	InvestmentsValue *float64

	// Liabilities
	TotalLiabilities     *float64 `gorm:"index"`
	TotalMonthlyPayments *float64

	// Financial Health Indicators
	DebtToIncomeRatio *float64 `gorm:"index"`
	AvailableCredit   *float64
	CreditUtilization *float64 `gorm:"index"` // Percentage of available credit used

	// Banking Behavior (from Open Finance)
	BankingRelationships int  // Number of active bank relationships
	AccountAgeAverage    *int // Average age in months
	HasCheckingAccount   bool
	HasSavingsAccount    bool
	HasInvestmentAccount bool
}
