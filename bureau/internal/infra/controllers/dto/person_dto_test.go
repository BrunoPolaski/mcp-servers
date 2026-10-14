package dto

import (
	"testing"

	"github.com/BrunoPolaski/bureau/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/bureau/internal/core/entities/value_objects"
	"gorm.io/gorm"
)

func TestNewPersonDTOWithoutAssociations(t *testing.T) {
	p := &entities.Person{
		Model:               gorm.Model{ID: 17},
		PersonalInformation: &entities.PersonalInformation{FullName: "Lucas Teixeira", Document: valueobjects.Document("12312312312")},
	}

	got := NewPersonDTO(p)

	if got.ID != 17 || got.PersonalInformation == nil {
		t.Fatalf("unexpected base DTO: %+v", got)
	}
	if got.CreditScore != nil || got.FinancialProfile != nil {
		t.Errorf("score and profile must stay nil when absent, got %+v / %+v", got.CreditScore, got.FinancialProfile)
	}
	if len(got.NegativeRecords) != 0 || len(got.FraudAlerts) != 0 || len(got.Debts) != 0 {
		t.Errorf("lists must stay empty when absent: %+v", got)
	}
}

func TestNewPersonDTOMapsBureauDimensions(t *testing.T) {
	income := 4200.0
	p := &entities.Person{
		Model:            gorm.Model{ID: 3},
		CreditScore:      &entities.CreditScore{Score: 640},
		FinancialProfile: &entities.FinancialProfile{DeclaredMonthlyIncome: &income},
		NegativeRecords:  []entities.NegativeRecord{{RecordType: "protest", Amount: 1500}, {RecordType: "spc", Amount: 300}},
	}

	got := NewPersonDTO(p)

	if got.CreditScore == nil || got.CreditScore.Score != 640 {
		t.Fatalf("credit score not mapped: %+v", got.CreditScore)
	}
	if got.FinancialProfile == nil || got.FinancialProfile.DeclaredMonthlyIncome == nil || *got.FinancialProfile.DeclaredMonthlyIncome != 4200 {
		t.Fatalf("financial profile not mapped: %+v", got.FinancialProfile)
	}
	if len(got.NegativeRecords) != 2 || got.NegativeRecords[0].Amount != 1500 || got.NegativeRecords[1].Amount != 300 {
		t.Fatalf("negative records not mapped in order: %+v", got.NegativeRecords)
	}
}

func TestNewPersonSummaryDTO(t *testing.T) {
	tests := []struct {
		name         string
		person       *entities.Person
		wantName     string
		wantDocument string
	}{
		{
			name:         "full record",
			person:       &entities.Person{Model: gorm.Model{ID: 7}, PersonalInformation: &entities.PersonalInformation{FullName: "Felipe Santos", Document: valueobjects.Document("11122233344")}},
			wantName:     "Felipe Santos",
			wantDocument: "11122233344",
		},
		{
			name:     "absent CPF stays empty",
			person:   &entities.Person{Model: gorm.Model{ID: 7}, PersonalInformation: &entities.PersonalInformation{FullName: "Felipe Santos"}},
			wantName: "Felipe Santos",
		},
		{
			name:   "no personal information",
			person: &entities.Person{Model: gorm.Model{ID: 7}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPersonSummaryDTO(tt.person)
			if got.ID != 7 || got.Name != tt.wantName || got.Document != tt.wantDocument {
				t.Errorf("got %+v, want name %q document %q", got, tt.wantName, tt.wantDocument)
			}
		})
	}
}
