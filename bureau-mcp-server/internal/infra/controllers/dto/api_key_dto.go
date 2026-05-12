package dto

import (
	"github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities"
)

type ApiKeyDTO struct {
	UUID      string `json:"uuid" validate:"required,uuid4" example:"123e4567-e89b-12d3-a456-426614174000"`
	Slug      string `json:"slug" validate:"required"`
	CreatedAt string `json:"created_at"`
}

func NewApiKeyDTO(entity *entities.ApiKey) *ApiKeyDTO {
	return &ApiKeyDTO{
		UUID:      entity.UUID,
		Slug:      entity.Slug,
		CreatedAt: entity.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (a ApiKeyDTO) ToEntity() *entities.ApiKey {
	return &entities.ApiKey{
		UUID: a.UUID,
		Slug: a.Slug,
	}
}
