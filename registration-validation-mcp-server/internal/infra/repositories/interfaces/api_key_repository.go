package interfaces

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/registration-validation-mcp-server/internal/core/entities"
)

type ApiKeyRepository interface {
	Create(ctx context.Context, apiKey *entities.ApiKey) *rest_err.RestErr
	GetById(ctx context.Context, uuid string) (*entities.ApiKey, *rest_err.RestErr)
	GetAll(ctx context.Context, limit, offset int, params map[string]any) ([]entities.ApiKey, int, *rest_err.RestErr)
	Delete(ctx context.Context, uuid string) *rest_err.RestErr
}
