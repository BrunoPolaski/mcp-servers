package dto

import (
	"testing"

	"github.com/BrunoPolaski/internal-registry/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/internal-registry/internal/core/entities/value_objects"
	"gorm.io/gorm"
)

func TestNewPersonSummaryDTO(t *testing.T) {
	p := &entities.Person{
		Model: gorm.Model{ID: 7},
		PersonalInformation: &entities.PersonalInformation{
			FullName: "Felipe Pereira Santos",
			Document: valueobjects.Document("11122233344"),
		},
	}
	got := NewPersonSummaryDTO(p)
	if got.ID != 7 || got.Name != "Felipe Pereira Santos" || got.Document != "11122233344" {
		t.Fatalf("summary inesperado: %+v", got)
	}
}

func TestNewPersonSummaryDTO_NilPersonalInformation(t *testing.T) {
	got := NewPersonSummaryDTO(&entities.Person{Model: gorm.Model{ID: 1}})
	if got.ID != 1 || got.Name != "" || got.Document != "" {
		t.Fatalf("esperado apenas ID preenchido, veio %+v", got)
	}
}
