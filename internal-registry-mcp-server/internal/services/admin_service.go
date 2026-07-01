package services

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry-mcp-server/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/internal-registry-mcp-server/internal/core/entities/value_objects"
	"github.com/BrunoPolaski/internal-registry-mcp-server/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/internal-registry-mcp-server/internal/infra/controllers/request"
	"github.com/BrunoPolaski/internal-registry-mcp-server/internal/infra/repositories"
	"github.com/BrunoPolaski/internal-registry-mcp-server/internal/infra/repositories/interfaces"
)

type AdminService struct {
	adminRepository interfaces.AdminRepository
	userRepository  interfaces.UserRepository
}

func NewAdminService(rf *repositories.RepositoryFactory) *AdminService {
	return &AdminService{
		adminRepository: rf.AdminRepository(),
		userRepository:  rf.UserRepository(),
	}
}

func (as *AdminService) Create(ctx context.Context, personalInformation *request.PersonalInformationRequest) (*entities.Admin, *rest_err.RestErr) {
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

	if user.UserType != valueobjects.UserTypeAdmin {
		return nil, rest_err.NewBadRequestError("user is not an admin")
	}

	if user.Admin != nil {
		return nil, rest_err.NewBadRequestError("user already has an admin profile")
	}

	user.Admin.PersonalInformation = personalInformation.ToEntity()

	user, err = as.userRepository.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return user.Admin, nil
}

func (as *AdminService) GetById(ctx context.Context, id uint) (*entities.Admin, *rest_err.RestErr) {
	admin, err := as.adminRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return admin, nil
}

func (as *AdminService) GetAll(ctx context.Context, limit, offset int, params map[string]any) (*dto.PaginatedResponse[dto.AdminDTO], *rest_err.RestErr) {
	admins, count, err := as.adminRepository.GetAll(ctx, limit, offset, params)
	if err != nil {
		return nil, err
	}

	paginatedAdmins := dto.NewPaginatedResponse(count, make([]*dto.AdminDTO, len(admins)))

	for i, admin := range admins {
		paginatedAdmins.Items[i] = dto.NewAdminDTO(&admin)
	}

	return paginatedAdmins, nil
}

func (as *AdminService) Delete(ctx context.Context, id uint) *rest_err.RestErr {
	return as.adminRepository.Delete(ctx, id)
}
