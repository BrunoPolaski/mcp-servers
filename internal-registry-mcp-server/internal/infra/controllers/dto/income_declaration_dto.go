package dto

import (
	"time"

	"github.com/BrunoPolaski/internal-registry-mcp-server/internal/core/entities"
)

type IncomeDeclarationDTO struct {
	ID              uint      `json:"id"`
	CreatedAt       string    `json:"created_at"`
	UpdatedAt       string    `json:"updated_at"`
	PersonID        uint      `json:"person_id"`
	DeclarationDate time.Time `json:"declaration_date"`
	IncomeType      string    `json:"income_type"`
	MonthlyAmount   float64   `json:"monthly_amount"`
	YearlyAmount    *float64  `json:"yearly_amount,omitempty"`
	Source          *string   `json:"source,omitempty"`
	Verified        bool      `json:"verified"`
	VerifiedBy      *string   `json:"verified_by,omitempty"`
	ProofFileID     *uint     `json:"proof_file_id,omitempty"`
	ProofFile       *FileDTO  `json:"proof_file,omitempty"`
}

func NewIncomeDeclarationDTO(entity *entities.IncomeDeclaration) *IncomeDeclarationDTO {
	dto := &IncomeDeclarationDTO{
		ID:              entity.ID,
		CreatedAt:       entity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       entity.UpdatedAt.Format("2006-01-02 15:04:05"),
		PersonID:        entity.PersonID,
		DeclarationDate: entity.DeclarationDate,
		IncomeType:      entity.IncomeType,
		MonthlyAmount:   entity.MonthlyAmount,
		YearlyAmount:    entity.YearlyAmount,
		Source:          entity.Source,
		Verified:        entity.Verified,
		VerifiedBy:      entity.VerifiedBy,
		ProofFileID:     entity.ProofFileID,
	}

	if entity.ProofFile != nil {
		dto.ProofFile = NewFileDTO(entity.ProofFile)
	}

	return dto
}

func (i IncomeDeclarationDTO) ToEntity() *entities.IncomeDeclaration {
	var file *entities.File
	if i.ProofFile != nil {
		file = i.ProofFile.ToEntity()
	}

	return &entities.IncomeDeclaration{
		PersonID:        i.PersonID,
		DeclarationDate: i.DeclarationDate,
		IncomeType:      i.IncomeType,
		MonthlyAmount:   i.MonthlyAmount,
		YearlyAmount:    i.YearlyAmount,
		Source:          i.Source,
		Verified:        i.Verified,
		VerifiedBy:      i.VerifiedBy,
		ProofFileID:     i.ProofFileID,
		ProofFile:       file,
	}
}
