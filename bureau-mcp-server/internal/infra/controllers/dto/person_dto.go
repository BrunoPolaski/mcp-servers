package dto

import (
	"time"

	"github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities/value_objects"
)

type PersonDTO struct {
	ID                    uint                    `json:"id"`
	CreatedAt             string                  `json:"created_at"`
	UpdatedAt             string                  `json:"updated_at"`
	PersonalInformationID uint                    `json:"personal_information_id"`
	PersonalInformation   *PersonalInformationDTO `json:"personal_information" validate:"required"`
	CreditScoreID         *uint                   `json:"credit_score_id,omitempty"`
	CreditScore           *CreditScoreDTO         `json:"credit_score,omitempty"`
	FinancialProfileID    *uint                   `json:"financial_profile_id,omitempty"`
	FinancialProfile      *FinancialProfileDTO    `json:"financial_profile,omitempty"`
	EmploymentRecords     []EmploymentRecordDTO   `json:"employment_records,omitempty"`
	CreditAccounts        []CreditAccountDTO      `json:"credit_accounts,omitempty"`
	CreditInquiries       []CreditInquiryDTO      `json:"credit_inquiries,omitempty"`
	PaymentHistories      []PaymentHistoryDTO     `json:"payment_histories,omitempty"`
	Debts                 []DebtDTO               `json:"debts,omitempty"`
	NegativeRecords       []NegativeRecordDTO     `json:"negative_records,omitempty"`
	LegalRecords          []LegalRecordDTO        `json:"legal_records,omitempty"`
	FraudAlerts           []FraudAlertDTO         `json:"fraud_alerts,omitempty"`
	RiskAssessments       []RiskAssessmentDTO     `json:"risk_assessments,omitempty"`
	RelatedPersons        []PersonRelationshipDTO `json:"related_persons,omitempty"`
	DataSources           []DataSourceDTO         `json:"data_sources,omitempty"`
	LastVerifiedAt        *time.Time              `json:"last_verified_at,omitempty"`
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

	dto.CreditScoreID = entity.CreditScoreID
	if entity.CreditScore != nil {
		dto.CreditScore = NewCreditScoreDTO(entity.CreditScore)
	}
	dto.FinancialProfileID = entity.FinancialProfileID
	if entity.FinancialProfile != nil {
		dto.FinancialProfile = NewFinancialProfileDTO(entity.FinancialProfile)
	}

	if len(entity.EmploymentRecords) > 0 {
		dto.EmploymentRecords = make([]EmploymentRecordDTO, 0, len(entity.EmploymentRecords))
		for _, record := range entity.EmploymentRecords {
			dto.EmploymentRecords = append(dto.EmploymentRecords, *NewEmploymentRecordDTO(&record))
		}
	}

	if len(entity.CreditAccounts) > 0 {
		dto.CreditAccounts = make([]CreditAccountDTO, 0, len(entity.CreditAccounts))
		for _, account := range entity.CreditAccounts {
			dto.CreditAccounts = append(dto.CreditAccounts, *NewCreditAccountDTO(&account))
		}
	}

	if len(entity.CreditInquiries) > 0 {
		dto.CreditInquiries = make([]CreditInquiryDTO, 0, len(entity.CreditInquiries))
		for _, inquiry := range entity.CreditInquiries {
			dto.CreditInquiries = append(dto.CreditInquiries, *NewCreditInquiryDTO(&inquiry))
		}
	}

	if len(entity.PaymentHistories) > 0 {
		dto.PaymentHistories = make([]PaymentHistoryDTO, 0, len(entity.PaymentHistories))
		for _, history := range entity.PaymentHistories {
			dto.PaymentHistories = append(dto.PaymentHistories, *NewPaymentHistoryDTO(&history))
		}
	}

	if len(entity.Debts) > 0 {
		dto.Debts = make([]DebtDTO, 0, len(entity.Debts))
		for _, debt := range entity.Debts {
			dto.Debts = append(dto.Debts, *NewDebtDTO(&debt))
		}
	}

	if len(entity.NegativeRecords) > 0 {
		dto.NegativeRecords = make([]NegativeRecordDTO, 0, len(entity.NegativeRecords))
		for _, record := range entity.NegativeRecords {
			dto.NegativeRecords = append(dto.NegativeRecords, *NewNegativeRecordDTO(&record))
		}
	}

	if len(entity.LegalRecords) > 0 {
		dto.LegalRecords = make([]LegalRecordDTO, 0, len(entity.LegalRecords))
		for _, record := range entity.LegalRecords {
			dto.LegalRecords = append(dto.LegalRecords, *NewLegalRecordDTO(&record))
		}
	}

	if len(entity.FraudAlerts) > 0 {
		dto.FraudAlerts = make([]FraudAlertDTO, 0, len(entity.FraudAlerts))
		for _, alert := range entity.FraudAlerts {
			dto.FraudAlerts = append(dto.FraudAlerts, *NewFraudAlertDTO(&alert))
		}
	}

	if len(entity.RiskAssessments) > 0 {
		dto.RiskAssessments = make([]RiskAssessmentDTO, 0, len(entity.RiskAssessments))
		for _, assessment := range entity.RiskAssessments {
			dto.RiskAssessments = append(dto.RiskAssessments, *NewRiskAssessmentDTO(&assessment))
		}
	}

	if len(entity.RelatedPersons) > 0 {
		dto.RelatedPersons = make([]PersonRelationshipDTO, 0, len(entity.RelatedPersons))
		for _, relation := range entity.RelatedPersons {
			dto.RelatedPersons = append(dto.RelatedPersons, *NewPersonRelationshipDTO(&relation))
		}
	}

	if len(entity.DataSources) > 0 {
		dto.DataSources = make([]DataSourceDTO, 0, len(entity.DataSources))
		for _, source := range entity.DataSources {
			dto.DataSources = append(dto.DataSources, *NewDataSourceDTO(&source))
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

	var creditScore *entities.CreditScore
	if c.CreditScore != nil {
		creditScore = c.CreditScore.ToEntity()
	}

	var financialProfile *entities.FinancialProfile
	if c.FinancialProfile != nil {
		financialProfile = c.FinancialProfile.ToEntity()
	}

	var employmentRecords []entities.EmploymentRecord
	if len(c.EmploymentRecords) > 0 {
		employmentRecords = make([]entities.EmploymentRecord, 0, len(c.EmploymentRecords))
		for _, record := range c.EmploymentRecords {
			employmentRecords = append(employmentRecords, *record.ToEntity())
		}
	}

	var creditAccounts []entities.CreditAccount
	if len(c.CreditAccounts) > 0 {
		creditAccounts = make([]entities.CreditAccount, 0, len(c.CreditAccounts))
		for _, account := range c.CreditAccounts {
			creditAccounts = append(creditAccounts, *account.ToEntity())
		}
	}

	var creditInquiries []entities.CreditInquiry
	if len(c.CreditInquiries) > 0 {
		creditInquiries = make([]entities.CreditInquiry, 0, len(c.CreditInquiries))
		for _, inquiry := range c.CreditInquiries {
			creditInquiries = append(creditInquiries, *inquiry.ToEntity())
		}
	}

	var paymentHistories []entities.PaymentHistory
	if len(c.PaymentHistories) > 0 {
		paymentHistories = make([]entities.PaymentHistory, 0, len(c.PaymentHistories))
		for _, history := range c.PaymentHistories {
			paymentHistories = append(paymentHistories, *history.ToEntity())
		}
	}

	var debts []entities.Debt
	if len(c.Debts) > 0 {
		debts = make([]entities.Debt, 0, len(c.Debts))
		for _, debt := range c.Debts {
			debts = append(debts, *debt.ToEntity())
		}
	}

	var negativeRecords []entities.NegativeRecord
	if len(c.NegativeRecords) > 0 {
		negativeRecords = make([]entities.NegativeRecord, 0, len(c.NegativeRecords))
		for _, record := range c.NegativeRecords {
			negativeRecords = append(negativeRecords, *record.ToEntity())
		}
	}

	var legalRecords []entities.LegalRecord
	if len(c.LegalRecords) > 0 {
		legalRecords = make([]entities.LegalRecord, 0, len(c.LegalRecords))
		for _, record := range c.LegalRecords {
			legalRecords = append(legalRecords, *record.ToEntity())
		}
	}

	var fraudAlerts []entities.FraudAlert
	if len(c.FraudAlerts) > 0 {
		fraudAlerts = make([]entities.FraudAlert, 0, len(c.FraudAlerts))
		for _, alert := range c.FraudAlerts {
			fraudAlerts = append(fraudAlerts, *alert.ToEntity())
		}
	}

	var riskAssessments []entities.RiskAssessment
	if len(c.RiskAssessments) > 0 {
		riskAssessments = make([]entities.RiskAssessment, 0, len(c.RiskAssessments))
		for _, assessment := range c.RiskAssessments {
			riskAssessments = append(riskAssessments, *assessment.ToEntity())
		}
	}

	var relatedPersons []entities.PersonRelationship
	if len(c.RelatedPersons) > 0 {
		relatedPersons = make([]entities.PersonRelationship, 0, len(c.RelatedPersons))
		for _, relation := range c.RelatedPersons {
			relatedPersons = append(relatedPersons, *relation.ToEntity())
		}
	}

	var dataSources []entities.DataSource
	if len(c.DataSources) > 0 {
		dataSources = make([]entities.DataSource, 0, len(c.DataSources))
		for _, source := range c.DataSources {
			dataSources = append(dataSources, *source.ToEntity())
		}
	}

	return &entities.Person{
		PersonalInformationID: c.PersonalInformationID,
		PersonalInformation:   pi,
		CreditScoreID:         c.CreditScoreID,
		CreditScore:           creditScore,
		FinancialProfileID:    c.FinancialProfileID,
		FinancialProfile:      financialProfile,
		EmploymentRecords:     employmentRecords,
		CreditAccounts:        creditAccounts,
		CreditInquiries:       creditInquiries,
		PaymentHistories:      paymentHistories,
		Debts:                 debts,
		NegativeRecords:       negativeRecords,
		LegalRecords:          legalRecords,
		FraudAlerts:           fraudAlerts,
		RiskAssessments:       riskAssessments,
		RelatedPersons:        relatedPersons,
		DataSources:           dataSources,
		LastVerifiedAt:        c.LastVerifiedAt,
	}
}

type ListPersonsDTO struct {
	PaginatedDTO
	Document valueobjects.Document `json:"document"`
}
