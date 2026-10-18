package tools

import (
	"context"
	"errors"
	"testing"

	"github.com/BrunoPolaski/bureau/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/bureau/internal/infra/thirdparty/logger"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestMain(m *testing.M) {
	logger.Init()
	m.Run()
}

func TestTracedStampsMetadataOnSuccess(t *testing.T) {
	handler := traced(func(context.Context, mcp.CallToolRequest, mcp.CallToolParams) (*dto.PersonDTO, error) {
		return &dto.PersonDTO{ID: 5}, nil
	})
	req := requestWith(map[string]any{"id": float64(5)})
	req.Params.Name = "get_customer_by_id"

	got, err := handler(context.Background(), req, mcp.CallToolParams{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Metadata == nil {
		t.Fatal("metadata not set")
	}
	if got.Metadata.Source != "bureau" || got.Metadata.Tool != "get_customer_by_id" {
		t.Errorf("unexpected metadata: %+v", got.Metadata)
	}
	if got.Metadata.TransactionID == "" || got.Metadata.QueriedAt == "" {
		t.Errorf("transaction id and query time must be set: %+v", got.Metadata)
	}
}

func TestTracedPassesErrorsThrough(t *testing.T) {
	want := errors.New("boom")
	handler := traced(func(context.Context, mcp.CallToolRequest, mcp.CallToolParams) (*dto.PersonDTO, error) {
		return nil, want
	})

	got, err := handler(context.Background(), requestWith(nil), mcp.CallToolParams{})

	if !errors.Is(err, want) || got != nil {
		t.Fatalf("got (%v, %v), want (nil, %v)", got, err, want)
	}
}

func TestTracedUsesNewTransactionIDPerCall(t *testing.T) {
	handler := traced(func(context.Context, mcp.CallToolRequest, mcp.CallToolParams) (*dto.PersonDTO, error) {
		return &dto.PersonDTO{}, nil
	})

	a, _ := handler(context.Background(), requestWith(nil), mcp.CallToolParams{})
	b, _ := handler(context.Background(), requestWith(nil), mcp.CallToolParams{})

	if a.Metadata.TransactionID == b.Metadata.TransactionID {
		t.Errorf("both calls got transaction id %q", a.Metadata.TransactionID)
	}
}
