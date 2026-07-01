package dto

import "github.com/BrunoPolaski/open-finance-mcp-server/internal/core/entities"

type AnalystDTO struct {
	ID                  uint                    `json:"id"`
	CreatedAt           string                  `json:"created_at"`
	UpdatedAt           string                  `json:"updated_at"`
	PersonalInformation *PersonalInformationDTO `json:"personal_information" validate:"required"`
}

func NewAnalystDTO(entity *entities.Analyst) *AnalystDTO {
	analyst := &AnalystDTO{
		ID:        entity.ID,
		CreatedAt: entity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: entity.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if entity.PersonalInformation != nil {
		analyst.PersonalInformation = NewPersonalInformationDTO(entity.PersonalInformation)
	}

	return analyst
}

func (e AnalystDTO) ToEntity() *entities.Analyst {
	analyst := &entities.Analyst{
		PersonalInformation: e.PersonalInformation.ToEntity(),
	}

	return analyst
}
