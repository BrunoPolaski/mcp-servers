package interfaces

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/open-finance/internal/core/entities"
)

type CashFlowAnalysisRepository interface {
	// GetByPersonID devolve as análises da pessoa ordenadas por analysis_date
	// decrescente. limit menor ou igual a zero devolve todas.
	GetByPersonID(ctx context.Context, personID uint, limit int) ([]entities.CashFlowAnalysis, *rest_err.RestErr)
}
