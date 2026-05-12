package dto

import (
	"github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities/value_objects"
)

type PersonalInformationDTO struct {
	ID       uint        `json:"id"`
	Name     string      `json:"name" validate:"required" example:"João da Silva"` // Full name
	Phone    string      `json:"phone" example:"11999999999"`                      // Phone in brazilian format, without country code
	Document string      `json:"document" example:"12345678909"`                   // Brazilian CPF
	Address  *AddressDTO `json:"address"`
	File     *FileDTO    `json:"file"`
}

func NewPersonalInformationDTO(entity *entities.PersonalInformation) *PersonalInformationDTO {
	piDTO := &PersonalInformationDTO{
		ID:       entity.ID,
		Name:     entity.Name,
		Phone:    entity.Phone.String(),
		Document: entity.Document.String(),
	}

	if entity.Address != nil {
		piDTO.Address = NewAddressDTO(entity.Address)
	}

	if entity.File != nil {
		piDTO.File = NewFileDTO(entity.File)
	}

	return piDTO
}

func (p PersonalInformationDTO) ToEntity() *entities.PersonalInformation {
	phone, err := valueobjects.NewPhoneNumber(p.Phone)
	if err != nil {
		return nil
	}

	pi := &entities.PersonalInformation{
		Name:     p.Name,
		Phone:    phone,
		Document: valueobjects.NewDocument(p.Document),
	}

	if p.Address != nil {
		pi.Address = p.Address.ToEntity()
	}

	if p.File != nil {
		pi.File = p.File.ToEntity()
	}

	return pi
}
