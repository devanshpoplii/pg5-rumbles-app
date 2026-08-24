package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	healthHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var body healthResponse
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body.Status != "ok" {
		t.Errorf("expected status \"ok\", got %q", body.Status)
	}
}

func TestVersionHandler(t *testing.T) {
	version = "test-1.2.3"

	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()

	versionHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var body versionResponse
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body.Version != "test-1.2.3" {
		t.Errorf("expected version \"test-1.2.3\", got %q", body.Version)
	}
}

func TestRootHandlerNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/does-not-exist", nil)
	rec := httptest.NewRecorder()

	rootHandler(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404 for unknown path, got %d", rec.Code)
	}
}
