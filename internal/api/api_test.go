package api

import (
	"net/http/httptest"
	"strings"
	"testing"

	"telemetry.local/drone/internal/service"
	"telemetry.local/drone/internal/store"
)

func TestHTTPHandlers(t *testing.T) {
	db, err := store.Open(t.TempDir() + "/api.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	server := NewServer(service.New(db, "api"))
	req := httptest.NewRequest("GET", "/health", nil)
	res := httptest.NewRecorder()
	server.Handler.ServeHTTP(res, req)
	if res.Code != 200 || !strings.Contains(res.Body.String(), "ok") {
		t.Fatalf("status=%d body=%s", res.Code, res.Body.String())
	}
}

func TestCLIParsing(t *testing.T) {
	options := ParseCLI([]string{"-command", "report", "-batch", "fixture"})
	if err := ValidateCLI(options); err != nil || options.Batch != "fixture" {
		t.Fatalf("options=%+v err=%v", options, err)
	}
}
