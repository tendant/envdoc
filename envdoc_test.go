package envdoc

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestRun_AllowlistMode(t *testing.T) {
	env := MapEnvReader{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	rules := []Rule{
		{Key: "DB_HOST", Required: true, Type: TypeString},
		{Key: "DB_PORT", Required: true, Type: TypeInt},
	}

	var buf bytes.Buffer
	report, err := Run(
		WithEnvReader(env),
		WithClock(fixedClock{t: time.Now()}),
		WithRules(rules),
		WithConfig(Config{Mode: ModeAllowlist}),
		WithOutput(&buf),
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(report.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(report.Results))
	}
	if report.Summary.Valid != 2 {
		t.Errorf("expected 2 valid, got %d", report.Summary.Valid)
	}
	if buf.Len() == 0 {
		t.Error("expected log output")
	}
}

func TestRun_FailFast(t *testing.T) {
	env := MapEnvReader{}
	rules := []Rule{
		{Key: "REQUIRED_VAR", Required: true},
	}

	var buf bytes.Buffer
	_, err := Run(
		WithEnvReader(env),
		WithClock(fixedClock{t: time.Now()}),
		WithRules(rules),
		WithConfig(Config{Mode: ModeAllowlist, FailFast: true}),
		WithOutput(&buf),
	)

	if err == nil {
		t.Fatal("expected fail-fast error")
	}
	if !strings.Contains(err.Error(), "fail-fast") {
		t.Errorf("expected fail-fast in error message, got: %v", err)
	}
	if !strings.Contains(err.Error(), "REQUIRED_VAR") {
		t.Errorf("expected REQUIRED_VAR in error message, got: %v", err)
	}
}

func TestCheckFailFast_NoError(t *testing.T) {
	report := &Report{
		Results: []VarResult{
			{Key: "A", Present: true, Valid: true, Required: true},
			{Key: "B", Present: true, Valid: true, Required: false},
		},
	}
	if err := CheckFailFast(report, true); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCheckFailFast_Disabled(t *testing.T) {
	report := &Report{
		Results: []VarResult{
			{Key: "A", Present: false, Valid: false, Required: true},
		},
	}
	if err := CheckFailFast(report, false); err != nil {
		t.Errorf("unexpected error when fail-fast disabled: %v", err)
	}
}

func TestCheckFailFast_MultipleFailures(t *testing.T) {
	report := &Report{
		Results: []VarResult{
			{Key: "A", Present: false, Valid: false, Required: true, Problems: []string{"required but not set"}},
			{Key: "B", Present: true, Valid: false, Required: true, Problems: []string{"not a valid int"}},
			{Key: "C", Present: true, Valid: true, Required: false},
		},
	}
	err := CheckFailFast(report, true)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "2 required variable(s) invalid") {
		t.Errorf("expected 2 failures, got: %v", err)
	}
}

func TestNew_DefaultConfig(t *testing.T) {
	env := MapEnvReader{
		"ENVDOC_FAIL_FAST": "true",
	}
	inspector := New(WithEnvReader(env))
	cfg := inspector.Config()
	if !cfg.FailFast {
		t.Error("expected fail-fast from env")
	}
}

func TestDumpAllMode_WithFingerprints(t *testing.T) {
	env := MapEnvReader{
		"APP_NAME":   "myapp",
		"SECRET_KEY": "mysecretvalue",
		"DB_HOST":    "localhost",
	}

	cfg := Config{
		Mode:               ModeDumpAll,
		DumpAll:            true,
		DumpAllFingerprint: true,
	}

	var buf bytes.Buffer
	report, err := Run(
		WithEnvReader(env),
		WithClock(fixedClock{t: time.Now()}),
		WithConfig(cfg),
		WithOutput(&buf),
	)
	if err != nil {
		t.Fatal(err)
	}

	// All env vars should be present
	if len(report.Results) < 3 {
		t.Errorf("expected at least 3 results, got %d", len(report.Results))
	}

	for _, r := range report.Results {
		if r.Key == "SECRET_KEY" {
			if !r.SecretLike {
				t.Error("expected SECRET_KEY to be secret-like")
			}
			// Secret-like vars should NOT have fingerprints (no explicit rule)
			if r.Fingerprint != "" {
				t.Error("expected no fingerprint for secret-like var without explicit rule")
			}
		}
		if r.Key == "APP_NAME" {
			if r.SecretLike {
				t.Error("expected APP_NAME to not be secret-like")
			}
			// Non-secret vars should have fingerprints in dump-all-fingerprint mode
			if r.Fingerprint == "" {
				t.Error("expected fingerprint for non-secret var in dump-all-fingerprint mode")
			}
		}
	}
}
