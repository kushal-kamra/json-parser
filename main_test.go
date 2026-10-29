package main

import "testing"

func TestValidateJSON(t *testing.T) {
	tests := []struct {
		input   string
		isValid bool
	}{
		{input: "{}", isValid: true},
		{input: " \t{\n\r} ", isValid: true},
		{input: "", isValid: false},
		{input: "[]", isValid: false},
		{input: "{ x }", isValid: false},
		{input: "{}{}", isValid: false},
		{input: "{} trailing", isValid: false},
		{input: "{\u00a0}", isValid: false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			isValid, errorMessage := ValidateJSON(tt.input)
			if isValid != tt.isValid {
				t.Errorf("ValidateJSON(%q) validity = %t, want %t", tt.input, isValid, tt.isValid)
			}
			if tt.isValid && errorMessage != "" {
				t.Errorf("ValidateJSON(%q) error message = %q, want empty", tt.input, errorMessage)
			}
			if !tt.isValid && errorMessage == "" {
				t.Errorf("ValidateJSON(%q) returned no error message for invalid input", tt.input)
			}
		})
	}
}
