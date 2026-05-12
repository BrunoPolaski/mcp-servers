package repositories

import (
	"context"
	"errors"

	"github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities"
	"github.com/BrunoPolaski/bureau-mcp-server/internal/infra/repositories/interfaces"
	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"gorm.io/gorm"
)

type gormAddressRepository struct {
	db *gorm.DB
}

func NewGormAddressRepository(db *gorm.DB) interfaces.AddressRepository {
	return &gormAddressRepository{
		db: db,
	}
}

func (g *gormAddressRepository) Create(ctx context.Context, a *entities.Address) *rest_err.RestErr {
	err := gorm.G[entities.Address](g.db).Create(ctx, a)
	if err != nil {
		return rest_err.NewInternalServerError("error while creating address")
	}
	return nil
}

func (g *gormAddressRepository) GetById(ctx context.Context, id uint) (*entities.Address, *rest_err.RestErr) {
	res, err := gorm.G[entities.Address](g.db).Where("id = ?", id).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, rest_err.NewNotFoundError("address not found")
		}
		return nil, rest_err.NewInternalServerError("error while fetching address")
	}
	return &res, nil
}

func (g *gormAddressRepository) GetAll(ctx context.Context, limit, offset int, params map[string]any) ([]entities.Address, int64, *rest_err.RestErr) {
	total, err := gorm.G[entities.Address](g.db).Where(params).Count(ctx, "id")
	if err != nil {
		return nil, 0, rest_err.NewInternalServerError("error while counting addresses")
	}

	addresses, err := gorm.G[entities.Address](g.db).Where(params).Limit(limit).Offset(offset).Find(ctx)
	if len(addresses) == 0 {
		return nil, 0, rest_err.NewNotFoundError("no addresses found")
	} else if err != nil {
		return nil, 0, rest_err.NewInternalServerError("error while fetching addresses")
	}
	return addresses, total, nil
}

func (g *gormAddressRepository) Delete(ctx context.Context, id uint) *rest_err.RestErr {
	affected, err := gorm.G[entities.Address](g.db).Where("id = ?", id).Delete(ctx)
	if affected == 0 {
		return rest_err.NewNotFoundError("address not found")
	} else if err != nil {
		return rest_err.NewInternalServerError("error while deleting address")
	}
	return nil
}
