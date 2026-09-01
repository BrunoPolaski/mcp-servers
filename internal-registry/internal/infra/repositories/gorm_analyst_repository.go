package repositories

import (
	"context"
	"errors"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
	"github.com/BrunoPolaski/internal-registry/internal/infra/repositories/interfaces"
	"gorm.io/gorm"
)

type gormAnalystRepository struct {
	db *gorm.DB
}

func NewGormAnalystRepository(db *gorm.DB) interfaces.AnalystRepository {
	return &gormAnalystRepository{
		db: db,
	}
}

func (g *gormAnalystRepository) GetById(ctx context.Context, id uint) (*entities.Analyst, *rest_err.RestErr) {
	res, err := gorm.G[entities.Analyst](g.db).
		Preload("PersonalInformation", nil).
		Preload("Businesses", nil).
		Where("id = ?", id).
		First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, rest_err.NewNotFoundError("analyst not found")
		}
		return nil, rest_err.NewInternalServerError("error while fetching analyst").WithCause(err)
	}
	return &res, nil
}

func (g *gormAnalystRepository) GetAll(ctx context.Context, limit, offset int, params map[string]any) ([]entities.Analyst, int64, *rest_err.RestErr) {
	total, err := gorm.G[entities.Analyst](g.db).Where(params).Count(ctx, "id")
	if err != nil {
		return nil, 0, rest_err.NewInternalServerError("error while counting analysts").WithCause(err)
	}

	analysts, err := gorm.G[entities.Analyst](g.db).Where(params).Limit(limit).Offset(offset).Find(ctx)
	if len(analysts) == 0 {
		return nil, 0, rest_err.NewNotFoundError("no analysts found")
	} else if err != nil {
		return nil, 0, rest_err.NewInternalServerError("error while fetching analysts").WithCause(err)
	}
	return analysts, total, nil
}

func (g *gormAnalystRepository) Delete(ctx context.Context, id uint) *rest_err.RestErr {
	affected, err := gorm.G[entities.Analyst](g.db).Where("id = ?", id).Delete(ctx)
	if affected == 0 {
		return rest_err.NewNotFoundError("analyst not found")
	} else if err != nil {
		return rest_err.NewInternalServerError("error while deleting analyst").WithCause(err)
	}
	return nil
}
