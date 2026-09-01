package dto

import (
	"time"

	"github.com/BrunoPolaski/registration-validation/internal/core/entities"
)

type ComplianceCheckDTO struct {
	ID               uint       `json:"id"`
	CreatedAt        string     `json:"created_at"`
	UpdatedAt        string     `json:"updated_at"`
	PersonID         uint       `json:"person_id"`
	CheckType        string     `json:"check_type"`
	CheckDate        time.Time  `json:"check_date"`
	Status           *string    `json:"status,omitempty"`
	Details          *string    `json:"details,omitempty"`
	IsPEP            bool       `json:"is_pep"`
	PEPDetails       *string    `json:"pep_details,omitempty"`
	OnSanctionsList  bool       `json:"on_sanctions_list"`
	SanctionsDetails *string    `json:"sanctions_details,omitempty"`
	ValidUntil       *time.Time `json:"valid_until,omitempty"`
}

func NewComplianceCheckDTO(entity *entities.ComplianceCheck) *ComplianceCheckDTO {
	return &ComplianceCheckDTO{
		ID:               entity.ID,
		CreatedAt:        entity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        entity.UpdatedAt.Format("2006-01-02 15:04:05"),
		PersonID:         entity.PersonID,
		CheckType:        entity.CheckType,
		CheckDate:        entity.CheckDate,
		Status:           entity.Status,
		Details:          entity.Details,
		IsPEP:            entity.IsPEP,
		PEPDetails:       entity.PEPDetails,
		OnSanctionsList:  entity.OnSanctionsList,
		SanctionsDetails: entity.SanctionsDetails,
		ValidUntil:       entity.ValidUntil,
	}
}

func (c ComplianceCheckDTO) ToEntity() *entities.ComplianceCheck {
	return &entities.ComplianceCheck{
		PersonID:         c.PersonID,
		CheckType:        c.CheckType,
		CheckDate:        c.CheckDate,
		Status:           c.Status,
		Details:          c.Details,
		IsPEP:            c.IsPEP,
		PEPDetails:       c.PEPDetails,
		OnSanctionsList:  c.OnSanctionsList,
		SanctionsDetails: c.SanctionsDetails,
		ValidUntil:       c.ValidUntil,
	}
}
