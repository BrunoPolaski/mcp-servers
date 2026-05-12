package mailer

import (
	"context"

	"github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities"
	"github.com/BrunoPolaski/go-rest-err/rest_err"
)

type Mailer interface {
	Send(ctx context.Context, mail *entities.Mail) *rest_err.RestErr
}
