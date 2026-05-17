package middlewares

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities"
	"github.com/BrunoPolaski/bureau-mcp-server/internal/helper/sessions"
	httphelper "github.com/BrunoPolaski/bureau-mcp-server/internal/interfaces/http"
	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/redis/go-redis/v9"
)

func SessionAuthMiddleware(rdb *redis.Client) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			var cookie *http.Cookie
			var err error

			cookieAsHeader := r.Header.Get(sessions.CookieName)
			if cookieAsHeader != "" {
				cookie = &http.Cookie{
					Name:  sessions.CookieName,
					Value: cookieAsHeader,
				}
			} else {
				cookie, err = r.Cookie(sessions.CookieName)
				if err != nil {
					httphelper.ErrorResponse(rest_err.NewUnauthorizedError("invalid session"), w)
					return
				}
			}

			data, err := rdb.Get(ctx, cookie.Value).Bytes()
			if errors.Is(err, redis.Nil) {
				httphelper.ErrorResponse(rest_err.NewUnauthorizedError("invalid session"), w)
				return
			} else if err != nil {
				httphelper.ErrorResponse(rest_err.NewInternalServerError("internal server error").WithCause(err), w)
				return
			}

			var sess entities.Session
			if err := json.Unmarshal(data, &sess); err != nil {
				httphelper.ErrorResponse(rest_err.NewInternalServerError("internal server error").WithCause(err), w)
				return
			}

			now := time.Now()

			if now.Sub(sess.LastActivity) > sessions.IdleTimeout {
				rdb.Del(ctx, cookie.Value)
				httphelper.ErrorResponse(rest_err.NewUnauthorizedError("session expired"), w)
				return
			}

			sess.LastActivity = now
			updated, err := json.Marshal(sess)
			if err != nil {
				httphelper.ErrorResponse(rest_err.NewInternalServerError("internal server error").WithCause(err), w)
				return
			}

			if err := rdb.Set(ctx, cookie.Value, updated, sessions.AbsoluteTimeout).Err(); err != nil {
				httphelper.ErrorResponse(rest_err.NewInternalServerError("internal server error").WithCause(err), w)
				return
			}

			ctx = context.WithValue(ctx, "user_id", sess.UserID)
			ctx = context.WithValue(ctx, "user_type", sess.UserType)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
