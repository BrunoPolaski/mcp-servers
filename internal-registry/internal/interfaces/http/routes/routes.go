package routes

import (
	"net/http"

	// _ "github.com/BrunoPolaski/internal-registry/docs"
	"github.com/BrunoPolaski/internal-registry/internal/infra/controllers"
	"github.com/BrunoPolaski/internal-registry/internal/infra/repositories"
	"github.com/BrunoPolaski/internal-registry/internal/infra/thirdparty"
	"github.com/BrunoPolaski/internal-registry/internal/infra/thirdparty/logger"
	"github.com/BrunoPolaski/internal-registry/internal/infra/thirdparty/validator"
	"github.com/BrunoPolaski/internal-registry/internal/interfaces/http/middlewares"
	"github.com/BrunoPolaski/internal-registry/internal/interfaces/http/routes/subroutes"

	// httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/time/rate"
)

func InitRoutes(mux *http.ServeMux) http.Handler {
	logger.Info("Setting up routes...")

	// Initialize the factory once
	tpf := thirdparty.NewThirdPartyFactory()
	rf := repositories.NewRepositoryFactory(tpf)

	// Initialize validator
	validator.InitValidator()

	// Create health controller
	healthController := controllers.NewHealthController(tpf.DB())

	// Register all subroutes with the factory
	subroutes.RegisterAuthRoutes(mux, tpf, rf)
	subroutes.RegisterUserRoutes(mux, tpf, rf)
	subroutes.RegisterAdminRoutes(mux, tpf, rf)
	subroutes.RegisterAddressRoutes(mux, tpf, rf)
	subroutes.RegisterAnalystRoutes(mux, tpf, rf)
	subroutes.RegisterPersonRoutes(mux, tpf, rf)

	// mux.Handle("/docs/swagger.json", http.FileServer(http.Dir("./docs")))
	// mux.Handle("/docs/", httpSwagger.WrapHandler)
	mux.Handle("/health", middlewares.HandlerChain(
		healthController.Check,
	))

	rateLimiter := middlewares.NewRateLimiter(rate.Limit(100), 200)
	rateLimiterMiddleware := middlewares.RateLimitMiddleware(rateLimiter)

	logger.Info("Routes set up successfully, listening on port 8080...")

	return rateLimiterMiddleware(mux)
}
