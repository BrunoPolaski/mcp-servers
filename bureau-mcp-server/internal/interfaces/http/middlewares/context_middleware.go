package middlewares

import (
	"context"
	"net/http"
	"time"
)

// ContextMiddleware adds a timeout context to all requests
// 30 seconds should be enough for most operations including:
// - Database queries with joins
// - Email sending
// - Redis operations
// - File uploads
func ContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
