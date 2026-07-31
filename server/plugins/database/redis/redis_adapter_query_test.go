package redis

import (
	"errors"
	"ivory/plugins/database"
	"slices"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    []string
		expectedErr error
	}{
		{name: "simple command", input: "get foo", expected: []string{"get", "foo"}},
		{name: "quoted value with spaces", input: `set greeting "hello world"`, expected: []string{"set", "greeting", "hello world"}},
		{name: "single quoted json", input: `set config '{"a": 1}'`, expected: []string{"set", "config", `{"a": 1}`}},
		{name: "extra whitespace tolerated", input: "  get   foo  \n", expected: []string{"get", "foo"}},
		{name: "empty input", input: "   ", expected: []string{}},
		{name: "unterminated quote", input: `get "foo`, expectedErr: ErrUnterminatedQuote},
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
			if !slices.Equal(tokens, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, tokens)
			}
		})
	}
}

func TestFormatValue(t *testing.T) {
	tests := []struct {
		name             string
		value            any
		expectedDataType string
		expectedValue    any
	}{
		{name: "nil", value: nil, expectedDataType: "text", expectedValue: "(nil)"},
		{name: "int64", value: int64(42), expectedDataType: "int8", expectedValue: int64(42)},
		{name: "float64", value: float64(1.5), expectedDataType: "float8", expectedValue: float64(1.5)},
		{name: "bool", value: true, expectedDataType: "bool", expectedValue: true},
		{name: "string", value: "hello", expectedDataType: "text", expectedValue: "hello"},
		{name: "bytes", value: []byte("hello"), expectedDataType: "text", expectedValue: "hello"},
		{name: "nested array", value: []interface{}{"a", int64(1)}, expectedDataType: "text", expectedValue: "[a, 1]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataType, value := formatValue(tt.value)
			if dataType != tt.expectedDataType {
				t.Errorf("expected data type %q, got %q", tt.expectedDataType, dataType)
			}
			if value != tt.expectedValue {
				t.Errorf("expected value %v, got %v", tt.expectedValue, value)
			}
		})
	}
}

func TestFormatReply(t *testing.T) {
	t.Run("array reply becomes index/value rows", func(t *testing.T) {
		fields, rows := formatReply([]interface{}{"a", "b", int64(3)})
		if len(fields) != 2 || fields[0].Name != "index" || fields[1].Name != "value" {
			t.Fatalf("unexpected fields: %+v", fields)
		}
		if len(rows) != 3 {
			t.Fatalf("expected 3 rows, got %d", len(rows))
		}
		if rows[2][0] != int64(2) || rows[2][1] != "3" {
			t.Errorf("unexpected row: %v", rows[2])
		}
	})

	t.Run("scalar reply becomes a single result row", func(t *testing.T) {
		fields, rows := formatReply("OK")
		if len(fields) != 1 || fields[0].Name != "result" {
			t.Fatalf("unexpected fields: %+v", fields)
		}
		if len(rows) != 1 || rows[0][0] != "OK" {
			t.Fatalf("unexpected rows: %v", rows)
		}
	})
}

func TestGetManyGetOneNotSupported(t *testing.T) {
	adapter := NewAdapter()
	ctx := database.Context{}

	if _, err := adapter.GetMany(ctx, "get a", nil); !errors.Is(err, ErrNotSupported) {
		t.Errorf("expected ErrNotSupported, got %v", err)
	}
	if _, err := adapter.GetOne(ctx, "get a"); !errors.Is(err, ErrNotSupported) {
		t.Errorf("expected ErrNotSupported, got %v", err)
	}
}

func TestGetFieldsEmptyCommand(t *testing.T) {
	adapter := NewAdapter()
	ctx := database.Context{}

	if _, err := adapter.GetFields(ctx, "   ", nil); !errors.Is(err, ErrEmptyCommand) {
		t.Errorf("expected ErrEmptyCommand, got %v", err)
	}
}
