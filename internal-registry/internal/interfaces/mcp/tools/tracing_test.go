package tools

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/BrunoPolaski/internal-registry/internal/infra/controllers/dto"
	"github.com/BrunoPolaski/internal-registry/internal/infra/thirdparty/logger"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestMain(m *testing.M) {
	logger.Init()
	m.Run()
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
