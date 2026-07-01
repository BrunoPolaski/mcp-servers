package middlewares

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry-mcp-server/internal/helper/sessions"
	"github.com/BrunoPolaski/internal-registry-mcp-server/internal/infra/repositories/interfaces"
	internaljwt "github.com/BrunoPolaski/internal-registry-mcp-server/internal/infra/thirdparty/jwt"
	httphelper "github.com/BrunoPolaski/internal-registry-mcp-server/internal/interfaces/http"
	"github.com/golang-jwt/jwt/v5"
)

func SessionAuthMiddleware(sessionRepo interfaces.SessionRepository) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			var cookie *http.Cookie
			var err error

			cookieAsHeader := r.Header.Get(sessions.CookieName)
			if cookieAsHeader != "" {
				cookie = &http.Cookie{Name: sessions.CookieName, Value: cookieAsHeader}
			} else {
				cookie, err = r.Cookie(sessions.CookieName)
				if err != nil {
					httphelper.ErrorResponse(rest_err.NewUnauthorizedError("invalid session"), w)
					return
				}
			}

			jwtAdapter := internaljwt.NewJWTAdapter()

			parsedToken, restErr := jwtAdapter.ParseToken(cookie.Value)
			if restErr != nil {
				httphelper.ErrorResponse(restErr, w)
				return
			}

			claims, ok := parsedToken.Claims.(jwt.MapClaims)
			if !ok || !parsedToken.Valid {
				httphelper.ErrorResponse(rest_err.NewUnauthorizedError("invalid token: %v", parsedToken.Claims), w)
				return
			}

			exp, err := claims.GetExpirationTime()
			if err != nil || exp.Before(time.Now()) {
				httphelper.ErrorResponse(rest_err.NewUnauthorizedError("token expired"), w)
				return
			}

			tidRaw, exists := claims["tid"]
			if !exists {
				httphelper.ErrorResponse(rest_err.NewUnauthorizedError("invalid token"), w)
				return
			}
			tokenID, ok := tidRaw.(string)
			if !ok {
				httphelper.ErrorResponse(rest_err.NewUnauthorizedError("invalid token"), w)
				return
			}

			data, repoErr := sessionRepo.GetById(ctx, tokenID)
			if repoErr != nil {
				if repoErr.Code == http.StatusNotFound {
					httphelper.ErrorResponse(rest_err.NewUnauthorizedError("invalid session"), w)
				} else {
					httphelper.ErrorResponse(rest_err.NewInternalServerError("internal server error"), w)
				}
				return
			}

			hash := sha256.Sum256([]byte(cookie.Value))
			if fmt.Sprintf("%x", hash[:]) != data.TokenHash {
				httphelper.ErrorResponse(rest_err.NewUnauthorizedError("invalid session"), w)
				return
			}

			if data.IsRevoked {
				httphelper.ErrorResponse(rest_err.NewUnauthorizedError("session revoked"), w)
				return
			}

			userID, _ := claims.GetSubject()
			ctx = context.WithValue(ctx, "user_id", userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
