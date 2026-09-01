package main

import (
	"net/http"
	"os"
	"time"

	"github.com/BrunoPolaski/open-finance/internal/infra/thirdparty/logger"
	internal_mcp "github.com/BrunoPolaski/open-finance/internal/interfaces/mcp"
	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

// @title Open Finance MCP Server
// @version 1.0
// @description This is the API documentation for the Open Finance MCP Server.
//
// license.name MIT
// @host localhost:8080
//
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-api-key
//
// @securityDefinitions.basic BasicAuth
//
// @securitydefinitions.apiKey CookieAuth
// @in cookie
// @name Sid
func main() {
	logger.Init()

	location, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		logger.Error("Error loading location", err)
	}
	time.Local = location

	if os.Getenv("ENV") == "local" {
		err = godotenv.Overload(".env")
		if err != nil {
			logger.Error("Error loading .env file", err)
		}
	}

	port := os.Getenv("MCP_PORT")
	if port == "" {
		port = "8080"
	}

	mcpServer := internal_mcp.InitMCPServer()
	httpServer := server.NewStreamableHTTPServer(
		mcpServer,
		server.WithHeartbeatInterval(30*time.Second),
	)

	logger.Info("Starting MCP server",
		zap.String("addr", ":"+port),
		zap.String("endpoint", "/mcp"),
	)

	http.ListenAndServe(":8080", httpServer)
}
