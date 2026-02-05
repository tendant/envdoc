package envdoc

import (
	"os"
	"strings"
	"time"
)

// Mode represents the operating mode.
type Mode string

const (
	ModeAllowlist Mode = "allowlist"
	ModeDumpAll   Mode = "dumpall"
)

// Config holds all ENVDOC_* configuration.
type Config struct {
	Mode               Mode
	EnableHTTP         bool
	FailFast           bool
	DumpAll            bool
	DumpAllFingerprint bool
	Token              string
	ExpiresAt          time.Time
	ListenAddr         string
}

// LoadConfig reads ENVDOC_* environment variables from the given EnvReader.
func LoadConfig(env EnvReader) Config {
	cfg := Config{
		Mode: ModeAllowlist,
	}

	if parseBool(env.Getenv("ENVDOC_DUMP_ALL")) {
		cfg.DumpAll = true
		cfg.Mode = ModeDumpAll
	}

	cfg.EnableHTTP = parseBool(env.Getenv("ENVDOC_ENABLE_HTTP"))
	cfg.FailFast = parseBool(env.Getenv("ENVDOC_FAIL_FAST"))
	cfg.DumpAllFingerprint = parseBool(env.Getenv("ENVDOC_DUMP_ALL_FINGERPRINT"))
	cfg.Token = env.Getenv("ENVDOC_TOKEN")
	cfg.ListenAddr = env.Getenv("ENVDOC_LISTEN_ADDR")

	if exp := env.Getenv("ENVDOC_EXPIRES_AT"); exp != "" {
		if t, err := time.Parse(time.RFC3339, exp); err == nil {
			cfg.ExpiresAt = t
		}
	}

	return cfg
}

func parseBool(s string) bool {
	return strings.EqualFold(s, "true") || s == "1"
}

// readFileOS reads a file from the OS filesystem.
func readFileOS(path string) ([]byte, error) {
	return os.ReadFile(path)
}
