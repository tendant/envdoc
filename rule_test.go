package envdoc

import (
	"os"
	"strings"
	"testing"
)

func TestLoadRules_Basic(t *testing.T) {
	data, err := os.ReadFile("testdata/basic_rules.yaml")
	if err != nil {
		t.Fatal(err)
	}
	rules, err := LoadRules(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 7 {
		t.Errorf("expected 7 rules, got %d", len(rules))
	}

	// Check first rule
	if rules[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", rules[0].Key)
	}
	if !rules[0].Required {
		t.Error("expected DB_HOST to be required")
	}
	if rules[0].Type != TypeString {
		t.Errorf("expected string type, got %s", rules[0].Type)
	}
	if rules[0].MinLen == nil || *rules[0].MinLen != 1 {
		t.Error("expected min_len=1")
	}

	// Check secret rule
	if rules[2].Key != "DB_PASSWORD" {
		t.Errorf("expected DB_PASSWORD, got %s", rules[2].Key)
	}
	if rules[2].Secret == nil || !*rules[2].Secret {
		t.Error("expected secret=true")
	}
	if rules[2].Fingerprint == nil || !*rules[2].Fingerprint {
		t.Error("expected fingerprint=true")
	}
}

func TestLoadRules_DuplicateKey(t *testing.T) {
	data, err := os.ReadFile("testdata/invalid_duplicate.yaml")
	if err != nil {
		t.Fatal(err)
	}
	_, err = LoadRules(data)
	if err == nil {
		t.Fatal("expected error for duplicate key")
	}
	if !strings.Contains(err.Error(), "duplicate key") {
		t.Errorf("expected 'duplicate key' error, got: %v", err)
	}
}

func TestLoadRules_UnknownType(t *testing.T) {
	data, err := os.ReadFile("testdata/invalid_type.yaml")
	if err != nil {
		t.Fatal(err)
	}
	_, err = LoadRules(data)
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
	if !strings.Contains(err.Error(), "unknown type") {
		t.Errorf("expected 'unknown type' error, got: %v", err)
	}
}

func TestLoadRules_InvalidRegex(t *testing.T) {
	data, err := os.ReadFile("testdata/invalid_regex.yaml")
	if err != nil {
		t.Fatal(err)
	}
	_, err = LoadRules(data)
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
	if !strings.Contains(err.Error(), "invalid regex") {
		t.Errorf("expected 'invalid regex' error, got: %v", err)
	}
}

func TestLoadRules_MinGreaterThanMax(t *testing.T) {
	data, err := os.ReadFile("testdata/invalid_minmax.yaml")
	if err != nil {
		t.Fatal(err)
	}
	_, err = LoadRules(data)
	if err == nil {
		t.Fatal("expected error for min > max")
	}
	if !strings.Contains(err.Error(), "min_len") {
		t.Errorf("expected min_len error, got: %v", err)
	}
}

func TestLoadRulesFile(t *testing.T) {
	rules, err := LoadRulesFile("testdata/basic_rules.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 7 {
		t.Errorf("expected 7 rules, got %d", len(rules))
	}
}

func TestLoadRulesFile_NotFound(t *testing.T) {
	_, err := LoadRulesFile("testdata/nonexistent.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
