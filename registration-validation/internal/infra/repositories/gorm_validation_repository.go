package repositories

import (
	"context"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/registration-validation/internal/core/entities"
	"github.com/BrunoPolaski/registration-validation/internal/infra/repositories/interfaces"
	"gorm.io/gorm"
)

type gormDocumentValidationRepository struct{ db *gorm.DB }

func NewGormDocumentValidationRepository(db *gorm.DB) interfaces.DocumentValidationRepository {
	return &gormDocumentValidationRepository{db: db}
}

func (g *gormDocumentValidationRepository) GetByPersonID(ctx context.Context, personID uint, documentType string) ([]entities.DocumentValidation, *rest_err.RestErr) {
	query := gorm.G[entities.DocumentValidation](g.db).Where("person_id = ?", personID)
	if documentType != "" {
		query = query.Where("document_type = ?", documentType)
	}

	validations, err := query.Order("validation_date DESC").Find(ctx)
	if err != nil {
		return nil, rest_err.NewInternalServerError("error while fetching document validations").WithCause(err)
	}
	return validations, nil
}

type gormFiscalRegularityRepository struct{ db *gorm.DB }

func NewGormFiscalRegularityRepository(db *gorm.DB) interfaces.FiscalRegularityRepository {
	return &gormFiscalRegularityRepository{db: db}
}

func (g *gormFiscalRegularityRepository) GetByPersonID(ctx context.Context, personID uint) ([]entities.FiscalRegularity, *rest_err.RestErr) {
	regularities, err := gorm.G[entities.FiscalRegularity](g.db).
		Where("person_id = ?", personID).
		Order("check_date DESC").
		Find(ctx)
	if err != nil {
		return nil, rest_err.NewInternalServerError("error while fetching fiscal regularities").WithCause(err)
	}
	return regularities, nil
}

type gormEmploymentLinkValidationRepository struct{ db *gorm.DB }

func NewGormEmploymentLinkValidationRepository(db *gorm.DB) interfaces.EmploymentLinkValidationRepository {
	return &gormEmploymentLinkValidationRepository{db: db}
}

func (g *gormEmploymentLinkValidationRepository) GetByPersonID(ctx context.Context, personID uint, status string) ([]entities.EmploymentLinkValidation, *rest_err.RestErr) {
	query := gorm.G[entities.EmploymentLinkValidation](g.db).Where("person_id = ?", personID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	links, err := query.Order("start_date DESC").Find(ctx)
	if err != nil {
		return nil, rest_err.NewInternalServerError("error while fetching employment link validations").WithCause(err)
	}
	return links, nil
}

type gormComplianceCheckRepository struct{ db *gorm.DB }

func NewGormComplianceCheckRepository(db *gorm.DB) interfaces.ComplianceCheckRepository {
	return &gormComplianceCheckRepository{db: db}
}

func (g *gormComplianceCheckRepository) GetByPersonID(ctx context.Context, personID uint, checkType string) ([]entities.ComplianceCheck, *rest_err.RestErr) {
	query := gorm.G[entities.ComplianceCheck](g.db).Where("person_id = ?", personID)
	if checkType != "" {
		query = query.Where("check_type = ?", checkType)
	}

	checks, err := query.Order("check_date DESC").Find(ctx)
	if err != nil {
		return nil, rest_err.NewInternalServerError("error while fetching compliance checks").WithCause(err)
	}
	return checks, nil
}
