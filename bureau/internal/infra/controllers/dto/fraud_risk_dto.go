package dto

import (
	"time"

	"github.com/BrunoPolaski/bureau/internal/core/entities"
)

type FraudAlertDTO struct {
	ID           uint       `json:"id"`
	CreatedAt    string     `json:"created_at"`
	UpdatedAt    string     `json:"updated_at"`
	PersonID     uint       `json:"person_id"`
	AlertType    string     `json:"alert_type"`
	Severity     *string    `json:"severity,omitempty"`
	Description  *string    `json:"description,omitempty"`
	DetectedDate *time.Time `json:"detected_date,omitempty"`
	Status       *string    `json:"status,omitempty"`
	ResolvedDate *time.Time `json:"resolved_date,omitempty"`
	ResolvedBy   *string    `json:"resolved_by,omitempty"`
	Notes        *string    `json:"notes,omitempty"`
}

func NewFraudAlertDTO(entity *entities.FraudAlert) *FraudAlertDTO {
	return &FraudAlertDTO{
		ID:           entity.ID,
		CreatedAt:    entity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    entity.UpdatedAt.Format("2006-01-02 15:04:05"),
		PersonID:     entity.PersonID,
		AlertType:    entity.AlertType,
		Severity:     entity.Severity,
		Description:  entity.Description,
		DetectedDate: entity.DetectedDate,
		Status:       entity.Status,
		ResolvedDate: entity.ResolvedDate,
		ResolvedBy:   entity.ResolvedBy,
		Notes:        entity.Notes,
	}
}

func (f FraudAlertDTO) ToEntity() *entities.FraudAlert {
	return &entities.FraudAlert{
		PersonID:     f.PersonID,
		AlertType:    f.AlertType,
		Severity:     f.Severity,
		Description:  f.Description,
		DetectedDate: f.DetectedDate,
		Status:       f.Status,
		ResolvedDate: f.ResolvedDate,
		ResolvedBy:   f.ResolvedBy,
		Notes:        f.Notes,
	}
}

type RiskAssessmentDTO struct {
	ID             uint      `json:"id"`
	CreatedAt      string    `json:"created_at"`
	UpdatedAt      string    `json:"updated_at"`
	PersonID       uint      `json:"person_id"`
	AssessmentDate time.Time `json:"assessment_date"`
	AssessmentType string    `json:"assessment_type"`
	RiskScore      *int      `json:"risk_score,omitempty"`
	RiskLevel      *string   `json:"risk_level,omitempty"`
	RiskFactors    *string   `json:"risk_factors,omitempty"`
	Recommendation *string   `json:"recommendation,omitempty"`
	ModelVersion   *string   `json:"model_version,omitempty"`
}

func NewRiskAssessmentDTO(entity *entities.RiskAssessment) *RiskAssessmentDTO {
	return &RiskAssessmentDTO{
		ID:             entity.ID,
		CreatedAt:      entity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      entity.UpdatedAt.Format("2006-01-02 15:04:05"),
		PersonID:       entity.PersonID,
		AssessmentDate: entity.AssessmentDate,
		AssessmentType: entity.AssessmentType,
		RiskScore:      entity.RiskScore,
		RiskLevel:      entity.RiskLevel,
		RiskFactors:    entity.RiskFactors,
		Recommendation: entity.Recommendation,
		ModelVersion:   entity.ModelVersion,
	}
}

func (r RiskAssessmentDTO) ToEntity() *entities.RiskAssessment {
	return &entities.RiskAssessment{
		PersonID:       r.PersonID,
		AssessmentDate: r.AssessmentDate,
		AssessmentType: r.AssessmentType,
		RiskScore:      r.RiskScore,
		RiskLevel:      r.RiskLevel,
		RiskFactors:    r.RiskFactors,
		Recommendation: r.Recommendation,
		ModelVersion:   r.ModelVersion,
	}
}
