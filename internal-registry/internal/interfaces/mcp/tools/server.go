package tools

import (
	"github.com/BrunoPolaski/internal-registry/internal/interfaces/http/middlewares"
	"github.com/BrunoPolaski/internal-registry/internal/services"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Server struct {
	userService    *services.UserService
	addressService *services.AddressService
	analystService *services.AnalystService
	personService  *services.PersonService
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
}
