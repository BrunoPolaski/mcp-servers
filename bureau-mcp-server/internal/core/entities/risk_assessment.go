package entities

import (
	"time"

	"gorm.io/gorm"
)

type RiskAssessment struct {
	gorm.Model
	PersonID       uint      `gorm:"not null;index"`
	AssessmentDate time.Time `gorm:"not null;index"`
	AssessmentType string    `gorm:"size:100;not null"` // credit, fraud, identity, compliance
	RiskScore      *int      `gorm:"index"`             // 0-100
	RiskLevel      *string   `gorm:"size:50;index"`     // low, medium, high, very_high
	RiskFactors    *string   `gorm:"type:json"`         // JSON array of risk factors
	Recommendation *string   `gorm:"type:text"`
	ModelVersion   *string   `gorm:"size:50"`
}
