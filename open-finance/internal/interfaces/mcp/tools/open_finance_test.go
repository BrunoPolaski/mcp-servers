package tools

import (
	"context"
	"testing"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/open-finance/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/open-finance/internal/services"
	"github.com/mark3labs/mcp-go/mcp"
)

type fakeOpenFinanceService struct {
	gotRef         services.CustomerRef
	gotAccountType string
	gotLimit       int
	gotType        string
	gotOnlyActive  bool
	err            *rest_err.RestErr
}

func (f *fakeOpenFinanceService) GetBankStatements(_ context.Context, ref services.CustomerRef, accountType string) (*dto.BankStatementsResultDTO, *rest_err.RestErr) {
	f.gotRef, f.gotAccountType = ref, accountType
	if f.err != nil {
		return nil, f.err
	}
	return &dto.BankStatementsResultDTO{CustomerID: 8, Items: []dto.BankStatementDTO{}}, nil
}

func (f *fakeOpenFinanceService) GetCashFlowAnalyses(_ context.Context, ref services.CustomerRef, limit int) (*dto.CashFlowAnalysesResultDTO, *rest_err.RestErr) {
	f.gotRef, f.gotLimit = ref, limit
	if f.err != nil {
		return nil, f.err
	}
	return &dto.CashFlowAnalysesResultDTO{CustomerID: 8, Items: []dto.CashFlowAnalysisDTO{}}, nil
}

func (f *fakeOpenFinanceService) GetRecurringTransactions(_ context.Context, ref services.CustomerRef, transactionType string, onlyActive bool) (*dto.RecurringTransactionsResultDTO, *rest_err.RestErr) {
	f.gotRef, f.gotType, f.gotOnlyActive = ref, transactionType, onlyActive
	if f.err != nil {
		return nil, f.err
	}
	return &dto.RecurringTransactionsResultDTO{CustomerID: 8, Items: []dto.RecurringTransactionDTO{}}, nil
}

func (f *fakeOpenFinanceService) GetDataSharingConsents(_ context.Context, ref services.CustomerRef) (*dto.DataSharingConsentsResultDTO, *rest_err.RestErr) {
	f.gotRef = ref
	if f.err != nil {
		return nil, f.err
	}
	return &dto.DataSharingConsentsResultDTO{CustomerID: 8, Items: []dto.DataSharingConsentDTO{}}, nil
}

func TestCustomerRefFrom(t *testing.T) {
	tests := []struct {
		name    string
		args    map[string]any
		wantErr bool
		want    services.CustomerRef
	}{
		{
			name: "por customer_id",
			args: map[string]any{"customer_id": float64(8)},
			want: services.CustomerRef{ID: 8},
		},
		{
			name: "por document",
			args: map[string]any{"document": "18979021232"},
			want: services.CustomerRef{Document: "18979021232"},
		},
		{
			name:    "nenhum dos dois",
			args:    map[string]any{},
			wantErr: true,
		},
		{
			name:    "customer_id negativo",
			args:    map[string]any{"customer_id": float64(-1)},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := customerRefFrom(requestWith(tt.args))

			if tt.wantErr {
				if err == nil {
					t.Fatal("esperado erro, veio nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if got != tt.want {
				t.Errorf("ref = %+v, esperado %+v", got, tt.want)
			}
		})
	}
}

func TestCustomerRefFromComOsDoisDeixaOServicoDecidir(t *testing.T) {
	// A recusa de referência ambígua é responsabilidade do serviço, que é quem
	// conhece a regra. O helper apenas transporta o que veio na chamada.
	got, err := customerRefFrom(requestWith(map[string]any{
		"customer_id": float64(8),
		"document":    "18979021232",
	}))

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if got.ID != 8 || got.Document != "18979021232" {
		t.Errorf("ref = %+v, esperado ambos preenchidos", got)
	}
}

func TestHandleGetBankStatements(t *testing.T) {
	tests := []struct {
		name            string
		args            map[string]any
		wantAccountType string
	}{
		{
			name:            "sem filtro de conta",
			args:            map[string]any{"customer_id": float64(8)},
			wantAccountType: "",
		},
		{
			name:            "com filtro de conta",
			args:            map[string]any{"customer_id": float64(8), "account_type": "savings"},
			wantAccountType: "savings",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &fakeOpenFinanceService{}
			s := &Server{openFinanceService: svc}

			got, err := s.HandleGetBankStatements(context.Background(), requestWith(tt.args), mcp.CallToolParams{})

			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if svc.gotAccountType != tt.wantAccountType {
				t.Errorf("account_type = %q, esperado %q", svc.gotAccountType, tt.wantAccountType)
			}
			if svc.gotRef.ID != 8 {
				t.Errorf("ref.ID = %d, esperado 8", svc.gotRef.ID)
			}
			if got.CustomerID != 8 {
				t.Errorf("CustomerID = %d, esperado 8", got.CustomerID)
			}
		})
	}
}

func TestHandleGetBankStatementsPropagaErro(t *testing.T) {
	s := &Server{openFinanceService: &fakeOpenFinanceService{err: rest_err.NewNotFoundError("person not found")}}

	if _, err := s.HandleGetBankStatements(context.Background(), requestWith(map[string]any{"customer_id": float64(99)}), mcp.CallToolParams{}); err == nil {
		t.Fatal("esperado erro, veio nil")
	}
}

func TestHandleGetCashFlowAnalysisUsaLimitePadraoUm(t *testing.T) {
	svc := &fakeOpenFinanceService{}
	s := &Server{openFinanceService: svc}

	if _, err := s.HandleGetCashFlowAnalysis(context.Background(), requestWith(map[string]any{"customer_id": float64(8)}), mcp.CallToolParams{}); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if svc.gotLimit != 1 {
		t.Errorf("limite = %d, esperado 1 (a análise mais recente)", svc.gotLimit)
	}
}

func TestHandleGetCashFlowAnalysisRecusaLimiteNegativo(t *testing.T) {
	s := &Server{openFinanceService: &fakeOpenFinanceService{}}

	if _, err := s.HandleGetCashFlowAnalysis(context.Background(), requestWith(map[string]any{
		"customer_id": float64(8),
		"limit":       float64(-1),
	}), mcp.CallToolParams{}); err == nil {
		t.Fatal("esperado erro, veio nil")
	}
}

func TestHandleGetRecurringTransactions(t *testing.T) {
	tests := []struct {
		name          string
		args          map[string]any
		wantType      string
		wantOnlyActiv bool
	}{
		{
			name:          "padrões",
			args:          map[string]any{"customer_id": float64(8)},
			wantType:      "",
			wantOnlyActiv: true,
		},
		{
			name:          "somente receitas, incluindo inativas",
			args:          map[string]any{"customer_id": float64(8), "transaction_type": "income", "only_active": false},
			wantType:      "income",
			wantOnlyActiv: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &fakeOpenFinanceService{}
			s := &Server{openFinanceService: svc}

			if _, err := s.HandleGetRecurringTransactions(context.Background(), requestWith(tt.args), mcp.CallToolParams{}); err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}

			if svc.gotType != tt.wantType {
				t.Errorf("transaction_type = %q, esperado %q", svc.gotType, tt.wantType)
			}
			if svc.gotOnlyActive != tt.wantOnlyActiv {
				t.Errorf("only_active = %v, esperado %v", svc.gotOnlyActive, tt.wantOnlyActiv)
			}
		})
	}
}

func TestHandleGetDataSharingConsents(t *testing.T) {
	svc := &fakeOpenFinanceService{}
	s := &Server{openFinanceService: svc}

	got, err := s.HandleGetDataSharingConsents(context.Background(), requestWith(map[string]any{"document": "18979021232"}), mcp.CallToolParams{})

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if svc.gotRef.Document != "18979021232" {
		t.Errorf("ref.Document = %q, esperado \"18979021232\"", svc.gotRef.Document)
	}
	if got.CustomerID != 8 {
		t.Errorf("CustomerID = %d, esperado 8", got.CustomerID)
	}
}

func TestNomesDasFerramentasPorDimensao(t *testing.T) {
	s := &Server{openFinanceService: &fakeOpenFinanceService{}}

	tests := []struct {
		got  string
		want string
	}{
		{s.GetBankStatementsTool().Name, "get_bank_statements"},
		{s.GetCashFlowAnalysisTool().Name, "get_cash_flow_analysis"},
		{s.GetRecurringTransactionsTool().Name, "get_recurring_transactions"},
		{s.GetDataSharingConsentsTool().Name, "get_data_sharing_consents"},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("nome = %q, esperado %q", tt.got, tt.want)
		}
	}
}
