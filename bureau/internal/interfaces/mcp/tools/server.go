package tools

import (
	"context"

	"github.com/BrunoPolaski/bureau/internal/core/entities"
	"github.com/BrunoPolaski/bureau/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/bureau/internal/interfaces/http/middlewares"
	"github.com/BrunoPolaski/bureau/internal/services"
	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type PersonService interface {
	GetById(ctx context.Context, id uint) (*entities.Person, *rest_err.RestErr)
	GetByDocument(ctx context.Context, document string) (*entities.Person, *rest_err.RestErr)
	GetAllSummary(ctx context.Context, limit, offset int, params map[string]any) (*dto.PaginatedResponse[dto.PersonSummaryDTO], *rest_err.RestErr)
}

type Server struct {
	userService    *services.UserService
	addressService *services.AddressService
	analystService *services.AnalystService
	personService  PersonService
}

func NewMCPServer(
	userService *services.UserService,
	addressService *services.AddressService,
	analystService *services.AnalystService,
	personService *services.PersonService,
) *server.MCPServer {
	s := &Server{
		userService:    userService,
		addressService: addressService,
		analystService: analystService,
		personService:  personService,
	}
	mcpSrv := server.NewMCPServer(
		"bureau-mcp",
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
}
