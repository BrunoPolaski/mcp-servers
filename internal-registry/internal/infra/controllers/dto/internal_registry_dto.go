package dto

import (
	"time"

	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
)

type CustomerRelationshipDTO struct {
	ID                 uint      `json:"id"`
	CreatedAt          string    `json:"created_at"`
	UpdatedAt          string    `json:"updated_at"`
	PersonID           uint      `json:"person_id"`
	CustomerSince      time.Time `json:"customer_since"`
	RelationshipMonths int       `json:"relationship_months"`
	Segment            string    `json:"segment,omitempty"`
	Branch             *string   `json:"branch,omitempty"`
	IsActive           bool      `json:"is_active"`
	ChurnRisk          *string   `json:"churn_risk,omitempty"`
	InternalScore      *int      `json:"internal_score,omitempty"`
}

func NewCustomerRelationshipDTO(entity *entities.CustomerRelationship) *CustomerRelationshipDTO {
	return &CustomerRelationshipDTO{
		ID:                 entity.ID,
		CreatedAt:          entity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:          entity.UpdatedAt.Format("2006-01-02 15:04:05"),
		PersonID:           entity.PersonID,
		CustomerSince:      entity.CustomerSince,
		RelationshipMonths: entity.RelationshipMonths,
		Segment:            entity.Segment,
		Branch:             entity.Branch,
		IsActive:           entity.IsActive,
		ChurnRisk:          entity.ChurnRisk,
		InternalScore:      entity.InternalScore,
	}
}

func (d CustomerRelationshipDTO) ToEntity() *entities.CustomerRelationship {
	return &entities.CustomerRelationship{
		PersonID:           d.PersonID,
		CustomerSince:      d.CustomerSince,
		RelationshipMonths: d.RelationshipMonths,
		Segment:            d.Segment,
		Branch:             d.Branch,
		IsActive:           d.IsActive,
		ChurnRisk:          d.ChurnRisk,
		InternalScore:      d.InternalScore,
	}
}

type ContractedProductDTO struct {
	ID             uint      `json:"id"`
	PersonID       uint      `json:"person_id"`
	ProductType    string    `json:"product_type"`
	ProductName    string    `json:"product_name"`
	ContractNumber string    `json:"contract_number,omitempty"`
	ContractedDate time.Time `json:"contracted_date"`
	Status         string    `json:"status"`
	Balance        *float64  `json:"balance,omitempty"`
	MonthlyValue   *float64  `json:"monthly_value,omitempty"`
}

func NewContractedProductDTO(entity *entities.ContractedProduct) *ContractedProductDTO {
	return &ContractedProductDTO{
		ID:             entity.ID,
		PersonID:       entity.PersonID,
		ProductType:    entity.ProductType,
		ProductName:    entity.ProductName,
		ContractNumber: entity.ContractNumber,
		ContractedDate: entity.ContractedDate,
		Status:         entity.Status,
		Balance:        entity.Balance,
		MonthlyValue:   entity.MonthlyValue,
	}
}

func (d ContractedProductDTO) ToEntity() *entities.ContractedProduct {
	return &entities.ContractedProduct{
		PersonID:       d.PersonID,
		ProductType:    d.ProductType,
		ProductName:    d.ProductName,
		ContractNumber: d.ContractNumber,
		ContractedDate: d.ContractedDate,
		Status:         d.Status,
		Balance:        d.Balance,
		MonthlyValue:   d.MonthlyValue,
	}
}

type InternalPaymentRecordDTO struct {
	ID                  uint       `json:"id"`
	PersonID            uint       `json:"person_id"`
	ContractedProductID *uint      `json:"contracted_product_id,omitempty"`
	ReferenceMonth      time.Time  `json:"reference_month"`
	DueDate             time.Time  `json:"due_date"`
	PaymentDate         *time.Time `json:"payment_date,omitempty"`
	AmountDue           float64    `json:"amount_due"`
	AmountPaid          float64    `json:"amount_paid"`
	Status              string     `json:"status"`
	DaysLate            int        `json:"days_late"`
}

func NewInternalPaymentRecordDTO(entity *entities.InternalPaymentRecord) *InternalPaymentRecordDTO {
	return &InternalPaymentRecordDTO{
		ID:                  entity.ID,
		PersonID:            entity.PersonID,
		ContractedProductID: entity.ContractedProductID,
		ReferenceMonth:      entity.ReferenceMonth,
		DueDate:             entity.DueDate,
		PaymentDate:         entity.PaymentDate,
		AmountDue:           entity.AmountDue,
		AmountPaid:          entity.AmountPaid,
		Status:              entity.Status,
		DaysLate:            entity.DaysLate,
	}
}

func (d InternalPaymentRecordDTO) ToEntity() *entities.InternalPaymentRecord {
	return &entities.InternalPaymentRecord{
		PersonID:            d.PersonID,
		ContractedProductID: d.ContractedProductID,
		ReferenceMonth:      d.ReferenceMonth,
		DueDate:             d.DueDate,
		PaymentDate:         d.PaymentDate,
		AmountDue:           d.AmountDue,
		AmountPaid:          d.AmountPaid,
		Status:              d.Status,
		DaysLate:            d.DaysLate,
	}
}

type PreApprovedLimitDTO struct {
	ID             uint      `json:"id"`
	PersonID       uint      `json:"person_id"`
	ProductType    string    `json:"product_type"`
	ApprovedAmount float64   `json:"approved_amount"`
	InterestRate   *float64  `json:"interest_rate,omitempty"`
	CalculatedDate time.Time `json:"calculated_date"`
	ValidUntil     time.Time `json:"valid_until"`
	PolicyVersion  string    `json:"policy_version,omitempty"`
	IsActive       bool      `json:"is_active"`
}

func NewPreApprovedLimitDTO(entity *entities.PreApprovedLimit) *PreApprovedLimitDTO {
	return &PreApprovedLimitDTO{
		ID:             entity.ID,
		PersonID:       entity.PersonID,
		ProductType:    entity.ProductType,
		ApprovedAmount: entity.ApprovedAmount,
		InterestRate:   entity.InterestRate,
		CalculatedDate: entity.CalculatedDate,
		ValidUntil:     entity.ValidUntil,
		PolicyVersion:  entity.PolicyVersion,
		IsActive:       entity.IsActive,
	}
}

func (d PreApprovedLimitDTO) ToEntity() *entities.PreApprovedLimit {
	return &entities.PreApprovedLimit{
		PersonID:       d.PersonID,
		ProductType:    d.ProductType,
		ApprovedAmount: d.ApprovedAmount,
		InterestRate:   d.InterestRate,
		CalculatedDate: d.CalculatedDate,
		ValidUntil:     d.ValidUntil,
		PolicyVersion:  d.PolicyVersion,
		IsActive:       d.IsActive,
	}
}
