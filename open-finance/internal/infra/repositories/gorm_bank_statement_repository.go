package repositories

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/open-finance/internal/core/entities"
	"github.com/BrunoPolaski/open-finance/internal/infra/repositories/interfaces"
	"gorm.io/gorm"
)

type gormBankStatementRepository struct {
	db *gorm.DB
}

func NewGormBankStatementRepository(db *gorm.DB) interfaces.BankStatementRepository {
	return &gormBankStatementRepository{db: db}
}

func (g *gormBankStatementRepository) GetByPersonID(ctx context.Context, personID uint, accountType string) ([]entities.BankStatement, *rest_err.RestErr) {
	query := gorm.G[entities.BankStatement](g.db).Where("person_id = ?", personID)
	if accountType != "" {
		query = query.Where("account_type = ?", accountType)
	}

	statements, err := query.Order("period_end DESC").Find(ctx)
	if err != nil {
		return nil, rest_err.NewInternalServerError("error while fetching bank statements").WithCause(err)
	}

	return statements, nil
}
