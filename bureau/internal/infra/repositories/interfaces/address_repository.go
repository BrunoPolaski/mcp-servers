package interfaces

import (
	"context"

	"github.com/BrunoPolaski/bureau/internal/core/entities"
	"github.com/BrunoPolaski/go-rest-err/rest_err"
)

type AddressRepository interface {
	Create(ctx context.Context, address *entities.Address) *rest_err.RestErr
	GetById(ctx context.Context, id uint) (*entities.Address, *rest_err.RestErr)
	GetAll(ctx context.Context, limit, offset int, params map[string]any) ([]entities.Address, int64, *rest_err.RestErr)
	Delete(ctx context.Context, id uint) *rest_err.RestErr
}
