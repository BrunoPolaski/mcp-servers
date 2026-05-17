package dto

import (
	"time"

	"github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities/value_objects"
)

type PersonDocumentDTO struct {
	ID             uint       `json:"id"`
	CreatedAt      string     `json:"created_at"`
	UpdatedAt      string     `json:"updated_at"`
	File           *FileDTO   `json:"file"`
	DocumentType   string     `json:"document_type"`
	IsVerified     bool       `json:"is_verified"`
	VerifiedAt     *time.Time `json:"verified_at,omitempty"`
	VerifiedBy     *string    `json:"verified_by,omitempty"`
	ExpirationDate *time.Time `json:"expiration_date,omitempty"`
}

func NewPersonDocumentDTO(entity *entities.PersonDocument) *PersonDocumentDTO {
	var file *FileDTO
	if entity.File != nil {
		file = NewFileDTO(entity.File)
	}

	return &PersonDocumentDTO{
		ID:             entity.ID,
		CreatedAt:      entity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      entity.UpdatedAt.Format("2006-01-02 15:04:05"),
		File:           file,
		DocumentType:   entity.DocumentType,
		IsVerified:     entity.IsVerified,
		VerifiedAt:     entity.VerifiedAt,
		VerifiedBy:     entity.VerifiedBy,
		ExpirationDate: entity.ExpirationDate,
	}
}

func (p PersonDocumentDTO) ToEntity() *entities.PersonDocument {
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

func phoneToStringPtr(phone *valueobjects.PhoneNumber) *string {
	if phone == nil {
		return nil
	}

	value := phone.String()
	return &value
}
