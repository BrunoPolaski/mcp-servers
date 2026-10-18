package tools

import (
	"context"
	"testing"

	"github.com/BrunoPolaski/bureau/internal/core/entities"
	valueobjects "github.com/BrunoPolaski/bureau/internal/core/entities/value_objects"
	"github.com/BrunoPolaski/bureau/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/mark3labs/mcp-go/mcp"
	"gorm.io/gorm"
)

type fakePersonService struct {
	person    *entities.Person
	summary   *dto.PaginatedResponse[dto.PersonSummaryDTO]
	err       *rest_err.RestErr
	gotID     uint
	gotDoc    string
	gotLimit  int
	gotOffset int
	gotParams map[string]any
}

func (f *fakePersonService) GetById(_ context.Context, id uint) (*entities.Person, *rest_err.RestErr) {
	f.gotID = id
	return f.person, f.err
}

func (f *fakePersonService) GetByDocument(_ context.Context, document string) (*entities.Person, *rest_err.RestErr) {
	f.gotDoc = document
	return f.person, f.err
}

func (f *fakePersonService) GetAllSummary(_ context.Context, limit, offset int, params map[string]any) (*dto.PaginatedResponse[dto.PersonSummaryDTO], *rest_err.RestErr) {
	f.gotLimit, f.gotOffset, f.gotParams = limit, offset, params
	return f.summary, f.err
}

func requestWith(args map[string]any) mcp.CallToolRequest {
	var r mcp.CallToolRequest
	r.Params.Arguments = args
	return r
}

func personEntity() *entities.Person {
	score := 812
	return &entities.Person{
		Model: gorm.Model{ID: 5},
		PersonalInformation: &entities.PersonalInformation{
			FullName: "Gabriela Ribeiro Barbosa",
			Document: valueobjects.Document("35509139404"),
		},
		CreditScore: &entities.CreditScore{Score: score},
	}
}

func TestHandleGetPersonByID(t *testing.T) {
	tests := []struct {
		name    string
		args    map[string]any
		wantErr bool
		wantID  uint
	}{
		{name: "valid id", args: map[string]any{"id": float64(5)}, wantID: 5},
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
			if svc.gotID != tt.wantID {
				t.Errorf("id passed = %d, want %d", svc.gotID, tt.wantID)
			}
			if got.ID != 5 || got.CreditScore == nil || got.CreditScore.Score != 812 {
				t.Errorf("unexpected DTO: %+v", got)
			}
		})
	}
}

func TestHandleGetPersonByIDPropagatesNotFound(t *testing.T) {
	svc := &fakePersonService{err: rest_err.NewNotFoundError("person not found")}
	s := &Server{personService: svc}

	_, err := s.HandleGetPersonByID(context.Background(), requestWith(map[string]any{"id": float64(99)}), mcp.CallToolParams{})

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestHandleGetPersonByDocument(t *testing.T) {
	tests := []struct {
		name    string
		args    map[string]any
		wantErr bool
		wantDoc string
	}{
		{name: "valid document", args: map[string]any{"document": "35509139404"}, wantDoc: "35509139404"},
		{name: "missing document", args: map[string]any{}, wantErr: true},
		{name: "empty document", args: map[string]any{"document": ""}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &fakePersonService{person: personEntity()}
			s := &Server{personService: svc}

			_, err := s.HandleGetPersonByDocument(context.Background(), requestWith(tt.args), mcp.CallToolParams{})

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if svc.gotDoc != tt.wantDoc {
				t.Errorf("document passed = %q, want %q", svc.gotDoc, tt.wantDoc)
			}
		})
	}
}

func TestHandleGetAllPersonsUsesDefaultsAndFilters(t *testing.T) {
	svc := &fakePersonService{
		summary: dto.NewPaginatedResponse(1, []*dto.PersonSummaryDTO{
			{ID: 5, Name: "Gabriela Ribeiro Barbosa", Document: "35509139404"},
		}),
	}
	s := &Server{personService: svc}

	got, err := s.HandleGetAllPersons(context.Background(), requestWith(map[string]any{
		"params": map[string]any{"id": float64(5)},
	}), mcp.CallToolParams{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if svc.gotLimit != 10 || svc.gotOffset != 0 {
		t.Errorf("defaults = (limit %d, offset %d), want (10, 0)", svc.gotLimit, svc.gotOffset)
	}
	if svc.gotParams == nil {
		t.Error("params not passed to the service")
	}
	if got.Total != 1 || len(got.Items) != 1 {
		t.Errorf("response = (total %d, %d items), want (1, 1)", got.Total, len(got.Items))
	}
}

func TestHandleGetAllPersonsRejectsNegativePagination(t *testing.T) {
	tests := []struct {
		name string
		args map[string]any
	}{
		{name: "negative limit", args: map[string]any{"limit": float64(-1)}},
		{name: "negative offset", args: map[string]any{"offset": float64(-1)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{personService: &fakePersonService{}}

			if _, err := s.HandleGetAllPersons(context.Background(), requestWith(tt.args), mcp.CallToolParams{}); err == nil {
				t.Fatal("expected an error, got nil")
			}
		})
	}
}

func TestToolNames(t *testing.T) {
	s := &Server{personService: &fakePersonService{}}

	tests := []struct {
		got  string
		want string
	}{
		{s.GetPersonByIDTool().Name, "get_customer_by_id"},
		{s.GetPersonByDocumentTool().Name, "get_customer_by_document"},
		{s.GetAllPersonsTool().Name, "get_all_customers"},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("name = %q, want %q", tt.got, tt.want)
		}
	}
}
