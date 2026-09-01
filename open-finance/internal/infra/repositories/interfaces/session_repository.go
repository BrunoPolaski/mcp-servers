package interfaces

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/open-finance/internal/core/entities"
)

type SessionRepository interface {
	Create(ctx context.Context, token *entities.Token) *rest_err.RestErr
	GetById(ctx context.Context, id string) (*entities.Token, *rest_err.RestErr)
	Delete(ctx context.Context, id string) *rest_err.RestErr
}
