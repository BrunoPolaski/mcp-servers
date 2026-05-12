package internal_mcp

import (
	"github.com/BrunoPolaski/bureau-mcp-server/internal/infra/repositories"
	"github.com/BrunoPolaski/bureau-mcp-server/internal/infra/thirdparty"
	"github.com/BrunoPolaski/bureau-mcp-server/internal/infra/thirdparty/logger"
	"github.com/BrunoPolaski/bureau-mcp-server/internal/infra/thirdparty/validator"
	"github.com/BrunoPolaski/bureau-mcp-server/internal/interfaces/mcp/tools"
	"github.com/BrunoPolaski/bureau-mcp-server/internal/services"
	"github.com/mark3labs/mcp-go/server"
)

func InitMCPServer() *server.MCPServer {
	logger.Info("Setting up MCP server...")

	tpf := thirdparty.NewThirdPartyFactory()
	rf := repositories.NewRepositoryFactory(tpf)
	validator.InitValidator()

	s := tools.NewMCPServer(
		services.NewUserService(rf),
		services.NewAddressService(rf),
		services.NewAnalystService(rf),
		services.NewPersonService(rf),
	)

	logger.Info("MCP server set up successfully.")

	return s
}
