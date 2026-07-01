package dto

import (
	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/registration-validation-mcp-server/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/registration-validation-mcp-server/internal/core/entities/value_objects"
)

type UserDTO struct {
	ID                            uint        `json:"id"`
	CreatedAt                     string      `json:"created_at"`
	UpdatedAt                     string      `json:"updated_at"`
	Email                         string      `json:"email" validate:"required,email" example:"user@example.com"`
	Password                      string      `json:"password,omitempty" validate:"required,password" example:"P@ssw0rd!"`
	UserType                      string      `json:"user_type" validate:"required,user_type" example:"person | owner | admin | analyst"`
	AdditionalInformationRequired bool        `json:"additional_information_required" example:"false"`
	Person                        *PersonDTO  `json:"person"`
	Admin                         *AdminDTO   `json:"admin"`
	Analyst                       *AnalystDTO `json:"analyst"`
}

func NewUserDTO(entity *entities.User) *UserDTO {
	userDTO := &UserDTO{
		ID:                            entity.ID,
		CreatedAt:                     entity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:                     entity.UpdatedAt.Format("2006-01-02 15:04:05"),
		Email:                         entity.Email.String(),
		UserType:                      entity.UserType.String(),
		AdditionalInformationRequired: entity.AdditionalInformationRequired(),
	}
	if entity.Person != nil {
		userDTO.Person = NewPersonDTO(entity.Person)
	}
	if entity.Admin != nil {
		userDTO.Admin = NewAdminDTO(entity.Admin)
	}
	if entity.Analyst != nil {
		userDTO.Analyst = NewAnalystDTO(entity.Analyst)
	}
	return userDTO
}

func (u UserDTO) ToEntity() (*entities.User, *rest_err.RestErr) {
	email := valueobjects.NewEmail(u.Email)
	password := valueobjects.NewPassword(u.Password)
	userType := valueobjects.NewUserType(u.UserType)

	hashedPassword, err := password.Hash()
	if err != nil {
		return nil, err
	}

	e := &entities.User{
		Email:    email,
		Password: hashedPassword,
		UserType: userType,
	}

	if u.Admin != nil {
		e.Admin = u.Admin.ToEntity()
	}
	if u.Person != nil {
		e.Person = u.Person.ToEntity()
	}
	if u.Analyst != nil {
		e.Analyst = u.Analyst.ToEntity()
	}

	return e, nil
}
