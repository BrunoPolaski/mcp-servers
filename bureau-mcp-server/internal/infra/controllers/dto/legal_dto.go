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
	ProcessNumber  *string    `json:"process_number,omitempty"`
	Court          *string    `json:"court,omitempty"`
	FilingDate     *time.Time `json:"filing_date,omitempty"`
	Status         *string    `json:"status,omitempty"`
	Amount         *float64   `json:"amount,omitempty"`
	Description    *string    `json:"description,omitempty"`
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
