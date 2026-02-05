package envdoc

import (
	"testing"
	"time"
)

func TestLoadConfig_Defaults(t *testing.T) {
	env := MapEnvReader{}
	cfg := LoadConfig(env)

	if cfg.Mode != ModeAllowlist {
		t.Errorf("expected allowlist mode, got %s", cfg.Mode)
	}
	if cfg.EnableHTTP {
		t.Error("expected EnableHTTP=false")
	}
	if cfg.FailFast {
		t.Error("expected FailFast=false")
	}
	if cfg.DumpAll {
		t.Error("expected DumpAll=false")
	}
	if cfg.Token != "" {
		t.Error("expected empty Token")
	}
}

func TestLoadConfig_DumpAll(t *testing.T) {
	env := MapEnvReader{
		"ENVDOC_DUMP_ALL":    "true",
		"ENVDOC_ENABLE_HTTP": "true",
	}
	cfg := LoadConfig(env)

	if cfg.Mode != ModeDumpAll {
		t.Errorf("expected dumpall mode, got %s", cfg.Mode)
	}
	if !cfg.DumpAll {
		t.Error("expected DumpAll=true")
	}
	if !cfg.EnableHTTP {
		t.Error("expected EnableHTTP=true")
	}
}

func TestLoadConfig_AllOptions(t *testing.T) {
	env := MapEnvReader{
		"ENVDOC_DUMP_ALL":             "true",
		"ENVDOC_ENABLE_HTTP":          "1",
		"ENVDOC_FAIL_FAST":            "true",
		"ENVDOC_DUMP_ALL_FINGERPRINT": "true",
		"ENVDOC_TOKEN":                "my-secret",
		"ENVDOC_EXPIRES_AT":           "2026-02-05T20:00:00Z",
		"ENVDOC_LISTEN_ADDR":          "0.0.0.0:8080",
	}
	cfg := LoadConfig(env)

	if !cfg.DumpAll {
		t.Error("expected DumpAll=true")
	}
	if !cfg.EnableHTTP {
		t.Error("expected EnableHTTP=true")
	}
	if !cfg.FailFast {
		t.Error("expected FailFast=true")
	}
	if !cfg.DumpAllFingerprint {
		t.Error("expected DumpAllFingerprint=true")
	}
	if cfg.Token != "my-secret" {
		t.Errorf("expected token 'my-secret', got %q", cfg.Token)
	}
	if cfg.ExpiresAt.IsZero() {
		t.Error("expected ExpiresAt to be set")
	}
	expected := time.Date(2026, 2, 5, 20, 0, 0, 0, time.UTC)
	if !cfg.ExpiresAt.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, cfg.ExpiresAt)
	}
	if cfg.ListenAddr != "0.0.0.0:8080" {
		t.Errorf("expected 0.0.0.0:8080, got %q", cfg.ListenAddr)
	}
}
