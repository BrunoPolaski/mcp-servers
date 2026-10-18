package services

import (
	"context"
	"testing"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/registration-validation/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/registration-validation/internal/core/entities/value_objects"
	"gorm.io/gorm"
)

type fakePersonRepository struct {
	byID       map[uint]*entities.Person
	byDocument map[string]*entities.Person
}

func (f *fakePersonRepository) GetById(_ context.Context, id uint) (*entities.Person, *rest_err.RestErr) {
	if p, ok := f.byID[id]; ok {
		return p, nil
	}
	return nil, rest_err.NewNotFoundError("person not found")
}

func (f *fakePersonRepository) GetByDocument(_ context.Context, doc string) (*entities.Person, *rest_err.RestErr) {
	if p, ok := f.byDocument[doc]; ok {
		return p, nil
	}
	return nil, rest_err.NewNotFoundError("person not found")
}

func (f *fakePersonRepository) GetAll(_ context.Context, _, _ int, _ map[string]any) ([]entities.Person, int64, *rest_err.RestErr) {
	return nil, 0, nil
}

func (f *fakePersonRepository) Delete(_ context.Context, _ uint) *rest_err.RestErr { return nil }

type fakeDocumentValidationRepository struct {
	gotPersonID uint
	gotType     string
	items       []entities.DocumentValidation
}

func (f *fakeDocumentValidationRepository) GetByPersonID(_ context.Context, personID uint, documentType string) ([]entities.DocumentValidation, *rest_err.RestErr) {
	f.gotPersonID, f.gotType = personID, documentType
	return f.items, nil
}

type fakeFiscalRegularityRepository struct {
	items []entities.FiscalRegularity
}

func (f *fakeFiscalRegularityRepository) GetByPersonID(_ context.Context, _ uint) ([]entities.FiscalRegularity, *rest_err.RestErr) {
	return f.items, nil
}

type fakeComplianceCheckRepository struct {
	gotType string
	items   []entities.ComplianceCheck
	err     *rest_err.RestErr
}

func (f *fakeComplianceCheckRepository) GetByPersonID(_ context.Context, _ uint, checkType string) ([]entities.ComplianceCheck, *rest_err.RestErr) {
	f.gotType = checkType
	return f.items, f.err
}

func person() *entities.Person {
	return &entities.Person{
		Model:               gorm.Model{ID: 12},
		PersonalInformation: &entities.PersonalInformation{FullName: "Rafael Moreira", Document: valueobjects.Document("35509139404")},
	}
}

func repoWith(p *entities.Person) *fakePersonRepository {
	return &fakePersonRepository{
		byID:       map[uint]*entities.Person{p.ID: p},
		byDocument: map[string]*entities.Person{p.PersonalInformation.Document.String(): p},
	}
}

func TestResolveCustomer(t *testing.T) {
	p := person()
	s := &RegistrationValidationService{personRepository: repoWith(p)}

	tests := []struct {
		name    string
		ref     CustomerRef
		wantErr bool
	}{
		{name: "by id", ref: CustomerRef{ID: 12}},
		{name: "by document", ref: CustomerRef{Document: "35509139404"}},
		{name: "both", ref: CustomerRef{ID: 12, Document: "35509139404"}, wantErr: true},
		{name: "neither", ref: CustomerRef{}, wantErr: true},
		{name: "unknown id", ref: CustomerRef{ID: 99}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.resolveCustomer(context.Background(), tt.ref)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.ID != 12 {
				t.Errorf("resolved person %d, want 12", got.ID)
			}
		})
	}
}

func TestGetDocumentValidationsResolvesDocumentAndFilters(t *testing.T) {
	p := person()
	docRepo := &fakeDocumentValidationRepository{items: []entities.DocumentValidation{
		{DocumentType: "cpf", ReceitaFederalStatus: "regular", IsValid: true},
		{DocumentType: "cpf", ReceitaFederalStatus: "suspended"},
	}}
	s := &RegistrationValidationService{personRepository: repoWith(p), documentValidationRepository: docRepo}

	got, err := s.GetDocumentValidations(context.Background(), CustomerRef{Document: "35509139404"}, "cpf")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if docRepo.gotPersonID != 12 || docRepo.gotType != "cpf" {
		t.Errorf("repository called with (%d, %q), want (12, \"cpf\")", docRepo.gotPersonID, docRepo.gotType)
	}
	if got.CustomerID != 12 || got.Document != "35509139404" {
		t.Errorf("wrong customer identification: %+v", got)
	}
	if len(got.Items) != 2 || got.Items[1].ReceitaFederalStatus != "suspended" {
		t.Errorf("items not mapped in order: %+v", got.Items)
	}
}

func TestNeverValidatedCustomerReturnsEmptyItems(t *testing.T) {
	p := person()
	s := &RegistrationValidationService{personRepository: repoWith(p), fiscalRegularityRepository: &fakeFiscalRegularityRepository{}}

	got, err := s.GetFiscalRegularities(context.Background(), CustomerRef{ID: 12})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Items == nil || len(got.Items) != 0 {
		t.Errorf("items must be an empty list, not nil, so the agent sees no records: %#v", got.Items)
	}
}

func TestGetComplianceChecksPropagatesRepositoryError(t *testing.T) {
	p := person()
	s := &RegistrationValidationService{
		personRepository:          repoWith(p),
		complianceCheckRepository: &fakeComplianceCheckRepository{err: rest_err.NewInternalServerError("db down")},
	}

	if _, err := s.GetComplianceChecks(context.Background(), CustomerRef{ID: 12}, "pep_screening"); err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestDocumentOfWithoutPersonalInformation(t *testing.T) {
	if got := documentOf(&entities.Person{}); got != "" {
		t.Errorf("documentOf = %q, want empty", got)
	}
}
