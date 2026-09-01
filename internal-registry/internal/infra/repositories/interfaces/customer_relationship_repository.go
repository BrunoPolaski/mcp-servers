package interfaces

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
)

type CustomerRelationshipRepository interface {
	// GetByPersonID devolve o relacionamento da pessoa, ou (nil, nil) quando ela
	// não é cliente da instituição.
	GetByPersonID(ctx context.Context, personID uint) (*entities.CustomerRelationship, *rest_err.RestErr)
}
