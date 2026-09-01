package mailer

import (
	"context"

	"github.com/BrunoPolaski/bureau/internal/core/entities"
	"github.com/BrunoPolaski/go-rest-err/rest_err"
)

type Mailer interface {
	Send(ctx context.Context, mail *entities.Mail) *rest_err.RestErr
}
