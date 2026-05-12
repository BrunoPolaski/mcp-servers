package services

import (
	"context"

	"github.com/BrunoPolaski/bureau-mcp-server/internal/core/entities"
	"github.com/BrunoPolaski/bureau-mcp-server/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/bureau-mcp-server/internal/infra/repositories"
	"github.com/BrunoPolaski/bureau-mcp-server/internal/infra/repositories/interfaces"
	"github.com/BrunoPolaski/go-rest-err/rest_err"
)

type UserService struct {
	userRepository interfaces.UserRepository
}

func NewUserService(rf *repositories.RepositoryFactory) *UserService {
	return &UserService{
		userRepository: rf.UserRepository(),
	}
}

func (us *UserService) GetById(ctx context.Context, id uint) (*entities.User, *rest_err.RestErr) {
	user, err := us.userRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserService) GetAll(ctx context.Context, limit, offset int, params map[string]any) (*dto.PaginatedResponse[dto.UserDTO], *rest_err.RestErr) {
	users, count, err := us.userRepository.GetAll(ctx, limit, offset, params)
	if err != nil {
		return nil, err
	}

	paginated := dto.NewPaginatedResponse(
		count,
		make([]*dto.UserDTO, len(users)),
	)

	for i := range users {
		paginated.Items[i] = dto.NewUserDTO(&users[i])
	}

	return paginated, nil
}

func (us *UserService) Delete(ctx context.Context, id uint) *rest_err.RestErr {
	return us.userRepository.Delete(ctx, id)
}
