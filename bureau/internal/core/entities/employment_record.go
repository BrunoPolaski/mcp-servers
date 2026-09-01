package entities

import (
	"time"

	"gorm.io/gorm"
)

type EmploymentRecord struct {
	gorm.Model
	PersonID           uint    `gorm:"not null;index"`
	EmployerName       string  `gorm:"size:255;not null"`
	EmployerDocument   *string `gorm:"size:14"` // CNPJ
	JobTitle           *string `gorm:"size:255"`
	EmploymentType     *string `gorm:"size:50"` // CLT, PJ, autonomous, unemployed
	Salary             *float64
	StartDate          *time.Time
	EndDate            *time.Time
	IsCurrent          bool    `gorm:"default:false;index"`
	VerificationStatus string  `gorm:"size:50;default:'unverified'"`
	DataSource         *string `gorm:"size:100"` // CAGED, RAIS, self-declared
}
