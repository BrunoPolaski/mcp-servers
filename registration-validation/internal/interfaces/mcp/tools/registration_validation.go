package tools

import (
	"context"
	"fmt"

	"github.com/BrunoPolaski/registration-validation/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/registration-validation/internal/services"
	"github.com/mark3labs/mcp-go/mcp"
)

const customerRefDescription = `Identify the customer by "customer_id" OR by "document" (CPF). Provide exactly one of them.`

func customerRefFrom(request mcp.CallToolRequest) (services.CustomerRef, error) {
	id := request.GetInt("customer_id", 0)
	if id < 0 {
		return services.CustomerRef{}, fmt.Errorf("invalid customer_id: must be a positive integer")
	}
	document := request.GetString("document", "")
	if id == 0 && document == "" {
		return services.CustomerRef{}, fmt.Errorf("informe customer_id ou document")
	}
	return services.CustomerRef{ID: uint(id), Document: document}, nil
}

func withCustomerRefParams() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithInteger("customer_id", mcp.Description("The ID of the customer")),
		mcp.WithString("document", mcp.Description("The document number (CPF) of the customer")),
	}
}

func (s *Server) GetDocumentValidationsTool() mcp.Tool {
	opts := append([]mcp.ToolOption{
		mcp.WithDescription(`
			Registration Validation: check the customer's CPF/CNPJ against the
			Receita Federal base. Each validation says whether the document is
			valid, its registration status (regular, pending, suspended, canceled,
			deceased) and whether name, birth date and biometrics matched the
			credit bureau and internal registry records.
			An empty list means the document was never validated.
			` + customerRefDescription),
		mcp.WithOutputSchema[dto.DocumentValidationsResultDTO](),
		mcp.WithString("document_type", mcp.Description("Optional filter by document type"),
			mcp.Enum("cpf", "cnpj")),
	}, withCustomerRefParams()...)
	return mcp.NewTool("get_document_validations", opts...)
}

func (s *Server) HandleGetDocumentValidations(ctx context.Context, request mcp.CallToolRequest, args mcp.CallToolParams) (*dto.DocumentValidationsResultDTO, error) {
	ref, err := customerRefFrom(request)
	if err != nil {
		return nil, err
	}

	result, restErr := s.registrationValidationService.GetDocumentValidations(ctx, ref,
		request.GetString("document_type", ""))
	if restErr != nil {
		return nil, restErr
	}
	return result, nil
}

func (s *Server) GetFiscalRegularityTool() mcp.Tool {
	opts := append([]mcp.ToolOption{
		mcp.WithDescription(`
			Registration Validation: fiscal regularity of the customer at the
			Receita Federal (Certidão Negativa de Débitos). Says whether there are
			outstanding tax debts, the CND status (regular, irregular, suspended),
			the certificate number and validity when one was issued, and the list
			of pending issues when it was denied.
			An empty list means fiscal regularity was never checked.
			` + customerRefDescription),
		mcp.WithOutputSchema[dto.FiscalRegularitiesResultDTO](),
	}, withCustomerRefParams()...)
	return mcp.NewTool("get_fiscal_regularity", opts...)
}

func (s *Server) HandleGetFiscalRegularity(ctx context.Context, request mcp.CallToolRequest, args mcp.CallToolParams) (*dto.FiscalRegularitiesResultDTO, error) {
	ref, err := customerRefFrom(request)
	if err != nil {
		return nil, err
	}

	result, restErr := s.registrationValidationService.GetFiscalRegularities(ctx, ref)
	if restErr != nil {
		return nil, restErr
	}
	return result, nil
}

func (s *Server) GetEmploymentLinksTool() mcp.Tool {
	opts := append([]mcp.ToolOption{
		mcp.WithDescription(`
			Registration Validation: employment relationships validated against
			eSocial, with employer name and CNPJ, employment type (CLT,
			estatutario, temporary, PJ, autonomous), status, start and end dates
			and whether the link was confirmed at the source.
			An empty list means eSocial has no employment record for this customer.
			` + customerRefDescription),
		mcp.WithOutputSchema[dto.EmploymentLinkValidationsResultDTO](),
		mcp.WithString("status", mcp.Description("Optional filter by link status"),
			mcp.Enum("active", "terminated")),
	}, withCustomerRefParams()...)
	return mcp.NewTool("get_employment_links", opts...)
}

func (s *Server) HandleGetEmploymentLinks(ctx context.Context, request mcp.CallToolRequest, args mcp.CallToolParams) (*dto.EmploymentLinkValidationsResultDTO, error) {
	ref, err := customerRefFrom(request)
	if err != nil {
		return nil, err
	}

	result, restErr := s.registrationValidationService.GetEmploymentLinks(ctx, ref,
		request.GetString("status", ""))
	if restErr != nil {
		return nil, restErr
	}
	return result, nil
}

func (s *Server) GetComplianceChecksTool() mcp.Tool {
	opts := append([]mcp.ToolOption{
		mcp.WithDescription(`
			Registration Validation: compliance screening of the customer, telling
			whether they are a politically exposed person (PEP), whether they appear
			on restrictive lists, the resulting status and how long the check
			remains valid.
			An empty list means the customer was never screened.
			` + customerRefDescription),
		mcp.WithOutputSchema[dto.ComplianceChecksResultDTO](),
		mcp.WithString("check_type", mcp.Description("Optional filter by check type"),
			mcp.Enum("kyc_full", "pep_screening", "sanctions_screening")),
	}, withCustomerRefParams()...)
	return mcp.NewTool("get_compliance_checks", opts...)
}

func (s *Server) HandleGetComplianceChecks(ctx context.Context, request mcp.CallToolRequest, args mcp.CallToolParams) (*dto.ComplianceChecksResultDTO, error) {
	ref, err := customerRefFrom(request)
	if err != nil {
		return nil, err
	}

	result, restErr := s.registrationValidationService.GetComplianceChecks(ctx, ref,
		request.GetString("check_type", ""))
	if restErr != nil {
		return nil, restErr
	}
	return result, nil
}
