package interfaces

import (
	"context"

	"github.com/BrunoPolaski/bureau/internal/core/entities"
	"github.com/BrunoPolaski/go-rest-err/rest_err"
)

type SessionRepository interface {
	Create(ctx context.Context, token *entities.Token) *rest_err.RestErr
	GetById(ctx context.Context, id string) (*entities.Token, *rest_err.RestErr)
	Delete(ctx context.Context, id string) *rest_err.RestErr
}
