package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/BrunoPolaski/bureau/internal/infra/thirdparty/logger"
	"github.com/BrunoPolaski/bureau/internal/interfaces/http/routes"
	"github.com/joho/godotenv"
)

// @title Bureau API
// @version 1.0
// @description This is the API documentation for the Bureau MCP Server.

// license.name MIT
// @host localhost:8080

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-api-key

// @securityDefinitions.basic BasicAuth

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

	mux := http.NewServeMux()

	log.Fatal(http.ListenAndServe(":8080", routes.InitRoutes(mux)))
}
