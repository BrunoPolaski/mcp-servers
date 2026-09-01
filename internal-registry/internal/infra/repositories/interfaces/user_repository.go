package interfaces

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
)

type UserRepository interface {
	Register(ctx context.Context, user *entities.User) (*entities.User, *rest_err.RestErr)
	GetById(ctx context.Context, id uint) (*entities.User, *rest_err.RestErr)
	GetByEmail(ctx context.Context, email string) (*entities.User, *rest_err.RestErr)
	GetByDocument(ctx context.Context, document string) (*entities.PersonalInformation, *rest_err.RestErr)
	GetAll(ctx context.Context, limit, offset int, params map[string]interface{}) ([]entities.User, int64, *rest_err.RestErr)
	Delete(ctx context.Context, id uint) *rest_err.RestErr
	Update(ctx context.Context, user *entities.User) (*entities.User, *rest_err.RestErr)
}
