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

	// Credit Bureau Specific
	CreditScore   *CreditScore
	CreditScoreID *uint `gorm:"index"`

	// Financial Profile
	FinancialProfile   *FinancialProfile
	FinancialProfileID *uint `gorm:"index"`

	// Employment & Income
	EmploymentRecords  []EmploymentRecord  `gorm:"foreignKey:PersonID"`
	IncomeDeclarations []IncomeDeclaration `gorm:"foreignKey:PersonID"`

	// Credit History
	CreditAccounts   []CreditAccount  `gorm:"foreignKey:PersonID"`
	CreditInquiries  []CreditInquiry  `gorm:"foreignKey:PersonID"`
	PaymentHistories []PaymentHistory `gorm:"foreignKey:PersonID"`

	// Debts & Obligations
	Debts           []Debt           `gorm:"foreignKey:PersonID"`
	NegativeRecords []NegativeRecord `gorm:"foreignKey:PersonID"` // Protests, bounced checks, etc.

	// Legal & Compliance
	LegalRecords     []LegalRecord     `gorm:"foreignKey:PersonID"`
	ComplianceChecks []ComplianceCheck `gorm:"foreignKey:PersonID"`

	// Fraud & Risk
	FraudAlerts     []FraudAlert     `gorm:"foreignKey:PersonID"`
	RiskAssessments []RiskAssessment `gorm:"foreignKey:PersonID"`

	// Relationships
	RelatedPersons []PersonRelationship `gorm:"foreignKey:PersonID"`

	// Audit
	DataSources      []DataSource `gorm:"many2many:person_data_sources;"` // Where data came from
	LastVerifiedAt   *time.Time
	ConsentStatus    string `gorm:"size:50;not null;default:'pending'"` // LGPD compliance
	ConsentGrantedAt *time.Time
}

func (Person) TableName() string {
	return "persons"
}
