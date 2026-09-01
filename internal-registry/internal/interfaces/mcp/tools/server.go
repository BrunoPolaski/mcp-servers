package tools

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
	"github.com/BrunoPolaski/internal-registry/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/internal-registry/internal/interfaces/http/middlewares"
	"github.com/BrunoPolaski/internal-registry/internal/services"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// PersonService e InternalRegistryService são declarados aqui, no consumidor, para
// que os handlers possam ser exercitados com dublês nos testes.
type PersonService interface {
	GetById(ctx context.Context, id uint) (*entities.Person, *rest_err.RestErr)
	GetByDocument(ctx context.Context, document string) (*entities.Person, *rest_err.RestErr)
	GetAllSummary(ctx context.Context, limit, offset int, params map[string]any) (*dto.PaginatedResponse[dto.PersonSummaryDTO], *rest_err.RestErr)
}

type InternalRegistryService interface {
	GetCustomerRelationship(ctx context.Context, ref services.CustomerRef) (*dto.CustomerRelationshipResultDTO, *rest_err.RestErr)
	GetContractedProducts(ctx context.Context, ref services.CustomerRef, productType, status string) (*dto.ContractedProductsResultDTO, *rest_err.RestErr)
	GetInternalPaymentRecords(ctx context.Context, ref services.CustomerRef, status string, productID *uint) (*dto.InternalPaymentRecordsResultDTO, *rest_err.RestErr)
	GetPreApprovedLimits(ctx context.Context, ref services.CustomerRef, onlyActive bool) (*dto.PreApprovedLimitsResultDTO, *rest_err.RestErr)
	GetIncomeDeclarations(ctx context.Context, ref services.CustomerRef, verifiedOnly bool) (*dto.IncomeDeclarationsResultDTO, *rest_err.RestErr)
}

type Server struct {
	userService             *services.UserService
	addressService          *services.AddressService
	analystService          *services.AnalystService
	personService           PersonService
	internalRegistryService InternalRegistryService
}

func NewMCPServer(
	userService *services.UserService,
	addressService *services.AddressService,
	analystService *services.AnalystService,
	personService PersonService,
	internalRegistryService InternalRegistryService,
) *server.MCPServer {
	s := &Server{
		userService:             userService,
		addressService:          addressService,
		analystService:          analystService,
		personService:           personService,
		internalRegistryService: internalRegistryService,
	}
	mcpSrv := server.NewMCPServer(
		"internal-registry-mcp",
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

	mcpSrv.AddTool(s.GetCustomerRelationshipTool(), mcp.NewStructuredToolHandler(s.HandleGetCustomerRelationship))
	mcpSrv.AddTool(s.GetContractedProductsTool(), mcp.NewStructuredToolHandler(s.HandleGetContractedProducts))
	mcpSrv.AddTool(s.GetInternalPaymentRecordsTool(), mcp.NewStructuredToolHandler(s.HandleGetInternalPaymentRecords))
	mcpSrv.AddTool(s.GetPreApprovedLimitsTool(), mcp.NewStructuredToolHandler(s.HandleGetPreApprovedLimits))
	mcpSrv.AddTool(s.GetIncomeDeclarationsTool(), mcp.NewStructuredToolHandler(s.HandleGetIncomeDeclarations))
}
