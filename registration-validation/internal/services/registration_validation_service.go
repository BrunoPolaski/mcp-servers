package services

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/registration-validation/internal/core/entities"
	"github.com/BrunoPolaski/registration-validation/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/registration-validation/internal/infra/repositories"
	"github.com/BrunoPolaski/registration-validation/internal/infra/repositories/interfaces"
)

// CustomerRef identifica o cliente de uma consulta por dimensão. Exatamente um
// dos dois campos deve estar preenchido.
type CustomerRef struct {
	ID       uint
	Document string
}

type RegistrationValidationService struct {
	personRepository                   interfaces.PersonRepository
	documentValidationRepository       interfaces.DocumentValidationRepository
	fiscalRegularityRepository         interfaces.FiscalRegularityRepository
	employmentLinkValidationRepository interfaces.EmploymentLinkValidationRepository
	complianceCheckRepository          interfaces.ComplianceCheckRepository
}

func NewRegistrationValidationService(rf *repositories.RepositoryFactory) *RegistrationValidationService {
	return &RegistrationValidationService{
		personRepository:                   rf.PersonRepository(),
		documentValidationRepository:       rf.DocumentValidationRepository(),
		fiscalRegularityRepository:         rf.FiscalRegularityRepository(),
		employmentLinkValidationRepository: rf.EmploymentLinkValidationRepository(),
		complianceCheckRepository:          rf.ComplianceCheckRepository(),
	}
}

func (s *RegistrationValidationService) resolveCustomer(ctx context.Context, ref CustomerRef) (*entities.Person, *rest_err.RestErr) {
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

func (s *RegistrationValidationService) GetDocumentValidations(ctx context.Context, ref CustomerRef, documentType string) (*dto.DocumentValidationsResultDTO, *rest_err.RestErr) {
	person, err := s.resolveCustomer(ctx, ref)
	if err != nil {
		return nil, err
	}
	validations, err := s.documentValidationRepository.GetByPersonID(ctx, person.ID, documentType)
	if err != nil {
		return nil, err
	}
	items := make([]dto.DocumentValidationDTO, 0, len(validations))
	for i := range validations {
		items = append(items, *dto.NewDocumentValidationDTO(&validations[i]))
	}
	return &dto.DocumentValidationsResultDTO{CustomerID: person.ID, Document: documentOf(person), Items: items}, nil
}

func (s *RegistrationValidationService) GetFiscalRegularities(ctx context.Context, ref CustomerRef) (*dto.FiscalRegularitiesResultDTO, *rest_err.RestErr) {
	person, err := s.resolveCustomer(ctx, ref)
	if err != nil {
		return nil, err
	}
	regularities, err := s.fiscalRegularityRepository.GetByPersonID(ctx, person.ID)
	if err != nil {
		return nil, err
	}
	items := make([]dto.FiscalRegularityDTO, 0, len(regularities))
	for i := range regularities {
		items = append(items, *dto.NewFiscalRegularityDTO(&regularities[i]))
	}
	return &dto.FiscalRegularitiesResultDTO{CustomerID: person.ID, Document: documentOf(person), Items: items}, nil
}

func (s *RegistrationValidationService) GetEmploymentLinks(ctx context.Context, ref CustomerRef, status string) (*dto.EmploymentLinkValidationsResultDTO, *rest_err.RestErr) {
	person, err := s.resolveCustomer(ctx, ref)
	if err != nil {
		return nil, err
	}
	links, err := s.employmentLinkValidationRepository.GetByPersonID(ctx, person.ID, status)
	if err != nil {
		return nil, err
	}
	items := make([]dto.EmploymentLinkValidationDTO, 0, len(links))
	for i := range links {
		items = append(items, *dto.NewEmploymentLinkValidationDTO(&links[i]))
	}
	return &dto.EmploymentLinkValidationsResultDTO{CustomerID: person.ID, Document: documentOf(person), Items: items}, nil
}

func (s *RegistrationValidationService) GetComplianceChecks(ctx context.Context, ref CustomerRef, checkType string) (*dto.ComplianceChecksResultDTO, *rest_err.RestErr) {
	person, err := s.resolveCustomer(ctx, ref)
	if err != nil {
		return nil, err
	}
	checks, err := s.complianceCheckRepository.GetByPersonID(ctx, person.ID, checkType)
	if err != nil {
		return nil, err
	}
	items := make([]dto.ComplianceCheckDTO, 0, len(checks))
	for i := range checks {
		items = append(items, *dto.NewComplianceCheckDTO(&checks[i]))
	}
	return &dto.ComplianceChecksResultDTO{CustomerID: person.ID, Document: documentOf(person), Items: items}, nil
}
