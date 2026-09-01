package repositories

import (
	"context"
	"errors"

	"github.com/BrunoPolaski/bureau/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/bureau/internal/core/entities/value_objects"
	"github.com/BrunoPolaski/bureau/internal/infra/repositories/interfaces"
	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type gormPersonRepository struct {
	db *gorm.DB
}

func NewGormPersonRepository(db *gorm.DB) interfaces.PersonRepository {
	return &gormPersonRepository{
		db: db,
	}
}

func (g *gormPersonRepository) GetById(ctx context.Context, id uint) (*entities.Person, *rest_err.RestErr) {
	res, err := gorm.G[*entities.Person](g.db).
		Preload("PersonalInformation", nil).
		Preload("CreditScore", nil).
		Preload("FinancialProfile", nil).
		Preload("EmploymentRecords", nil).
		Preload("CreditAccounts", nil).
		Preload("CreditInquiries", nil).
		Preload("PaymentHistories", nil).
		Preload("Debts", nil).
		Preload("NegativeRecords", nil).
		Preload("LegalRecords", nil).
		Preload("FraudAlerts", nil).
		Preload("RiskAssessments", nil).
		Preload("RelatedPersons", nil).
		Preload("DataSources", nil).
		Where("id = ?", id).
		First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, rest_err.NewNotFoundError("person not found")
		}
		return nil, rest_err.NewInternalServerError("error while fetching person").WithCause(err)
	}
	return res, nil
}

func (g *gormPersonRepository) GetByDocument(ctx context.Context, document string) (*entities.Person, *rest_err.RestErr) {
	res, err := gorm.G[*entities.Person](g.db).
		Joins(
			clause.JoinTarget{Association: "PersonalInformation"},
			func(db gorm.JoinBuilder, joinTable, curTable clause.Table) error {
				db.Where(&entities.PersonalInformation{Document: valueobjects.Document(document)})
				return nil
			},
		).
		Preload("CreditScore", nil).
		Preload("FinancialProfile", nil).
		Preload("EmploymentRecords", nil).
		Preload("CreditAccounts", nil).
		Preload("CreditInquiries", nil).
		Preload("PaymentHistories", nil).
		Preload("Debts", nil).
		Preload("NegativeRecords", nil).
		Preload("LegalRecords", nil).
		Preload("FraudAlerts", nil).
		Preload("RiskAssessments", nil).
		Preload("RelatedPersons", nil).
		Preload("DataSources", nil).
		First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, rest_err.NewNotFoundError("person not found")
		}
		return nil, rest_err.NewInternalServerError("error while fetching person").WithCause(err)
	}
	return res, nil
}

func (g *gormPersonRepository) GetAll(ctx context.Context, limit, offset int, params map[string]any) ([]entities.Person, int64, *rest_err.RestErr) {
	total, err := gorm.G[entities.Person](g.db).
		Where(params).
		Count(ctx, "id")
	if err != nil {
		return nil, 0, rest_err.NewInternalServerError("error while counting persons").WithCause(err)
	}

	query := gorm.G[entities.Person](g.db).
		Preload("PersonalInformation", nil).
		Where(params)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	persons, err := query.Find(ctx)
	if len(persons) == 0 {
		return nil, 0, rest_err.NewNotFoundError("no persons found")
	} else if err != nil {
		return nil, 0, rest_err.NewInternalServerError("error while fetching persons").WithCause(err)
	}
	return persons, total, nil
}

func (g *gormPersonRepository) Delete(ctx context.Context, id uint) *rest_err.RestErr {
	affected, err := gorm.G[entities.Person](g.db).Where("id = ?", id).Delete(ctx)
	if affected == 0 {
		return rest_err.NewNotFoundError("person not found")
	} else if err != nil {
		return rest_err.NewInternalServerError("error while deleting person").WithCause(err)
	}
	return nil
}
