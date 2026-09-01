package entities

import (
	"time"

	"gorm.io/gorm"
)

type Person struct {
	gorm.Model

	// Core Identity
	PersonalInformationID uint `gorm:"not null;index"`
	PersonalInformation   *PersonalInformation

	// Open Finance Domain
	BankAccountProfile   *BankAccountProfile
	BankAccountProfileID *uint `gorm:"index"`

	BankStatements        []BankStatement        `gorm:"foreignKey:PersonID"`
	CashFlowAnalyses      []CashFlowAnalysis     `gorm:"foreignKey:PersonID"`
	RecurringTransactions []RecurringTransaction `gorm:"foreignKey:PersonID"`
	DataSharingConsents   []DataSharingConsent   `gorm:"foreignKey:PersonID"`

	// Audit
	LastVerifiedAt *time.Time
}

func (Person) TableName() string {
	return "persons"
}
