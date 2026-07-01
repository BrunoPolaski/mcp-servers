package dto

import (
	"time"

	"github.com/BrunoPolaski/registration-validation-mcp-server/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/registration-validation-mcp-server/internal/core/entities/value_objects"
)

type PersonDTO struct {
	ID                        uint                          `json:"id"`
	CreatedAt                 string                        `json:"created_at"`
	UpdatedAt                 string                        `json:"updated_at"`
	PersonalInformationID     uint                          `json:"personal_information_id"`
	PersonalInformation       *PersonalInformationDTO       `json:"personal_information" validate:"required"`
	DocumentValidations       []DocumentValidationDTO       `json:"document_validations,omitempty"`
	FiscalRegularities        []FiscalRegularityDTO         `json:"fiscal_regularities,omitempty"`
	EmploymentLinkValidations []EmploymentLinkValidationDTO `json:"employment_link_validations,omitempty"`
	ComplianceChecks          []ComplianceCheckDTO          `json:"compliance_checks,omitempty"`
	LastVerifiedAt            *time.Time                    `json:"last_verified_at,omitempty"`
}

func NewPersonDTO(entity *entities.Person) *PersonDTO {
	dto := &PersonDTO{
		ID:        entity.ID,
		CreatedAt: entity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: entity.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if entity.PersonalInformation != nil {
		dto.PersonalInformation = NewPersonalInformationDTO(entity.PersonalInformation)
	}
	dto.PersonalInformationID = entity.PersonalInformationID

	if len(entity.DocumentValidations) > 0 {
		dto.DocumentValidations = make([]DocumentValidationDTO, 0, len(entity.DocumentValidations))
		for _, v := range entity.DocumentValidations {
			dto.DocumentValidations = append(dto.DocumentValidations, *NewDocumentValidationDTO(&v))
		}
	}

	if len(entity.FiscalRegularities) > 0 {
		dto.FiscalRegularities = make([]FiscalRegularityDTO, 0, len(entity.FiscalRegularities))
		for _, r := range entity.FiscalRegularities {
			dto.FiscalRegularities = append(dto.FiscalRegularities, *NewFiscalRegularityDTO(&r))
		}
	}

	if len(entity.EmploymentLinkValidations) > 0 {
		dto.EmploymentLinkValidations = make([]EmploymentLinkValidationDTO, 0, len(entity.EmploymentLinkValidations))
		for _, e := range entity.EmploymentLinkValidations {
			dto.EmploymentLinkValidations = append(dto.EmploymentLinkValidations, *NewEmploymentLinkValidationDTO(&e))
		}
	}

	if len(entity.ComplianceChecks) > 0 {
		dto.ComplianceChecks = make([]ComplianceCheckDTO, 0, len(entity.ComplianceChecks))
		for _, c := range entity.ComplianceChecks {
			dto.ComplianceChecks = append(dto.ComplianceChecks, *NewComplianceCheckDTO(&c))
		}
	}

	dto.LastVerifiedAt = entity.LastVerifiedAt

	return dto
}

func (c PersonDTO) ToEntity() *entities.Person {
	var pi *entities.PersonalInformation
	if c.PersonalInformation != nil {
		pi = c.PersonalInformation.ToEntity()
	}

	var documentValidations []entities.DocumentValidation
	if len(c.DocumentValidations) > 0 {
		documentValidations = make([]entities.DocumentValidation, 0, len(c.DocumentValidations))
		for _, v := range c.DocumentValidations {
			documentValidations = append(documentValidations, *v.ToEntity())
		}
	}

	var fiscalRegularities []entities.FiscalRegularity
	if len(c.FiscalRegularities) > 0 {
		fiscalRegularities = make([]entities.FiscalRegularity, 0, len(c.FiscalRegularities))
		for _, r := range c.FiscalRegularities {
			fiscalRegularities = append(fiscalRegularities, *r.ToEntity())
		}
	}

	var employmentLinkValidations []entities.EmploymentLinkValidation
	if len(c.EmploymentLinkValidations) > 0 {
		employmentLinkValidations = make([]entities.EmploymentLinkValidation, 0, len(c.EmploymentLinkValidations))
		for _, e := range c.EmploymentLinkValidations {
			employmentLinkValidations = append(employmentLinkValidations, *e.ToEntity())
		}
	}

	var complianceChecks []entities.ComplianceCheck
	if len(c.ComplianceChecks) > 0 {
		complianceChecks = make([]entities.ComplianceCheck, 0, len(c.ComplianceChecks))
		for _, cc := range c.ComplianceChecks {
			complianceChecks = append(complianceChecks, *cc.ToEntity())
		}
	}

	return &entities.Person{
		PersonalInformationID:     c.PersonalInformationID,
		PersonalInformation:       pi,
		DocumentValidations:       documentValidations,
		FiscalRegularities:        fiscalRegularities,
		EmploymentLinkValidations: employmentLinkValidations,
		ComplianceChecks:          complianceChecks,
		LastVerifiedAt:            c.LastVerifiedAt,
	}
}

type ListPersonsDTO struct {
	PaginatedDTO
	Document valueobjects.Document `json:"document"`
}
