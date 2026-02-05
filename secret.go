package envdoc

import "regexp"

// secretPattern matches environment variable names that likely contain secrets.
var secretPattern = regexp.MustCompile(
	`(?i)(password|passwd|pwd|secret|token|apikey|api[_\-]?key|private[_\-]?key|ssh|pem|cert|jwt|session|cookie|authorization|bearer|kms|encrypt|signing)`,
)

// IsSecretLike returns true if the key name matches common secret patterns.
func IsSecretLike(key string) bool {
	return secretPattern.MatchString(key)
}

// classifySecretLike determines if a variable is secret-like.
// Rule.Secret explicitly overrides the heuristic.
func classifySecretLike(key string, rule Rule) bool {
	if rule.Secret != nil {
		return *rule.Secret
	}
	return IsSecretLike(key)
}
