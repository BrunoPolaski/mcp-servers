package subroutes

import (
	"net/http"

	valueobjects "github.com/BrunoPolaski/open-finance/internal/core/entities/value_objects"
	"github.com/BrunoPolaski/open-finance/internal/infra/controllers"
	"github.com/BrunoPolaski/open-finance/internal/infra/repositories"
	"github.com/BrunoPolaski/open-finance/internal/infra/thirdparty"
	"github.com/BrunoPolaski/open-finance/internal/interfaces/http/middlewares"
	"github.com/BrunoPolaski/open-finance/internal/services"
)

func RegisterPersonRoutes(s *http.ServeMux, tpf *thirdparty.ThirdPartyFactory, rf *repositories.RepositoryFactory) {
	personController := controllers.NewPersonController(
		services.NewPersonService(
			rf,
		),
	)

	s.Handle("POST /person", middlewares.HandlerChain(
		personController.Create,
		middlewares.SessionAuthMiddleware(rf.SessionRepository()),
	))

	s.Handle("GET /person/{id}", middlewares.HandlerChain(
		personController.GetById,
		middlewares.SessionAuthMiddleware(rf.SessionRepository()),
	))

	s.Handle("GET /person/document/{document}", middlewares.HandlerChain(
		personController.GetByDocument,
		middlewares.SessionAuthMiddleware(rf.SessionRepository()),
	))

	s.Handle("GET /person", middlewares.HandlerChain(
		personController.GetAll,
		middlewares.SessionAuthMiddleware(rf.SessionRepository()),
		middlewares.UserTypeMiddleware(valueobjects.UserTypeAnalyst),
	))

	s.Handle("DELETE /person/{id}", middlewares.HandlerChain(
		personController.Delete,
		middlewares.SessionAuthMiddleware(rf.SessionRepository()),
		middlewares.UserTypeMiddleware(valueobjects.UserTypeAdmin),
	))
}
