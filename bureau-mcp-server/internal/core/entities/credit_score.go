package entities

import (
	"time"

	"gorm.io/gorm"
)

type CreditScore struct {
	gorm.Model
	PersonID    uint      `gorm:"not null;uniqueIndex:idx_person_score_date"`
	Score       int       `gorm:"not null;index"` // 0-1000 scale (Serasa-like)
	ScoreDate   time.Time `gorm:"not null;uniqueIndex:idx_person_score_date"`
	ScoreModel  string    `gorm:"size:50;not null"` // Which scoring model
	ScoreReason string    `gorm:"type:text"`        // Factors affecting score

	// Score Breakdown
	PaymentHistory  int `gorm:"not null"` // Weight percentage
	CreditUsage     int `gorm:"not null"`
	CreditAge       int `gorm:"not null"`
	CreditMix       int `gorm:"not null"`
	RecentInquiries int `gorm:"not null"`

	// Risk Classification
	RiskLevel          string  `gorm:"size:50;not null"` // low, medium, high, very_high
	DefaultProbability float64 `gorm:"not null"`         // Probability of default
}
