package node

import (
	"testing"
)

func TestEscapeInterpolatedValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "plain value is unchanged", input: "postgres", expected: "postgres"},
		{name: "embedded single quote is escaped", input: "it's a test", expected: `it\'s a test`},
		{name: "embedded double quote is escaped", input: `pa"ss`, expected: `pa\"ss`},
		{name: "embedded backslash is escaped first", input: `pa\ss`, expected: `pa\\ss`},
		{name: "backslash and quote together", input: `pa\'ss`, expected: `pa\\\'ss`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := escapeInterpolatedValue(tt.input); got != tt.expected {
				t.Errorf("escapeInterpolatedValue(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
