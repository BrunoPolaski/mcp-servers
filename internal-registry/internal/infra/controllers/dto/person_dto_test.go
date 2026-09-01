package dto

import (
	"testing"

	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/internal-registry/internal/core/entities/value_objects"
	"gorm.io/gorm"
)

func TestNewPersonDTO_NilAssociationsNoPanic(t *testing.T) {
	// Pessoa sem nenhuma dimensão interna (não-cliente): não pode dar panic e as
	// listas ficam vazias/omitidas.
	p := &entities.Person{
		Model:               gorm.Model{ID: 6},
		PersonalInformation: &entities.PersonalInformation{FullName: "Henrique Almeida Ribeiro", Document: valueobjects.Document("99988877766")},
	}
	got := NewPersonDTO(p)
	if got.ID != 6 || got.PersonalInformation == nil {
		t.Fatalf("DTO base inesperado: %+v", got)
	}
	if got.CustomerRelationship != nil {
		t.Errorf("esperado relationship nil para não-cliente")
	}
	if len(got.ContractedProducts) != 0 || len(got.InternalPaymentRecords) != 0 ||
		len(got.PreApprovedLimits) != 0 || len(got.IncomeDeclarations) != 0 {
		t.Errorf("esperado dimensões vazias, veio %+v", got)
	}
}

func TestNewPersonDTO_FullAssociations(t *testing.T) {
	score := 740
	p := &entities.Person{
		Model:                gorm.Model{ID: 2},
		PersonalInformation:  &entities.PersonalInformation{FullName: "Henrique Martins Barbosa", Document: valueobjects.Document("22233344455")},
		CustomerRelationship: &entities.CustomerRelationship{RelationshipMonths: 96, Segment: "retail", InternalScore: &score},
		ContractedProducts:   []entities.ContractedProduct{{ProductType: "credit_card", ProductName: "Cartão Gold", Status: "active"}},
		InternalPaymentRecords: []entities.InternalPaymentRecord{{Status: "on_time", AmountDue: 500, AmountPaid: 500}},
		PreApprovedLimits:    []entities.PreApprovedLimit{{ProductType: "personal_loan", ApprovedAmount: 20000, IsActive: true}},
	}
	got := NewPersonDTO(p)
	if got.CustomerRelationship == nil || got.CustomerRelationship.RelationshipMonths != 96 {
		t.Fatalf("relationship não mapeado: %+v", got.CustomerRelationship)
	}
	if len(got.ContractedProducts) != 1 || got.ContractedProducts[0].ProductType != "credit_card" {
		t.Fatalf("produtos não mapeados: %+v", got.ContractedProducts)
	}
	if len(got.PreApprovedLimits) != 1 || got.PreApprovedLimits[0].ApprovedAmount != 20000 {
		t.Fatalf("limites não mapeados: %+v", got.PreApprovedLimits)
	}
}
