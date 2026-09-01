package interfaces

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
)

type PreApprovedLimitRepository interface {
	// GetByPersonID devolve os limites pré-aprovados da pessoa. onlyActive
	// restringe a is_active = true.
	GetByPersonID(ctx context.Context, personID uint, onlyActive bool) ([]entities.PreApprovedLimit, *rest_err.RestErr)
}
