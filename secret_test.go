package envdoc

import "testing"

func TestIsSecretLike(t *testing.T) {
	secrets := []string{
		"DB_PASSWORD", "API_KEY", "JWT_TOKEN", "PRIVATE_KEY",
		"SSH_KEY", "BEARER_TOKEN", "SESSION_ID", "COOKIE_SECRET",
		"AUTHORIZATION_HEADER", "KMS_KEY", "ENCRYPT_KEY", "SIGNING_KEY",
		"aws_secret_access_key", "api-key", "PEM_FILE",
	}
	for _, k := range secrets {
		if !IsSecretLike(k) {
			t.Errorf("expected %q to be secret-like", k)
		}
	}

	nonSecrets := []string{
		"DB_HOST", "DB_PORT", "APP_NAME", "LOG_LEVEL",
		"FEATURE_FLAGS", "MAX_RETRIES", "HOME",
	}
	for _, k := range nonSecrets {
		if IsSecretLike(k) {
			t.Errorf("expected %q to NOT be secret-like", k)
		}
	}
}

func TestClassifySecretLike_Override(t *testing.T) {
	// Explicit secret=true overrides heuristic
	if !classifySecretLike("DB_HOST", Rule{Secret: boolPtr(true)}) {
		t.Error("expected secret=true to override heuristic")
	}

	// Explicit secret=false overrides heuristic
	if classifySecretLike("DB_PASSWORD", Rule{Secret: boolPtr(false)}) {
		t.Error("expected secret=false to override heuristic")
	}

	// No override: falls back to heuristic
	if !classifySecretLike("DB_PASSWORD", Rule{}) {
		t.Error("expected heuristic to classify DB_PASSWORD as secret-like")
	}
	if classifySecretLike("DB_HOST", Rule{}) {
		t.Error("expected heuristic to classify DB_HOST as not secret-like")
	}
}
