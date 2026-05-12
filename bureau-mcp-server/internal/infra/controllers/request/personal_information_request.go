package request

import (
	"github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities/value_objects"
)

type PersonalInformationRequest struct {
	ID       uint            `json:"id"`
	Name     string          `json:"name" validate:"required" example:"João da Silva"`             // Full name
	Phone    string          `json:"phone" validate:"required,phone_number" example:"11999999999"` // Phone in brazilian format, without country code
	Document string          `json:"document" validate:"required,document" example:"12345678909"`  // Brazilian CPF
	Address  *AddressRequest `json:"address" validate:"required"`
	File     *FileRequest    `json:"file"`
}

func (p PersonalInformationRequest) ToEntity() *entities.PersonalInformation {
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

type AddressRequest struct {
	Id         uint    `json:"id"`
	Street     string  `json:"street" validate:"required"`
	Number     string  `json:"number" validate:"required" example:"1234 | S/N"`
	City       string  `json:"city" validate:"required"`
	State      string  `json:"state" validate:"required" example:"SC"`
	ZipCode    string  `json:"zip_code" validate:"required" example:"12345-678"`
	Complement *string `json:"complement"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

func (a AddressRequest) ToEntity() *entities.Address {
	return &entities.Address{
		Street:     a.Street,
		Number:     a.Number,
		City:       a.City,
		State:      a.State,
		ZipCode:    a.ZipCode,
		Complement: a.Complement,
	}
}

type FileRequest struct {
	ID           uint   `json:"id"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	OriginalName string `json:"original_name" validate:"required"`
	Name         string `json:"name"`
	URL          string `json:"url" validate:"required,url" example:"https://example.com/file.jpg"`
	MimeType     string `json:"mime_type"`
}

func (f FileRequest) ToEntity() *entities.File {
	return &entities.File{
		OriginalName: f.OriginalName,
		Name:         f.Name,
		URL:          f.URL,
		MimeType:     f.MimeType,
	}
}
