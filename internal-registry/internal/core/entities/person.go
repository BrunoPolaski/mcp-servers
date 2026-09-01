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

	// Internal Registry Domain
	CustomerRelationship   *CustomerRelationship
	CustomerRelationshipID *uint `gorm:"index"`

	ContractedProducts     []ContractedProduct     `gorm:"foreignKey:PersonID"`
	InternalPaymentRecords []InternalPaymentRecord `gorm:"foreignKey:PersonID"`
	PreApprovedLimits      []PreApprovedLimit      `gorm:"foreignKey:PersonID"`
	IncomeDeclarations     []IncomeDeclaration     `gorm:"foreignKey:PersonID"`

	// Audit
	LastVerifiedAt *time.Time
}

func (Person) TableName() string {
	return "persons"
}
