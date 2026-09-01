package services

import (
	"context"

	"github.com/BrunoPolaski/bureau/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/bureau/internal/core/entities/value_objects"
	"github.com/BrunoPolaski/bureau/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/bureau/internal/infra/controllers/request"
	"github.com/BrunoPolaski/bureau/internal/infra/repositories"
	"github.com/BrunoPolaski/bureau/internal/infra/repositories/interfaces"
	"github.com/BrunoPolaski/go-rest-err/rest_err"
)

type AnalystService struct {
	analystRepository interfaces.AnalystRepository
	userRepository    interfaces.UserRepository
}

func NewAnalystService(
	rf *repositories.RepositoryFactory,
) *AnalystService {
	return &AnalystService{
		analystRepository: rf.AnalystRepository(),
		userRepository:    rf.UserRepository(),
	}
}

func (as *AnalystService) Create(ctx context.Context, personalInformation *request.PersonalInformationRequest) (*entities.Analyst, *rest_err.RestErr) {
	_, err := as.userRepository.GetByDocument(ctx, personalInformation.Document)
	if err == nil {
		return nil, rest_err.NewBadRequestError("document already exists")
	}

	uid, ok := ctx.Value("user_id").(uint)
	if !ok {
		return nil, rest_err.NewUnauthorizedError("invalid user id in token")
	}

	user, err := as.userRepository.GetById(ctx, uid)
	if err != nil {
		return nil, err
	}

	if user.UserType != valueobjects.UserTypeAnalyst {
		return nil, rest_err.NewBadRequestError("user is not an analyst")
	}

	if user.Analyst != nil {
		return nil, rest_err.NewBadRequestError("user already has an analyst profile")
	}

	user.Analyst.PersonalInformation = personalInformation.ToEntity()

	user, err = as.userRepository.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return user.Analyst, nil
}

func (as *AnalystService) GetById(ctx context.Context, id uint) (*entities.Analyst, *rest_err.RestErr) {
	analyst, err := as.analystRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return analyst, nil
}

func (as *AnalystService) GetAll(ctx context.Context, limit, offset int, params map[string]any) (*dto.PaginatedResponse[dto.AnalystDTO], *rest_err.RestErr) {
	analysts, count, err := as.analystRepository.GetAll(ctx, limit, offset, params)
	if err != nil {
		return nil, err
	}

	paginatedAnalysts := dto.NewPaginatedResponse(count, make([]*dto.AnalystDTO, len(analysts)))

	for i, analyst := range analysts {
		paginatedAnalysts.Items[i] = dto.NewAnalystDTO(&analyst)
	}

	return paginatedAnalysts, nil
}

func (as *AnalystService) Delete(ctx context.Context, id uint) *rest_err.RestErr {
	return as.analystRepository.Delete(ctx, id)
}
