package tools

import (
	"context"
	"testing"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/internal-registry/internal/services"
	"github.com/mark3labs/mcp-go/mcp"
)

type fakeInternalRegistryService struct {
	gotRef        services.CustomerRef
	gotType       string
	gotStatus     string
	gotOnlyActive bool
	err           *rest_err.RestErr
}

func (f *fakeInternalRegistryService) GetCustomerRelationship(_ context.Context, ref services.CustomerRef) (*dto.CustomerRelationshipResultDTO, *rest_err.RestErr) {
	f.gotRef = ref
	if f.err != nil {
		return nil, f.err
	}
	return &dto.CustomerRelationshipResultDTO{CustomerID: ref.ID, Document: ref.Document}, nil
}
func (f *fakeInternalRegistryService) GetContractedProducts(_ context.Context, ref services.CustomerRef, productType, status string) (*dto.ContractedProductsResultDTO, *rest_err.RestErr) {
	f.gotRef, f.gotType, f.gotStatus = ref, productType, status
	if f.err != nil {
		return nil, f.err
	}
	return &dto.ContractedProductsResultDTO{CustomerID: ref.ID, Document: ref.Document}, nil
}
func (f *fakeInternalRegistryService) GetInternalPaymentRecords(_ context.Context, ref services.CustomerRef, status string, productID *uint) (*dto.InternalPaymentRecordsResultDTO, *rest_err.RestErr) {
	f.gotRef, f.gotStatus = ref, status
	if f.err != nil {
		return nil, f.err
	}
	return &dto.InternalPaymentRecordsResultDTO{CustomerID: ref.ID, Document: ref.Document}, nil
}
func (f *fakeInternalRegistryService) GetPreApprovedLimits(_ context.Context, ref services.CustomerRef, onlyActive bool) (*dto.PreApprovedLimitsResultDTO, *rest_err.RestErr) {
	f.gotRef, f.gotOnlyActive = ref, onlyActive
	if f.err != nil {
		return nil, f.err
	}
	return &dto.PreApprovedLimitsResultDTO{CustomerID: ref.ID, Document: ref.Document}, nil
}
func (f *fakeInternalRegistryService) GetIncomeDeclarations(_ context.Context, ref services.CustomerRef, verifiedOnly bool) (*dto.IncomeDeclarationsResultDTO, *rest_err.RestErr) {
	f.gotRef = ref
	if f.err != nil {
		return nil, f.err
	}
	return &dto.IncomeDeclarationsResultDTO{CustomerID: ref.ID, Document: ref.Document}, nil
}

func reqWith(args map[string]any) mcp.CallToolRequest {
	var r mcp.CallToolRequest
	r.Params.Arguments = args
	return r
}

func TestHandleGetContractedProducts_PassesRefAndFilters(t *testing.T) {
	svc := &fakeInternalRegistryService{}
	s := &Server{internalRegistryService: svc}
	_, err := s.HandleGetContractedProducts(context.Background(),
		reqWith(map[string]any{"customer_id": float64(5), "product_type": "credit_card", "status": "active"}), mcp.CallToolParams{})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if svc.gotRef.ID != 5 || svc.gotType != "credit_card" || svc.gotStatus != "active" {
		t.Fatalf("ref/filtros não repassados: %+v", svc)
	}
}

func TestHandleGetContractedProducts_NoIdentifier(t *testing.T) {
	s := &Server{internalRegistryService: &fakeInternalRegistryService{}}
	_, err := s.HandleGetContractedProducts(context.Background(), reqWith(map[string]any{}), mcp.CallToolParams{})
	if err == nil {
		t.Fatal("esperado erro quando nem customer_id nem document são informados")
	}
}

func TestHandleGetPreApprovedLimits_DefaultsOnlyActiveTrue(t *testing.T) {
	svc := &fakeInternalRegistryService{}
	s := &Server{internalRegistryService: svc}
	_, err := s.HandleGetPreApprovedLimits(context.Background(),
		reqWith(map[string]any{"document": "11122233344"}), mcp.CallToolParams{})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if !svc.gotOnlyActive || svc.gotRef.Document != "11122233344" {
		t.Fatalf("only_active default ou ref errados: %+v", svc)
	}
}
