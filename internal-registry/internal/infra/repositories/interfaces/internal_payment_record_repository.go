package interfaces

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
)

type InternalPaymentRecordRepository interface {
	// GetByPersonID devolve os registros de pagamento interno da pessoa, do mais
	// recente para o mais antigo. status vazio e productID nil não filtram.
	GetByPersonID(ctx context.Context, personID uint, status string, productID *uint) ([]entities.InternalPaymentRecord, *rest_err.RestErr)
}
