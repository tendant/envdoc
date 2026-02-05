package envdoc

import "testing"

func TestFingerprintValue(t *testing.T) {
	fp := FingerprintValue("hello")
	if len(fp) != 8 {
		t.Errorf("expected 8 hex chars, got %d: %q", len(fp), fp)
	}

	// SHA-256 of "hello" starts with 2cf24dba
	if fp != "2cf24dba" {
		t.Errorf("expected 2cf24dba, got %q", fp)
	}

	// Same input gives same fingerprint
	if FingerprintValue("hello") != FingerprintValue("hello") {
		t.Error("expected deterministic fingerprint")
	}

	// Different input gives different fingerprint
	if FingerprintValue("hello") == FingerprintValue("world") {
		t.Error("expected different fingerprints for different inputs")
	}
}

func TestShouldFingerprint(t *testing.T) {
	// Explicit rule fingerprint=true always wins
	if !shouldFingerprint(true, Rule{Fingerprint: boolPtr(true)}, false) {
		t.Error("explicit fingerprint=true should override secret-like")
	}

	// Explicit rule fingerprint=false always wins
	if shouldFingerprint(false, Rule{Fingerprint: boolPtr(false)}, true) {
		t.Error("explicit fingerprint=false should override dump-all")
	}

	// Secret-like without explicit rule: no fingerprint
	if shouldFingerprint(true, Rule{}, false) {
		t.Error("secret-like var should not be fingerprinted by default")
	}

	// Non-secret with dump-all-fingerprint: yes
	if !shouldFingerprint(false, Rule{}, true) {
		t.Error("non-secret var should be fingerprinted in dump-all-fingerprint mode")
	}

	// Non-secret without dump-all: no
	if shouldFingerprint(false, Rule{}, false) {
		t.Error("non-secret var should not be fingerprinted without dump-all-fingerprint")
	}
}
