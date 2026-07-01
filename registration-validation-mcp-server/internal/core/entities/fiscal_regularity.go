package entities

import (
	"time"

	"gorm.io/gorm"
)

// FiscalRegularity represents the fiscal regularity status (Certidão Negativa
// de Débitos) verified against the Receita Federal.
type FiscalRegularity struct {
	gorm.Model
	PersonID  uint      `gorm:"not null;index"`
	CheckDate time.Time `gorm:"not null;index"`

	HasDebts  bool   `gorm:"default:false;index"`
	CNDStatus string `gorm:"size:50;not null"` // regular, irregular, suspended

	CNDNumber     *string `gorm:"size:100"`
	CNDIssueDate  *time.Time
	CNDValidUntil *time.Time

	PendingIssues *string `gorm:"type:json"`
}
