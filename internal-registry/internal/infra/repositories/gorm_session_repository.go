package repositories

import (
	"context"
	"errors"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
	"github.com/BrunoPolaski/internal-registry/internal/infra/repositories/interfaces"
	"gorm.io/gorm"
)

type gormSessionRepository struct {
	db *gorm.DB
}

func NewGormSessionRepository(db *gorm.DB) interfaces.SessionRepository {
	return &gormSessionRepository{
		db: db,
	}
}

func (g *gormSessionRepository) Create(ctx context.Context, t *entities.Token) *rest_err.RestErr {
	err := gorm.G[entities.Token](g.db).Create(ctx, t)
	if err != nil {
		return rest_err.NewInternalServerError("error while creating token").WithCause(err)
	}
	return nil
}

func (g *gormSessionRepository) GetById(ctx context.Context, id string) (*entities.Token, *rest_err.RestErr) {
	res, err := gorm.G[entities.Token](g.db).Where("id = ?", id).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, rest_err.NewNotFoundError("token not found")
		}
		return nil, rest_err.NewInternalServerError("error while fetching token").WithCause(err)
	}
	return &res, nil
}

func (g *gormSessionRepository) GetAll(ctx context.Context, limit, offset int, params map[string]any) ([]entities.Token, int, *rest_err.RestErr) {
	tokens, err := gorm.G[entities.Token](g.db).Where(params).Limit(limit).Offset(offset).Find(ctx)
	if len(tokens) == 0 {
		return nil, 0, rest_err.NewNotFoundError("no tokens found")
	} else if err != nil {
		return nil, 0, rest_err.NewInternalServerError("error while fetching tokens").WithCause(err)
	}
	return tokens, len(tokens), nil
}

func (g *gormSessionRepository) Delete(ctx context.Context, id string) *rest_err.RestErr {
	affected, err := gorm.G[entities.Token](g.db).Where("id = ?", id).Delete(ctx)
	if affected == 0 {
		return rest_err.NewNotFoundError("token not found")
	} else if err != nil {
		return rest_err.NewInternalServerError("error while deleting token").WithCause(err)
	}
	return nil
}
