package dto

import (
	"time"

	"github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities"
)

type LegalRecordDTO struct {
	ID             uint       `json:"id"`
	CreatedAt      string     `json:"created_at"`
	UpdatedAt      string     `json:"updated_at"`
	PersonID       uint       `json:"person_id"`
	RecordType     string     `json:"record_type"`
	ProcessNumber  string     `json:"process_number"`
	Court          string     `json:"court"`
	FilingDate     time.Time  `json:"filing_date"`
	Status         string     `json:"status"`
	Amount         *float64   `json:"amount,omitempty"`
	Description    string     `json:"description"`
	Resolution     *string    `json:"resolution,omitempty"`
	ResolutionDate *time.Time `json:"resolution_date,omitempty"`
}

func NewLegalRecordDTO(entity *entities.LegalRecord) *LegalRecordDTO {
	return &LegalRecordDTO{
		ID:             entity.ID,
		CreatedAt:      entity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      entity.UpdatedAt.Format("2006-01-02 15:04:05"),
		PersonID:       entity.PersonID,
		RecordType:     entity.RecordType,
		ProcessNumber:  entity.ProcessNumber,
		Court:          entity.Court,
		FilingDate:     entity.FilingDate,
		Status:         entity.Status,
		Amount:         entity.Amount,
		Description:    entity.Description,
		Resolution:     entity.Resolution,
		ResolutionDate: entity.ResolutionDate,
	}
}

func (l LegalRecordDTO) ToEntity() *entities.LegalRecord {
	return &entities.LegalRecord{
		PersonID:       l.PersonID,
		RecordType:     l.RecordType,
		ProcessNumber:  l.ProcessNumber,
		Court:          l.Court,
		FilingDate:     l.FilingDate,
		Status:         l.Status,
		Amount:         l.Amount,
		Description:    l.Description,
		Resolution:     l.Resolution,
		ResolutionDate: l.ResolutionDate,
	}
}

type ComplianceCheckDTO struct {
	ID               uint       `json:"id"`
	CreatedAt        string     `json:"created_at"`
	UpdatedAt        string     `json:"updated_at"`
	PersonID         uint       `json:"person_id"`
	CheckType        string     `json:"check_type"`
	CheckDate        time.Time  `json:"check_date"`
	Status           string     `json:"status"`
	Details          string     `json:"details"`
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
