package tools

import (
	"context"
	"time"

	"github.com/BrunoPolaski/registration-validation/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/registration-validation/internal/infra/thirdparty/logger"
	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

// source identifica esta fonte de dados nos metadados e nos logs.
const source = "registration-validation"

type traceable interface {
	SetMetadata(dto.ResponseMetadata)
}

// traced registra cada invocacao de ferramenta no log e carimba a resposta com
// os metadados de rastreabilidade. O transaction_id aparece nos dois lugares,
// de modo que log e resposta possam ser cruzados na analise dos experimentos.
func traced[T interface {
	traceable
	comparable
}](handler func(context.Context, mcp.CallToolRequest, mcp.CallToolParams) (T, error),
) func(context.Context, mcp.CallToolRequest, mcp.CallToolParams) (T, error) {
	return func(ctx context.Context, request mcp.CallToolRequest, args mcp.CallToolParams) (T, error) {
		start := time.Now()
		transactionID := uuid.NewString()
		tool := request.Params.Name

		result, err := handler(ctx, request, args)

		tags := []zap.Field{
			zap.String("source", source),
			zap.String("tool", tool),
			zap.String("transaction_id", transactionID),
			zap.Any("arguments", request.GetArguments()),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
		}

		var zero T
		if err != nil || result == zero {
			logger.Error("mcp tool call", err, tags...)
			return result, err
		}

		result.SetMetadata(dto.ResponseMetadata{
			Source:        source,
			Tool:          tool,
			QueriedAt:     start.UTC().Format(time.RFC3339),
			TransactionID: transactionID,
		})

		logger.Info("mcp tool call", tags...)
		return result, nil
	}
}
