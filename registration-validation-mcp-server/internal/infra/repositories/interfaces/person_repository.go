package interfaces

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/registration-validation-mcp-server/internal/core/entities"
)

type PersonRepository interface {
	GetById(ctx context.Context, id uint) (*entities.Person, *rest_err.RestErr)
	GetByDocument(ctx context.Context, document string) (*entities.Person, *rest_err.RestErr)
	GetAll(ctx context.Context, limit, offset int, params map[string]any) ([]entities.Person, int64, *rest_err.RestErr)
	Delete(ctx context.Context, id uint) *rest_err.RestErr
}
