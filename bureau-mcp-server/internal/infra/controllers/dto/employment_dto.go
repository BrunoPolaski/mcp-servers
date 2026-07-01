package dto

import (
	"time"

	"github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities"
)

type EmploymentRecordDTO struct {
	ID                 uint       `json:"id"`
	CreatedAt          string     `json:"created_at"`
	UpdatedAt          string     `json:"updated_at"`
	PersonID           uint       `json:"person_id"`
	EmployerName       string     `json:"employer_name"`
	EmployerDocument   *string    `json:"employer_document,omitempty"`
	JobTitle           *string    `json:"job_title,omitempty"`
	EmploymentType     *string    `json:"employment_type,omitempty"`
	Salary             *float64   `json:"salary,omitempty"`
	StartDate          *time.Time `json:"start_date,omitempty"`
	EndDate            *time.Time `json:"end_date,omitempty"`
	IsCurrent          bool       `json:"is_current"`
	VerificationStatus string     `json:"verification_status"`
	DataSource         *string    `json:"data_source,omitempty"`
}

func NewEmploymentRecordDTO(entity *entities.EmploymentRecord) *EmploymentRecordDTO {
	return &EmploymentRecordDTO{
		ID:                 entity.ID,
		CreatedAt:          entity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:          entity.UpdatedAt.Format("2006-01-02 15:04:05"),
		PersonID:           entity.PersonID,
		EmployerName:       entity.EmployerName,
		EmployerDocument:   entity.EmployerDocument,
		JobTitle:           entity.JobTitle,
		EmploymentType:     entity.EmploymentType,
		Salary:             entity.Salary,
		StartDate:          entity.StartDate,
		EndDate:            entity.EndDate,
		IsCurrent:          entity.IsCurrent,
		VerificationStatus: entity.VerificationStatus,
		DataSource:         entity.DataSource,
	}
}

func (e EmploymentRecordDTO) ToEntity() *entities.EmploymentRecord {
	return &entities.EmploymentRecord{
		PersonID:           e.PersonID,
		EmployerName:       e.EmployerName,
		EmployerDocument:   e.EmployerDocument,
		JobTitle:           e.JobTitle,
		EmploymentType:     e.EmploymentType,
		Salary:             e.Salary,
		StartDate:          e.StartDate,
		EndDate:            e.EndDate,
		IsCurrent:          e.IsCurrent,
		VerificationStatus: e.VerificationStatus,
		DataSource:         e.DataSource,
	}
}
