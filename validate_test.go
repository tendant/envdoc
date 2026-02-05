package envdoc

import "testing"

func TestValidateVar_TypeString(t *testing.T) {
	problems := ValidateVar("hello", Rule{Type: TypeString})
	if len(problems) != 0 {
		t.Errorf("expected no problems, got %v", problems)
	}
}

func TestValidateVar_TypeInt(t *testing.T) {
	problems := ValidateVar("42", Rule{Type: TypeInt})
	if len(problems) != 0 {
		t.Errorf("expected no problems, got %v", problems)
	}

	problems = ValidateVar("notanint", Rule{Type: TypeInt})
	if len(problems) != 1 {
		t.Errorf("expected 1 problem, got %v", problems)
	}
}

func TestValidateVar_TypeBool(t *testing.T) {
	for _, v := range []string{"true", "false", "1", "0", "TRUE"} {
		problems := ValidateVar(v, Rule{Type: TypeBool})
		if len(problems) != 0 {
			t.Errorf("expected no problems for %q, got %v", v, problems)
		}
	}

	problems := ValidateVar("maybe", Rule{Type: TypeBool})
	if len(problems) != 1 {
		t.Errorf("expected 1 problem, got %v", problems)
	}
}

func TestValidateVar_TypeDuration(t *testing.T) {
	problems := ValidateVar("5s", Rule{Type: TypeDuration})
	if len(problems) != 0 {
		t.Errorf("expected no problems, got %v", problems)
	}

	problems = ValidateVar("notduration", Rule{Type: TypeDuration})
	if len(problems) != 1 {
		t.Errorf("expected 1 problem, got %v", problems)
	}
}

func TestValidateVar_TypeURL(t *testing.T) {
	problems := ValidateVar("https://example.com", Rule{Type: TypeURL})
	if len(problems) != 0 {
		t.Errorf("expected no problems, got %v", problems)
	}

	problems = ValidateVar("not-a-url", Rule{Type: TypeURL})
	if len(problems) != 1 {
		t.Errorf("expected 1 problem, got %v", problems)
	}
}

func TestValidateVar_TypeJSON(t *testing.T) {
	problems := ValidateVar(`{"key":"value"}`, Rule{Type: TypeJSON})
	if len(problems) != 0 {
		t.Errorf("expected no problems, got %v", problems)
	}

	problems = ValidateVar(`{invalid`, Rule{Type: TypeJSON})
	if len(problems) != 1 {
		t.Errorf("expected 1 problem, got %v", problems)
	}
}

func TestValidateVar_MinLen(t *testing.T) {
	min := 5
	problems := ValidateVar("hi", Rule{MinLen: &min})
	if len(problems) != 1 {
		t.Errorf("expected 1 problem, got %v", problems)
	}

	problems = ValidateVar("hello", Rule{MinLen: &min})
	if len(problems) != 0 {
		t.Errorf("expected no problems, got %v", problems)
	}
}

func TestValidateVar_MaxLen(t *testing.T) {
	max := 3
	problems := ValidateVar("hello", Rule{MaxLen: &max})
	if len(problems) != 1 {
		t.Errorf("expected 1 problem, got %v", problems)
	}

	problems = ValidateVar("hi", Rule{MaxLen: &max})
	if len(problems) != 0 {
		t.Errorf("expected no problems, got %v", problems)
	}
}

func TestValidateVar_Regex(t *testing.T) {
	problems := ValidateVar("abc123", Rule{Regex: `^[a-z]+\d+$`})
	if len(problems) != 0 {
		t.Errorf("expected no problems, got %v", problems)
	}

	problems = ValidateVar("ABC", Rule{Regex: `^[a-z]+$`})
	if len(problems) != 1 {
		t.Errorf("expected 1 problem, got %v", problems)
	}
}

func TestValidateVar_Allowed(t *testing.T) {
	problems := ValidateVar("prod", Rule{Allowed: []string{"dev", "staging", "prod"}})
	if len(problems) != 0 {
		t.Errorf("expected no problems, got %v", problems)
	}

	problems = ValidateVar("test", Rule{Allowed: []string{"dev", "staging", "prod"}})
	if len(problems) != 1 {
		t.Errorf("expected 1 problem, got %v", problems)
	}
}

func TestDetectWhitespace(t *testing.T) {
	if detectWhitespace("hello") {
		t.Error("expected false for no whitespace")
	}
	if !detectWhitespace(" hello") {
		t.Error("expected true for leading space")
	}
	if !detectWhitespace("hello ") {
		t.Error("expected true for trailing space")
	}
	if !detectWhitespace("\thello") {
		t.Error("expected true for leading tab")
	}
}
