package entities

import (
	"time"

	"gorm.io/gorm"
)

type ComplianceCheck struct {
	gorm.Model
	PersonID  uint      `gorm:"not null;index"`
	CheckType string    `gorm:"size:100;not null;index"` // receita_federal, pep, sanctions, watchlist
	CheckDate time.Time `gorm:"not null;index"`
	Status    string    `gorm:"size:50;not null"` // clear, flagged, review_required
	Details   string    `gorm:"type:json"`        // JSON with check results

	// PEP (Politically Exposed Person)
	IsPEP      bool    `gorm:"default:false;index"`
	PEPDetails *string `gorm:"type:text"`

	// Sanctions/Watchlist
	OnSanctionsList  bool    `gorm:"default:false;index"`
	SanctionsDetails *string `gorm:"type:text"`

	ValidUntil *time.Time
}
