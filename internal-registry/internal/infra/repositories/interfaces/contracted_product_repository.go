package interfaces

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
)

type ContractedProductRepository interface {
	// GetByPersonID devolve os produtos contratados da pessoa. productType e
	// status vazios não filtram.
	GetByPersonID(ctx context.Context, personID uint, productType, status string) ([]entities.ContractedProduct, *rest_err.RestErr)
}
