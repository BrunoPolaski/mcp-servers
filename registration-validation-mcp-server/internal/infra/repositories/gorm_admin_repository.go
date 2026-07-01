package repositories

import (
	"context"
	"errors"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/registration-validation-mcp-server/internal/core/entities"
	"github.com/BrunoPolaski/registration-validation-mcp-server/internal/infra/repositories/interfaces"
	"gorm.io/gorm"
)

type gormAdminRepository struct {
	db *gorm.DB
}

func NewGormAdminRepository(db *gorm.DB) interfaces.AdminRepository {
	return &gormAdminRepository{
		db: db,
	}
}

func (g *gormAdminRepository) GetById(ctx context.Context, id uint) (*entities.Admin, *rest_err.RestErr) {
	res, err := gorm.G[entities.Admin](g.db).
		Preload("PersonalInformation", nil).
		Where("id = ?", id).
		First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, rest_err.NewNotFoundError("admin not found")
		}
		return nil, rest_err.NewInternalServerError("error while fetching admin").WithCause(err)
	}
	return &res, nil
}

func (g *gormAdminRepository) GetAll(ctx context.Context, limit, offset int, params map[string]any) ([]entities.Admin, int64, *rest_err.RestErr) {
	total, err := gorm.G[entities.Admin](g.db).Where(params).Count(ctx, "id")
	if err != nil {
		return nil, 0, rest_err.NewInternalServerError("error while counting admins").WithCause(err)
	}

	admins, err := gorm.G[entities.Admin](g.db).Where(params).Limit(limit).Offset(offset).Find(ctx)
	if len(admins) == 0 {
		return nil, 0, rest_err.NewNotFoundError("no admins found")
	} else if err != nil {
		return nil, 0, rest_err.NewInternalServerError("error while fetching admins").WithCause(err)
	}
	return admins, total, nil
}

func (g *gormAdminRepository) Delete(ctx context.Context, id uint) *rest_err.RestErr {
	affected, err := gorm.G[entities.Admin](g.db).Where("id = ?", id).Delete(ctx)
	if affected == 0 {
		return rest_err.NewNotFoundError("admin not found")
	} else if err != nil {
		return rest_err.NewInternalServerError("error while deleting admin").WithCause(err)
	}
	return nil
}
