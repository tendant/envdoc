package envdoc

import (
	"bytes"
	"strings"
	"testing"
)

func TestLogReport(t *testing.T) {
	report := &Report{
		Results: []VarResult{
			{Key: "DB_HOST", Present: true, Length: 9, Valid: true},
			{Key: "DB_PORT", Present: true, Length: 4, Valid: true, Required: true},
			{Key: "MISSING", Present: false, Valid: false, Required: true, Problems: []string{"required but not set"}},
			{Key: "DB_PASSWORD", Present: true, Length: 32, Valid: true, Fingerprint: "9f2c1a2b", SecretLike: true},
			{Key: "PADDED", Present: true, Length: 7, Valid: true, Trimmed: true},
		},
	}

	var buf bytes.Buffer
	LogReport(&buf, report)
	output := buf.String()

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 5 {
		t.Errorf("expected 5 lines, got %d", len(lines))
	}

	// Check DB_HOST line
	if !strings.Contains(lines[0], "key=DB_HOST") {
		t.Errorf("expected key=DB_HOST in line: %s", lines[0])
	}
	if !strings.Contains(lines[0], "present=true") {
		t.Errorf("expected present=true in line: %s", lines[0])
	}
	if !strings.Contains(lines[0], "len=9") {
		t.Errorf("expected len=9 in line: %s", lines[0])
	}

	// Check required field
	if !strings.Contains(lines[1], "required=true") {
		t.Errorf("expected required=true in line: %s", lines[1])
	}

	// Check missing var
	if !strings.Contains(lines[2], "present=false") {
		t.Errorf("expected present=false in line: %s", lines[2])
	}
	if !strings.Contains(lines[2], "valid=false") {
		t.Errorf("expected valid=false in line: %s", lines[2])
	}

	// Check fingerprint
	if !strings.Contains(lines[3], "fp=9f2c1a2b") {
		t.Errorf("expected fp=9f2c1a2b in line: %s", lines[3])
	}

	// Check trimmed
	if !strings.Contains(lines[4], "trimmed=true") {
		t.Errorf("expected trimmed=true in line: %s", lines[4])
	}
}
