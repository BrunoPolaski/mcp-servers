package repositories

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
	"github.com/BrunoPolaski/internal-registry/internal/infra/repositories/interfaces"
	"gorm.io/gorm"
)

type gormInternalPaymentRecordRepository struct {
	db *gorm.DB
}

func NewGormInternalPaymentRecordRepository(db *gorm.DB) interfaces.InternalPaymentRecordRepository {
	return &gormInternalPaymentRecordRepository{db: db}
}

func (g *gormInternalPaymentRecordRepository) GetByPersonID(ctx context.Context, personID uint, status string, productID *uint) ([]entities.InternalPaymentRecord, *rest_err.RestErr) {
	query := gorm.G[entities.InternalPaymentRecord](g.db).Where("person_id = ?", personID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if productID != nil {
		query = query.Where("contracted_product_id = ?", *productID)
	}

	records, err := query.Order("reference_month DESC").Find(ctx)
	if err != nil {
		return nil, rest_err.NewInternalServerError("error while fetching internal payment records").WithCause(err)
	}
	return records, nil
}
