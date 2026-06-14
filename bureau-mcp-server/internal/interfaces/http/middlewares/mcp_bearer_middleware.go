package middlewares

import (
	"net/http"
	"os"
)

func MCPBearerMiddleware(next http.Handler) http.Handler {
	token := os.Getenv("MCP_AUTH_TOKEN")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+token {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
