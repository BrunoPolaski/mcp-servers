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
	ScoreModel  *string   `gorm:"size:50"`   // Which scoring model
	ScoreReason *string   `gorm:"type:text"` // Factors affecting score

	// Score Breakdown
	PaymentHistory  *int // Weight percentage
	CreditUsage     *int
	CreditAge       *int
	CreditMix       *int
	RecentInquiries *int

	// Risk Classification
	RiskLevel          *string  `gorm:"size:50"` // low, medium, high, very_high
	DefaultProbability *float64 // Probability of default
}
