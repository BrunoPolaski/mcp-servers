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

func RegisterAnalystRoutes(s *http.ServeMux, tpf *thirdparty.ThirdPartyFactory, rf *repositories.RepositoryFactory) {
	analystController := controllers.NewAnalystController(
		services.NewAnalystService(
			rf,
		),
	)

	s.Handle("POST /analyst", middlewares.HandlerChain(
		analystController.Create,
		middlewares.SessionAuthMiddleware(tpf.Redis()),
	))

	s.Handle("GET /analyst/{id}", middlewares.HandlerChain(
		analystController.GetById,
		middlewares.SessionAuthMiddleware(tpf.Redis()),
	))

	s.Handle("GET /analyst", middlewares.HandlerChain(
		analystController.GetAll,
		middlewares.SessionAuthMiddleware(tpf.Redis()),
	))

	s.Handle("DELETE /analyst/{id}", middlewares.HandlerChain(
		analystController.Delete,
		middlewares.SessionAuthMiddleware(tpf.Redis()),
		middlewares.UserTypeMiddleware(valueobjects.UserTypeAdmin),
	))
}
