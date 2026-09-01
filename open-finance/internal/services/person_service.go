package services

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/open-finance/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/open-finance/internal/core/entities/value_objects"
	"github.com/BrunoPolaski/open-finance/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/open-finance/internal/infra/controllers/request"
	"github.com/BrunoPolaski/open-finance/internal/infra/repositories"
	"github.com/BrunoPolaski/open-finance/internal/infra/repositories/interfaces"
)

type PersonService struct {
	personRepository interfaces.PersonRepository
	userRepository   interfaces.UserRepository
}

func NewPersonService(rf *repositories.RepositoryFactory) *PersonService {
	return &PersonService{
		personRepository: rf.PersonRepository(),
		userRepository:   rf.UserRepository(),
	}
}

func (as *PersonService) Create(ctx context.Context, personalInformation *request.PersonalInformationRequest) (*entities.Person, *rest_err.RestErr) {
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

	if user.UserType != valueobjects.UserTypePerson {
		return nil, rest_err.NewBadRequestError("user is not a person")
	}

	if user.Person != nil {
		return nil, rest_err.NewBadRequestError("user already has a person profile")
	}

	user.Person.PersonalInformation = personalInformation.ToEntity()

	user, err = as.userRepository.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return user.Person, nil
}

func (as *PersonService) GetById(ctx context.Context, id uint) (*entities.Person, *rest_err.RestErr) {
	person, err := as.personRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return person, nil
}

func (as *PersonService) GetByDocument(ctx context.Context, document string) (*entities.Person, *rest_err.RestErr) {
	person, err := as.personRepository.GetByDocument(ctx, document)
	if err != nil {
		return nil, err
	}

	return person, nil
}

func (as *PersonService) GetAll(ctx context.Context, limit, offset int, params map[string]any) (*dto.PaginatedResponse[dto.PersonDTO], *rest_err.RestErr) {
	persons, count, err := as.personRepository.GetAll(ctx, limit, offset, params)
	if err != nil {
		return nil, err
	}

	paginatedPersones := dto.NewPaginatedResponse(count, make([]*dto.PersonDTO, len(persons)))

	for i, person := range persons {
		paginatedPersones.Items[i] = dto.NewPersonDTO(&person)
	}

	return paginatedPersones, nil
}

func (as *PersonService) Delete(ctx context.Context, id uint) *rest_err.RestErr {
	return as.personRepository.Delete(ctx, id)
}
