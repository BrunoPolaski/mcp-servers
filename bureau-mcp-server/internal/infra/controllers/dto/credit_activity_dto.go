package dto

import (
	"time"

	"github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities"
)

type CreditAccountDTO struct {
	ID                uint       `json:"id"`
	CreatedAt         string     `json:"created_at"`
	UpdatedAt         string     `json:"updated_at"`
	PersonID          uint       `json:"person_id"`
	AccountType       string     `json:"account_type"`
	Creditor          string     `json:"creditor"`
	CreditorDocument  string     `json:"creditor_document"`
	AccountNumber     string     `json:"account_number"`
	OpenedDate        time.Time  `json:"opened_date"`
	ClosedDate        *time.Time `json:"closed_date,omitempty"`
	Status            string     `json:"status"`
	CreditLimit       *float64   `json:"credit_limit,omitempty"`
	CurrentBalance    float64    `json:"current_balance"`
	AvailableCredit   *float64   `json:"available_credit,omitempty"`
	OriginalAmount    *float64   `json:"original_amount,omitempty"`
	RemainingAmount   *float64   `json:"remaining_amount,omitempty"`
	InterestRate      *float64   `json:"interest_rate,omitempty"`
	MonthlyPayment    *float64   `json:"monthly_payment,omitempty"`
	PaymentDueDay     *int       `json:"payment_due_day,omitempty"`
	NumberOfPayments  *int       `json:"number_of_payments,omitempty"`
	RemainingPayments *int       `json:"remaining_payments,omitempty"`
	PaymentStatus     string     `json:"payment_status"`
	DaysLate          int        `json:"days_late"`
	HighestDaysLate   int        `json:"highest_days_late"`
	TimesLate30Days   int        `json:"times_late_30_days"`
	TimesLate60Days   int        `json:"times_late_60_days"`
	TimesLate90Days   int        `json:"times_late_90_days"`
	LastReportedDate  time.Time  `json:"last_reported_date"`
}

func NewCreditAccountDTO(entity *entities.CreditAccount) *CreditAccountDTO {
	return &CreditAccountDTO{
		ID:                entity.ID,
		CreatedAt:         entity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:         entity.UpdatedAt.Format("2006-01-02 15:04:05"),
		PersonID:          entity.PersonID,
		AccountType:       entity.AccountType,
		Creditor:          entity.Creditor,
		CreditorDocument:  entity.CreditorDocument,
		AccountNumber:     entity.AccountNumber,
		OpenedDate:        entity.OpenedDate,
		ClosedDate:        entity.ClosedDate,
		Status:            entity.Status,
		CreditLimit:       entity.CreditLimit,
		CurrentBalance:    entity.CurrentBalance,
		AvailableCredit:   entity.AvailableCredit,
		OriginalAmount:    entity.OriginalAmount,
		RemainingAmount:   entity.RemainingAmount,
		InterestRate:      entity.InterestRate,
		MonthlyPayment:    entity.MonthlyPayment,
		PaymentDueDay:     entity.PaymentDueDay,
		NumberOfPayments:  entity.NumberOfPayments,
		RemainingPayments: entity.RemainingPayments,
		PaymentStatus:     entity.PaymentStatus,
		DaysLate:          entity.DaysLate,
		HighestDaysLate:   entity.HighestDaysLate,
		TimesLate30Days:   entity.TimesLate30Days,
		TimesLate60Days:   entity.TimesLate60Days,
		TimesLate90Days:   entity.TimesLate90Days,
		LastReportedDate:  entity.LastReportedDate,
	}
}

func (c CreditAccountDTO) ToEntity() *entities.CreditAccount {
	return &entities.CreditAccount{
		PersonID:          c.PersonID,
		AccountType:       c.AccountType,
		Creditor:          c.Creditor,
		CreditorDocument:  c.CreditorDocument,
		AccountNumber:     c.AccountNumber,
		OpenedDate:        c.OpenedDate,
		ClosedDate:        c.ClosedDate,
		Status:            c.Status,
		CreditLimit:       c.CreditLimit,
		CurrentBalance:    c.CurrentBalance,
		AvailableCredit:   c.AvailableCredit,
		OriginalAmount:    c.OriginalAmount,
		RemainingAmount:   c.RemainingAmount,
		InterestRate:      c.InterestRate,
		MonthlyPayment:    c.MonthlyPayment,
		PaymentDueDay:     c.PaymentDueDay,
		NumberOfPayments:  c.NumberOfPayments,
		RemainingPayments: c.RemainingPayments,
		PaymentStatus:     c.PaymentStatus,
		DaysLate:          c.DaysLate,
		HighestDaysLate:   c.HighestDaysLate,
		TimesLate30Days:   c.TimesLate30Days,
		TimesLate60Days:   c.TimesLate60Days,
		TimesLate90Days:   c.TimesLate90Days,
		LastReportedDate:  c.LastReportedDate,
	}
}

type CreditInquiryDTO struct {
	ID               uint      `json:"id"`
	CreatedAt        string    `json:"created_at"`
	UpdatedAt        string    `json:"updated_at"`
	PersonID         uint      `json:"person_id"`
	InquiryDate      time.Time `json:"inquiry_date"`
	InquiryType      string    `json:"inquiry_type"`
	Creditor         string    `json:"creditor"`
	CreditorDocument string    `json:"creditor_document"`
	Purpose          string    `json:"purpose"`
	Amount           *float64  `json:"amount,omitempty"`
	Result           *string   `json:"result,omitempty"`
}

func NewCreditInquiryDTO(entity *entities.CreditInquiry) *CreditInquiryDTO {
	return &CreditInquiryDTO{
		ID:               entity.ID,
		CreatedAt:        entity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        entity.UpdatedAt.Format("2006-01-02 15:04:05"),
		PersonID:         entity.PersonID,
		InquiryDate:      entity.InquiryDate,
		InquiryType:      entity.InquiryType,
		Creditor:         entity.Creditor,
		CreditorDocument: entity.CreditorDocument,
		Purpose:          entity.Purpose,
		Amount:           entity.Amount,
		Result:           entity.Result,
	}
}

func (c CreditInquiryDTO) ToEntity() *entities.CreditInquiry {
	return &entities.CreditInquiry{
		PersonID:         c.PersonID,
		InquiryDate:      c.InquiryDate,
		InquiryType:      c.InquiryType,
		Creditor:         c.Creditor,
		CreditorDocument: c.CreditorDocument,
		Purpose:          c.Purpose,
		Amount:           c.Amount,
		Result:           c.Result,
	}
}

type PaymentHistoryDTO struct {
	ID              uint      `json:"id"`
	CreatedAt       string    `json:"created_at"`
	UpdatedAt       string    `json:"updated_at"`
	PersonID        uint      `json:"person_id"`
	CreditAccountID *uint     `json:"credit_account_id,omitempty"`
	DebtID          *uint     `json:"debt_id,omitempty"`
	PaymentDate     time.Time `json:"payment_date"`
	DueDate         time.Time `json:"due_date"`
	Amount          float64   `json:"amount"`
	AmountDue       float64   `json:"amount_due"`
	Status          string    `json:"status"`
	DaysLate        int       `json:"days_late"`
}

func NewPaymentHistoryDTO(entity *entities.PaymentHistory) *PaymentHistoryDTO {
	return &PaymentHistoryDTO{
		ID:              entity.ID,
		CreatedAt:       entity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       entity.UpdatedAt.Format("2006-01-02 15:04:05"),
		PersonID:        entity.PersonID,
		CreditAccountID: entity.CreditAccountID,
		DebtID:          entity.DebtID,
		PaymentDate:     entity.PaymentDate,
		DueDate:         entity.DueDate,
		Amount:          entity.Amount,
		AmountDue:       entity.AmountDue,
		Status:          entity.Status,
		DaysLate:        entity.DaysLate,
	}
}

func (p PaymentHistoryDTO) ToEntity() *entities.PaymentHistory {
	return &entities.PaymentHistory{
		PersonID:        p.PersonID,
		CreditAccountID: p.CreditAccountID,
		DebtID:          p.DebtID,
		PaymentDate:     p.PaymentDate,
		DueDate:         p.DueDate,
		Amount:          p.Amount,
		AmountDue:       p.AmountDue,
		Status:          p.Status,
		DaysLate:        p.DaysLate,
	}
}
