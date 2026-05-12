package dto

import (
	"github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities/value_objects"
)

type PersonDTO struct {
	ID                  uint                    `json:"id"`
	CreatedAt           string                  `json:"created_at"`
	UpdatedAt           string                  `json:"updated_at"`
	PersonalInformation *PersonalInformationDTO `json:"personal_information" validate:"required"`
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

	return dto
}

func (c PersonDTO) ToEntity() *entities.Person {
	pi := c.PersonalInformation.ToEntity()

	return &entities.Person{
		PersonalInformation: pi,
	}
}

type ListPersonsDTO struct {
	PaginatedDTO
	Document valueobjects.Document `json:"document"`
}
