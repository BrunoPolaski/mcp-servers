package middlewares

import (
	"net/http"
	"time"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	internaljwt "github.com/BrunoPolaski/open-finance-mcp-server/internal/infra/thirdparty/jwt"
	"github.com/BrunoPolaski/open-finance-mcp-server/internal/infra/thirdparty/logger"
	httphelper "github.com/BrunoPolaski/open-finance-mcp-server/internal/interfaces/http"
	"github.com/golang-jwt/jwt/v5"
)

func BearerAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		jwtAdapter := internaljwt.NewJWTAdapter()
		token, restErr := jwtAdapter.TrimPrefix(header)
		if restErr != nil {
			httphelper.ErrorResponse(restErr, w)
			return
		}

		parsedToken, restErr := jwtAdapter.ParseToken(token)
		if restErr != nil {
			httphelper.ErrorResponse(restErr, w)
			return
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok || !parsedToken.Valid {
			httphelper.ErrorResponse(rest_err.NewUnauthorizedError("invalid token"), w)
			return
		}

		exp, err := claims.GetExpirationTime()
		if err != nil {
			logger.Error("failed to get expiration time from claims: %s", err)
			httphelper.ErrorResponse(rest_err.NewUnauthorizedError("invalid token"), w)
			return
		}

		if exp.Before(time.Now()) {
			httphelper.ErrorResponse(rest_err.NewUnauthorizedError("token expired"), w)
			return
		}

		next.ServeHTTP(w, r)
	})
}
