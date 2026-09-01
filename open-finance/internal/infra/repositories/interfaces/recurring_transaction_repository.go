package interfaces

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/open-finance/internal/core/entities"
)

type RecurringTransactionRepository interface {
	// GetByPersonID devolve as transações recorrentes da pessoa.
	// transactionType vazio não filtra; onlyActive restringe a is_active = true.
	GetByPersonID(ctx context.Context, personID uint, transactionType string, onlyActive bool) ([]entities.RecurringTransaction, *rest_err.RestErr)
}
