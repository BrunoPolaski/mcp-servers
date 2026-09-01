package repositories

import (
	"context"
	"errors"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/open-finance/internal/core/entities"
	"github.com/BrunoPolaski/open-finance/internal/infra/repositories/interfaces"
	"gorm.io/gorm"
)

type gormApiKeyRepository struct {
	db *gorm.DB
}

func NewGormApiKeyRepository(db *gorm.DB) interfaces.ApiKeyRepository {
	return &gormApiKeyRepository{
		db: db,
	}
}

func (g *gormApiKeyRepository) Create(ctx context.Context, a *entities.ApiKey) *rest_err.RestErr {
	err := gorm.G[entities.ApiKey](g.db).Create(ctx, a)
	if err != nil {
		return rest_err.NewInternalServerError("error while creating api key").WithCause(err)
	}
	return nil
}

func (g *gormApiKeyRepository) GetById(ctx context.Context, uuid string) (*entities.ApiKey, *rest_err.RestErr) {
	res, err := gorm.G[entities.ApiKey](g.db).Where("uuid = ?", uuid).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, rest_err.NewNotFoundError("api key not found")
		}
		return nil, rest_err.NewInternalServerError("error while fetching api key").WithCause(err)
	}
	return &res, nil
}

func (g *gormApiKeyRepository) GetAll(ctx context.Context, limit, offset int, params map[string]any) ([]entities.ApiKey, int, *rest_err.RestErr) {
	apiKeys, err := gorm.G[entities.ApiKey](g.db).Where(params).Limit(limit).Offset(offset).Find(ctx)
	if len(apiKeys) == 0 {
		return nil, 0, rest_err.NewNotFoundError("no api keys found")
	} else if err != nil {
		return nil, 0, rest_err.NewInternalServerError("error while fetching api keys").WithCause(err)
	}
	return apiKeys, len(apiKeys), nil
}

func (g *gormApiKeyRepository) Delete(ctx context.Context, uuid string) *rest_err.RestErr {
	affected, err := gorm.G[entities.ApiKey](g.db).Where("id = ?", uuid).Delete(ctx)
	if affected == 0 {
		return rest_err.NewNotFoundError("api key not found")
	} else if err != nil {
		return rest_err.NewInternalServerError("error while deleting api key").WithCause(err)
	}
	return nil
}
