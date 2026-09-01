package services

import (
	"context"
	"testing"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/internal-registry/internal/core/entities/value_objects"
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

type fakeContractedProductRepository struct {
	gotType, gotStatus string
	products           []entities.ContractedProduct
}

func (f *fakeContractedProductRepository) GetByPersonID(_ context.Context, _ uint, productType, status string) ([]entities.ContractedProduct, *rest_err.RestErr) {
	f.gotType, f.gotStatus = productType, status
	return f.products, nil
}

func personFixture() *entities.Person {
	return &entities.Person{
		Model:               gorm.Model{ID: 5},
		PersonalInformation: &entities.PersonalInformation{FullName: "Gabriela Ribeiro Barbosa", Document: valueobjects.Document("35509139404")},
	}
}

func TestResolveCustomer_Ambiguous(t *testing.T) {
	s := &InternalRegistryService{personRepository: &fakePersonRepository{}}
	_, err := s.resolveCustomer(context.Background(), CustomerRef{ID: 1, Document: "x"})
	if err == nil {
		t.Fatal("esperado erro para customer_id e document simultâneos")
	}
}

func TestResolveCustomer_Empty(t *testing.T) {
	s := &InternalRegistryService{personRepository: &fakePersonRepository{}}
	_, err := s.resolveCustomer(context.Background(), CustomerRef{})
	if err == nil {
		t.Fatal("esperado erro quando nenhum identificador é informado")
	}
}

func TestGetContractedProducts_ResolvesDocumentAndFilters(t *testing.T) {
	person := personFixture()
	prodRepo := &fakeContractedProductRepository{products: []entities.ContractedProduct{{ProductType: "credit_card", Status: "active"}}}
	s := &InternalRegistryService{
		personRepository:            &fakePersonRepository{byDocument: map[string]*entities.Person{"35509139404": person}},
		contractedProductRepository: prodRepo,
	}

	res, err := s.GetContractedProducts(context.Background(), CustomerRef{Document: "35509139404"}, "credit_card", "active")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if res.CustomerID != 5 || res.Document != "35509139404" {
		t.Fatalf("identificação errada no resultado: %+v", res)
	}
	if prodRepo.gotType != "credit_card" || prodRepo.gotStatus != "active" {
		t.Fatalf("filtros não repassados: %+v", prodRepo)
	}
	if len(res.Items) != 1 {
		t.Fatalf("esperado 1 item, veio %d", len(res.Items))
	}
}
