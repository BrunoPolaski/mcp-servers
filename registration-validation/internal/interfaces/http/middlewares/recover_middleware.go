package middlewares

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/registration-validation/internal/infra/thirdparty/logger"
	httphelper "github.com/BrunoPolaski/registration-validation/internal/interfaces/http"
)

const stackTraceFilter = "github.com/BrunoPolaski/registration-validation"

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Error("Panic recovered", nil)
				logger.Error(fmt.Sprintf("Error: %v", rec), nil)
				logger.Error(fmt.Sprintf("Path: %s %s", r.Method, r.URL.Path), nil)
				logger.Error("----------------------STACKTRACE------------------------------", nil)
				stack := strings.Split(string(debug.Stack()), "\n")
				for i, line := range stack {
					if strings.Contains(line, stackTraceFilter) {
						logger.Error(fmt.Sprintf("%s at line %s", line, strings.Split(stack[i+1], ":")[1]), nil)
					}
				}

				httphelper.JSONWithStatus(rest_err.NewInternalServerError("internal server error"), http.StatusInternalServerError, w)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
