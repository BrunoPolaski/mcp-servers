package entities

import (
	"time"

	"gorm.io/gorm"
)

type LegalRecord struct {
	gorm.Model
	PersonID       uint    `gorm:"not null;index"`
	RecordType     string  `gorm:"size:100;not null;index"` // lawsuit, bankruptcy, criminal, tax_lien
	ProcessNumber  *string `gorm:"size:100;index"`
	Court          *string `gorm:"size:255"`
	FilingDate     *time.Time
	Status         *string `gorm:"size:50"` // active, closed, appealed
	Amount         *float64
	Description    *string `gorm:"type:text"`
	Resolution     *string `gorm:"type:text"`
	ResolutionDate *time.Time
}
