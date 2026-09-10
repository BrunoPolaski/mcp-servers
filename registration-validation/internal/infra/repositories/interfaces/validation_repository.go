package interfaces

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/registration-validation/internal/core/entities"
)

type DocumentValidationRepository interface {
	// GetByPersonID devolve as validações de documento da pessoa. documentType
	// vazio não filtra.
	GetByPersonID(ctx context.Context, personID uint, documentType string) ([]entities.DocumentValidation, *rest_err.RestErr)
}

type FiscalRegularityRepository interface {
	// GetByPersonID devolve as consultas de regularidade fiscal da pessoa.
	GetByPersonID(ctx context.Context, personID uint) ([]entities.FiscalRegularity, *rest_err.RestErr)
}

type EmploymentLinkValidationRepository interface {
	// GetByPersonID devolve os vínculos validados no eSocial. status vazio não
	// filtra.
	GetByPersonID(ctx context.Context, personID uint, status string) ([]entities.EmploymentLinkValidation, *rest_err.RestErr)
}

type ComplianceCheckRepository interface {
	// GetByPersonID devolve as verificações de compliance da pessoa. checkType
	// vazio não filtra.
	GetByPersonID(ctx context.Context, personID uint, checkType string) ([]entities.ComplianceCheck, *rest_err.RestErr)
}
