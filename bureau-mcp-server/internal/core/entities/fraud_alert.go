package entities

import (
	"time"

	"gorm.io/gorm"
)

type FraudAlert struct {
	gorm.Model
	PersonID     uint      `gorm:"not null;index"`
	AlertType    string    `gorm:"size:100;not null;index"` // identity_theft, suspicious_activity, document_fraud
	Severity     string    `gorm:"size:50;not null;index"`  // low, medium, high, critical
	Description  string    `gorm:"type:text;not null"`
	DetectedDate time.Time `gorm:"not null;index"`
	Status       string    `gorm:"size:50;not null"` // active, investigating, resolved, false_positive
	ResolvedDate *time.Time
	ResolvedBy   *string `gorm:"size:255"`
	Notes        *string `gorm:"type:text"`
}
