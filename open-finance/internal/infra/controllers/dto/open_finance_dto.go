package dto

import (
	"time"

	"github.com/BrunoPolaski/open-finance/internal/core/entities"
)

type BankAccountProfileDTO struct {
	ID                   uint      `json:"id"`
	CreatedAt            string    `json:"created_at"`
	UpdatedAt            string    `json:"updated_at"`
	PersonID             uint      `json:"person_id"`
	ProfileDate          time.Time `json:"profile_date"`
	BankingRelationships int       `json:"banking_relationships"`
	AccountAgeAverage    *int      `json:"account_age_average,omitempty"`
	HasCheckingAccount   bool      `json:"has_checking_account"`
	HasSavingsAccount    bool      `json:"has_savings_account"`
	HasInvestmentAccount bool      `json:"has_investment_account"`
	InvestmentsValue     *float64  `json:"investments_value,omitempty"`
}

func NewBankAccountProfileDTO(entity *entities.BankAccountProfile) *BankAccountProfileDTO {
	return &BankAccountProfileDTO{
		ID:                   entity.ID,
		CreatedAt:            entity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:            entity.UpdatedAt.Format("2006-01-02 15:04:05"),
		PersonID:             entity.PersonID,
		ProfileDate:          entity.ProfileDate,
		BankingRelationships: entity.BankingRelationships,
		AccountAgeAverage:    entity.AccountAgeAverage,
		HasCheckingAccount:   entity.HasCheckingAccount,
		HasSavingsAccount:    entity.HasSavingsAccount,
		HasInvestmentAccount: entity.HasInvestmentAccount,
		InvestmentsValue:     entity.InvestmentsValue,
	}
}

func (d BankAccountProfileDTO) ToEntity() *entities.BankAccountProfile {
	return &entities.BankAccountProfile{
		PersonID:             d.PersonID,
		ProfileDate:          d.ProfileDate,
		BankingRelationships: d.BankingRelationships,
		AccountAgeAverage:    d.AccountAgeAverage,
		HasCheckingAccount:   d.HasCheckingAccount,
		HasSavingsAccount:    d.HasSavingsAccount,
		HasInvestmentAccount: d.HasInvestmentAccount,
		InvestmentsValue:     d.InvestmentsValue,
	}
}

type BankStatementDTO struct {
	ID                  uint      `json:"id"`
	PersonID            uint      `json:"person_id"`
	Institution         string    `json:"institution"`
	InstitutionDocument string    `json:"institution_document,omitempty"`
	AccountType         string    `json:"account_type"`
	PeriodStart         time.Time `json:"period_start"`
	PeriodEnd           time.Time `json:"period_end"`
	OpeningBalance      float64   `json:"opening_balance"`
	ClosingBalance      float64   `json:"closing_balance"`
	TotalCredits        float64   `json:"total_credits"`
	TotalDebits         float64   `json:"total_debits"`
	TransactionCount    int       `json:"transaction_count"`
	Currency            string    `json:"currency"`
}

func NewBankStatementDTO(entity *entities.BankStatement) *BankStatementDTO {
	return &BankStatementDTO{
		ID:                  entity.ID,
		PersonID:            entity.PersonID,
		Institution:         entity.Institution,
		InstitutionDocument: entity.InstitutionDocument,
		AccountType:         entity.AccountType,
		PeriodStart:         entity.PeriodStart,
		PeriodEnd:           entity.PeriodEnd,
		OpeningBalance:      entity.OpeningBalance,
		ClosingBalance:      entity.ClosingBalance,
		TotalCredits:        entity.TotalCredits,
		TotalDebits:         entity.TotalDebits,
		TransactionCount:    entity.TransactionCount,
		Currency:            entity.Currency,
	}
}

func (d BankStatementDTO) ToEntity() *entities.BankStatement {
	return &entities.BankStatement{
		PersonID:            d.PersonID,
		Institution:         d.Institution,
		InstitutionDocument: d.InstitutionDocument,
		AccountType:         d.AccountType,
		PeriodStart:         d.PeriodStart,
		PeriodEnd:           d.PeriodEnd,
		OpeningBalance:      d.OpeningBalance,
		ClosingBalance:      d.ClosingBalance,
		TotalCredits:        d.TotalCredits,
		TotalDebits:         d.TotalDebits,
		TransactionCount:    d.TransactionCount,
		Currency:            d.Currency,
	}
}

type CashFlowAnalysisDTO struct {
	ID                    uint      `json:"id"`
	PersonID              uint      `json:"person_id"`
	AnalysisDate          time.Time `json:"analysis_date"`
	PeriodDays            int       `json:"period_days"`
	AverageMonthlyInflow  float64   `json:"average_monthly_inflow"`
	AverageMonthlyOutflow float64   `json:"average_monthly_outflow"`
	NetCashFlow           float64   `json:"net_cash_flow"`
	InflowVolatility      *float64  `json:"inflow_volatility,omitempty"`
	NegativeBalanceDays   int       `json:"negative_balance_days"`
	HasRecurringIncome    bool      `json:"has_recurring_income"`
}

func NewCashFlowAnalysisDTO(entity *entities.CashFlowAnalysis) *CashFlowAnalysisDTO {
	return &CashFlowAnalysisDTO{
		ID:                    entity.ID,
		PersonID:              entity.PersonID,
		AnalysisDate:          entity.AnalysisDate,
		PeriodDays:            entity.PeriodDays,
		AverageMonthlyInflow:  entity.AverageMonthlyInflow,
		AverageMonthlyOutflow: entity.AverageMonthlyOutflow,
		NetCashFlow:           entity.NetCashFlow,
		InflowVolatility:      entity.InflowVolatility,
		NegativeBalanceDays:   entity.NegativeBalanceDays,
		HasRecurringIncome:    entity.HasRecurringIncome,
	}
}

func (d CashFlowAnalysisDTO) ToEntity() *entities.CashFlowAnalysis {
	return &entities.CashFlowAnalysis{
		PersonID:              d.PersonID,
		AnalysisDate:          d.AnalysisDate,
		PeriodDays:            d.PeriodDays,
		AverageMonthlyInflow:  d.AverageMonthlyInflow,
		AverageMonthlyOutflow: d.AverageMonthlyOutflow,
		NetCashFlow:           d.NetCashFlow,
		InflowVolatility:      d.InflowVolatility,
		NegativeBalanceDays:   d.NegativeBalanceDays,
		HasRecurringIncome:    d.HasRecurringIncome,
	}
}

type RecurringTransactionDTO struct {
	ID                 uint      `json:"id"`
	PersonID           uint      `json:"person_id"`
	TransactionType    string    `json:"transaction_type"`
	Category           string    `json:"category"`
	Description        string    `json:"description,omitempty"`
	Amount             float64   `json:"amount"`
	Frequency          string    `json:"frequency"`
	Counterparty       *string   `json:"counterparty,omitempty"`
	FirstDetectedDate  time.Time `json:"first_detected_date"`
	LastOccurrenceDate time.Time `json:"last_occurrence_date"`
	IsActive           bool      `json:"is_active"`
}

func NewRecurringTransactionDTO(entity *entities.RecurringTransaction) *RecurringTransactionDTO {
	return &RecurringTransactionDTO{
		ID:                 entity.ID,
		PersonID:           entity.PersonID,
		TransactionType:    entity.TransactionType,
		Category:           entity.Category,
		Description:        entity.Description,
		Amount:             entity.Amount,
		Frequency:          entity.Frequency,
		Counterparty:       entity.Counterparty,
		FirstDetectedDate:  entity.FirstDetectedDate,
		LastOccurrenceDate: entity.LastOccurrenceDate,
		IsActive:           entity.IsActive,
	}
}

func (d RecurringTransactionDTO) ToEntity() *entities.RecurringTransaction {
	return &entities.RecurringTransaction{
		PersonID:           d.PersonID,
		TransactionType:    d.TransactionType,
		Category:           d.Category,
		Description:        d.Description,
		Amount:             d.Amount,
		Frequency:          d.Frequency,
		Counterparty:       d.Counterparty,
		FirstDetectedDate:  d.FirstDetectedDate,
		LastOccurrenceDate: d.LastOccurrenceDate,
		IsActive:           d.IsActive,
	}
}

type DataSharingConsentDTO struct {
	ID          uint       `json:"id"`
	PersonID    uint       `json:"person_id"`
	ConsentID   string     `json:"consent_id"`
	Institution string     `json:"institution"`
	Status      string     `json:"status"`
	Scope       string     `json:"scope,omitempty"`
	GrantedAt   time.Time  `json:"granted_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	RevokedAt   *time.Time `json:"revoked_at,omitempty"`
}

func NewDataSharingConsentDTO(entity *entities.DataSharingConsent) *DataSharingConsentDTO {
	return &DataSharingConsentDTO{
		ID:          entity.ID,
		PersonID:    entity.PersonID,
		ConsentID:   entity.ConsentID,
		Institution: entity.Institution,
		Status:      entity.Status,
		Scope:       entity.Scope,
		GrantedAt:   entity.GrantedAt,
		ExpiresAt:   entity.ExpiresAt,
		RevokedAt:   entity.RevokedAt,
	}
}

func (d DataSharingConsentDTO) ToEntity() *entities.DataSharingConsent {
	return &entities.DataSharingConsent{
		PersonID:    d.PersonID,
		ConsentID:   d.ConsentID,
		Institution: d.Institution,
		Status:      d.Status,
		Scope:       d.Scope,
		GrantedAt:   d.GrantedAt,
		ExpiresAt:   d.ExpiresAt,
		RevokedAt:   d.RevokedAt,
	}
}
