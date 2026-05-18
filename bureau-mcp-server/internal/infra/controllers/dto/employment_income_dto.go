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
	EmployerDocument   string     `json:"employer_document"`
	JobTitle           string     `json:"job_title"`
	EmploymentType     string     `json:"employment_type"`
	Salary             *float64   `json:"salary,omitempty"`
	StartDate          time.Time  `json:"start_date"`
	EndDate            *time.Time `json:"end_date,omitempty"`
	IsCurrent          bool       `json:"is_current"`
	VerificationStatus string     `json:"verification_status"`
	DataSource         string     `json:"data_source"`
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

type IncomeDeclarationDTO struct {
	ID              uint      `json:"id"`
	CreatedAt       string    `json:"created_at"`
	UpdatedAt       string    `json:"updated_at"`
	PersonID        uint      `json:"person_id"`
	DeclarationDate time.Time `json:"declaration_date"`
	IncomeType      string    `json:"income_type"`
	MonthlyAmount   float64   `json:"monthly_amount"`
	YearlyAmount    *float64  `json:"yearly_amount,omitempty"`
	Source          string    `json:"source"`
	Verified        bool      `json:"verified"`
	VerifiedBy      *string   `json:"verified_by,omitempty"`
	ProofFileID     *uint     `json:"proof_file_id,omitempty"`
	ProofFile       *FileDTO  `json:"proof_file,omitempty"`
}

func NewIncomeDeclarationDTO(entity *entities.IncomeDeclaration) *IncomeDeclarationDTO {
	dto := &IncomeDeclarationDTO{
		ID:              entity.ID,
		CreatedAt:       entity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       entity.UpdatedAt.Format("2006-01-02 15:04:05"),
		PersonID:        entity.PersonID,
		DeclarationDate: entity.DeclarationDate,
		IncomeType:      entity.IncomeType,
		MonthlyAmount:   entity.MonthlyAmount,
		YearlyAmount:    entity.YearlyAmount,
		Source:          entity.Source,
		Verified:        entity.Verified,
		VerifiedBy:      entity.VerifiedBy,
		ProofFileID:     entity.ProofFileID,
	}

	if entity.ProofFile != nil {
		dto.ProofFile = NewFileDTO(entity.ProofFile)
	}

	return dto
}

func (i IncomeDeclarationDTO) ToEntity() *entities.IncomeDeclaration {
	var file *entities.File
	if i.ProofFile != nil {
		file = i.ProofFile.ToEntity()
	}

	return &entities.IncomeDeclaration{
		PersonID:        i.PersonID,
		DeclarationDate: i.DeclarationDate,
		IncomeType:      i.IncomeType,
		MonthlyAmount:   i.MonthlyAmount,
		YearlyAmount:    i.YearlyAmount,
		Source:          i.Source,
		Verified:        i.Verified,
		VerifiedBy:      i.VerifiedBy,
		ProofFileID:     i.ProofFileID,
		ProofFile:       file,
	}
}
