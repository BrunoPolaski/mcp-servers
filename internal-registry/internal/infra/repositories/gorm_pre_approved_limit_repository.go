package repositories

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
	"github.com/BrunoPolaski/internal-registry/internal/infra/repositories/interfaces"
	"gorm.io/gorm"
)

type gormPreApprovedLimitRepository struct {
	db *gorm.DB
}

func NewGormPreApprovedLimitRepository(db *gorm.DB) interfaces.PreApprovedLimitRepository {
	return &gormPreApprovedLimitRepository{db: db}
}

func (g *gormPreApprovedLimitRepository) GetByPersonID(ctx context.Context, personID uint, onlyActive bool) ([]entities.PreApprovedLimit, *rest_err.RestErr) {
	query := gorm.G[entities.PreApprovedLimit](g.db).Where("person_id = ?", personID)
	if onlyActive {
		query = query.Where("is_active = ?", true)
	}

	limits, err := query.Order("calculated_date DESC").Find(ctx)
	if err != nil {
		return nil, rest_err.NewInternalServerError("error while fetching pre approved limits").WithCause(err)
	}
	return limits, nil
}
