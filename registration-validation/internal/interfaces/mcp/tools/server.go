package tools

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/registration-validation/internal/core/entities"
	"github.com/BrunoPolaski/registration-validation/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/registration-validation/internal/interfaces/http/middlewares"
	"github.com/BrunoPolaski/registration-validation/internal/services"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type PersonService interface {
	GetById(ctx context.Context, id uint) (*entities.Person, *rest_err.RestErr)
	GetByDocument(ctx context.Context, document string) (*entities.Person, *rest_err.RestErr)
	GetAll(ctx context.Context, limit, offset int, params map[string]any) (*dto.PaginatedResponse[dto.PersonDTO], *rest_err.RestErr)
}

type RegistrationValidationService interface {
	GetDocumentValidations(ctx context.Context, ref services.CustomerRef, documentType string) (*dto.DocumentValidationsResultDTO, *rest_err.RestErr)
	GetFiscalRegularities(ctx context.Context, ref services.CustomerRef) (*dto.FiscalRegularitiesResultDTO, *rest_err.RestErr)
	GetEmploymentLinks(ctx context.Context, ref services.CustomerRef, status string) (*dto.EmploymentLinkValidationsResultDTO, *rest_err.RestErr)
	GetComplianceChecks(ctx context.Context, ref services.CustomerRef, checkType string) (*dto.ComplianceChecksResultDTO, *rest_err.RestErr)
}

type Server struct {
	userService    *services.UserService
	addressService *services.AddressService
	analystService *services.AnalystService
	personService  PersonService

	registrationValidationService RegistrationValidationService
}

func NewMCPServer(
	userService *services.UserService,
	addressService *services.AddressService,
	analystService *services.AnalystService,
	personService *services.PersonService,
	registrationValidationService *services.RegistrationValidationService,
) *server.MCPServer {
	s := &Server{
		userService:    userService,
		addressService: addressService,
		analystService: analystService,
		personService:  personService,

		registrationValidationService: registrationValidationService,
	}
	mcpSrv := server.NewMCPServer(
		"registration-validation-mcp",
		"1.0.0",
		server.WithRecovery(),
		server.WithOutputSchemaValidation(),
		server.WithToolHandlerMiddleware(middlewares.MCPLogMiddleware),
	)
	s.registerTools(mcpSrv)
	return mcpSrv
}

func (s *Server) registerTools(mcpSrv *server.MCPServer) {
	mcpSrv.AddTool(s.GetPersonByIDTool(), mcp.NewStructuredToolHandler(traced(s.HandleGetPersonByID)))
	mcpSrv.AddTool(s.GetPersonByDocumentTool(), mcp.NewStructuredToolHandler(traced(s.HandleGetPersonByDocument)))
	mcpSrv.AddTool(s.GetAllPersonsTool(), mcp.NewStructuredToolHandler(traced(s.HandleGetAllPersons)))

	mcpSrv.AddTool(s.GetDocumentValidationsTool(), mcp.NewStructuredToolHandler(traced(s.HandleGetDocumentValidations)))
	mcpSrv.AddTool(s.GetFiscalRegularityTool(), mcp.NewStructuredToolHandler(traced(s.HandleGetFiscalRegularity)))
	mcpSrv.AddTool(s.GetEmploymentLinksTool(), mcp.NewStructuredToolHandler(traced(s.HandleGetEmploymentLinks)))
	mcpSrv.AddTool(s.GetComplianceChecksTool(), mcp.NewStructuredToolHandler(traced(s.HandleGetComplianceChecks)))
}
