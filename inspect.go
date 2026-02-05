package envdoc

import (
	"strings"
	"time"
)

// VarResult holds the inspection result for a single environment variable.
type VarResult struct {
	Key         string   `json:"key"`
	Present     bool     `json:"present"`
	Length      int      `json:"length"`
	Required    bool     `json:"required"`
	Valid       bool     `json:"valid"`
	Problems    []string `json:"problems,omitempty"`
	SecretLike  bool     `json:"secret_like"`
	Fingerprint string   `json:"fingerprint,omitempty"`
	Trimmed     bool     `json:"trimmed"`
}

// Summary holds aggregate counts.
type Summary struct {
	Total    int `json:"total"`
	Present  int `json:"present"`
	Valid    int `json:"valid"`
	Required int `json:"required"`
	Missing  int `json:"missing"`
}

// Report is the complete inspection output.
type Report struct {
	Timestamp time.Time   `json:"timestamp"`
	Mode      string      `json:"mode"`
	Results   []VarResult `json:"results"`
	Summary   Summary     `json:"summary"`
}

// inspect performs the core inspection logic.
func inspect(env EnvReader, clock Clock, rules []Rule, cfg Config) *Report {
	report := &Report{
		Timestamp: clock.Now(),
		Mode:      string(cfg.Mode),
	}
	if report.Mode == "" {
		report.Mode = "allowlist"
	}

	// Build rule map for quick lookup
	ruleMap := make(map[string]Rule)
	for _, r := range rules {
		ruleMap[r.Key] = r
	}

	// Determine which keys to inspect
	var keys []string
	if cfg.DumpAll {
		// Dump-all mode: inspect all env vars
		for _, pair := range env.Environ() {
			k, _, _ := strings.Cut(pair, "=")
			keys = append(keys, k)
		}
		// Also include rule keys that may not be set
		for _, r := range rules {
			if _, ok := ruleMap[r.Key]; ok {
				found := false
				for _, k := range keys {
					if k == r.Key {
						found = true
						break
					}
				}
				if !found {
					keys = append(keys, r.Key)
				}
			}
		}
	} else {
		// Allow-list mode: only rule-defined keys
		for _, r := range rules {
			keys = append(keys, r.Key)
		}
	}

	for _, key := range keys {
		vr := inspectVar(env, key, ruleMap[key], cfg)
		report.Results = append(report.Results, vr)

		// Update summary
		report.Summary.Total++
		if vr.Present {
			report.Summary.Present++
		}
		if vr.Valid {
			report.Summary.Valid++
		}
		if vr.Required {
			report.Summary.Required++
			if !vr.Present {
				report.Summary.Missing++
			}
		}
	}

	return report
}

// inspectVar inspects a single environment variable.
func inspectVar(env EnvReader, key string, rule Rule, cfg Config) VarResult {
	vr := VarResult{
		Key:      key,
		Required: rule.Required,
		Valid:    true,
	}

	value, present := env.LookupEnv(key)
	vr.Present = present

	if !present {
		if rule.Required {
			vr.Valid = false
			vr.Problems = append(vr.Problems, "required but not set")
		}
		// Classify secret-like even when not present
		vr.SecretLike = classifySecretLike(key, rule)
		return vr
	}

	vr.Length = len(value)
	vr.Trimmed = detectWhitespace(value)
	vr.SecretLike = classifySecretLike(key, rule)

	// Run validation if rule has any constraints
	if rule.Key != "" {
		problems := ValidateVar(value, rule)
		if len(problems) > 0 {
			vr.Valid = false
			vr.Problems = problems
		}
	}

	// Fingerprint decision
	if shouldFingerprint(vr.SecretLike, rule, cfg.DumpAllFingerprint) {
		vr.Fingerprint = FingerprintValue(value)
	}

	return vr
}
