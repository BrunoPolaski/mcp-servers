package dto

import (
	"github.com/BrunoPolaski/registration-validation-mcp-server/internal/core/entities"
)

type FileDTO struct {
	ID           uint   `json:"id"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	OriginalName string `json:"original_name" validate:"required"`
	Name         string `json:"name"`
	URL          string `json:"url" validate:"required,url" example:"https://example.com/file.jpg"`
	MimeType     string `json:"mime_type"`
}

func NewFileDTO(entity *entities.File) *FileDTO {
	return &FileDTO{
		ID:           entity.ID,
		CreatedAt:    entity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    entity.UpdatedAt.Format("2006-01-02 15:04:05"),
		OriginalName: entity.OriginalName,
		Name:         entity.Name,
		URL:          entity.URL,
		MimeType:     entity.MimeType,
	}
}

func (f FileDTO) ToEntity() *entities.File {
	return &entities.File{
		OriginalName: f.OriginalName,
		Name:         f.Name,
		URL:          f.URL,
		MimeType:     f.MimeType,
	}
}
