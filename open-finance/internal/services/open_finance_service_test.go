package services

import (
	"context"
	"testing"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/open-finance/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/open-finance/internal/core/entities/value_objects"
	"gorm.io/gorm"
)

type fakePersonRepository struct {
	byID       map[uint]*entities.Person
	byDocument map[string]*entities.Person
}

func (f *fakePersonRepository) GetById(_ context.Context, id uint) (*entities.Person, *rest_err.RestErr) {
	if p, ok := f.byID[id]; ok {
		return p, nil
	}
	return nil, rest_err.NewNotFoundError("person not found")
}

func (f *fakePersonRepository) GetByDocument(_ context.Context, document string) (*entities.Person, *rest_err.RestErr) {
	if p, ok := f.byDocument[document]; ok {
		return p, nil
	}
	return nil, rest_err.NewNotFoundError("person not found")
}

func (f *fakePersonRepository) GetAll(_ context.Context, _, _ int, _ map[string]any) ([]entities.Person, int64, *rest_err.RestErr) {
	return nil, 0, nil
}

func (f *fakePersonRepository) Delete(_ context.Context, _ uint) *rest_err.RestErr { return nil }

type fakeBankStatementRepository struct {
	gotPersonID    uint
	gotAccountType string
	statements     []entities.BankStatement
	err            *rest_err.RestErr
}

func (f *fakeBankStatementRepository) GetByPersonID(_ context.Context, personID uint, accountType string) ([]entities.BankStatement, *rest_err.RestErr) {
	f.gotPersonID = personID
	f.gotAccountType = accountType
	return f.statements, f.err
}

type fakeCashFlowAnalysisRepository struct {
	gotLimit int
	analyses []entities.CashFlowAnalysis
}

func (f *fakeCashFlowAnalysisRepository) GetByPersonID(_ context.Context, _ uint, limit int) ([]entities.CashFlowAnalysis, *rest_err.RestErr) {
	f.gotLimit = limit
	return f.analyses, nil
}

type fakeRecurringTransactionRepository struct {
	gotType       string
	gotOnlyActive bool
	transactions  []entities.RecurringTransaction
}

func (f *fakeRecurringTransactionRepository) GetByPersonID(_ context.Context, _ uint, transactionType string, onlyActive bool) ([]entities.RecurringTransaction, *rest_err.RestErr) {
	f.gotType = transactionType
	f.gotOnlyActive = onlyActive
	return f.transactions, nil
}

type fakeDataSharingConsentRepository struct {
	consents []entities.DataSharingConsent
}

func (f *fakeDataSharingConsentRepository) GetByPersonID(_ context.Context, _ uint) ([]entities.DataSharingConsent, *rest_err.RestErr) {
	return f.consents, nil
}

func personFixture() *entities.Person {
	return &entities.Person{
		Model: gorm.Model{ID: 8},
		PersonalInformation: &entities.PersonalInformation{
			FullName: "Lucas Martins Souza",
			Document: valueobjects.Document("18979021232"),
		},
	}
}

func newTestService(bs *fakeBankStatementRepository) (*OpenFinanceService, *fakeCashFlowAnalysisRepository, *fakeRecurringTransactionRepository, *fakeDataSharingConsentRepository) {
	person := personFixture()
	cf := &fakeCashFlowAnalysisRepository{}
	rt := &fakeRecurringTransactionRepository{}
	dc := &fakeDataSharingConsentRepository{}

	svc := &OpenFinanceService{
		personRepository: &fakePersonRepository{
			byID:       map[uint]*entities.Person{8: person},
			byDocument: map[string]*entities.Person{"18979021232": person},
		},
		bankStatementRepository:        bs,
		cashFlowAnalysisRepository:     cf,
		recurringTransactionRepository: rt,
		dataSharingConsentRepository:   dc,
	}

	return svc, cf, rt, dc
}

func TestResolveCustomerRef(t *testing.T) {
	tests := []struct {
		name         string
		ref          CustomerRef
		wantErr      bool
		wantMessage  string
		wantCustomer uint
	}{
		{name: "por id", ref: CustomerRef{ID: 8}, wantCustomer: 8},
		{name: "por documento", ref: CustomerRef{Document: "18979021232"}, wantCustomer: 8},
		{
			name:        "os dois informados",
			ref:         CustomerRef{ID: 8, Document: "18979021232"},
			wantErr:     true,
			wantMessage: "informe apenas um entre customer_id e document",
		},
		{
			name:        "nenhum informado",
			ref:         CustomerRef{},
			wantErr:     true,
			wantMessage: "informe customer_id ou document",
		},
		{
			name:        "id inexistente",
			ref:         CustomerRef{ID: 99},
			wantErr:     true,
			wantMessage: "person not found",
		},
		{
			name:        "documento inexistente",
			ref:         CustomerRef{Document: "00000000000"},
			wantErr:     true,
			wantMessage: "person not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, _, _, _ := newTestService(&fakeBankStatementRepository{})

			person, err := svc.resolveCustomer(context.Background(), tt.ref)

			if tt.wantErr {
				if err == nil {
					t.Fatal("esperado erro, veio nil")
				}
				if err.Message != tt.wantMessage {
					t.Errorf("mensagem = %q, esperado %q", err.Message, tt.wantMessage)
				}
				return
			}

			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if person.ID != tt.wantCustomer {
				t.Errorf("ID = %d, esperado %d", person.ID, tt.wantCustomer)
			}
		})
	}
}

func TestGetBankStatementsRepassaFiltroEIdentificacao(t *testing.T) {
	bs := &fakeBankStatementRepository{
		statements: []entities.BankStatement{{Institution: "Banco Sintético Gama", AccountType: "checking"}},
	}
	svc, _, _, _ := newTestService(bs)

	got, err := svc.GetBankStatements(context.Background(), CustomerRef{Document: "18979021232"}, "checking")

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if bs.gotPersonID != 8 {
		t.Errorf("person_id repassado = %d, esperado 8", bs.gotPersonID)
	}
	if bs.gotAccountType != "checking" {
		t.Errorf("account_type repassado = %q, esperado \"checking\"", bs.gotAccountType)
	}
	if got.CustomerID != 8 || got.Document != "18979021232" {
		t.Errorf("identificação = (%d, %q), esperado (8, \"18979021232\")", got.CustomerID, got.Document)
	}
	if len(got.Items) != 1 || got.Items[0].Institution != "Banco Sintético Gama" {
		t.Error("extratos não foram mapeados para o DTO")
	}
}

func TestGetBankStatementsPropagaErroDoRepositorio(t *testing.T) {
	bs := &fakeBankStatementRepository{err: rest_err.NewInternalServerError("boom")}
	svc, _, _, _ := newTestService(bs)

	_, err := svc.GetBankStatements(context.Background(), CustomerRef{ID: 8}, "")

	if err == nil {
		t.Fatal("esperado erro, veio nil")
	}
	if err.Message != "boom" {
		t.Errorf("mensagem = %q, esperado \"boom\"", err.Message)
	}
}

func TestGetBankStatementsSemRegistrosDevolveListaVazia(t *testing.T) {
	svc, _, _, _ := newTestService(&fakeBankStatementRepository{})

	got, err := svc.GetBankStatements(context.Background(), CustomerRef{ID: 8}, "")

	if err != nil {
		t.Fatalf("ausência de extrato não é erro, mas veio: %v", err)
	}
	if got.Items == nil {
		t.Error("Items deve ser lista vazia, nunca nil, para serializar como [] no JSON")
	}
	if len(got.Items) != 0 {
		t.Errorf("Items = %d, esperado 0", len(got.Items))
	}
}

func TestGetCashFlowAnalysesRepassaLimite(t *testing.T) {
	svc, cf, _, _ := newTestService(&fakeBankStatementRepository{})
	cf.analyses = []entities.CashFlowAnalysis{{NetCashFlow: 150}}

	got, err := svc.GetCashFlowAnalyses(context.Background(), CustomerRef{ID: 8}, 1)

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if cf.gotLimit != 1 {
		t.Errorf("limite repassado = %d, esperado 1", cf.gotLimit)
	}
	if len(got.Items) != 1 || got.Items[0].NetCashFlow != 150 {
		t.Error("análises não foram mapeadas para o DTO")
	}
}

func TestGetRecurringTransactionsRepassaFiltros(t *testing.T) {
	svc, _, rt, _ := newTestService(&fakeBankStatementRepository{})
	rt.transactions = []entities.RecurringTransaction{{TransactionType: "income", Amount: 4200}}

	got, err := svc.GetRecurringTransactions(context.Background(), CustomerRef{ID: 8}, "income", true)

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if rt.gotType != "income" {
		t.Errorf("tipo repassado = %q, esperado \"income\"", rt.gotType)
	}
	if !rt.gotOnlyActive {
		t.Error("only_active repassado = false, esperado true")
	}
	if len(got.Items) != 1 || got.Items[0].Amount != 4200 {
		t.Error("transações não foram mapeadas para o DTO")
	}
}

func TestGetDataSharingConsents(t *testing.T) {
	svc, _, _, dc := newTestService(&fakeBankStatementRepository{})
	dc.consents = []entities.DataSharingConsent{{ConsentID: "urn:openfinance:consent:008", Status: "granted"}}

	got, err := svc.GetDataSharingConsents(context.Background(), CustomerRef{ID: 8})

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(got.Items) != 1 || got.Items[0].Status != "granted" {
		t.Error("consentimentos não foram mapeados para o DTO")
	}
}
