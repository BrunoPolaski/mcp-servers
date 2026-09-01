package middlewares

import (
	"net/http"

	"github.com/mark3labs/mcp-go/server"
)

func HandlerChain(h http.HandlerFunc, middlewares ...Middleware) http.Handler {
	var handler http.Handler = h

	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	handler = LogMiddleware(handler)
	handler = RecoverMiddleware(handler)

	return handler
}

func MCPHandlerChain(h server.ToolHandlerFunc, middlewares ...func(server.ToolHandlerFunc) server.ToolHandlerFunc) server.ToolHandlerFunc {
	var handler server.ToolHandlerFunc = h

	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}
