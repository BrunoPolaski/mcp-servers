package dto

import (
	"time"

	"github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities"
)

type CreditScoreDTO struct {
	ID                 uint      `json:"id"`
	CreatedAt          string    `json:"created_at"`
	UpdatedAt          string    `json:"updated_at"`
	PersonID           uint      `json:"person_id"`
	Score              int       `json:"score"`
	ScoreDate          time.Time `json:"score_date"`
	ScoreModel         string    `json:"score_model"`
	ScoreReason        string    `json:"score_reason"`
	PaymentHistory     int       `json:"payment_history"`
	CreditUsage        int       `json:"credit_usage"`
	CreditAge          int       `json:"credit_age"`
	CreditMix          int       `json:"credit_mix"`
	RecentInquiries    int       `json:"recent_inquiries"`
	RiskLevel          string    `json:"risk_level"`
	DefaultProbability float64   `json:"default_probability"`
}

func NewCreditScoreDTO(entity *entities.CreditScore) *CreditScoreDTO {
	return &CreditScoreDTO{
		ID:                 entity.ID,
		CreatedAt:          entity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:          entity.UpdatedAt.Format("2006-01-02 15:04:05"),
		PersonID:           entity.PersonID,
		Score:              entity.Score,
		ScoreDate:          entity.ScoreDate,
		ScoreModel:         entity.ScoreModel,
		ScoreReason:        entity.ScoreReason,
		PaymentHistory:     entity.PaymentHistory,
		CreditUsage:        entity.CreditUsage,
		CreditAge:          entity.CreditAge,
		CreditMix:          entity.CreditMix,
		RecentInquiries:    entity.RecentInquiries,
		RiskLevel:          entity.RiskLevel,
		DefaultProbability: entity.DefaultProbability,
	}
}

func (c CreditScoreDTO) ToEntity() *entities.CreditScore {
	return &entities.CreditScore{
		PersonID:           c.PersonID,
		Score:              c.Score,
		ScoreDate:          c.ScoreDate,
		ScoreModel:         c.ScoreModel,
		ScoreReason:        c.ScoreReason,
		PaymentHistory:     c.PaymentHistory,
		CreditUsage:        c.CreditUsage,
		CreditAge:          c.CreditAge,
		CreditMix:          c.CreditMix,
		RecentInquiries:    c.RecentInquiries,
		RiskLevel:          c.RiskLevel,
		DefaultProbability: c.DefaultProbability,
	}
}
