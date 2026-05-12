package dto

import (
	"github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities"
)

type AdminDTO struct {
	ID                  uint                    `json:"id"`
	CreatedAt           string                  `json:"created_at"`
	UpdatedAt           string                  `json:"updated_at"`
	PersonalInformation *PersonalInformationDTO `json:"personal_information" validate:"required"`
}

func NewAdminDTO(entity *entities.Admin) *AdminDTO {
	dto := &AdminDTO{
		ID:        entity.ID,
		CreatedAt: entity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: entity.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if entity.PersonalInformation != nil {
		dto.PersonalInformation = NewPersonalInformationDTO(entity.PersonalInformation)
	}

	return dto
}

func (a AdminDTO) ToEntity() *entities.Admin {
	return &entities.Admin{
		PersonalInformation: a.PersonalInformation.ToEntity(),
	}
}
