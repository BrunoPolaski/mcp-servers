package services

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/registration-validation/internal/core/entities"
	"github.com/BrunoPolaski/registration-validation/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/registration-validation/internal/infra/repositories"
	"github.com/BrunoPolaski/registration-validation/internal/infra/repositories/interfaces"
)

type AddressService struct {
	addressRepository interfaces.AddressRepository
}

func NewAddressService(rf *repositories.RepositoryFactory) *AddressService {
	return &AddressService{
		addressRepository: rf.AddressRepository(),
	}
}

func (as *AddressService) Create(ctx context.Context, addr *dto.AddressDTO) (*entities.Address, *rest_err.RestErr) {
	domain := addr.ToEntity()

	err := as.addressRepository.Create(ctx, domain)
	if err != nil {
		return nil, err
	}

	return domain, nil
}

func (as *AddressService) GetById(ctx context.Context, id uint) (*entities.Address, *rest_err.RestErr) {
	address, err := as.addressRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func (as *AddressService) GetAll(ctx context.Context, limit, offset int, params map[string]any) (*dto.PaginatedResponse[dto.AddressDTO], *rest_err.RestErr) {
	addresses, count, err := as.addressRepository.GetAll(ctx, limit, offset, params)
	if err != nil {
		return nil, err
	}

	paginatedAddresses := dto.NewPaginatedResponse(count, make([]*dto.AddressDTO, len(addresses)))

	for i, addr := range addresses {
		paginatedAddresses.Items[i] = dto.NewAddressDTO(&addr)
	}

	return paginatedAddresses, nil
}

func (as *AddressService) Delete(ctx context.Context, id uint) *rest_err.RestErr {
	return as.addressRepository.Delete(ctx, id)
}
