package entities

import (
	"time"

	"gorm.io/gorm"
)

// CashFlowAnalysis holds the cash-flow indicators derived from the bank
// statements shared through Open Finance.
type CashFlowAnalysis struct {
	gorm.Model
	PersonID     uint      `gorm:"not null;index"`
	AnalysisDate time.Time `gorm:"not null;index"`
	PeriodDays   int       `gorm:"not null"` // Window analyzed, e.g. 90

	AverageMonthlyInflow  float64 `gorm:"not null"`
	AverageMonthlyOutflow float64 `gorm:"not null"`
	NetCashFlow           float64 `gorm:"not null;index"`
	InflowVolatility      *float64

	NegativeBalanceDays int  `gorm:"default:0"`
	HasRecurringIncome  bool `gorm:"default:false;index"`
}
