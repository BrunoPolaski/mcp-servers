package repositories

import (
	"context"
	"errors"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
	"github.com/BrunoPolaski/internal-registry/internal/infra/repositories/interfaces"
	"gorm.io/gorm"
)

type gormCustomerRelationshipRepository struct {
	db *gorm.DB
}

func NewGormCustomerRelationshipRepository(db *gorm.DB) interfaces.CustomerRelationshipRepository {
	return &gormCustomerRelationshipRepository{db: db}
}

func (g *gormCustomerRelationshipRepository) GetByPersonID(ctx context.Context, personID uint) (*entities.CustomerRelationship, *rest_err.RestErr) {
	rel, err := gorm.G[entities.CustomerRelationship](g.db).Where("person_id = ?", personID).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // não-cliente: ausência é dado válido
		}
		return nil, rest_err.NewInternalServerError("error while fetching customer relationship").WithCause(err)
	}
	return &rel, nil
}
