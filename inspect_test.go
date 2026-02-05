package envdoc

import (
	"testing"
	"time"
)

func TestInspect_AllowlistMode(t *testing.T) {
	env := MapEnvReader{
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"DB_PASSWORD": "super-secret-password-long",
	}

	rules := []Rule{
		{Key: "DB_HOST", Required: true, Type: TypeString, MinLen: intPtr(1)},
		{Key: "DB_PORT", Required: true, Type: TypeInt},
		{Key: "DB_PASSWORD", Required: true, Type: TypeString, Secret: boolPtr(true), MinLen: intPtr(16), Fingerprint: boolPtr(true)},
		{Key: "MISSING_VAR", Required: false},
	}

	cfg := Config{Mode: ModeAllowlist}
	clock := fixedClock{t: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}

	report := inspect(env, clock, rules, cfg)

	if report.Mode != "allowlist" {
		t.Errorf("expected allowlist mode, got %s", report.Mode)
	}
	if len(report.Results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(report.Results))
	}

	// DB_HOST
	r := report.Results[0]
	if !r.Present || !r.Valid || r.Key != "DB_HOST" {
		t.Errorf("DB_HOST: present=%t valid=%t key=%s", r.Present, r.Valid, r.Key)
	}
	if r.Length != 9 {
		t.Errorf("DB_HOST length: expected 9, got %d", r.Length)
	}

	// DB_PORT
	r = report.Results[1]
	if !r.Present || !r.Valid {
		t.Errorf("DB_PORT: present=%t valid=%t", r.Present, r.Valid)
	}

	// DB_PASSWORD
	r = report.Results[2]
	if !r.Present || !r.Valid {
		t.Errorf("DB_PASSWORD: present=%t valid=%t", r.Present, r.Valid)
	}
	if !r.SecretLike {
		t.Error("expected DB_PASSWORD to be secret-like")
	}
	if r.Fingerprint == "" {
		t.Error("expected DB_PASSWORD to have a fingerprint (rule says fingerprint=true)")
	}

	// MISSING_VAR
	r = report.Results[3]
	if r.Present {
		t.Error("expected MISSING_VAR to not be present")
	}
	if !r.Valid {
		t.Error("expected MISSING_VAR to be valid (not required)")
	}

	// Summary
	if report.Summary.Total != 4 {
		t.Errorf("expected total 4, got %d", report.Summary.Total)
	}
	if report.Summary.Present != 3 {
		t.Errorf("expected present 3, got %d", report.Summary.Present)
	}
	if report.Summary.Required != 3 {
		t.Errorf("expected required 3, got %d", report.Summary.Required)
	}
}

func TestInspect_MissingRequired(t *testing.T) {
	env := MapEnvReader{}
	rules := []Rule{
		{Key: "REQUIRED_VAR", Required: true, Type: TypeString},
	}
	cfg := Config{Mode: ModeAllowlist}
	clock := fixedClock{t: time.Now()}

	report := inspect(env, clock, rules, cfg)

	r := report.Results[0]
	if r.Present {
		t.Error("expected not present")
	}
	if r.Valid {
		t.Error("expected not valid for missing required var")
	}
	if len(r.Problems) != 1 {
		t.Errorf("expected 1 problem, got %d", len(r.Problems))
	}
	if report.Summary.Missing != 1 {
		t.Errorf("expected 1 missing, got %d", report.Summary.Missing)
	}
}

func TestInspect_DumpAllMode(t *testing.T) {
	env := MapEnvReader{
		"DB_HOST":    "localhost",
		"APP_NAME":   "myapp",
		"SECRET_KEY": "s3cr3t",
	}
	cfg := Config{Mode: ModeDumpAll, DumpAll: true}
	clock := fixedClock{t: time.Now()}

	report := inspect(env, clock, nil, cfg)

	if report.Mode != "dumpall" {
		t.Errorf("expected dumpall mode, got %s", report.Mode)
	}

	// Should have at least 3 results (all env vars)
	if len(report.Results) < 3 {
		t.Errorf("expected at least 3 results, got %d", len(report.Results))
	}

	// Find SECRET_KEY and verify it's classified as secret-like
	found := false
	for _, r := range report.Results {
		if r.Key == "SECRET_KEY" {
			found = true
			if !r.SecretLike {
				t.Error("expected SECRET_KEY to be secret-like")
			}
		}
	}
	if !found {
		t.Error("expected to find SECRET_KEY in results")
	}
}

func TestInspect_ValidationFails(t *testing.T) {
	env := MapEnvReader{
		"PORT": "notanumber",
	}
	rules := []Rule{
		{Key: "PORT", Required: true, Type: TypeInt},
	}
	cfg := Config{Mode: ModeAllowlist}
	clock := fixedClock{t: time.Now()}

	report := inspect(env, clock, rules, cfg)

	r := report.Results[0]
	if r.Valid {
		t.Error("expected invalid for non-int value")
	}
	if len(r.Problems) == 0 {
		t.Error("expected problems for non-int value")
	}
}

func TestInspect_WhitespaceDetection(t *testing.T) {
	env := MapEnvReader{
		"PADDED": " hello ",
	}
	rules := []Rule{
		{Key: "PADDED"},
	}
	cfg := Config{}
	clock := fixedClock{t: time.Now()}

	report := inspect(env, clock, rules, cfg)

	if !report.Results[0].Trimmed {
		t.Error("expected trimmed=true for value with whitespace")
	}
}
