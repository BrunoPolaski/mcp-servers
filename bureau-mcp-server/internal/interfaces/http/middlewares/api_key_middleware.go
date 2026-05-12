package middlewares

import (
	"net/http"

	"github.com/BrunoPolaski/bureau-mcp-server/internal/infra/repositories/interfaces"
	"github.com/BrunoPolaski/bureau-mcp-server/internal/infra/thirdparty/logger"
	httphelper "github.com/BrunoPolaski/bureau-mcp-server/internal/interfaces/http"
	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func ApiKeyMiddleware(apiKeyRepository interfaces.ApiKeyRepository) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("x-api-key")

			logger.Info("Validating API key", zap.String("api_key", apiKey))

			if _, err := uuid.Parse(apiKey); err != nil {
				httphelper.ErrorResponse(rest_err.NewUnauthorizedError("invalid api key"), w)
				return
			}

			if _, err := apiKeyRepository.GetById(r.Context(), apiKey); err != nil {
				httphelper.ErrorResponse(err, w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
