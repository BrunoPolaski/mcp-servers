package dto

import (
	"github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities"
)

type AddressDTO struct {
	Id         uint    `json:"id"`
	Street     *string `json:"street,omitempty"`
	Number     *string `json:"number,omitempty" example:"1234 | S/N"`
	City       *string `json:"city,omitempty"`
	State      *string `json:"state,omitempty" example:"SC"`
	ZipCode    *string `json:"zip_code,omitempty" example:"12345-678"`
	Complement *string `json:"complement,omitempty"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

func NewAddressDTO(entity *entities.Address) *AddressDTO {
	return &AddressDTO{
		Id:         entity.ID,
		Street:     entity.Street,
		Number:     entity.Number,
		City:       entity.City,
		State:      entity.State,
		ZipCode:    entity.ZipCode,
		Complement: entity.Complement,
		CreatedAt:  entity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  entity.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (a AddressDTO) ToEntity() *entities.Address {
	return &entities.Address{
		Street:     a.Street,
		Number:     a.Number,
		City:       a.City,
		State:      a.State,
		ZipCode:    a.ZipCode,
		Complement: a.Complement,
	}
}
