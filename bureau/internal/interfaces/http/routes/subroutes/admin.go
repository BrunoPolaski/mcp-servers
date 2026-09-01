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

func RegisterAdminRoutes(s *http.ServeMux, tpf *thirdparty.ThirdPartyFactory, rf *repositories.RepositoryFactory) {
	adminController := controllers.NewAdminController(
		services.NewAdminService(
			rf,
		),
	)

	s.Handle("POST /admin", middlewares.HandlerChain(
		adminController.Create,
		middlewares.SessionAuthMiddleware(rf.SessionRepository()),
		middlewares.UserTypeMiddleware(valueobjects.UserTypeAdmin),
	))

	s.Handle("GET /admin/{id}", middlewares.HandlerChain(
		adminController.GetById,
		middlewares.SessionAuthMiddleware(rf.SessionRepository()),
		middlewares.UserTypeMiddleware(valueobjects.UserTypeAdmin),
	))

	s.Handle("GET /admin", middlewares.HandlerChain(
		adminController.GetAll,
		middlewares.SessionAuthMiddleware(rf.SessionRepository()),
		middlewares.UserTypeMiddleware(valueobjects.UserTypeAdmin),
	))

	s.Handle("DELETE /admin/{id}", middlewares.HandlerChain(
		adminController.Delete,
		middlewares.SessionAuthMiddleware(rf.SessionRepository()),
		middlewares.UserTypeMiddleware(valueobjects.UserTypeAdmin),
	))
}
