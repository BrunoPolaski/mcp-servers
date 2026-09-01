package middlewares

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/BrunoPolaski/registration-validation/internal/infra/thirdparty/logger"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		logger.Info("[REQUEST]",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Any("headers", r.Header),
		)

		start := time.Now()
		next.ServeHTTP(recorder, r)
		duration := time.Since(start)

		jsonBody := map[string]any{}
		if err := json.Unmarshal(recorder.body.Bytes(), &jsonBody); err != nil {
			jsonBody = map[string]any{
				"raw": recorder.body.String(),
			}
		}

		logger.Info("[RESPONSE]",
			zap.Int("status", recorder.statusCode),
			zap.String("duration", duration.String()),
			zap.Any("body", jsonBody),
		)
	})
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (r *responseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func MCPLogMiddleware(thf server.ToolHandlerFunc) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logger.Info("[MCP REQUEST]",
			zap.String("name", request.Params.Name),
			zap.Any("arguments", request.Params.Arguments),
		)

		start := time.Now()
		result, err := thf(ctx, request)
		duration := time.Since(start)

		if err != nil {
			logger.Error("[MCP RESPONSE ERROR]",
				err,
				zap.String("name", request.Params.Name),
				zap.String("duration", duration.String()),
				zap.Error(err),
			)
		} else {
			logger.Info("[MCP RESPONSE]",
				zap.String("name", request.Params.Name),
				zap.String("duration", duration.String()),
				zap.Any("result", result),
			)
		}

		return result, err
	}
}
