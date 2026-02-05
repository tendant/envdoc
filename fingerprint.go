package envdoc

import (
	"crypto/sha256"
	"fmt"
)

// FingerprintValue returns the first 8 hex characters of the SHA-256 hash of value.
func FingerprintValue(value string) string {
	h := sha256.Sum256([]byte(value))
	return fmt.Sprintf("%x", h[:4])
}

// shouldFingerprint decides whether to compute a fingerprint for a variable.
// Logic:
//  1. If the rule explicitly sets Fingerprint, use that.
//  2. If the variable is secret-like, do NOT fingerprint (unless rule says so).
//  3. If DumpAllFingerprint is enabled, fingerprint everything non-secret.
func shouldFingerprint(secretLike bool, rule Rule, dumpAllFingerprint bool) bool {
	// Explicit rule override
	if rule.Fingerprint != nil {
		return *rule.Fingerprint
	}
	// Secret-like vars don't get fingerprinted by default
	if secretLike {
		return false
	}
	// In dump-all fingerprint mode, fingerprint all non-secret vars
	if dumpAllFingerprint {
		return true
	}
	return false
}
