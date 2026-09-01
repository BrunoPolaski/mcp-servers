package services

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/open-finance/internal/core/entities"
	"github.com/BrunoPolaski/open-finance/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/open-finance/internal/infra/repositories"
	"github.com/BrunoPolaski/open-finance/internal/infra/repositories/interfaces"
)

// CustomerRef identifica o cliente de uma consulta por dimensão. Exatamente um
// dos dois campos deve estar preenchido.
type CustomerRef struct {
	ID       uint
	Document string
}

type OpenFinanceService struct {
	personRepository               interfaces.PersonRepository
	bankStatementRepository        interfaces.BankStatementRepository
	cashFlowAnalysisRepository     interfaces.CashFlowAnalysisRepository
	recurringTransactionRepository interfaces.RecurringTransactionRepository
	dataSharingConsentRepository   interfaces.DataSharingConsentRepository
}

func NewOpenFinanceService(rf *repositories.RepositoryFactory) *OpenFinanceService {
	return &OpenFinanceService{
		personRepository:               rf.PersonRepository(),
		bankStatementRepository:        rf.BankStatementRepository(),
		cashFlowAnalysisRepository:     rf.CashFlowAnalysisRepository(),
		recurringTransactionRepository: rf.RecurringTransactionRepository(),
		dataSharingConsentRepository:   rf.DataSharingConsentRepository(),
	}
}

// resolveCustomer traduz a referência recebida pela ferramenta em uma pessoa
// existente, recusando referências ambíguas ou vazias.
func (s *OpenFinanceService) resolveCustomer(ctx context.Context, ref CustomerRef) (*entities.Person, *rest_err.RestErr) {
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

func (s *OpenFinanceService) GetBankStatements(ctx context.Context, ref CustomerRef, accountType string) (*dto.BankStatementsResultDTO, *rest_err.RestErr) {
	person, err := s.resolveCustomer(ctx, ref)
	if err != nil {
		return nil, err
	}

	statements, err := s.bankStatementRepository.GetByPersonID(ctx, person.ID, accountType)
	if err != nil {
		return nil, err
	}

	items := make([]dto.BankStatementDTO, 0, len(statements))
	for i := range statements {
		items = append(items, *dto.NewBankStatementDTO(&statements[i]))
	}

	return &dto.BankStatementsResultDTO{
		CustomerID: person.ID,
		Document:   documentOf(person),
		Items:      items,
	}, nil
}

func (s *OpenFinanceService) GetCashFlowAnalyses(ctx context.Context, ref CustomerRef, limit int) (*dto.CashFlowAnalysesResultDTO, *rest_err.RestErr) {
	person, err := s.resolveCustomer(ctx, ref)
	if err != nil {
		return nil, err
	}

	analyses, err := s.cashFlowAnalysisRepository.GetByPersonID(ctx, person.ID, limit)
	if err != nil {
		return nil, err
	}

	items := make([]dto.CashFlowAnalysisDTO, 0, len(analyses))
	for i := range analyses {
		items = append(items, *dto.NewCashFlowAnalysisDTO(&analyses[i]))
	}

	return &dto.CashFlowAnalysesResultDTO{
		CustomerID: person.ID,
		Document:   documentOf(person),
		Items:      items,
	}, nil
}

func (s *OpenFinanceService) GetRecurringTransactions(ctx context.Context, ref CustomerRef, transactionType string, onlyActive bool) (*dto.RecurringTransactionsResultDTO, *rest_err.RestErr) {
	person, err := s.resolveCustomer(ctx, ref)
	if err != nil {
		return nil, err
	}

	transactions, err := s.recurringTransactionRepository.GetByPersonID(ctx, person.ID, transactionType, onlyActive)
	if err != nil {
		return nil, err
	}

	items := make([]dto.RecurringTransactionDTO, 0, len(transactions))
	for i := range transactions {
		items = append(items, *dto.NewRecurringTransactionDTO(&transactions[i]))
	}

	return &dto.RecurringTransactionsResultDTO{
		CustomerID: person.ID,
		Document:   documentOf(person),
		Items:      items,
	}, nil
}

func (s *OpenFinanceService) GetDataSharingConsents(ctx context.Context, ref CustomerRef) (*dto.DataSharingConsentsResultDTO, *rest_err.RestErr) {
	person, err := s.resolveCustomer(ctx, ref)
	if err != nil {
		return nil, err
	}

	consents, err := s.dataSharingConsentRepository.GetByPersonID(ctx, person.ID)
	if err != nil {
		return nil, err
	}

	items := make([]dto.DataSharingConsentDTO, 0, len(consents))
	for i := range consents {
		items = append(items, *dto.NewDataSharingConsentDTO(&consents[i]))
	}

	return &dto.DataSharingConsentsResultDTO{
		CustomerID: person.ID,
		Document:   documentOf(person),
		Items:      items,
	}, nil
}
