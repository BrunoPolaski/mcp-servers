package repositories

import (
	"testing"

	"github.com/BrunoPolaski/open-finance/internal/infra/repositories/interfaces"
)

// As asserções abaixo falham em compilação se alguma implementação divergir da
// interface, que é o erro mais provável ao adicionar um repositório novo.
var (
	_ interfaces.BankStatementRepository        = (*gormBankStatementRepository)(nil)
	_ interfaces.CashFlowAnalysisRepository     = (*gormCashFlowAnalysisRepository)(nil)
	_ interfaces.RecurringTransactionRepository = (*gormRecurringTransactionRepository)(nil)
	_ interfaces.DataSharingConsentRepository   = (*gormDataSharingConsentRepository)(nil)
)

func TestFactoryExposesOpenFinanceRepositories(t *testing.T) {
	f := &RepositoryFactory{}

	if f.BankStatementRepository() != nil {
		t.Error("esperado nil na factory zerada; o getter deve apenas devolver o campo")
	}
	if f.CashFlowAnalysisRepository() != nil {
		t.Error("esperado nil na factory zerada; o getter deve apenas devolver o campo")
	}
	if f.RecurringTransactionRepository() != nil {
		t.Error("esperado nil na factory zerada; o getter deve apenas devolver o campo")
	}
	if f.DataSharingConsentRepository() != nil {
		t.Error("esperado nil na factory zerada; o getter deve apenas devolver o campo")
	}
}
