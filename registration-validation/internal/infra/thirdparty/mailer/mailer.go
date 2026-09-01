package mailer

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/registration-validation/internal/core/entities"
)

type Mailer interface {
	Send(ctx context.Context, mail *entities.Mail) *rest_err.RestErr
}
