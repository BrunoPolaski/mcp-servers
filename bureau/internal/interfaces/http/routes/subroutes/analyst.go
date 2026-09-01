package subroutes

import (
	"net/http"

	valueobjects "github.com/BrunoPolaski/bureau/internal/core/entities/value_objects"
	"github.com/BrunoPolaski/bureau/internal/infra/controllers"
	"github.com/BrunoPolaski/bureau/internal/infra/repositories"
	"github.com/BrunoPolaski/bureau/internal/infra/thirdparty"
	"github.com/BrunoPolaski/bureau/internal/interfaces/http/middlewares"
	"github.com/BrunoPolaski/bureau/internal/services"
)

func RegisterAnalystRoutes(s *http.ServeMux, tpf *thirdparty.ThirdPartyFactory, rf *repositories.RepositoryFactory) {
	analystController := controllers.NewAnalystController(
		services.NewAnalystService(
			rf,
		),
	)

	s.Handle("POST /analyst", middlewares.HandlerChain(
		analystController.Create,
		middlewares.SessionAuthMiddleware(rf.SessionRepository()),
	))

	s.Handle("GET /analyst/{id}", middlewares.HandlerChain(
		analystController.GetById,
		middlewares.SessionAuthMiddleware(rf.SessionRepository()),
	))

	s.Handle("GET /analyst", middlewares.HandlerChain(
		analystController.GetAll,
		middlewares.SessionAuthMiddleware(rf.SessionRepository()),
	))

	s.Handle("DELETE /analyst/{id}", middlewares.HandlerChain(
		analystController.Delete,
		middlewares.SessionAuthMiddleware(rf.SessionRepository()),
		middlewares.UserTypeMiddleware(valueobjects.UserTypeAdmin),
	))
}
