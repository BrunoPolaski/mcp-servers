package middlewares

import (
	"net/http"

	valueobjects "github.com/BrunoPolaski/bureau/internal/core/entities/value_objects"
	httphelper "github.com/BrunoPolaski/bureau/internal/interfaces/http"
	"github.com/BrunoPolaski/go-rest-err/rest_err"
)

func UserTypeMiddleware(wanted valueobjects.UserType) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userType := r.Context().Value("user_type").(valueobjects.UserType)

			if userType.GreaterOrEqualThan(wanted) {
				next.ServeHTTP(w, r)
				return
			}

			httphelper.ErrorResponse(rest_err.NewForbiddenError("forbidden"), w) // Hide existence of resource
		})
	}
}
