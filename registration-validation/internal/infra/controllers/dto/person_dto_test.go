package dto

import (
	"testing"

	"github.com/BrunoPolaski/registration-validation/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/registration-validation/internal/core/entities/value_objects"
	"gorm.io/gorm"
)

func TestNewPersonDTOWithoutValidations(t *testing.T) {
	p := &entities.Person{
		Model:               gorm.Model{ID: 12},
		PersonalInformation: &entities.PersonalInformation{FullName: "Rafael Moreira", Document: valueobjects.Document("35509139404")},
	}

	got := NewPersonDTO(p)

	if got.ID != 12 || got.PersonalInformation == nil {
		t.Fatalf("unexpected base DTO: %+v", got)
	}
	if len(got.DocumentValidations) != 0 || len(got.FiscalRegularities) != 0 ||
		len(got.EmploymentLinkValidations) != 0 || len(got.ComplianceChecks) != 0 {
		t.Errorf("validation lists must stay empty for a customer never validated: %+v", got)
	}
}

func TestNewPersonDTOMapsEveryValidation(t *testing.T) {
	p := &entities.Person{
		Model:                     gorm.Model{ID: 3},
		DocumentValidations:       []entities.DocumentValidation{{DocumentType: "cpf", ReceitaFederalStatus: "regular"}},
		FiscalRegularities:        []entities.FiscalRegularity{{CNDStatus: "irregular", HasDebts: true}},
		EmploymentLinkValidations: []entities.EmploymentLinkValidation{{EmployerName: "Acme Ltda", Status: "active"}},
		ComplianceChecks:          []entities.ComplianceCheck{{CheckType: "pep_screening", IsPEP: true}},
	}

	got := NewPersonDTO(p)

	if len(got.DocumentValidations) != 1 || got.DocumentValidations[0].ReceitaFederalStatus != "regular" {
		t.Errorf("document validations not mapped: %+v", got.DocumentValidations)
	}
	if len(got.FiscalRegularities) != 1 || got.FiscalRegularities[0].CNDStatus != "irregular" {
		t.Errorf("fiscal regularities not mapped: %+v", got.FiscalRegularities)
	}
	if len(got.EmploymentLinkValidations) != 1 || got.EmploymentLinkValidations[0].EmployerName != "Acme Ltda" {
		t.Errorf("employment links not mapped: %+v", got.EmploymentLinkValidations)
	}
	if len(got.ComplianceChecks) != 1 || !got.ComplianceChecks[0].IsPEP {
		t.Errorf("compliance checks not mapped: %+v", got.ComplianceChecks)
	}
}
