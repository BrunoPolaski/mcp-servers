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

	// Registration Validation Domain
	DocumentValidations       []DocumentValidation       `gorm:"foreignKey:PersonID"`
	FiscalRegularities        []FiscalRegularity         `gorm:"foreignKey:PersonID"`
	EmploymentLinkValidations []EmploymentLinkValidation `gorm:"foreignKey:PersonID"`
	ComplianceChecks          []ComplianceCheck          `gorm:"foreignKey:PersonID"`

	// Audit
	LastVerifiedAt *time.Time
}

func (Person) TableName() string {
	return "persons"
}
