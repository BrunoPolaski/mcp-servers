package tools

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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

func TestWithRunIDStoresHeaderInContext(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	r.Header.Set("X-Run-Id", "5-1")

	got, _ := WithRunID(context.Background(), r).Value(runIDKey{}).(string)

	if got != "5-1" {
		t.Errorf("run id in context = %q, want %q", got, "5-1")
	}
}

func TestWithRunIDLeavesContextWithoutHeader(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/mcp", nil)

	if got := WithRunID(context.Background(), r).Value(runIDKey{}); got != nil {
		t.Errorf("run id in context = %v, want none", got)
	}
}

func TestTracedLogsRunID(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	r.Header.Set("X-Run-Id", "5-1")
	ctx := WithRunID(context.Background(), r)

	out := captureLog(t, func() { callTraced(ctx) })

	if !strings.Contains(out, `"run_id":"5-1"`) {
		t.Errorf("log has no run_id 5-1: %s", out)
	}
}

func TestTracedOmitsRunIDWithoutHeader(t *testing.T) {
	out := captureLog(t, func() { callTraced(context.Background()) })

	if !strings.Contains(out, `"transaction_id"`) {
		t.Fatalf("no trace record logged: %s", out)
	}
	if strings.Contains(out, `"run_id"`) {
		t.Errorf("log has run_id without header: %s", out)
	}
}

func callTraced(ctx context.Context) {
	handler := traced(func(context.Context, mcp.CallToolRequest, mcp.CallToolParams) (*dto.PersonDTO, error) {
		return &dto.PersonDTO{}, nil
	})
	handler(ctx, requestWith(nil), mcp.CallToolParams{})
}

func captureLog(t *testing.T, fn func()) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "log")
	if err != nil {
		t.Fatal(err)
	}
	stdout := os.Stdout
	os.Stdout = f
	logger.Init()
	os.Stdout = stdout
	t.Cleanup(logger.Init)

	fn()

	data, err := os.ReadFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
