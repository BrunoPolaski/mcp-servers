package repositories

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/open-finance/internal/core/entities"
	"github.com/BrunoPolaski/open-finance/internal/infra/repositories/interfaces"
	"gorm.io/gorm"
)

type gormDataSharingConsentRepository struct {
	db *gorm.DB
}

func NewGormDataSharingConsentRepository(db *gorm.DB) interfaces.DataSharingConsentRepository {
	return &gormDataSharingConsentRepository{db: db}
}

func (g *gormDataSharingConsentRepository) GetByPersonID(ctx context.Context, personID uint) ([]entities.DataSharingConsent, *rest_err.RestErr) {
	consents, err := gorm.G[entities.DataSharingConsent](g.db).
		Where("person_id = ?", personID).
		Order("granted_at DESC").
		Find(ctx)
	if err != nil {
		return nil, rest_err.NewInternalServerError("error while fetching data sharing consents").WithCause(err)
	}

	return consents, nil
}
