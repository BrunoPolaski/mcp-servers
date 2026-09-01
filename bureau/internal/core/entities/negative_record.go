package entities

import (
	"time"

	"gorm.io/gorm"
)

type NegativeRecord struct {
	gorm.Model
	PersonID         uint    `gorm:"not null;index"`
	RecordType       string  `gorm:"size:100;not null;index"` // protest, bounced_check, spc, serasa, refin
	Creditor         *string `gorm:"size:255"`
	CreditorDocument *string `gorm:"size:14"`

	Amount         float64   `gorm:"not null"`
	InclusionDate  time.Time `gorm:"not null;index"`
	ContractNumber *string   `gorm:"size:100"`

	Status        *string `gorm:"size:50;index"` // active, removed, disputed, expired
	RemovalDate   *time.Time
	RemovalReason *string `gorm:"size:255"`

	// Legal Info
	ProcessNumber *string `gorm:"size:100"` // For protests
	Notary        *string `gorm:"size:255"` // Cartório

	IsDisputed    bool `gorm:"default:false"`
	DisputeDate   *time.Time
	DisputeReason *string `gorm:"type:text"`
}
