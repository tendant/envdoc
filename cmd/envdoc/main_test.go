package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestCLI_WithRules(t *testing.T) {
	// Build the binary
	binPath := t.TempDir() + "/envdoc"
	cmd := exec.Command("go", "build", "-o", binPath, ".")
	cmd.Dir = "."
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}

	// Run with rules file and set required env vars
	run := exec.Command(binPath, "-rules", "../../testdata/basic_rules.yaml")
	run.Env = append(os.Environ(),
		"DB_HOST=localhost",
		"DB_PORT=5432",
		"DB_PASSWORD=super-secret-password-long-enough",
	)
	out, err := run.CombinedOutput()
	if err != nil {
		t.Fatalf("unexpected error: %v\n%s", err, out)
	}

	output := string(out)
	if !strings.Contains(output, "key=DB_HOST") {
		t.Errorf("expected DB_HOST in output: %s", output)
	}
	if !strings.Contains(output, "key=DB_PORT") {
		t.Errorf("expected DB_PORT in output: %s", output)
	}
}

func TestCLI_FailFast(t *testing.T) {
	binPath := t.TempDir() + "/envdoc"
	cmd := exec.Command("go", "build", "-o", binPath, ".")
	cmd.Dir = "."
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}

	// Run with fail-fast and missing required var
	run := exec.Command(binPath, "-rules", "../../testdata/basic_rules.yaml")
	run.Env = append(os.Environ(), "ENVDOC_FAIL_FAST=true")
	out, err := run.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit")
	}

	output := string(out)
	if !strings.Contains(output, "fail-fast") {
		t.Errorf("expected fail-fast in output: %s", output)
	}
}

func TestCLI_NoRules(t *testing.T) {
	binPath := t.TempDir() + "/envdoc"
	cmd := exec.Command("go", "build", "-o", binPath, ".")
	cmd.Dir = "."
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}

	// Run without rules - should work fine (no vars to check)
	run := exec.Command(binPath)
	run.Env = os.Environ()
	if out, err := run.CombinedOutput(); err != nil {
		t.Fatalf("unexpected error: %v\n%s", err, out)
	}
}

func TestCLI_InvalidRulesPath(t *testing.T) {
	binPath := t.TempDir() + "/envdoc"
	cmd := exec.Command("go", "build", "-o", binPath, ".")
	cmd.Dir = "."
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}

	run := exec.Command(binPath, "-rules", "/nonexistent/path.yaml")
	run.Env = os.Environ()
	out, err := run.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit for invalid rules path")
	}
	if !strings.Contains(string(out), "envdoc:") {
		t.Errorf("expected error message: %s", out)
	}
}
