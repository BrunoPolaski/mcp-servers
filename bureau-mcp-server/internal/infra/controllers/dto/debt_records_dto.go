package dto

import (
	"time"

	"github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities"
)

type DebtDTO struct {
	ID               uint       `json:"id"`
	CreatedAt        string     `json:"created_at"`
	UpdatedAt        string     `json:"updated_at"`
	PersonID         uint       `json:"person_id"`
	DebtType         string     `json:"debt_type"`
	Creditor         string     `json:"creditor"`
	CreditorDocument string     `json:"creditor_document"`
	OriginalAmount   float64    `json:"original_amount"`
	CurrentAmount    float64    `json:"current_amount"`
	InterestRate     *float64   `json:"interest_rate,omitempty"`
	Fees             *float64   `json:"fees,omitempty"`
	OriginDate       time.Time  `json:"origin_date"`
	DueDate          time.Time  `json:"due_date"`
	Status           string     `json:"status"`
	InCollection     bool       `json:"in_collection"`
	CollectionDate   *time.Time `json:"collection_date,omitempty"`
	CollectionAgency *string    `json:"collection_agency,omitempty"`
	SettlementAmount *float64   `json:"settlement_amount,omitempty"`
	SettlementDate   *time.Time `json:"settlement_date,omitempty"`
}

func NewDebtDTO(entity *entities.Debt) *DebtDTO {
	return &DebtDTO{
		ID:               entity.ID,
		CreatedAt:        entity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        entity.UpdatedAt.Format("2006-01-02 15:04:05"),
		PersonID:         entity.PersonID,
		DebtType:         entity.DebtType,
		Creditor:         entity.Creditor,
		CreditorDocument: entity.CreditorDocument,
		OriginalAmount:   entity.OriginalAmount,
		CurrentAmount:    entity.CurrentAmount,
		InterestRate:     entity.InterestRate,
		Fees:             entity.Fees,
		OriginDate:       entity.OriginDate,
		DueDate:          entity.DueDate,
		Status:           entity.Status,
		InCollection:     entity.InCollection,
		CollectionDate:   entity.CollectionDate,
		CollectionAgency: entity.CollectionAgency,
		SettlementAmount: entity.SettlementAmount,
		SettlementDate:   entity.SettlementDate,
	}
}

func (d DebtDTO) ToEntity() *entities.Debt {
	return &entities.Debt{
		PersonID:         d.PersonID,
		DebtType:         d.DebtType,
		Creditor:         d.Creditor,
		CreditorDocument: d.CreditorDocument,
		OriginalAmount:   d.OriginalAmount,
		CurrentAmount:    d.CurrentAmount,
		InterestRate:     d.InterestRate,
		Fees:             d.Fees,
		OriginDate:       d.OriginDate,
		DueDate:          d.DueDate,
		Status:           d.Status,
		InCollection:     d.InCollection,
		CollectionDate:   d.CollectionDate,
		CollectionAgency: d.CollectionAgency,
		SettlementAmount: d.SettlementAmount,
		SettlementDate:   d.SettlementDate,
	}
}

type NegativeRecordDTO struct {
	ID               uint       `json:"id"`
	CreatedAt        string     `json:"created_at"`
	UpdatedAt        string     `json:"updated_at"`
	PersonID         uint       `json:"person_id"`
	RecordType       string     `json:"record_type"`
	Creditor         string     `json:"creditor"`
	CreditorDocument string     `json:"creditor_document"`
	Amount           float64    `json:"amount"`
	InclusionDate    time.Time  `json:"inclusion_date"`
	ContractNumber   *string    `json:"contract_number,omitempty"`
	Status           string     `json:"status"`
	RemovalDate      *time.Time `json:"removal_date,omitempty"`
	RemovalReason    *string    `json:"removal_reason,omitempty"`
	ProcessNumber    *string    `json:"process_number,omitempty"`
	Notary           *string    `json:"notary,omitempty"`
	IsDisputed       bool       `json:"is_disputed"`
	DisputeDate      *time.Time `json:"dispute_date,omitempty"`
	DisputeReason    *string    `json:"dispute_reason,omitempty"`
}

func NewNegativeRecordDTO(entity *entities.NegativeRecord) *NegativeRecordDTO {
	return &NegativeRecordDTO{
		ID:               entity.ID,
		CreatedAt:        entity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        entity.UpdatedAt.Format("2006-01-02 15:04:05"),
		PersonID:         entity.PersonID,
		RecordType:       entity.RecordType,
		Creditor:         entity.Creditor,
		CreditorDocument: entity.CreditorDocument,
		Amount:           entity.Amount,
		InclusionDate:    entity.InclusionDate,
		ContractNumber:   entity.ContractNumber,
		Status:           entity.Status,
		RemovalDate:      entity.RemovalDate,
		RemovalReason:    entity.RemovalReason,
		ProcessNumber:    entity.ProcessNumber,
		Notary:           entity.Notary,
		IsDisputed:       entity.IsDisputed,
		DisputeDate:      entity.DisputeDate,
		DisputeReason:    entity.DisputeReason,
	}
}

func (n NegativeRecordDTO) ToEntity() *entities.NegativeRecord {
	return &entities.NegativeRecord{
		PersonID:         n.PersonID,
		RecordType:       n.RecordType,
		Creditor:         n.Creditor,
		CreditorDocument: n.CreditorDocument,
		Amount:           n.Amount,
		InclusionDate:    n.InclusionDate,
		ContractNumber:   n.ContractNumber,
		Status:           n.Status,
		RemovalDate:      n.RemovalDate,
		RemovalReason:    n.RemovalReason,
		ProcessNumber:    n.ProcessNumber,
		Notary:           n.Notary,
		IsDisputed:       n.IsDisputed,
		DisputeDate:      n.DisputeDate,
		DisputeReason:    n.DisputeReason,
	}
}
