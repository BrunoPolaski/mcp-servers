package dto

import (
	"time"

	"github.com/BrunoPolaski/bureau/internal/core/entities"
)

type FinancialProfileDTO struct {
	ID                     uint      `json:"id"`
	CreatedAt              string    `json:"created_at"`
	UpdatedAt              string    `json:"updated_at"`
	PersonID               uint      `json:"person_id"`
	ProfileDate            time.Time `json:"profile_date"`
	DeclaredMonthlyIncome  *float64  `json:"declared_monthly_income,omitempty"`
	EstimatedMonthlyIncome *float64  `json:"estimated_monthly_income,omitempty"`
	IncomeSource           *string   `json:"income_source,omitempty"`
	TotalAssets            *float64  `json:"total_assets,omitempty"`
	RealEstateValue        *float64  `json:"real_estate_value,omitempty"`
	VehiclesValue          *float64  `json:"vehicles_value,omitempty"`
	TotalLiabilities       *float64  `json:"total_liabilities,omitempty"`
	TotalMonthlyPayments   *float64  `json:"total_monthly_payments,omitempty"`
	DebtToIncomeRatio      *float64  `json:"debt_to_income_ratio,omitempty"`
	AvailableCredit        *float64  `json:"available_credit,omitempty"`
	CreditUtilization      *float64  `json:"credit_utilization,omitempty"`
}

func NewFinancialProfileDTO(entity *entities.FinancialProfile) *FinancialProfileDTO {
	return &FinancialProfileDTO{
		ID:                     entity.ID,
		CreatedAt:              entity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:              entity.UpdatedAt.Format("2006-01-02 15:04:05"),
		PersonID:               entity.PersonID,
		ProfileDate:            entity.ProfileDate,
		DeclaredMonthlyIncome:  entity.DeclaredMonthlyIncome,
		EstimatedMonthlyIncome: entity.EstimatedMonthlyIncome,
		IncomeSource:           entity.IncomeSource,
		TotalAssets:            entity.TotalAssets,
		RealEstateValue:        entity.RealEstateValue,
		VehiclesValue:          entity.VehiclesValue,
		TotalLiabilities:       entity.TotalLiabilities,
		TotalMonthlyPayments:   entity.TotalMonthlyPayments,
		DebtToIncomeRatio:      entity.DebtToIncomeRatio,
		AvailableCredit:        entity.AvailableCredit,
		CreditUtilization:      entity.CreditUtilization,
	}
}

func (f FinancialProfileDTO) ToEntity() *entities.FinancialProfile {
	return &entities.FinancialProfile{
		PersonID:               f.PersonID,
		ProfileDate:            f.ProfileDate,
		DeclaredMonthlyIncome:  f.DeclaredMonthlyIncome,
		EstimatedMonthlyIncome: f.EstimatedMonthlyIncome,
		IncomeSource:           f.IncomeSource,
		TotalAssets:            f.TotalAssets,
		RealEstateValue:        f.RealEstateValue,
		VehiclesValue:          f.VehiclesValue,
		TotalLiabilities:       f.TotalLiabilities,
		TotalMonthlyPayments:   f.TotalMonthlyPayments,
		DebtToIncomeRatio:      f.DebtToIncomeRatio,
		AvailableCredit:        f.AvailableCredit,
		CreditUtilization:      f.CreditUtilization,
	}
}
