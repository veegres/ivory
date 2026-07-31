package zookeeper

import (
	"errors"
	"ivory/plugins/database"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    []string
		expectedErr error
	}{
		{name: "simple command", input: "ls /service", expected: []string{"ls", "/service"}},
		{name: "quoted value with spaces", input: `create /motd "hello world"`, expected: []string{"create", "/motd", "hello world"}},
		{name: "extra whitespace tolerated", input: "  get   /a  \n", expected: []string{"get", "/a"}},
		{name: "empty input", input: "   ", expected: []string{}},
		{name: "unterminated quote", input: `get "/a`, expectedErr: ErrUnterminatedQuote},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := tokenize(tt.input)
			if tt.expectedErr != nil {
				if !errors.Is(err, tt.expectedErr) {
					t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if len(tokens) != len(tt.expected) {
				t.Fatalf("expected %v, got %v", tt.expected, tokens)
			}
			for i := range tokens {
				if tokens[i] != tt.expected[i] {
					t.Errorf("expected %v, got %v", tt.expected, tokens)
				}
			}
		})
	}
}

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		expected    command
		expectedErr error
	}{
		{name: "ls", query: "ls /service", expected: command{Verb: verbLs, Path: "/service"}},
		{name: "get", query: "get /service/config", expected: command{Verb: verbGet, Path: "/service/config"}},
		{name: "exists", query: "exists /service", expected: command{Verb: verbExists, Path: "/service"}},
		{name: "delete", query: "delete /service", expected: command{Verb: verbDelete, Path: "/service"}},
		{name: "create", query: `create /motd "hello world"`, expected: command{Verb: verbCreate, Path: "/motd", Data: "hello world"}},
		{name: "set", query: "set /service/config enabled", expected: command{Verb: verbSet, Path: "/service/config", Data: "enabled"}},
		{name: "empty query", query: "   ", expectedErr: ErrEmptyCommand},
		{name: "unknown command", query: "watch /a", expectedErr: ErrUnknownCommand},
		{name: "ls missing path", query: "ls", expectedErr: ErrMissingArgument},
		{name: "create missing data", query: "create /a", expectedErr: ErrMissingArgument},
		{name: "get extra argument", query: "get /a extra", expectedErr: ErrUnexpectedArgument},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := parseCommand(tt.query)
			if tt.expectedErr != nil {
				if !errors.Is(err, tt.expectedErr) {
					t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if *parsed != tt.expected {
				t.Errorf("expected %+v, got %+v", tt.expected, *parsed)
			}
		})
	}
}

func TestGetManyGetOneNotSupported(t *testing.T) {
	adapter := NewAdapter()
	ctx := database.Context{}

	if _, err := adapter.GetMany(ctx, "ls /", nil); !errors.Is(err, ErrNotSupported) {
		t.Errorf("expected ErrNotSupported, got %v", err)
	}
	if _, err := adapter.GetOne(ctx, "ls /"); !errors.Is(err, ErrNotSupported) {
		t.Errorf("expected ErrNotSupported, got %v", err)
	}
}
