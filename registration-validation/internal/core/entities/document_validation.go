package entities

import (
	"time"

	"gorm.io/gorm"
)

// DocumentValidation holds the result of validating a person's CPF/CNPJ
// against the Receita Federal public base.
type DocumentValidation struct {
	gorm.Model
	PersonID       uint      `gorm:"not null;index"`
	ValidationDate time.Time `gorm:"not null;index"`
	DocumentNumber string    `gorm:"size:14;not null;index"`
	DocumentType   string    `gorm:"size:10;not null"` // cpf, cnpj

	ReceitaFederalStatus string `gorm:"size:50;not null"` // regular, pending, suspended, canceled, deceased
	IsValid              bool   `gorm:"default:false;index"`
	NameMatches          bool   `gorm:"default:false"`
	BirthDateMatches     bool   `gorm:"default:false"`
	BiometricValidated   bool   `gorm:"default:false"`

	Source      string  `gorm:"size:100;default:'receita_federal'"`
	RawResponse *string `gorm:"type:json"`
}
