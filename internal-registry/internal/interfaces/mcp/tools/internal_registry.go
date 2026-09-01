package tools

import (
	"context"
	"fmt"

	"github.com/BrunoPolaski/internal-registry/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/internal-registry/internal/services"
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

func (s *Server) GetCustomerRelationshipTool() mcp.Tool {
	opts := append([]mcp.ToolOption{
		mcp.WithDescription(`
			Internal Registry: get the customer's relationship with the institution:
			how long they have been a customer (customer_since, relationship_months),
			segment, branch, whether the relationship is active, the internal behavioral
			score and the churn risk.
			A customer with no relationship record is not a client of this institution;
			treat the relationship criteria as missing data rather than as zero.
			` + customerRefDescription),
		mcp.WithOutputSchema[dto.CustomerRelationshipResultDTO](),
	}, withCustomerRefParams()...)
	return mcp.NewTool("get_customer_relationship", opts...)
}

func (s *Server) HandleGetCustomerRelationship(ctx context.Context, request mcp.CallToolRequest, args mcp.CallToolParams) (*dto.CustomerRelationshipResultDTO, error) {
	ref, err := customerRefFrom(request)
	if err != nil {
		return nil, err
	}
	result, restErr := s.internalRegistryService.GetCustomerRelationship(ctx, ref)
	if restErr != nil {
		return nil, restErr
	}
	return result, nil
}

func (s *Server) GetContractedProductsTool() mcp.Tool {
	opts := append([]mcp.ToolOption{
		mcp.WithDescription(`
			Internal Registry: get the products the customer has contracted with the
			institution (checking account, credit card, loan, insurance, investment),
			each with its status, balance and monthly value.
			` + customerRefDescription),
		mcp.WithOutputSchema[dto.ContractedProductsResultDTO](),
		mcp.WithString("product_type", mcp.Description("Optional filter by product type"),
			mcp.Enum("checking_account", "credit_card", "loan", "insurance", "investment")),
		mcp.WithString("status", mcp.Description("Optional filter by status"),
			mcp.Enum("active", "closed", "suspended")),
	}, withCustomerRefParams()...)
	return mcp.NewTool("get_contracted_products", opts...)
}

func (s *Server) HandleGetContractedProducts(ctx context.Context, request mcp.CallToolRequest, args mcp.CallToolParams) (*dto.ContractedProductsResultDTO, error) {
	ref, err := customerRefFrom(request)
	if err != nil {
		return nil, err
	}
	result, restErr := s.internalRegistryService.GetContractedProducts(ctx, ref,
		request.GetString("product_type", ""), request.GetString("status", ""))
	if restErr != nil {
		return nil, restErr
	}
	return result, nil
}

func (s *Server) GetInternalPaymentRecordsTool() mcp.Tool {
	opts := append([]mcp.ToolOption{
		mcp.WithDescription(`
			Internal Registry: get the customer's internal payment history for products
			contracted with the institution. Each record reports the reference month,
			due date, payment date, amount due and paid, status (on_time, late, missed,
			partial) and days late.
			` + customerRefDescription),
		mcp.WithOutputSchema[dto.InternalPaymentRecordsResultDTO](),
		mcp.WithString("status", mcp.Description("Optional filter by status"),
			mcp.Enum("on_time", "late", "missed", "partial")),
		mcp.WithInteger("product_id", mcp.Description("Optional filter by contracted product id"), mcp.Min(1)),
	}, withCustomerRefParams()...)
	return mcp.NewTool("get_internal_payment_records", opts...)
}

func (s *Server) HandleGetInternalPaymentRecords(ctx context.Context, request mcp.CallToolRequest, args mcp.CallToolParams) (*dto.InternalPaymentRecordsResultDTO, error) {
	ref, err := customerRefFrom(request)
	if err != nil {
		return nil, err
	}
	var productID *uint
	if pid := request.GetInt("product_id", 0); pid > 0 {
		v := uint(pid)
		productID = &v
	}
	result, restErr := s.internalRegistryService.GetInternalPaymentRecords(ctx, ref,
		request.GetString("status", ""), productID)
	if restErr != nil {
		return nil, restErr
	}
	return result, nil
}

func (s *Server) GetPreApprovedLimitsTool() mcp.Tool {
	opts := append([]mcp.ToolOption{
		mcp.WithDescription(`
			Internal Registry: get the customer's pre-approved credit limits granted by
			the institution's internal policies. Each limit reports the product type,
			approved amount, interest rate, calculation date, validity and whether it is
			active.
			` + customerRefDescription + `
			By default returns only active limits. Set "only_active" to false for the
			full history.`),
		mcp.WithOutputSchema[dto.PreApprovedLimitsResultDTO](),
		mcp.WithBoolean("only_active", mcp.Description("Return only active limits. Defaults to true"), mcp.DefaultBool(true)),
	}, withCustomerRefParams()...)
	return mcp.NewTool("get_pre_approved_limits", opts...)
}

func (s *Server) HandleGetPreApprovedLimits(ctx context.Context, request mcp.CallToolRequest, args mcp.CallToolParams) (*dto.PreApprovedLimitsResultDTO, error) {
	ref, err := customerRefFrom(request)
	if err != nil {
		return nil, err
	}
	result, restErr := s.internalRegistryService.GetPreApprovedLimits(ctx, ref, request.GetBool("only_active", true))
	if restErr != nil {
		return nil, restErr
	}
	return result, nil
}

func (s *Server) GetIncomeDeclarationsTool() mcp.Tool {
	opts := append([]mcp.ToolOption{
		mcp.WithDescription(`
			Internal Registry: get the income the customer declared to the institution
			during onboarding/relationship. Each declaration reports the type, monthly
			and yearly amount, source, and whether it was verified.
			` + customerRefDescription + `
			Set "verified_only" to true to return only verified declarations.`),
		mcp.WithOutputSchema[dto.IncomeDeclarationsResultDTO](),
		mcp.WithBoolean("verified_only", mcp.Description("Return only verified declarations. Defaults to false"), mcp.DefaultBool(false)),
	}, withCustomerRefParams()...)
	return mcp.NewTool("get_income_declarations", opts...)
}

func (s *Server) HandleGetIncomeDeclarations(ctx context.Context, request mcp.CallToolRequest, args mcp.CallToolParams) (*dto.IncomeDeclarationsResultDTO, error) {
	ref, err := customerRefFrom(request)
	if err != nil {
		return nil, err
	}
	result, restErr := s.internalRegistryService.GetIncomeDeclarations(ctx, ref, request.GetBool("verified_only", false))
	if restErr != nil {
		return nil, restErr
	}
	return result, nil
}
