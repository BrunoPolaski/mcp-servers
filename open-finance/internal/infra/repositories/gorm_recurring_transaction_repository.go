package repositories

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/open-finance/internal/core/entities"
	"github.com/BrunoPolaski/open-finance/internal/infra/repositories/interfaces"
	"gorm.io/gorm"
)

type gormRecurringTransactionRepository struct {
	db *gorm.DB
}

func NewGormRecurringTransactionRepository(db *gorm.DB) interfaces.RecurringTransactionRepository {
	return &gormRecurringTransactionRepository{db: db}
}

func (g *gormRecurringTransactionRepository) GetByPersonID(ctx context.Context, personID uint, transactionType string, onlyActive bool) ([]entities.RecurringTransaction, *rest_err.RestErr) {
	query := gorm.G[entities.RecurringTransaction](g.db).Where("person_id = ?", personID)
	if transactionType != "" {
		query = query.Where("transaction_type = ?", transactionType)
	}
	if onlyActive {
		query = query.Where("is_active = ?", true)
	}

	transactions, err := query.Order("last_occurrence_date DESC").Find(ctx)
	if err != nil {
		return nil, rest_err.NewInternalServerError("error while fetching recurring transactions").WithCause(err)
	}

	return transactions, nil
}
