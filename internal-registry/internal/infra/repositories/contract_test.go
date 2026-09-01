package repositories

import (
	"testing"

	"github.com/BrunoPolaski/internal-registry/internal/infra/repositories/interfaces"
)

// Asserções de compilação: cada construtor devolve a interface esperada.
var (
	_ interfaces.CustomerRelationshipRepository  = (*gormCustomerRelationshipRepository)(nil)
	_ interfaces.ContractedProductRepository     = (*gormContractedProductRepository)(nil)
	_ interfaces.InternalPaymentRecordRepository = (*gormInternalPaymentRecordRepository)(nil)
	_ interfaces.PreApprovedLimitRepository      = (*gormPreApprovedLimitRepository)(nil)
	_ interfaces.IncomeDeclarationRepository     = (*gormIncomeDeclarationRepository)(nil)
)

func TestRepositoryFactoryWiresDimensions(t *testing.T) {
	// Guarda contra getters esquecidos: o factory precisa expor os cinco.
	var f RepositoryFactory
	_ = f.CustomerRelationshipRepository
	_ = f.ContractedProductRepository
	_ = f.InternalPaymentRecordRepository
	_ = f.PreApprovedLimitRepository
	_ = f.IncomeDeclarationRepository
}
