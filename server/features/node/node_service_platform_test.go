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

func TestService_normalizeDatabaseOptions(t *testing.T) {
	s := &Service{}
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single line",
			input:    "--name test",
			expected: "--name test",
		},
		{
			name:     "multiple lines",
			input:    "--name test\n--restart always",
			expected: "--name test --restart always",
		},
		{
			name:     "multiple lines with spaces and tabs",
			input:    "--name test \n\t --restart always  \r\n -p 80:80",
			expected: "--name test --restart always -p 80:80",
		},
		{
			name:     "quoted strings with spaces",
			input:    "-e SCOPE=\"my cluster\"\n-e NAME=test",
			expected: "-e SCOPE=\"my cluster\" -e NAME=test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.normalizeDatabaseOptions(tt.input)
			if got != tt.expected {
				t.Errorf("normalizeDatabaseOptions() = %q, want %q", got, tt.expected)
			}
		})
	}
}
