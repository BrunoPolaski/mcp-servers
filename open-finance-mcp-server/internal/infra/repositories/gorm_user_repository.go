package repositories

import (
	"context"
	"errors"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/open-finance-mcp-server/internal/core/entities"
	"github.com/BrunoPolaski/open-finance-mcp-server/internal/infra/repositories/interfaces"
	"gorm.io/gorm"
)

type gormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &gormUserRepository{
		db: db,
	}
}

func (g *gormUserRepository) Register(ctx context.Context, user *entities.User) (*entities.User, *rest_err.RestErr) {
	err := gorm.G[entities.User](g.db).Create(ctx, user)
	if err != nil {
		return nil, rest_err.NewInternalServerError("error while creating user").WithCause(err)
	}
	return user, nil
}

func (g *gormUserRepository) GetById(ctx context.Context, id uint) (*entities.User, *rest_err.RestErr) {
	res, err := gorm.G[*entities.User](g.db).
		Where("id = ?", id).
		Preload("Admin.PersonalInformation", nil).
		Preload("Analyst.PersonalInformation", nil).
		Preload("Person.PersonalInformation", nil).
		First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, rest_err.NewNotFoundError("user not found")
		}
		return nil, rest_err.NewInternalServerError("user not found").WithCause(err)
	}
	return res, nil
}

func (g *gormUserRepository) GetByDocument(ctx context.Context, document string) (*entities.PersonalInformation, *rest_err.RestErr) {
	personalInformation, err := gorm.G[*entities.PersonalInformation](g.db).Where("document = ?", document).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, rest_err.NewNotFoundError("user not found")
		}
		return nil, rest_err.NewInternalServerError("internal server error while fetching user").WithCause(err)
	}
	return personalInformation, nil
}

func (g *gormUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, *rest_err.RestErr) {
	usr, err := gorm.G[*entities.User](g.db).
		Preload("Admin.PersonalInformation", nil).
		Preload("Analyst.PersonalInformation", nil).
		Preload("Person.PersonalInformation", nil).
		Where("email = ?", email).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, rest_err.NewNotFoundError("user not found")
		}
		return nil, rest_err.NewInternalServerError("internal server error while fetching user").WithCause(err)
	}
	return usr, nil
}

func (g *gormUserRepository) GetAll(ctx context.Context, limit, offset int, params map[string]interface{}) ([]entities.User, int64, *rest_err.RestErr) {
	count, err := gorm.G[entities.User](g.db).Where(params).Count(ctx, "id")
	if err != nil {
		return nil, 0, rest_err.NewInternalServerError("error while counting users").WithCause(err)
	}

	users, err := gorm.G[entities.User](g.db).Find(ctx)
	if len(users) == 0 {
		return nil, 0, rest_err.NewNotFoundError("no users found")
	} else if err != nil {
		return nil, 0, rest_err.NewInternalServerError("error while fetching users").WithCause(err)
	}
	return users, count, nil
}

func (g *gormUserRepository) Delete(ctx context.Context, id uint) *rest_err.RestErr {
	affected, err := gorm.G[entities.User](g.db).Where("id = ?", id).Delete(ctx)
	if affected == 0 {
		return rest_err.NewNotFoundError("user not found")
	} else if err != nil {
		return rest_err.NewInternalServerError("error while deleting user").WithCause(err)
	}
	return nil
}

func (g *gormUserRepository) Update(ctx context.Context, user *entities.User) (*entities.User, *rest_err.RestErr) {
	affected, err := gorm.G[*entities.User](g.db).
		Where("id = ?", user.ID).
		Updates(ctx, user)

	if err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return nil, rest_err.NewBadRequestError("document already exists")
		}
		if errors.Is(err, gorm.ErrCheckConstraintViolated) {
			return nil, rest_err.NewBadRequestError("email already exists")
		}
		return nil, rest_err.NewInternalServerError("error while updating user").WithCause(err)
	}
	if affected == 0 {
		return nil, rest_err.NewNotFoundError("user not found")
	}
	return user, nil
}
