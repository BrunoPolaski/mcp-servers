package tools

import (
	"context"
	"fmt"

	"github.com/BrunoPolaski/open-finance/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/open-finance/internal/services"
	"github.com/mark3labs/mcp-go/mcp"
)

const customerRefDescription = `Identify the customer by "customer_id" OR by "document" (CPF). Provide exactly one of them.`

// customerRefFrom apenas transporta a identificação recebida. A recusa de
// referência ambígua ou vazia cabe ao serviço, que é onde a regra vive.
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

func (s *Server) GetBankStatementsTool() mcp.Tool {
	opts := append(
		[]mcp.ToolOption{
			mcp.WithDescription(
				`
				Open Finance: get the customer's bank statements for the last 90 days.
				Each statement covers one month of one account at one institution and
				reports opening and closing balances, total credits and debits, and the
				number of transactions.
				` + customerRefDescription + `
				Example usage:
				{
					"document": "12345678900",
					"account_type": "checking"
				}
				`,
			),
			mcp.WithOutputSchema[dto.BankStatementsResultDTO](),
			mcp.WithString(
				"account_type",
				mcp.Description("Optional filter by account type"),
				mcp.Enum("checking", "savings", "payment"),
			),
		},
		withCustomerRefParams()...,
	)

	return mcp.NewTool("get_bank_statements", opts...)
}

func (s *Server) HandleGetBankStatements(ctx context.Context, request mcp.CallToolRequest, args mcp.CallToolParams) (*dto.BankStatementsResultDTO, error) {
	ref, err := customerRefFrom(request)
	if err != nil {
		return nil, err
	}

	result, restErr := s.openFinanceService.GetBankStatements(ctx, ref, request.GetString("account_type", ""))
	if restErr != nil {
		return nil, restErr
	}

	return result, nil
}

func (s *Server) GetCashFlowAnalysisTool() mcp.Tool {
	opts := append(
		[]mcp.ToolOption{
			mcp.WithDescription(
				`
				Open Finance: get the customer's cash flow analysis derived from the
				shared bank statements. Reports average monthly inflow and outflow, net
				cash flow, inflow volatility, days with a negative balance, and whether
				a recurring income was detected.
				` + customerRefDescription + `
				By default returns only the most recent analysis. Set "limit" to 0 to
				retrieve the full history.
				Example usage:
				{
					"customer_id": 123
				}
				`,
			),
			mcp.WithOutputSchema[dto.CashFlowAnalysesResultDTO](),
			mcp.WithInteger(
				"limit",
				mcp.Description("How many analyses to return, most recent first. 0 returns all"),
				mcp.DefaultNumber(1),
				mcp.Min(0),
			),
		},
		withCustomerRefParams()...,
	)

	return mcp.NewTool("get_cash_flow_analysis", opts...)
}

func (s *Server) HandleGetCashFlowAnalysis(ctx context.Context, request mcp.CallToolRequest, args mcp.CallToolParams) (*dto.CashFlowAnalysesResultDTO, error) {
	ref, err := customerRefFrom(request)
	if err != nil {
		return nil, err
	}

	limit := request.GetInt("limit", 1)
	if limit < 0 {
		return nil, fmt.Errorf("invalid limit: must be a non-negative integer")
	}

	result, restErr := s.openFinanceService.GetCashFlowAnalyses(ctx, ref, limit)
	if restErr != nil {
		return nil, restErr
	}

	return result, nil
}

func (s *Server) GetRecurringTransactionsTool() mcp.Tool {
	opts := append(
		[]mcp.ToolOption{
			mcp.WithDescription(
				`
				Open Finance: get the customer's recurring incomes and fixed expenses
				identified from the shared transaction history. Each entry reports its
				type, category, amount, frequency, counterparty and whether it is still
				active.
				` + customerRefDescription + `
				Example usage:
				{
					"customer_id": 123,
					"transaction_type": "income"
				}
				`,
			),
			mcp.WithOutputSchema[dto.RecurringTransactionsResultDTO](),
			mcp.WithString(
				"transaction_type",
				mcp.Description("Optional filter by transaction type"),
				mcp.Enum("income", "expense"),
			),
			mcp.WithBoolean(
				"only_active",
				mcp.Description("Return only transactions still active. Defaults to true"),
				mcp.DefaultBool(true),
			),
		},
		withCustomerRefParams()...,
	)

	return mcp.NewTool("get_recurring_transactions", opts...)
}

func (s *Server) HandleGetRecurringTransactions(ctx context.Context, request mcp.CallToolRequest, args mcp.CallToolParams) (*dto.RecurringTransactionsResultDTO, error) {
	ref, err := customerRefFrom(request)
	if err != nil {
		return nil, err
	}

	result, restErr := s.openFinanceService.GetRecurringTransactions(
		ctx,
		ref,
		request.GetString("transaction_type", ""),
		request.GetBool("only_active", true),
	)
	if restErr != nil {
		return nil, restErr
	}

	return result, nil
}

func (s *Server) GetDataSharingConsentsTool() mcp.Tool {
	opts := append(
		[]mcp.ToolOption{
			mcp.WithDescription(
				`
				Open Finance: get the customer's data sharing consents. Each consent
				reports the institution, its status (granted, revoked, expired,
				awaiting), the authorized scope, and the grant, expiry and revocation
				timestamps.
				A customer without an active consent has no Open Finance data available;
				treat the corresponding criteria as missing data rather than as zero.
				` + customerRefDescription + `
				Example usage:
				{
					"document": "12345678900"
				}
				`,
			),
			mcp.WithOutputSchema[dto.DataSharingConsentsResultDTO](),
		},
		withCustomerRefParams()...,
	)

	return mcp.NewTool("get_data_sharing_consents", opts...)
}

func (s *Server) HandleGetDataSharingConsents(ctx context.Context, request mcp.CallToolRequest, args mcp.CallToolParams) (*dto.DataSharingConsentsResultDTO, error) {
	ref, err := customerRefFrom(request)
	if err != nil {
		return nil, err
	}

	result, restErr := s.openFinanceService.GetDataSharingConsents(ctx, ref)
	if restErr != nil {
		return nil, restErr
	}

	return result, nil
}
