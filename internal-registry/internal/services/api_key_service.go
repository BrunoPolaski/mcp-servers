package services

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
	"github.com/BrunoPolaski/internal-registry/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/internal-registry/internal/infra/repositories"
	"github.com/BrunoPolaski/internal-registry/internal/infra/repositories/interfaces"
	"github.com/google/uuid"
)

type ApiKeyService struct {
	apiKeyRepository interfaces.ApiKeyRepository
}

func NewApiKeyService(rf *repositories.RepositoryFactory) *ApiKeyService {
	return &ApiKeyService{
		apiKeyRepository: rf.ApiKeyRepository(),
	}
}

func (us *ApiKeyService) Create(ctx context.Context, req *dto.ApiKeyDTO) (*entities.ApiKey, *rest_err.RestErr) {
	existing, _ := us.apiKeyRepository.GetById(ctx, req.UUID)
	if existing != nil {
		return nil, rest_err.NewBadRequestError("api key already exists")
	}

	uuid := uuid.NewString()

	apiKey := &entities.ApiKey{
		UUID: uuid,
		Slug: req.Slug,
	}

	err := us.apiKeyRepository.Create(ctx, apiKey)
	if err != nil {
		return nil, err
	}

	return apiKey, nil
}

func (us *ApiKeyService) GetById(ctx context.Context, uuid string) (*entities.ApiKey, *rest_err.RestErr) {
	user, err := us.apiKeyRepository.GetById(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *ApiKeyService) GetAll(ctx context.Context, limit, offset int, filters map[string]any) (*dto.PaginatedResponse[dto.ApiKeyDTO], *rest_err.RestErr) {
	apiKeys, count, err := us.apiKeyRepository.GetAll(ctx, limit, offset, filters)
	if err != nil {
		return nil, err
	}

	var paginated dto.PaginatedResponse[dto.ApiKeyDTO]
	for _, apiKey := range apiKeys {
		paginated.Items = append(paginated.Items, dto.NewApiKeyDTO(&apiKey))
	}
	paginated.Total = int64(count)

	return &paginated, nil
}

func (us *ApiKeyService) Delete(ctx context.Context, uuid string) *rest_err.RestErr {
	return us.apiKeyRepository.Delete(ctx, uuid)
}
