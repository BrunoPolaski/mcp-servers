package repositories

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
	"github.com/BrunoPolaski/internal-registry/internal/infra/repositories/interfaces"
	"gorm.io/gorm"
)

type gormContractedProductRepository struct {
	db *gorm.DB
}

func NewGormContractedProductRepository(db *gorm.DB) interfaces.ContractedProductRepository {
	return &gormContractedProductRepository{db: db}
}

func (g *gormContractedProductRepository) GetByPersonID(ctx context.Context, personID uint, productType, status string) ([]entities.ContractedProduct, *rest_err.RestErr) {
	query := gorm.G[entities.ContractedProduct](g.db).Where("person_id = ?", personID)
	if productType != "" {
		query = query.Where("product_type = ?", productType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	products, err := query.Order("contracted_date DESC").Find(ctx)
	if err != nil {
		return nil, rest_err.NewInternalServerError("error while fetching contracted products").WithCause(err)
	}
	return products, nil
}
