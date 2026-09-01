package subroutes

import (
	"net/http"

	valueobjects "github.com/BrunoPolaski/internal-registry/internal/core/entities/value_objects"
	"github.com/BrunoPolaski/internal-registry/internal/infra/controllers"
	"github.com/BrunoPolaski/internal-registry/internal/infra/repositories"
	"github.com/BrunoPolaski/internal-registry/internal/infra/thirdparty"
	"github.com/BrunoPolaski/internal-registry/internal/interfaces/http/middlewares"
	"github.com/BrunoPolaski/internal-registry/internal/services"
)

func RegisterUserRoutes(s *http.ServeMux, tpf *thirdparty.ThirdPartyFactory, rf *repositories.RepositoryFactory) {
	userController := controllers.NewUserController(
		services.NewUserService(
			rf,
		),
	)

	s.Handle("GET /user", middlewares.HandlerChain(
		userController.GetAll,
		middlewares.SessionAuthMiddleware(rf.SessionRepository()),
	))

	s.Handle("GET /user/{id}", middlewares.HandlerChain(
		userController.GetById,
		middlewares.SessionAuthMiddleware(rf.SessionRepository()),
	))

	s.Handle("DELETE /user/{id}", middlewares.HandlerChain(
		userController.Delete,
		middlewares.SessionAuthMiddleware(rf.SessionRepository()),
		middlewares.UserTypeMiddleware(valueobjects.UserTypeAdmin),
	))
}
