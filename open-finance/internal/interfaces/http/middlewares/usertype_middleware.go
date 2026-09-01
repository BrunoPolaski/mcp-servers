package middlewares

import (
	"net/http"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	valueobjects "github.com/BrunoPolaski/open-finance/internal/core/entities/value_objects"
	httphelper "github.com/BrunoPolaski/open-finance/internal/interfaces/http"
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
