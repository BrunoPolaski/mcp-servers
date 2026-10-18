package tools

import (
	"context"
	"testing"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/registration-validation/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/registration-validation/internal/core/entities/value_objects"
	"github.com/BrunoPolaski/registration-validation/internal/infra/controllers/dto"
	"github.com/mark3labs/mcp-go/mcp"
	"gorm.io/gorm"
)

type fakePersonService struct {
	person    *entities.Person
	page      *dto.PaginatedResponse[dto.PersonDTO]
	err       *rest_err.RestErr
	gotID     uint
	gotDoc    string
	gotLimit  int
	gotOffset int
}

func (f *fakePersonService) GetById(_ context.Context, id uint) (*entities.Person, *rest_err.RestErr) {
	f.gotID = id
	return f.person, f.err
}

func (f *fakePersonService) GetByDocument(_ context.Context, document string) (*entities.Person, *rest_err.RestErr) {
	f.gotDoc = document
	return f.person, f.err
}

func (f *fakePersonService) GetAll(_ context.Context, limit, offset int, _ map[string]any) (*dto.PaginatedResponse[dto.PersonDTO], *rest_err.RestErr) {
	f.gotLimit, f.gotOffset = limit, offset
	return f.page, f.err
}

func requestWith(args map[string]any) mcp.CallToolRequest {
	var r mcp.CallToolRequest
	r.Params.Arguments = args
	return r
}

func personEntity() *entities.Person {
	return &entities.Person{
		Model: gorm.Model{ID: 5},
		PersonalInformation: &entities.PersonalInformation{
			FullName: "Gabriela Ribeiro Barbosa",
			Document: valueobjects.Document("35509139404"),
		},
		DocumentValidations: []entities.DocumentValidation{{DocumentType: "cpf", ReceitaFederalStatus: "regular", IsValid: true}},
	}
}

func TestHandleGetPersonByID(t *testing.T) {
	tests := []struct {
		name    string
		args    map[string]any
		wantErr bool
	}{
		{name: "valid id", args: map[string]any{"id": float64(5)}},
		{name: "missing id", args: map[string]any{}, wantErr: true},
		{name: "zero id", args: map[string]any{"id": float64(0)}, wantErr: true},
		{name: "negative id", args: map[string]any{"id": float64(-3)}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &fakePersonService{person: personEntity()}
			s := &Server{personService: svc}

			got, err := s.HandleGetPersonByID(context.Background(), requestWith(tt.args), mcp.CallToolParams{})

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if svc.gotID != 5 || len(got.DocumentValidations) != 1 || got.DocumentValidations[0].ReceitaFederalStatus != "regular" {
				t.Errorf("unexpected DTO or id: id=%d dto=%+v", svc.gotID, got)
			}
		})
	}
}

func TestHandleGetPersonByDocumentRejectsEmpty(t *testing.T) {
	s := &Server{personService: &fakePersonService{person: personEntity()}}

	if _, err := s.HandleGetPersonByDocument(context.Background(), requestWith(map[string]any{"document": ""}), mcp.CallToolParams{}); err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestHandleGetPersonByDocumentPropagatesNotFound(t *testing.T) {
	svc := &fakePersonService{err: rest_err.NewNotFoundError("person not found")}
	s := &Server{personService: svc}

	_, err := s.HandleGetPersonByDocument(context.Background(), requestWith(map[string]any{"document": "00000000000"}), mcp.CallToolParams{})

	if err == nil || svc.gotDoc != "00000000000" {
		t.Fatalf("expected not found for the given document, got err=%v doc=%q", err, svc.gotDoc)
	}
}

func TestHandleGetAllPersons(t *testing.T) {
	tests := []struct {
		name       string
		args       map[string]any
		wantErr    bool
		wantLimit  int
		wantOffset int
	}{
		{name: "defaults", args: map[string]any{}, wantLimit: 10},
		{name: "explicit pagination", args: map[string]any{"limit": float64(0), "offset": float64(20)}, wantLimit: 0, wantOffset: 20},
		{name: "negative limit", args: map[string]any{"limit": float64(-1)}, wantErr: true},
		{name: "negative offset", args: map[string]any{"offset": float64(-1)}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &fakePersonService{page: dto.NewPaginatedResponse(0, []*dto.PersonDTO{})}
			s := &Server{personService: svc}

			_, err := s.HandleGetAllPersons(context.Background(), requestWith(tt.args), mcp.CallToolParams{})

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if svc.gotLimit != tt.wantLimit || svc.gotOffset != tt.wantOffset {
				t.Errorf("pagination = (%d, %d), want (%d, %d)", svc.gotLimit, svc.gotOffset, tt.wantLimit, tt.wantOffset)
			}
		})
	}
}
