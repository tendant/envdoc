package envdoc

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_Basic(t *testing.T) {
	env := MapEnvReader{
		"DB_HOST": "localhost",
	}
	rules := []Rule{
		{Key: "DB_HOST", Required: true, Type: TypeString},
	}

	inspector := New(
		WithEnvReader(env),
		WithClock(fixedClock{t: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}),
		WithRules(rules),
		WithConfig(Config{Mode: ModeAllowlist}),
	)

	handler := inspector.Handler()

	req := httptest.NewRequest(http.MethodGet, "/debug/env", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	var report Report
	if err := json.NewDecoder(rec.Body).Decode(&report); err != nil {
		t.Fatal(err)
	}
	if len(report.Results) != 1 {
		t.Errorf("expected 1 result, got %d", len(report.Results))
	}
	if report.Results[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", report.Results[0].Key)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	inspector := New(
		WithEnvReader(MapEnvReader{}),
		WithConfig(Config{}),
	)
	handler := inspector.Handler()

	req := httptest.NewRequest(http.MethodPost, "/debug/env", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}

func TestHandler_TokenAuth(t *testing.T) {
	inspector := New(
		WithEnvReader(MapEnvReader{}),
		WithConfig(Config{Token: "my-token"}),
	)
	handler := inspector.Handler()

	// No token
	req := httptest.NewRequest(http.MethodGet, "/debug/env", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 without token, got %d", rec.Code)
	}

	// Wrong token
	req = httptest.NewRequest(http.MethodGet, "/debug/env", nil)
	req.Header.Set("Authorization", "Bearer wrong-token")
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 with wrong token, got %d", rec.Code)
	}

	// Correct token
	req = httptest.NewRequest(http.MethodGet, "/debug/env", nil)
	req.Header.Set("Authorization", "Bearer my-token")
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 with correct token, got %d", rec.Code)
	}
}

func TestHandler_Expiry(t *testing.T) {
	expiresAt := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

	// Before expiry
	inspector := New(
		WithEnvReader(MapEnvReader{}),
		WithClock(fixedClock{t: time.Date(2026, 1, 1, 11, 0, 0, 0, time.UTC)}),
		WithConfig(Config{ExpiresAt: expiresAt}),
	)
	handler := inspector.Handler()

	req := httptest.NewRequest(http.MethodGet, "/debug/env", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 before expiry, got %d", rec.Code)
	}

	// After expiry
	inspector = New(
		WithEnvReader(MapEnvReader{}),
		WithClock(fixedClock{t: time.Date(2026, 1, 1, 13, 0, 0, 0, time.UTC)}),
		WithConfig(Config{ExpiresAt: expiresAt}),
	)
	handler = inspector.Handler()

	req = httptest.NewRequest(http.MethodGet, "/debug/env", nil)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusGone {
		t.Errorf("expected 410 after expiry, got %d", rec.Code)
	}
}
