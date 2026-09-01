package interfaces

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/open-finance/internal/core/entities"
)

type BankStatementRepository interface {
	// GetByPersonID devolve os extratos da pessoa, do mais recente para o mais
	// antigo. accountType vazio não filtra por tipo de conta.
	GetByPersonID(ctx context.Context, personID uint, accountType string) ([]entities.BankStatement, *rest_err.RestErr)
}
