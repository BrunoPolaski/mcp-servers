package repositories

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
	"github.com/BrunoPolaski/internal-registry/internal/infra/repositories/interfaces"
	"gorm.io/gorm"
)

type gormIncomeDeclarationRepository struct {
	db *gorm.DB
}

func NewGormIncomeDeclarationRepository(db *gorm.DB) interfaces.IncomeDeclarationRepository {
	return &gormIncomeDeclarationRepository{db: db}
}

func (g *gormIncomeDeclarationRepository) GetByPersonID(ctx context.Context, personID uint, verifiedOnly bool) ([]entities.IncomeDeclaration, *rest_err.RestErr) {
	query := gorm.G[entities.IncomeDeclaration](g.db).Where("person_id = ?", personID)
	if verifiedOnly {
		query = query.Where("verified = ?", true)
	}

	declarations, err := query.Order("declaration_date DESC").Find(ctx)
	if err != nil {
		return nil, rest_err.NewInternalServerError("error while fetching income declarations").WithCause(err)
	}
	return declarations, nil
}
