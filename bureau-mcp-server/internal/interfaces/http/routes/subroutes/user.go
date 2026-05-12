package subroutes

import (
	"net/http"

	valueobjects "github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities/value_objects"
	"github.com/BrunoPolaski/bureau-mcp-server/internal/infra/controllers"
	"github.com/BrunoPolaski/bureau-mcp-server/internal/infra/repositories"
	"github.com/BrunoPolaski/bureau-mcp-server/internal/infra/thirdparty"
	"github.com/BrunoPolaski/bureau-mcp-server/internal/interfaces/http/middlewares"
	"github.com/BrunoPolaski/bureau-mcp-server/internal/services"
)

func RegisterUserRoutes(s *http.ServeMux, tpf *thirdparty.ThirdPartyFactory, rf *repositories.RepositoryFactory) {
	userController := controllers.NewUserController(
		services.NewUserService(
			rf,
		),
	)

	s.Handle("GET /user", middlewares.HandlerChain(
		userController.GetAll,
		middlewares.SessionAuthMiddleware(tpf.Redis()),
	))

	s.Handle("GET /user/{id}", middlewares.HandlerChain(
		userController.GetById,
		middlewares.SessionAuthMiddleware(tpf.Redis()),
	))

	s.Handle("DELETE /user/{id}", middlewares.HandlerChain(
		userController.Delete,
		middlewares.SessionAuthMiddleware(tpf.Redis()),
		middlewares.UserTypeMiddleware(valueobjects.UserTypeAdmin),
	))
}
