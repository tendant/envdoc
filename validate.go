package envdoc

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ValidateVar validates a value against a Rule and returns a list of problems.
// An empty return means the value is valid.
func ValidateVar(value string, rule Rule) []string {
	var problems []string

	// Type check
	if rule.Type != "" {
		if err := checkType(value, rule.Type); err != nil {
			problems = append(problems, err.Error())
		}
	}

	// Length checks
	if rule.MinLen != nil && len(value) < *rule.MinLen {
		problems = append(problems, fmt.Sprintf("length %d < min_len %d", len(value), *rule.MinLen))
	}
	if rule.MaxLen != nil && len(value) > *rule.MaxLen {
		problems = append(problems, fmt.Sprintf("length %d > max_len %d", len(value), *rule.MaxLen))
	}

	// Regex check
	if rule.Regex != "" {
		re, err := regexp.Compile(rule.Regex)
		if err == nil && !re.MatchString(value) {
			problems = append(problems, fmt.Sprintf("does not match regex %q", rule.Regex))
		}
	}

	// Allowed values
	if len(rule.Allowed) > 0 {
		found := false
		for _, a := range rule.Allowed {
			if value == a {
				found = true
				break
			}
		}
		if !found {
			problems = append(problems, fmt.Sprintf("value not in allowed set [%s]", strings.Join(rule.Allowed, ", ")))
		}
	}

	return problems
}

// checkType validates a string value against the expected VarType.
func checkType(value string, typ VarType) error {
	switch typ {
	case TypeString:
		return nil
	case TypeInt:
		if _, err := strconv.ParseInt(value, 10, 64); err != nil {
			return fmt.Errorf("not a valid int")
		}
	case TypeBool:
		if _, err := strconv.ParseBool(value); err != nil {
			return fmt.Errorf("not a valid bool")
		}
	case TypeDuration:
		if _, err := time.ParseDuration(value); err != nil {
			return fmt.Errorf("not a valid duration")
		}
	case TypeURL:
		u, err := url.Parse(value)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return fmt.Errorf("not a valid url")
		}
	case TypeJSON:
		if !json.Valid([]byte(value)) {
			return fmt.Errorf("not valid json")
		}
	}
	return nil
}

// detectWhitespace checks if a value has leading or trailing whitespace.
func detectWhitespace(value string) bool {
	return value != strings.TrimSpace(value)
}
