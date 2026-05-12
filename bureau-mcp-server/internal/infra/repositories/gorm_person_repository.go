package repositories

import (
	"context"
	"errors"

	"github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities"
	"github.com/BrunoPolaski/bureau-mcp-server/internal/infra/repositories/interfaces"
	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"gorm.io/gorm"
)

type gormPersonRepository struct {
	db *gorm.DB
}

func NewGormPersonRepository(db *gorm.DB) interfaces.PersonRepository {
	return &gormPersonRepository{
		db: db,
	}
}

func (g *gormPersonRepository) GetById(ctx context.Context, id uint) (*entities.Person, *rest_err.RestErr) {
	res, err := gorm.G[*entities.Person](g.db).
		Preload("PersonalInformation", nil).
		Where("id = ?", id).
		First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, rest_err.NewNotFoundError("person not found")
		}
		return nil, rest_err.NewInternalServerError("error while fetching person")
	}
	return res, nil
}

func (g *gormPersonRepository) GetAll(ctx context.Context, limit, offset int, params map[string]any) ([]entities.Person, int64, *rest_err.RestErr) {
	total, err := gorm.G[entities.Person](g.db).Where(params).Count(ctx, "id")
	if err != nil {
		return nil, 0, rest_err.NewInternalServerError("error while counting persons")
	}

	persons, err := gorm.G[entities.Person](g.db).Where(params).Limit(limit).Offset(offset).Find(ctx)
	if len(persons) == 0 {
		return nil, 0, rest_err.NewNotFoundError("no persons found")
	} else if err != nil {
		return nil, 0, rest_err.NewInternalServerError("error while fetching persons")
	}
	return persons, total, nil
}

func (g *gormPersonRepository) Delete(ctx context.Context, id uint) *rest_err.RestErr {
	affected, err := gorm.G[entities.Person](g.db).Where("id = ?", id).Delete(ctx)
	if affected == 0 {
		return rest_err.NewNotFoundError("person not found")
	} else if err != nil {
		return rest_err.NewInternalServerError("error while deleting person")
	}
	return nil
}
