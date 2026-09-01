package interfaces

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
)

type AnalystRepository interface {
	GetById(ctx context.Context, id uint) (*entities.Analyst, *rest_err.RestErr)
	GetAll(ctx context.Context, limit, offset int, params map[string]any) ([]entities.Analyst, int64, *rest_err.RestErr)
	Delete(ctx context.Context, id uint) *rest_err.RestErr
}
