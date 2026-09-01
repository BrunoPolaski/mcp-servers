package interfaces

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
)

type IncomeDeclarationRepository interface {
	// GetByPersonID devolve as declarações de renda da pessoa, da mais recente
	// para a mais antiga. verifiedOnly restringe a verified = true.
	GetByPersonID(ctx context.Context, personID uint, verifiedOnly bool) ([]entities.IncomeDeclaration, *rest_err.RestErr)
}
