package repositories

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/open-finance/internal/core/entities"
	"github.com/BrunoPolaski/open-finance/internal/infra/repositories/interfaces"
	"gorm.io/gorm"
)

type gormCashFlowAnalysisRepository struct {
	db *gorm.DB
}

func NewGormCashFlowAnalysisRepository(db *gorm.DB) interfaces.CashFlowAnalysisRepository {
	return &gormCashFlowAnalysisRepository{db: db}
}

func (g *gormCashFlowAnalysisRepository) GetByPersonID(ctx context.Context, personID uint, limit int) ([]entities.CashFlowAnalysis, *rest_err.RestErr) {
	query := gorm.G[entities.CashFlowAnalysis](g.db).
		Where("person_id = ?", personID).
		Order("analysis_date DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}

	analyses, err := query.Find(ctx)
	if err != nil {
		return nil, rest_err.NewInternalServerError("error while fetching cash flow analyses").WithCause(err)
	}

	return analyses, nil
}
