package repositories

import (
	"github.com/BrunoPolaski/internal-registry/internal/infra/repositories/interfaces"
	"github.com/BrunoPolaski/internal-registry/internal/infra/thirdparty"
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

	customerRelationshipRepository  interfaces.CustomerRelationshipRepository
	contractedProductRepository     interfaces.ContractedProductRepository
	internalPaymentRecordRepository interfaces.InternalPaymentRecordRepository
	preApprovedLimitRepository      interfaces.PreApprovedLimitRepository
	incomeDeclarationRepository     interfaces.IncomeDeclarationRepository
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

		customerRelationshipRepository:  NewGormCustomerRelationshipRepository(tpf.DB()),
		contractedProductRepository:     NewGormContractedProductRepository(tpf.DB()),
		internalPaymentRecordRepository: NewGormInternalPaymentRecordRepository(tpf.DB()),
		preApprovedLimitRepository:      NewGormPreApprovedLimitRepository(tpf.DB()),
		incomeDeclarationRepository:     NewGormIncomeDeclarationRepository(tpf.DB()),
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

func (f *RepositoryFactory) CustomerRelationshipRepository() interfaces.CustomerRelationshipRepository {
	return f.customerRelationshipRepository
}

func (f *RepositoryFactory) ContractedProductRepository() interfaces.ContractedProductRepository {
	return f.contractedProductRepository
}

func (f *RepositoryFactory) InternalPaymentRecordRepository() interfaces.InternalPaymentRecordRepository {
	return f.internalPaymentRecordRepository
}

func (f *RepositoryFactory) PreApprovedLimitRepository() interfaces.PreApprovedLimitRepository {
	return f.preApprovedLimitRepository
}

func (f *RepositoryFactory) IncomeDeclarationRepository() interfaces.IncomeDeclarationRepository {
	return f.incomeDeclarationRepository
}
