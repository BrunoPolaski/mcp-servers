package middlewares

import (
	"net/http"
	"strings"

	httphelper "github.com/BrunoPolaski/bureau/internal/interfaces/http"
	"github.com/BrunoPolaski/go-rest-err/rest_err"
)

func BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
			httphelper.ErrorResponse(
				rest_err.NewUnauthorizedError("basic auth header not found or invalid"),
				w,
			)
			return
		}
		next.ServeHTTP(w, r)
	})
}
