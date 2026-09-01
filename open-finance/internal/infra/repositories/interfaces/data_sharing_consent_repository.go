package interfaces

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/open-finance/internal/core/entities"
)

type DataSharingConsentRepository interface {
	GetByPersonID(ctx context.Context, personID uint) ([]entities.DataSharingConsent, *rest_err.RestErr)
}
