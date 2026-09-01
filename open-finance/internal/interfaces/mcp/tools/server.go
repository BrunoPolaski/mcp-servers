package tools

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/open-finance/internal/core/entities"
	"github.com/BrunoPolaski/open-finance/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/open-finance/internal/interfaces/http/middlewares"
	"github.com/BrunoPolaski/open-finance/internal/services"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// PersonService e OpenFinanceService são declarados aqui, no consumidor, para
// que os handlers possam ser exercitados com dublês nos testes.
type PersonService interface {
	GetById(ctx context.Context, id uint) (*entities.Person, *rest_err.RestErr)
	GetByDocument(ctx context.Context, document string) (*entities.Person, *rest_err.RestErr)
	GetAllSummary(ctx context.Context, limit, offset int, params map[string]any) (*dto.PaginatedResponse[dto.PersonSummaryDTO], *rest_err.RestErr)
}

type OpenFinanceService interface {
	GetBankStatements(ctx context.Context, ref services.CustomerRef, accountType string) (*dto.BankStatementsResultDTO, *rest_err.RestErr)
	GetCashFlowAnalyses(ctx context.Context, ref services.CustomerRef, limit int) (*dto.CashFlowAnalysesResultDTO, *rest_err.RestErr)
	GetRecurringTransactions(ctx context.Context, ref services.CustomerRef, transactionType string, onlyActive bool) (*dto.RecurringTransactionsResultDTO, *rest_err.RestErr)
	GetDataSharingConsents(ctx context.Context, ref services.CustomerRef) (*dto.DataSharingConsentsResultDTO, *rest_err.RestErr)
}

type Server struct {
	userService        *services.UserService
	addressService     *services.AddressService
	analystService     *services.AnalystService
	personService      PersonService
	openFinanceService OpenFinanceService
}

func NewMCPServer(
	userService *services.UserService,
	addressService *services.AddressService,
	analystService *services.AnalystService,
	personService PersonService,
	openFinanceService OpenFinanceService,
) *server.MCPServer {
	s := &Server{
		userService:        userService,
		addressService:     addressService,
		analystService:     analystService,
		personService:      personService,
		openFinanceService: openFinanceService,
	}
	mcpSrv := server.NewMCPServer(
		"open-finance-mcp",
		"1.0.0",
		server.WithRecovery(),
		server.WithOutputSchemaValidation(),
		server.WithToolHandlerMiddleware(middlewares.MCPLogMiddleware),
	)
	s.registerTools(mcpSrv)
	return mcpSrv
}

func (s *Server) registerTools(mcpSrv *server.MCPServer) {
	mcpSrv.AddTool(s.GetPersonByIDTool(), mcp.NewStructuredToolHandler(s.HandleGetPersonByID))
	mcpSrv.AddTool(s.GetPersonByDocumentTool(), mcp.NewStructuredToolHandler(s.HandleGetPersonByDocument))
	mcpSrv.AddTool(s.GetAllPersonsTool(), mcp.NewStructuredToolHandler(s.HandleGetAllPersons))

	mcpSrv.AddTool(s.GetBankStatementsTool(), mcp.NewStructuredToolHandler(s.HandleGetBankStatements))
	mcpSrv.AddTool(s.GetCashFlowAnalysisTool(), mcp.NewStructuredToolHandler(s.HandleGetCashFlowAnalysis))
	mcpSrv.AddTool(s.GetRecurringTransactionsTool(), mcp.NewStructuredToolHandler(s.HandleGetRecurringTransactions))
	mcpSrv.AddTool(s.GetDataSharingConsentsTool(), mcp.NewStructuredToolHandler(s.HandleGetDataSharingConsents))
}
