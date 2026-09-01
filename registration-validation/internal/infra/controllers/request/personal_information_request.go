package request

import (
	"time"

	"github.com/BrunoPolaski/registration-validation/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/registration-validation/internal/core/entities/value_objects"
)

type PersonalInformationRequest struct {
	ID                   uint                    `json:"id"`
	FullName             string                  `json:"full_name" validate:"required" example:"João da Silva"`
	MotherName           *string                 `json:"mother_name,omitempty"`
	BirthDate            *time.Time              `json:"birth_date,omitempty"`
	Gender               *string                 `json:"gender,omitempty"`
	Nationality          *string                 `json:"nationality,omitempty"`
	MaritalStatus        *string                 `json:"marital_status,omitempty"`
	Document             string                  `json:"document" validate:"required,document" example:"12345678909"`
	RG                   *string                 `json:"rg,omitempty"`
	RGIssuer             *string                 `json:"rg_issuer,omitempty"`
	RGIssueDate          *time.Time              `json:"rg_issue_date,omitempty"`
	VoterID              *string                 `json:"voter_id,omitempty"`
	WorkCard             *string                 `json:"work_card,omitempty"`
	PrimaryPhone         *string                 `json:"primary_phone,omitempty" validate:"omitempty,phone_number" example:"11999999999"`
	SecondaryPhone       *string                 `json:"secondary_phone,omitempty"`
	Email                *string                 `json:"email,omitempty" validate:"omitempty,email"`
	AlternativeEmail     *string                 `json:"alternative_email,omitempty"`
	Addresses            []PersonAddressRequest  `json:"addresses,omitempty"`
	ProfilePhoto         *FileRequest            `json:"profile_photo,omitempty"`
	Documents            []PersonDocumentRequest `json:"documents,omitempty"`
	DocumentValidated    bool                    `json:"document_validated"`
	EmailVerified        bool                    `json:"email_verified"`
	PhoneVerified        bool                    `json:"phone_verified"`
	BiometricValidated   bool                    `json:"biometric_validated"`
	ReceitaFederalStatus string                  `json:"receita_federal_status"`
}

func (p PersonalInformationRequest) ToEntity() *entities.PersonalInformation {
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

type PersonAddressRequest struct {
	AddressID uint `json:"address_id" validate:"required"`
}

func (p PersonAddressRequest) ToEntity() entities.PersonAddress {
	return entities.PersonAddress{
		AddressID: p.AddressID,
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

type PersonDocumentRequest struct {
	ID             uint         `json:"id"`
	File           *FileRequest `json:"file"`
	DocumentType   string       `json:"document_type" validate:"required"`
	IsVerified     bool         `json:"is_verified"`
	VerifiedAt     *time.Time   `json:"verified_at,omitempty"`
	VerifiedBy     *string      `json:"verified_by,omitempty"`
	ExpirationDate *time.Time   `json:"expiration_date,omitempty"`
}

func (p PersonDocumentRequest) ToEntity() *entities.PersonDocument {
	var file *entities.File
	if p.File != nil {
		file = p.File.ToEntity()
	}

	return &entities.PersonDocument{
		File:           file,
		DocumentType:   p.DocumentType,
		IsVerified:     p.IsVerified,
		VerifiedAt:     p.VerifiedAt,
		VerifiedBy:     p.VerifiedBy,
		ExpirationDate: p.ExpirationDate,
	}
}
