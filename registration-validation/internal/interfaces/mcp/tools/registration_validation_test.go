package tools

import (
	"context"
	"testing"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/registration-validation/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/registration-validation/internal/services"
	"github.com/mark3labs/mcp-go/mcp"
)

type fakeValidationService struct {
	gotRef    services.CustomerRef
	gotFilter string
	err       *rest_err.RestErr
}

func (f *fakeValidationService) GetDocumentValidations(_ context.Context, ref services.CustomerRef, documentType string) (*dto.DocumentValidationsResultDTO, *rest_err.RestErr) {
	f.gotRef, f.gotFilter = ref, documentType
	if f.err != nil {
		return nil, f.err
	}
	return &dto.DocumentValidationsResultDTO{CustomerID: ref.ID, Document: ref.Document}, nil
}

func (f *fakeValidationService) GetFiscalRegularities(_ context.Context, ref services.CustomerRef) (*dto.FiscalRegularitiesResultDTO, *rest_err.RestErr) {
	f.gotRef = ref
	if f.err != nil {
		return nil, f.err
	}
	return &dto.FiscalRegularitiesResultDTO{CustomerID: ref.ID, Document: ref.Document}, nil
}

func (f *fakeValidationService) GetEmploymentLinks(_ context.Context, ref services.CustomerRef, status string) (*dto.EmploymentLinkValidationsResultDTO, *rest_err.RestErr) {
	f.gotRef, f.gotFilter = ref, status
	if f.err != nil {
		return nil, f.err
	}
	return &dto.EmploymentLinkValidationsResultDTO{CustomerID: ref.ID, Document: ref.Document}, nil
}

func (f *fakeValidationService) GetComplianceChecks(_ context.Context, ref services.CustomerRef, checkType string) (*dto.ComplianceChecksResultDTO, *rest_err.RestErr) {
	f.gotRef, f.gotFilter = ref, checkType
	if f.err != nil {
		return nil, f.err
	}
	return &dto.ComplianceChecksResultDTO{CustomerID: ref.ID, Document: ref.Document}, nil
}

func TestCustomerRefFrom(t *testing.T) {
	tests := []struct {
		name    string
		args    map[string]any
		want    services.CustomerRef
		wantErr bool
	}{
		{name: "by id", args: map[string]any{"customer_id": float64(5)}, want: services.CustomerRef{ID: 5}},
		{name: "by document", args: map[string]any{"document": "35509139404"}, want: services.CustomerRef{Document: "35509139404"}},
		{name: "both are passed through for the service to reject", args: map[string]any{"customer_id": float64(5), "document": "x"}, want: services.CustomerRef{ID: 5, Document: "x"}},
		{name: "neither", args: map[string]any{}, wantErr: true},
		{name: "negative id", args: map[string]any{"customer_id": float64(-1)}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := customerRefFrom(requestWith(tt.args))

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestValidationHandlersPassRefAndFilter(t *testing.T) {
	type handler func(*Server, context.Context, mcp.CallToolRequest) error
	tests := []struct {
		name      string
		filterArg string
		filterVal string
		call      handler
	}{
		{
			name: "document validations", filterArg: "document_type", filterVal: "cpf",
			call: func(s *Server, ctx context.Context, r mcp.CallToolRequest) error {
				_, err := s.HandleGetDocumentValidations(ctx, r, mcp.CallToolParams{})
				return err
			},
		},
		{
			name: "fiscal regularity",
			call: func(s *Server, ctx context.Context, r mcp.CallToolRequest) error {
				_, err := s.HandleGetFiscalRegularity(ctx, r, mcp.CallToolParams{})
				return err
			},
		},
		{
			name: "employment links", filterArg: "status", filterVal: "active",
			call: func(s *Server, ctx context.Context, r mcp.CallToolRequest) error {
				_, err := s.HandleGetEmploymentLinks(ctx, r, mcp.CallToolParams{})
				return err
			},
		},
		{
			name: "compliance checks", filterArg: "check_type", filterVal: "pep_screening",
			call: func(s *Server, ctx context.Context, r mcp.CallToolRequest) error {
				_, err := s.HandleGetComplianceChecks(ctx, r, mcp.CallToolParams{})
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &fakeValidationService{}
			s := &Server{registrationValidationService: svc}
			args := map[string]any{"document": "35509139404"}
			if tt.filterArg != "" {
				args[tt.filterArg] = tt.filterVal
			}

			if err := tt.call(s, context.Background(), requestWith(args)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if svc.gotRef.Document != "35509139404" || svc.gotFilter != tt.filterVal {
				t.Errorf("ref or filter not passed: ref=%+v filter=%q", svc.gotRef, svc.gotFilter)
			}
		})
	}
}

func TestValidationHandlerWithoutCustomerSkipsService(t *testing.T) {
	svc := &fakeValidationService{}
	s := &Server{registrationValidationService: svc}

	_, err := s.HandleGetComplianceChecks(context.Background(), requestWith(map[string]any{}), mcp.CallToolParams{})

	if err == nil {
		t.Fatal("expected an error when neither customer_id nor document is given")
	}
	if svc.gotRef != (services.CustomerRef{}) {
		t.Errorf("service should not be called, got ref %+v", svc.gotRef)
	}
}

func TestValidationHandlerPropagatesServiceError(t *testing.T) {
	s := &Server{registrationValidationService: &fakeValidationService{err: rest_err.NewNotFoundError("person not found")}}

	if _, err := s.HandleGetFiscalRegularity(context.Background(), requestWith(map[string]any{"customer_id": float64(99)}), mcp.CallToolParams{}); err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestToolNames(t *testing.T) {
	s := &Server{}

	tests := []struct {
		got  string
		want string
	}{
		{s.GetPersonByIDTool().Name, "get_customer_by_id"},
		{s.GetPersonByDocumentTool().Name, "get_customer_by_document"},
		{s.GetAllPersonsTool().Name, "get_all_customers"},
		{s.GetDocumentValidationsTool().Name, "get_document_validations"},
		{s.GetFiscalRegularityTool().Name, "get_fiscal_regularity"},
		{s.GetEmploymentLinksTool().Name, "get_employment_links"},
		{s.GetComplianceChecksTool().Name, "get_compliance_checks"},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("name = %q, want %q", tt.got, tt.want)
		}
	}
}
