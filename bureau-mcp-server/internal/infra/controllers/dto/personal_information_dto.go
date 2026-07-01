package dto

import (
	"time"

	"github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities/value_objects"
)

type PersonalInformationDTO struct {
	ID               uint                `json:"id"`
	FullName         string              `json:"full_name" validate:"required" example:"João da Silva"`
	MotherName       *string             `json:"mother_name,omitempty"`
	BirthDate        *time.Time          `json:"birth_date,omitempty"`
	Gender           *string             `json:"gender,omitempty"`
	Nationality      *string             `json:"nationality,omitempty"`
	MaritalStatus    *string             `json:"marital_status,omitempty"`
	Document         string              `json:"document" example:"12345678909"`
	RG               *string             `json:"rg,omitempty"`
	RGIssuer         *string             `json:"rg_issuer,omitempty"`
	RGIssueDate      *time.Time          `json:"rg_issue_date,omitempty"`
	VoterID          *string             `json:"voter_id,omitempty"`
	WorkCard         *string             `json:"work_card,omitempty"`
	PrimaryPhone     *string             `json:"primary_phone,omitempty" example:"11999999999"`
	SecondaryPhone   *string             `json:"secondary_phone,omitempty"`
	Email            *string             `json:"email,omitempty"`
	AlternativeEmail *string             `json:"alternative_email,omitempty"`
	Addresses        []PersonAddressDTO  `json:"addresses,omitempty"`
	ProfilePhoto     *FileDTO            `json:"profile_photo,omitempty"`
	Documents        []PersonDocumentDTO `json:"documents,omitempty"`
	EmailVerified    bool                `json:"email_verified"`
	PhoneVerified    bool                `json:"phone_verified"`
}

func NewPersonalInformationDTO(entity *entities.PersonalInformation) *PersonalInformationDTO {
	piDTO := &PersonalInformationDTO{
		ID:               entity.ID,
		FullName:         entity.FullName,
		MotherName:       entity.MotherName,
		BirthDate:        entity.BirthDate,
		Gender:           entity.Gender,
		Nationality:      entity.Nationality,
		MaritalStatus:    entity.MaritalStatus,
		Document:         entity.Document.String(),
		RG:               entity.RG,
		RGIssuer:         entity.RGIssuer,
		RGIssueDate:      entity.RGIssueDate,
		VoterID:          entity.VoterID,
		WorkCard:         entity.WorkCard,
		PrimaryPhone:     phoneToStringPtr(entity.PrimaryPhone),
		SecondaryPhone:   phoneToStringPtr(entity.SecondaryPhone),
		Email:            entity.Email,
		AlternativeEmail: entity.AlternativeEmail,
		EmailVerified:    entity.EmailVerified,
		PhoneVerified:    entity.PhoneVerified,
	}

	if len(entity.Addresses) > 0 {
		piDTO.Addresses = make([]PersonAddressDTO, 0, len(entity.Addresses))
		for _, address := range entity.Addresses {
			piDTO.Addresses = append(piDTO.Addresses, NewPersonAddressDTO(&address))
		}
	}

	if entity.ProfilePhoto != nil {
		piDTO.ProfilePhoto = NewFileDTO(entity.ProfilePhoto)
	}

	if len(entity.Documents) > 0 {
		piDTO.Documents = make([]PersonDocumentDTO, 0, len(entity.Documents))
		for _, doc := range entity.Documents {
			piDTO.Documents = append(piDTO.Documents, *NewPersonDocumentDTO(&doc))
		}
	}

	return piDTO
}

func (p PersonalInformationDTO) ToEntity() *entities.PersonalInformation {
	var primaryPhone *valueobjects.PhoneNumber
	if p.PrimaryPhone != nil {
		phone, err := valueobjects.NewPhoneNumber(*p.PrimaryPhone)
		if err != nil {
			return nil
		}
		primaryPhone = &phone
	}

	var secondaryPhone *valueobjects.PhoneNumber
	if p.SecondaryPhone != nil {
		phone, err := valueobjects.NewPhoneNumber(*p.SecondaryPhone)
		if err != nil {
			return nil
		}
		secondaryPhone = &phone
	}

	pi := &entities.PersonalInformation{
		FullName:         p.FullName,
		MotherName:       p.MotherName,
		BirthDate:        p.BirthDate,
		Gender:           p.Gender,
		Nationality:      p.Nationality,
		MaritalStatus:    p.MaritalStatus,
		Document:         valueobjects.NewDocument(p.Document),
		RG:               p.RG,
		RGIssuer:         p.RGIssuer,
		RGIssueDate:      p.RGIssueDate,
		VoterID:          p.VoterID,
		WorkCard:         p.WorkCard,
		PrimaryPhone:     primaryPhone,
		SecondaryPhone:   secondaryPhone,
		Email:            p.Email,
		AlternativeEmail: p.AlternativeEmail,
		EmailVerified:    p.EmailVerified,
		PhoneVerified:    p.PhoneVerified,
	}

	if len(p.Addresses) > 0 {
		pi.Addresses = make([]entities.PersonAddress, 0, len(p.Addresses))
		for _, address := range p.Addresses {
			pi.Addresses = append(pi.Addresses, address.ToEntity())
		}
	}

	if p.ProfilePhoto != nil {
		pi.ProfilePhoto = p.ProfilePhoto.ToEntity()
	}

	if len(p.Documents) > 0 {
		pi.Documents = make([]entities.PersonDocument, 0, len(p.Documents))
		for _, doc := range p.Documents {
			if entityDoc := doc.ToEntity(); entityDoc != nil {
				pi.Documents = append(pi.Documents, *entityDoc)
			}
		}
	}

	return pi
}

type PersonAddressDTO struct {
	ID        uint `json:"id"`
	AddressID uint `json:"address_id"`
}

func NewPersonAddressDTO(entity *entities.PersonAddress) PersonAddressDTO {
	return PersonAddressDTO{
		ID:        entity.ID,
		AddressID: entity.AddressID,
	}
}

func (p PersonAddressDTO) ToEntity() entities.PersonAddress {
	return entities.PersonAddress{
		AddressID: p.AddressID,
	}
}
