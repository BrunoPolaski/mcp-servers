package subroutes

import (
	"net/http"

	"github.com/BrunoPolaski/registration-validation/internal/infra/controllers"
	"github.com/BrunoPolaski/registration-validation/internal/infra/repositories"
	"github.com/BrunoPolaski/registration-validation/internal/infra/thirdparty"
	"github.com/BrunoPolaski/registration-validation/internal/interfaces/http/middlewares"
	"github.com/BrunoPolaski/registration-validation/internal/services"
)

func RegisterAuthRoutes(s *http.ServeMux, tpf *thirdparty.ThirdPartyFactory, rf *repositories.RepositoryFactory) {
	authController := controllers.NewAuthController(
		services.NewAuthService(
			rf,
			tpf,
		),
	)

	s.Handle("POST /auth/signin",
		middlewares.HandlerChain(
			authController.SignIn,
			middlewares.BasicAuthMiddleware,
			middlewares.ApiKeyMiddleware(rf.ApiKeyRepository()),
		),
	)

	s.Handle("POST /auth/register",
		middlewares.HandlerChain(
			authController.Register,
			middlewares.ApiKeyMiddleware(rf.ApiKeyRepository()),
		),
	)
}
