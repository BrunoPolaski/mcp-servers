package subroutes

import (
	"net/http"

	"github.com/BrunoPolaski/bureau/internal/infra/controllers"
	"github.com/BrunoPolaski/bureau/internal/infra/repositories"
	"github.com/BrunoPolaski/bureau/internal/infra/thirdparty"
	"github.com/BrunoPolaski/bureau/internal/interfaces/http/middlewares"
	"github.com/BrunoPolaski/bureau/internal/services"
)

func RegisterAddressRoutes(s *http.ServeMux, tpf *thirdparty.ThirdPartyFactory, rf *repositories.RepositoryFactory) {
	addressController := controllers.NewAddressController(
		services.NewAddressService(
			rf,
		),
	)

	s.Handle("POST /address", middlewares.HandlerChain(
		addressController.Create,
		middlewares.SessionAuthMiddleware(rf.SessionRepository()),
	))

	s.Handle("GET /address/{id}", middlewares.HandlerChain(
		addressController.GetById,
		middlewares.SessionAuthMiddleware(rf.SessionRepository()),
	))

	s.Handle("GET /address", middlewares.HandlerChain(
		addressController.GetAll,
		middlewares.SessionAuthMiddleware(rf.SessionRepository()),
	))

	s.Handle("DELETE /address/{id}", middlewares.HandlerChain(
		addressController.Delete,
		middlewares.SessionAuthMiddleware(rf.SessionRepository()),
	))
}
