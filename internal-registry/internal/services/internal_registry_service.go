package services

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
	"github.com/BrunoPolaski/internal-registry/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/internal-registry/internal/infra/repositories"
	"github.com/BrunoPolaski/internal-registry/internal/infra/repositories/interfaces"
)

// CustomerRef identifica o cliente de uma consulta por dimensão. Exatamente um
// dos dois campos deve estar preenchido.
type CustomerRef struct {
	ID       uint
	Document string
}

type InternalRegistryService struct {
	personRepository                interfaces.PersonRepository
	customerRelationshipRepository  interfaces.CustomerRelationshipRepository
	contractedProductRepository     interfaces.ContractedProductRepository
	internalPaymentRecordRepository interfaces.InternalPaymentRecordRepository
	preApprovedLimitRepository      interfaces.PreApprovedLimitRepository
	incomeDeclarationRepository     interfaces.IncomeDeclarationRepository
}

func NewInternalRegistryService(rf *repositories.RepositoryFactory) *InternalRegistryService {
	return &InternalRegistryService{
		personRepository:                rf.PersonRepository(),
		customerRelationshipRepository:  rf.CustomerRelationshipRepository(),
		contractedProductRepository:     rf.ContractedProductRepository(),
		internalPaymentRecordRepository: rf.InternalPaymentRecordRepository(),
		preApprovedLimitRepository:      rf.PreApprovedLimitRepository(),
		incomeDeclarationRepository:     rf.IncomeDeclarationRepository(),
	}
}

func (s *InternalRegistryService) resolveCustomer(ctx context.Context, ref CustomerRef) (*entities.Person, *rest_err.RestErr) {
	hasID := ref.ID > 0
	hasDocument := ref.Document != ""
	switch {
	case hasID && hasDocument:
		return nil, rest_err.NewBadRequestError("informe apenas um entre customer_id e document")
	case !hasID && !hasDocument:
		return nil, rest_err.NewBadRequestError("informe customer_id ou document")
	case hasID:
		return s.personRepository.GetById(ctx, ref.ID)
	default:
		return s.personRepository.GetByDocument(ctx, ref.Document)
	}
}

func documentOf(person *entities.Person) string {
	if person.PersonalInformation == nil {
		return ""
	}
	return person.PersonalInformation.Document.String()
}

func (s *InternalRegistryService) GetCustomerRelationship(ctx context.Context, ref CustomerRef) (*dto.CustomerRelationshipResultDTO, *rest_err.RestErr) {
	person, err := s.resolveCustomer(ctx, ref)
	if err != nil {
		return nil, err
	}
	rel, err := s.customerRelationshipRepository.GetByPersonID(ctx, person.ID)
	if err != nil {
		return nil, err
	}
	result := &dto.CustomerRelationshipResultDTO{CustomerID: person.ID, Document: documentOf(person)}
	if rel != nil {
		result.Relationship = dto.NewCustomerRelationshipDTO(rel)
	}
	return result, nil
}

func (s *InternalRegistryService) GetContractedProducts(ctx context.Context, ref CustomerRef, productType, status string) (*dto.ContractedProductsResultDTO, *rest_err.RestErr) {
	person, err := s.resolveCustomer(ctx, ref)
	if err != nil {
		return nil, err
	}
	products, err := s.contractedProductRepository.GetByPersonID(ctx, person.ID, productType, status)
	if err != nil {
		return nil, err
	}
	items := make([]dto.ContractedProductDTO, 0, len(products))
	for i := range products {
		items = append(items, *dto.NewContractedProductDTO(&products[i]))
	}
	return &dto.ContractedProductsResultDTO{CustomerID: person.ID, Document: documentOf(person), Items: items}, nil
}

func (s *InternalRegistryService) GetInternalPaymentRecords(ctx context.Context, ref CustomerRef, status string, productID *uint) (*dto.InternalPaymentRecordsResultDTO, *rest_err.RestErr) {
	person, err := s.resolveCustomer(ctx, ref)
	if err != nil {
		return nil, err
	}
	records, err := s.internalPaymentRecordRepository.GetByPersonID(ctx, person.ID, status, productID)
	if err != nil {
		return nil, err
	}
	items := make([]dto.InternalPaymentRecordDTO, 0, len(records))
	for i := range records {
		items = append(items, *dto.NewInternalPaymentRecordDTO(&records[i]))
	}
	return &dto.InternalPaymentRecordsResultDTO{CustomerID: person.ID, Document: documentOf(person), Items: items}, nil
}

func (s *InternalRegistryService) GetPreApprovedLimits(ctx context.Context, ref CustomerRef, onlyActive bool) (*dto.PreApprovedLimitsResultDTO, *rest_err.RestErr) {
	person, err := s.resolveCustomer(ctx, ref)
	if err != nil {
		return nil, err
	}
	limits, err := s.preApprovedLimitRepository.GetByPersonID(ctx, person.ID, onlyActive)
	if err != nil {
		return nil, err
	}
	items := make([]dto.PreApprovedLimitDTO, 0, len(limits))
	for i := range limits {
		items = append(items, *dto.NewPreApprovedLimitDTO(&limits[i]))
	}
	return &dto.PreApprovedLimitsResultDTO{CustomerID: person.ID, Document: documentOf(person), Items: items}, nil
}

func (s *InternalRegistryService) GetIncomeDeclarations(ctx context.Context, ref CustomerRef, verifiedOnly bool) (*dto.IncomeDeclarationsResultDTO, *rest_err.RestErr) {
	person, err := s.resolveCustomer(ctx, ref)
	if err != nil {
		return nil, err
	}
	declarations, err := s.incomeDeclarationRepository.GetByPersonID(ctx, person.ID, verifiedOnly)
	if err != nil {
		return nil, err
	}
	items := make([]dto.IncomeDeclarationDTO, 0, len(declarations))
	for i := range declarations {
		items = append(items, *dto.NewIncomeDeclarationDTO(&declarations[i]))
	}
	return &dto.IncomeDeclarationsResultDTO{CustomerID: person.ID, Document: documentOf(person), Items: items}, nil
}
