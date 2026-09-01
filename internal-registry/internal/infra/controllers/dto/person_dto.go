package dto

import (
	"time"

	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/internal-registry/internal/core/entities/value_objects"
)

type PersonDTO struct {
	ID                     uint                       `json:"id"`
	CreatedAt              string                     `json:"created_at"`
	UpdatedAt              string                     `json:"updated_at"`
	PersonalInformationID  uint                       `json:"personal_information_id"`
	PersonalInformation    *PersonalInformationDTO    `json:"personal_information" validate:"required"`
	CustomerRelationshipID *uint                      `json:"customer_relationship_id,omitempty"`
	CustomerRelationship   *CustomerRelationshipDTO   `json:"customer_relationship,omitempty"`
	ContractedProducts     []ContractedProductDTO     `json:"contracted_products,omitempty"`
	InternalPaymentRecords []InternalPaymentRecordDTO `json:"internal_payment_records,omitempty"`
	PreApprovedLimits      []PreApprovedLimitDTO      `json:"pre_approved_limits,omitempty"`
	IncomeDeclarations     []IncomeDeclarationDTO     `json:"income_declarations,omitempty"`
	LastVerifiedAt         *time.Time                 `json:"last_verified_at,omitempty"`
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

	dto.CustomerRelationshipID = entity.CustomerRelationshipID
	if entity.CustomerRelationship != nil {
		dto.CustomerRelationship = NewCustomerRelationshipDTO(entity.CustomerRelationship)
	}

	if len(entity.ContractedProducts) > 0 {
		dto.ContractedProducts = make([]ContractedProductDTO, 0, len(entity.ContractedProducts))
		for _, p := range entity.ContractedProducts {
			dto.ContractedProducts = append(dto.ContractedProducts, *NewContractedProductDTO(&p))
		}
	}

	if len(entity.InternalPaymentRecords) > 0 {
		dto.InternalPaymentRecords = make([]InternalPaymentRecordDTO, 0, len(entity.InternalPaymentRecords))
		for _, r := range entity.InternalPaymentRecords {
			dto.InternalPaymentRecords = append(dto.InternalPaymentRecords, *NewInternalPaymentRecordDTO(&r))
		}
	}

	if len(entity.PreApprovedLimits) > 0 {
		dto.PreApprovedLimits = make([]PreApprovedLimitDTO, 0, len(entity.PreApprovedLimits))
		for _, l := range entity.PreApprovedLimits {
			dto.PreApprovedLimits = append(dto.PreApprovedLimits, *NewPreApprovedLimitDTO(&l))
		}
	}

	if len(entity.IncomeDeclarations) > 0 {
		dto.IncomeDeclarations = make([]IncomeDeclarationDTO, 0, len(entity.IncomeDeclarations))
		for _, d := range entity.IncomeDeclarations {
			dto.IncomeDeclarations = append(dto.IncomeDeclarations, *NewIncomeDeclarationDTO(&d))
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

	var customerRelationship *entities.CustomerRelationship
	if c.CustomerRelationship != nil {
		customerRelationship = c.CustomerRelationship.ToEntity()
	}

	var contractedProducts []entities.ContractedProduct
	if len(c.ContractedProducts) > 0 {
		contractedProducts = make([]entities.ContractedProduct, 0, len(c.ContractedProducts))
		for _, p := range c.ContractedProducts {
			contractedProducts = append(contractedProducts, *p.ToEntity())
		}
	}

	var internalPaymentRecords []entities.InternalPaymentRecord
	if len(c.InternalPaymentRecords) > 0 {
		internalPaymentRecords = make([]entities.InternalPaymentRecord, 0, len(c.InternalPaymentRecords))
		for _, r := range c.InternalPaymentRecords {
			internalPaymentRecords = append(internalPaymentRecords, *r.ToEntity())
		}
	}

	var preApprovedLimits []entities.PreApprovedLimit
	if len(c.PreApprovedLimits) > 0 {
		preApprovedLimits = make([]entities.PreApprovedLimit, 0, len(c.PreApprovedLimits))
		for _, l := range c.PreApprovedLimits {
			preApprovedLimits = append(preApprovedLimits, *l.ToEntity())
		}
	}

	var incomeDeclarations []entities.IncomeDeclaration
	if len(c.IncomeDeclarations) > 0 {
		incomeDeclarations = make([]entities.IncomeDeclaration, 0, len(c.IncomeDeclarations))
		for _, d := range c.IncomeDeclarations {
			incomeDeclarations = append(incomeDeclarations, *d.ToEntity())
		}
	}

	return &entities.Person{
		PersonalInformationID:  c.PersonalInformationID,
		PersonalInformation:    pi,
		CustomerRelationshipID: c.CustomerRelationshipID,
		CustomerRelationship:   customerRelationship,
		ContractedProducts:     contractedProducts,
		InternalPaymentRecords: internalPaymentRecords,
		PreApprovedLimits:      preApprovedLimits,
		IncomeDeclarations:     incomeDeclarations,
		LastVerifiedAt:         c.LastVerifiedAt,
	}
}

type ListPersonsDTO struct {
	PaginatedDTO
	Document valueobjects.Document `json:"document"`
}
