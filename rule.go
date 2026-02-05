package envdoc

import (
	"fmt"
	"regexp"

	"gopkg.in/yaml.v3"
)

// VarType represents the expected type of an environment variable value.
type VarType string

const (
	TypeString   VarType = "string"
	TypeInt      VarType = "int"
	TypeBool     VarType = "bool"
	TypeDuration VarType = "duration"
	TypeURL      VarType = "url"
	TypeJSON     VarType = "json"
)

var validTypes = map[VarType]bool{
	TypeString: true, TypeInt: true, TypeBool: true,
	TypeDuration: true, TypeURL: true, TypeJSON: true,
}

// Rule defines validation for a single environment variable.
type Rule struct {
	Key         string   `yaml:"key"`
	Required    bool     `yaml:"required"`
	Type        VarType  `yaml:"type"`
	MinLen      *int     `yaml:"min_len,omitempty"`
	MaxLen      *int     `yaml:"max_len,omitempty"`
	Regex       string   `yaml:"regex,omitempty"`
	Allowed     []string `yaml:"allowed,omitempty"`
	Secret      *bool    `yaml:"secret,omitempty"`
	Fingerprint *bool    `yaml:"fingerprint,omitempty"`
}

// RuleSet is the top-level YAML structure.
type RuleSet struct {
	Rules []Rule `yaml:"rules"`
}

// LoadRules parses YAML bytes into a slice of Rules.
func LoadRules(data []byte) ([]Rule, error) {
	var rs RuleSet
	if err := yaml.Unmarshal(data, &rs); err != nil {
		return nil, fmt.Errorf("envdoc: parsing rules: %w", err)
	}
	if err := validateRules(rs.Rules); err != nil {
		return nil, err
	}
	return rs.Rules, nil
}

// LoadRulesFile reads and parses a YAML rules file.
func LoadRulesFile(path string) ([]Rule, error) {
	data, err := readFile(path)
	if err != nil {
		return nil, fmt.Errorf("envdoc: reading rules file: %w", err)
	}
	return LoadRules(data)
}

// validateRules checks rules for duplicate keys, unknown types, invalid regex, and min>max.
func validateRules(rules []Rule) error {
	seen := make(map[string]bool)
	for idx, r := range rules {
		if r.Key == "" {
			return fmt.Errorf("envdoc: rule[%d]: key is required", idx)
		}
		if seen[r.Key] {
			return fmt.Errorf("envdoc: rule[%d]: duplicate key %q", idx, r.Key)
		}
		seen[r.Key] = true

		if r.Type != "" && !validTypes[r.Type] {
			return fmt.Errorf("envdoc: rule[%d] (%s): unknown type %q", idx, r.Key, r.Type)
		}

		if r.Regex != "" {
			if _, err := regexp.Compile(r.Regex); err != nil {
				return fmt.Errorf("envdoc: rule[%d] (%s): invalid regex: %w", idx, r.Key, err)
			}
		}

		if r.MinLen != nil && r.MaxLen != nil && *r.MinLen > *r.MaxLen {
			return fmt.Errorf("envdoc: rule[%d] (%s): min_len (%d) > max_len (%d)", idx, r.Key, *r.MinLen, *r.MaxLen)
		}
	}
	return nil
}

// readFile is a helper to read a file. Separated for testability.
var readFile = readFileOS
