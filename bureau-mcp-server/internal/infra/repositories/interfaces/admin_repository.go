package interfaces

import (
	"context"

	"github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities"
	"github.com/BrunoPolaski/go-rest-err/rest_err"
)

type AdminRepository interface {
	GetById(ctx context.Context, id uint) (*entities.Admin, *rest_err.RestErr)
	GetAll(ctx context.Context, limit, offset int, params map[string]any) ([]entities.Admin, int64, *rest_err.RestErr)
	Delete(ctx context.Context, id uint) *rest_err.RestErr
}
