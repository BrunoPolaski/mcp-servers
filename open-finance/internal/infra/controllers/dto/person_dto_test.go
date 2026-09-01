package dto

import (
	"testing"

	"github.com/BrunoPolaski/open-finance/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/open-finance/internal/core/entities/value_objects"
	"gorm.io/gorm"
)

func TestNewPersonSummaryDTO(t *testing.T) {
	tests := []struct {
		name         string
		entity       *entities.Person
		wantID       uint
		wantName     string
		wantDocument string
	}{
		{
			name: "pessoa com informações pessoais",
			entity: &entities.Person{
				Model: gorm.Model{ID: 7},
				PersonalInformation: &entities.PersonalInformation{
					FullName: "Igor Souza Martins",
					Document: valueobjects.Document("40161087990"),
				},
			},
			wantID:       7,
			wantName:     "Igor Souza Martins",
			wantDocument: "40161087990",
		},
		{
			name:         "pessoa sem informações pessoais não deve causar panic",
			entity:       &entities.Person{Model: gorm.Model{ID: 3}},
			wantID:       3,
			wantName:     "",
			wantDocument: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPersonSummaryDTO(tt.entity)

			if got.ID != tt.wantID {
				t.Errorf("ID = %d, esperado %d", got.ID, tt.wantID)
			}
			if got.Name != tt.wantName {
				t.Errorf("Name = %q, esperado %q", got.Name, tt.wantName)
			}
			if got.Document != tt.wantDocument {
				t.Errorf("Document = %q, esperado %q", got.Document, tt.wantDocument)
			}
		})
	}
}

func TestNewPersonDTOComAssociacoesVazias(t *testing.T) {
	got := NewPersonDTO(&entities.Person{Model: gorm.Model{ID: 1}})

	if got.ID != 1 {
		t.Errorf("ID = %d, esperado 1", got.ID)
	}
	if got.PersonalInformation != nil {
		t.Error("PersonalInformation deveria ser nil quando a entidade não a carrega")
	}
	if got.BankAccountProfile != nil {
		t.Error("BankAccountProfile deveria ser nil quando a entidade não o carrega")
	}
	if len(got.BankStatements) != 0 {
		t.Errorf("BankStatements = %d itens, esperado 0", len(got.BankStatements))
	}
	if len(got.DataSharingConsents) != 0 {
		t.Errorf("DataSharingConsents = %d itens, esperado 0", len(got.DataSharingConsents))
	}
}

func TestNewPersonDTOAgregaDimensoes(t *testing.T) {
	entity := &entities.Person{
		Model: gorm.Model{ID: 2},
		PersonalInformation: &entities.PersonalInformation{
			FullName: "Henrique Martins Barbosa",
			Document: valueobjects.Document("15393772025"),
		},
		BankAccountProfile: &entities.BankAccountProfile{BankingRelationships: 3},
		BankStatements:     []entities.BankStatement{{Institution: "Banco Sintético Beta"}},
		CashFlowAnalyses:   []entities.CashFlowAnalysis{{NetCashFlow: 1500}},
		RecurringTransactions: []entities.RecurringTransaction{
			{TransactionType: "income", Amount: 6950},
		},
		DataSharingConsents: []entities.DataSharingConsent{{Status: "granted"}},
	}

	got := NewPersonDTO(entity)

	if got.PersonalInformation == nil || got.PersonalInformation.FullName != "Henrique Martins Barbosa" {
		t.Error("PersonalInformation não foi mapeada")
	}
	if got.BankAccountProfile == nil || got.BankAccountProfile.BankingRelationships != 3 {
		t.Error("BankAccountProfile não foi mapeado")
	}
	if len(got.BankStatements) != 1 || got.BankStatements[0].Institution != "Banco Sintético Beta" {
		t.Error("BankStatements não foram mapeados")
	}
	if len(got.CashFlowAnalyses) != 1 || got.CashFlowAnalyses[0].NetCashFlow != 1500 {
		t.Error("CashFlowAnalyses não foram mapeadas")
	}
	if len(got.RecurringTransactions) != 1 || got.RecurringTransactions[0].Amount != 6950 {
		t.Error("RecurringTransactions não foram mapeadas")
	}
	if len(got.DataSharingConsents) != 1 || got.DataSharingConsents[0].Status != "granted" {
		t.Error("DataSharingConsents não foram mapeados")
	}
}
