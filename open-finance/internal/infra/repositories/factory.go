package repositories

import (
	"github.com/BrunoPolaski/open-finance/internal/infra/repositories/interfaces"
	"github.com/BrunoPolaski/open-finance/internal/infra/thirdparty"
)

type RepositoryFactory struct {
	tpf *thirdparty.ThirdPartyFactory

	addressRepository interfaces.AddressRepository
	adminRepository   interfaces.AdminRepository
	apiKeyRepository  interfaces.ApiKeyRepository
	personRepository  interfaces.PersonRepository
	analystRepository interfaces.AnalystRepository
	userRepository    interfaces.UserRepository
	sessionRepository interfaces.SessionRepository

	bankStatementRepository        interfaces.BankStatementRepository
	cashFlowAnalysisRepository     interfaces.CashFlowAnalysisRepository
	recurringTransactionRepository interfaces.RecurringTransactionRepository
	dataSharingConsentRepository   interfaces.DataSharingConsentRepository
}

func NewRepositoryFactory(tpf *thirdparty.ThirdPartyFactory) *RepositoryFactory {
	return &RepositoryFactory{
		tpf:               tpf,
		addressRepository: NewGormAddressRepository(tpf.DB()),
		adminRepository:   NewGormAdminRepository(tpf.DB()),
		apiKeyRepository:  NewGormApiKeyRepository(tpf.DB()),
		personRepository:  NewGormPersonRepository(tpf.DB()),
		analystRepository: NewGormAnalystRepository(tpf.DB()),
		userRepository:    NewGormUserRepository(tpf.DB()),
		sessionRepository: NewGormSessionRepository(tpf.DB()),

		bankStatementRepository:        NewGormBankStatementRepository(tpf.DB()),
		cashFlowAnalysisRepository:     NewGormCashFlowAnalysisRepository(tpf.DB()),
		recurringTransactionRepository: NewGormRecurringTransactionRepository(tpf.DB()),
		dataSharingConsentRepository:   NewGormDataSharingConsentRepository(tpf.DB()),
	}
}

func (f *RepositoryFactory) AddressRepository() interfaces.AddressRepository {
	return f.addressRepository
}

func (f *RepositoryFactory) AdminRepository() interfaces.AdminRepository {
	return f.adminRepository
}

func (f *RepositoryFactory) ApiKeyRepository() interfaces.ApiKeyRepository {
	return f.apiKeyRepository
}

func (f *RepositoryFactory) PersonRepository() interfaces.PersonRepository {
	return f.personRepository
}

func (f *RepositoryFactory) AnalystRepository() interfaces.AnalystRepository {
	return f.analystRepository
}

func (f *RepositoryFactory) UserRepository() interfaces.UserRepository {
	return f.userRepository
}

func (f *RepositoryFactory) SessionRepository() interfaces.SessionRepository {
	return f.sessionRepository
}

func (f *RepositoryFactory) BankStatementRepository() interfaces.BankStatementRepository {
	return f.bankStatementRepository
}

func (f *RepositoryFactory) CashFlowAnalysisRepository() interfaces.CashFlowAnalysisRepository {
	return f.cashFlowAnalysisRepository
}

func (f *RepositoryFactory) RecurringTransactionRepository() interfaces.RecurringTransactionRepository {
	return f.recurringTransactionRepository
}

func (f *RepositoryFactory) DataSharingConsentRepository() interfaces.DataSharingConsentRepository {
	return f.dataSharingConsentRepository
}
