package dto

import (
	"time"

	"github.com/BrunoPolaski/open-finance/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/open-finance/internal/core/entities/value_objects"
)

type PersonDTO struct {
	ID                    uint                      `json:"id"`
	CreatedAt             string                    `json:"created_at"`
	UpdatedAt             string                    `json:"updated_at"`
	PersonalInformationID uint                      `json:"personal_information_id"`
	PersonalInformation   *PersonalInformationDTO   `json:"personal_information" validate:"required"`
	BankAccountProfileID  *uint                     `json:"bank_account_profile_id,omitempty"`
	BankAccountProfile    *BankAccountProfileDTO    `json:"bank_account_profile,omitempty"`
	BankStatements        []BankStatementDTO        `json:"bank_statements,omitempty"`
	CashFlowAnalyses      []CashFlowAnalysisDTO     `json:"cash_flow_analyses,omitempty"`
	RecurringTransactions []RecurringTransactionDTO `json:"recurring_transactions,omitempty"`
	DataSharingConsents   []DataSharingConsentDTO   `json:"data_sharing_consents,omitempty"`
	LastVerifiedAt        *time.Time                `json:"last_verified_at,omitempty"`
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

	dto.BankAccountProfileID = entity.BankAccountProfileID
	if entity.BankAccountProfile != nil {
		dto.BankAccountProfile = NewBankAccountProfileDTO(entity.BankAccountProfile)
	}

	if len(entity.BankStatements) > 0 {
		dto.BankStatements = make([]BankStatementDTO, 0, len(entity.BankStatements))
		for _, s := range entity.BankStatements {
			dto.BankStatements = append(dto.BankStatements, *NewBankStatementDTO(&s))
		}
	}

	if len(entity.CashFlowAnalyses) > 0 {
		dto.CashFlowAnalyses = make([]CashFlowAnalysisDTO, 0, len(entity.CashFlowAnalyses))
		for _, c := range entity.CashFlowAnalyses {
			dto.CashFlowAnalyses = append(dto.CashFlowAnalyses, *NewCashFlowAnalysisDTO(&c))
		}
	}

	if len(entity.RecurringTransactions) > 0 {
		dto.RecurringTransactions = make([]RecurringTransactionDTO, 0, len(entity.RecurringTransactions))
		for _, r := range entity.RecurringTransactions {
			dto.RecurringTransactions = append(dto.RecurringTransactions, *NewRecurringTransactionDTO(&r))
		}
	}

	if len(entity.DataSharingConsents) > 0 {
		dto.DataSharingConsents = make([]DataSharingConsentDTO, 0, len(entity.DataSharingConsents))
		for _, c := range entity.DataSharingConsents {
			dto.DataSharingConsents = append(dto.DataSharingConsents, *NewDataSharingConsentDTO(&c))
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

	var bankAccountProfile *entities.BankAccountProfile
	if c.BankAccountProfile != nil {
		bankAccountProfile = c.BankAccountProfile.ToEntity()
	}

	var bankStatements []entities.BankStatement
	if len(c.BankStatements) > 0 {
		bankStatements = make([]entities.BankStatement, 0, len(c.BankStatements))
		for _, s := range c.BankStatements {
			bankStatements = append(bankStatements, *s.ToEntity())
		}
	}

	var cashFlowAnalyses []entities.CashFlowAnalysis
	if len(c.CashFlowAnalyses) > 0 {
		cashFlowAnalyses = make([]entities.CashFlowAnalysis, 0, len(c.CashFlowAnalyses))
		for _, a := range c.CashFlowAnalyses {
			cashFlowAnalyses = append(cashFlowAnalyses, *a.ToEntity())
		}
	}

	var recurringTransactions []entities.RecurringTransaction
	if len(c.RecurringTransactions) > 0 {
		recurringTransactions = make([]entities.RecurringTransaction, 0, len(c.RecurringTransactions))
		for _, r := range c.RecurringTransactions {
			recurringTransactions = append(recurringTransactions, *r.ToEntity())
		}
	}

	var dataSharingConsents []entities.DataSharingConsent
	if len(c.DataSharingConsents) > 0 {
		dataSharingConsents = make([]entities.DataSharingConsent, 0, len(c.DataSharingConsents))
		for _, s := range c.DataSharingConsents {
			dataSharingConsents = append(dataSharingConsents, *s.ToEntity())
		}
	}

	return &entities.Person{
		PersonalInformationID: c.PersonalInformationID,
		PersonalInformation:   pi,
		BankAccountProfileID:  c.BankAccountProfileID,
		BankAccountProfile:    bankAccountProfile,
		BankStatements:        bankStatements,
		CashFlowAnalyses:      cashFlowAnalyses,
		RecurringTransactions: recurringTransactions,
		DataSharingConsents:   dataSharingConsents,
		LastVerifiedAt:        c.LastVerifiedAt,
	}
}

type ListPersonsDTO struct {
	PaginatedDTO
	Document valueobjects.Document `json:"document"`
}

// PersonSummaryDTO é uma projeção enxuta de um cliente, usada pela listagem.
// Expõe apenas o necessário para identificá-lo, de modo que o chamador seja
// obrigado a buscar o cadastro completo por get_customer_by_id ou
// get_customer_by_document antes de qualquer análise.
type PersonSummaryDTO struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Document string `json:"document"`
}

func NewPersonSummaryDTO(entity *entities.Person) *PersonSummaryDTO {
	dto := &PersonSummaryDTO{
		ID: entity.ID,
	}

	if entity.PersonalInformation != nil {
		dto.Name = entity.PersonalInformation.FullName
		dto.Document = entity.PersonalInformation.Document.String()
	}

	return dto
}
