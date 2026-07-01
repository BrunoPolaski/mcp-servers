package dto

import (
	"time"

	"github.com/BrunoPolaski/registration-validation-mcp-server/internal/core/entities"
)

type DocumentValidationDTO struct {
	ID                   uint      `json:"id"`
	PersonID             uint      `json:"person_id"`
	ValidationDate       time.Time `json:"validation_date"`
	DocumentNumber       string    `json:"document_number"`
	DocumentType         string    `json:"document_type"`
	ReceitaFederalStatus string    `json:"receita_federal_status"`
	IsValid              bool      `json:"is_valid"`
	NameMatches          bool      `json:"name_matches"`
	BirthDateMatches     bool      `json:"birth_date_matches"`
	BiometricValidated   bool      `json:"biometric_validated"`
	Source               string    `json:"source"`
	RawResponse          *string   `json:"raw_response,omitempty"`
}

func NewDocumentValidationDTO(entity *entities.DocumentValidation) *DocumentValidationDTO {
	return &DocumentValidationDTO{
		ID:                   entity.ID,
		PersonID:             entity.PersonID,
		ValidationDate:       entity.ValidationDate,
		DocumentNumber:       entity.DocumentNumber,
		DocumentType:         entity.DocumentType,
		ReceitaFederalStatus: entity.ReceitaFederalStatus,
		IsValid:              entity.IsValid,
		NameMatches:          entity.NameMatches,
		BirthDateMatches:     entity.BirthDateMatches,
		BiometricValidated:   entity.BiometricValidated,
		Source:               entity.Source,
		RawResponse:          entity.RawResponse,
	}
}

func (d DocumentValidationDTO) ToEntity() *entities.DocumentValidation {
	return &entities.DocumentValidation{
		PersonID:             d.PersonID,
		ValidationDate:       d.ValidationDate,
		DocumentNumber:       d.DocumentNumber,
		DocumentType:         d.DocumentType,
		ReceitaFederalStatus: d.ReceitaFederalStatus,
		IsValid:              d.IsValid,
		NameMatches:          d.NameMatches,
		BirthDateMatches:     d.BirthDateMatches,
		BiometricValidated:   d.BiometricValidated,
		Source:               d.Source,
		RawResponse:          d.RawResponse,
	}
}

type FiscalRegularityDTO struct {
	ID            uint       `json:"id"`
	PersonID      uint       `json:"person_id"`
	CheckDate     time.Time  `json:"check_date"`
	HasDebts      bool       `json:"has_debts"`
	CNDStatus     string     `json:"cnd_status"`
	CNDNumber     *string    `json:"cnd_number,omitempty"`
	CNDIssueDate  *time.Time `json:"cnd_issue_date,omitempty"`
	CNDValidUntil *time.Time `json:"cnd_valid_until,omitempty"`
	PendingIssues *string    `json:"pending_issues,omitempty"`
}

func NewFiscalRegularityDTO(entity *entities.FiscalRegularity) *FiscalRegularityDTO {
	return &FiscalRegularityDTO{
		ID:            entity.ID,
		PersonID:      entity.PersonID,
		CheckDate:     entity.CheckDate,
		HasDebts:      entity.HasDebts,
		CNDStatus:     entity.CNDStatus,
		CNDNumber:     entity.CNDNumber,
		CNDIssueDate:  entity.CNDIssueDate,
		CNDValidUntil: entity.CNDValidUntil,
		PendingIssues: entity.PendingIssues,
	}
}

func (d FiscalRegularityDTO) ToEntity() *entities.FiscalRegularity {
	return &entities.FiscalRegularity{
		PersonID:      d.PersonID,
		CheckDate:     d.CheckDate,
		HasDebts:      d.HasDebts,
		CNDStatus:     d.CNDStatus,
		CNDNumber:     d.CNDNumber,
		CNDIssueDate:  d.CNDIssueDate,
		CNDValidUntil: d.CNDValidUntil,
		PendingIssues: d.PendingIssues,
	}
}

type EmploymentLinkValidationDTO struct {
	ID               uint       `json:"id"`
	PersonID         uint       `json:"person_id"`
	ValidationDate   time.Time  `json:"validation_date"`
	EmployerName     string     `json:"employer_name"`
	EmployerDocument string     `json:"employer_document,omitempty"`
	EmploymentType   string     `json:"employment_type,omitempty"`
	Status           string     `json:"status"`
	StartDate        time.Time  `json:"start_date"`
	EndDate          *time.Time `json:"end_date,omitempty"`
	Source           string     `json:"source"`
	Verified         bool       `json:"verified"`
}

func NewEmploymentLinkValidationDTO(entity *entities.EmploymentLinkValidation) *EmploymentLinkValidationDTO {
	return &EmploymentLinkValidationDTO{
		ID:               entity.ID,
		PersonID:         entity.PersonID,
		ValidationDate:   entity.ValidationDate,
		EmployerName:     entity.EmployerName,
		EmployerDocument: entity.EmployerDocument,
		EmploymentType:   entity.EmploymentType,
		Status:           entity.Status,
		StartDate:        entity.StartDate,
		EndDate:          entity.EndDate,
		Source:           entity.Source,
		Verified:         entity.Verified,
	}
}

func (d EmploymentLinkValidationDTO) ToEntity() *entities.EmploymentLinkValidation {
	return &entities.EmploymentLinkValidation{
		PersonID:         d.PersonID,
		ValidationDate:   d.ValidationDate,
		EmployerName:     d.EmployerName,
		EmployerDocument: d.EmployerDocument,
		EmploymentType:   d.EmploymentType,
		Status:           d.Status,
		StartDate:        d.StartDate,
		EndDate:          d.EndDate,
		Source:           d.Source,
		Verified:         d.Verified,
	}
}
