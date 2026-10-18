package services

import (
	"context"
	"testing"

	"github.com/BrunoPolaski/bureau/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/bureau/internal/core/entities/value_objects"
	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"gorm.io/gorm"
)

type fakePersonRepository struct {
	persons   []entities.Person
	total     int64
	err       *rest_err.RestErr
	gotLimit  int
	gotOffset int
	gotParams map[string]any
}

func (f *fakePersonRepository) GetById(_ context.Context, id uint) (*entities.Person, *rest_err.RestErr) {
	for i := range f.persons {
		if f.persons[i].ID == id {
			return &f.persons[i], nil
		}
	}
	return nil, rest_err.NewNotFoundError("person not found")
}

func (f *fakePersonRepository) GetByDocument(_ context.Context, document string) (*entities.Person, *rest_err.RestErr) {
	for i := range f.persons {
		if pi := f.persons[i].PersonalInformation; pi != nil && pi.Document.String() == document {
			return &f.persons[i], nil
		}
	}
	return nil, rest_err.NewNotFoundError("person not found")
}

func (f *fakePersonRepository) GetAll(_ context.Context, limit, offset int, params map[string]any) ([]entities.Person, int64, *rest_err.RestErr) {
	f.gotLimit, f.gotOffset, f.gotParams = limit, offset, params
	return f.persons, f.total, f.err
}

func (f *fakePersonRepository) Delete(_ context.Context, _ uint) *rest_err.RestErr { return nil }

func persons() []entities.Person {
	return []entities.Person{
		{Model: gorm.Model{ID: 1}, PersonalInformation: &entities.PersonalInformation{FullName: "Ana Souza", Document: valueobjects.Document("11122233344")}},
		{Model: gorm.Model{ID: 2}, PersonalInformation: &entities.PersonalInformation{FullName: "Bruno Lima", Document: valueobjects.Document("55566677788")}},
	}
}

func TestGetAllSummaryMapsEveryPersonAndKeepsTotal(t *testing.T) {
	repo := &fakePersonRepository{persons: persons(), total: 30}
	s := &PersonService{personRepository: repo}
	params := map[string]any{"name": "Ana"}

	got, err := s.GetAllSummary(context.Background(), 2, 4, params)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.gotLimit != 2 || repo.gotOffset != 4 || repo.gotParams["name"] != "Ana" {
		t.Errorf("pagination or filters not passed: %+v", repo)
	}
	if got.Total != 30 || len(got.Items) != 2 {
		t.Fatalf("response = (total %d, %d items), want (30, 2)", got.Total, len(got.Items))
	}
	if got.Items[0].ID != 1 || got.Items[1].Name != "Bruno Lima" || got.Items[1].Document != "55566677788" {
		t.Errorf("items mapped in the wrong order or with wrong fields: %+v, %+v", got.Items[0], got.Items[1])
	}
}

func TestGetAllSummaryPropagatesRepositoryError(t *testing.T) {
	s := &PersonService{personRepository: &fakePersonRepository{err: rest_err.NewInternalServerError("db down")}}

	if _, err := s.GetAllSummary(context.Background(), 10, 0, nil); err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestGetByDocumentNotFound(t *testing.T) {
	s := &PersonService{personRepository: &fakePersonRepository{persons: persons()}}

	if _, err := s.GetByDocument(context.Background(), "00000000000"); err == nil {
		t.Fatal("expected not found, got nil")
	}
}
