package redis

import (
	"errors"
	"ivory/plugins/database"
	"slices"
	"testing"
)

func TestDatabaseCount(t *testing.T) {
	tests := []struct {
		name     string
		settings map[string]string
		expected int
	}{
		{name: "parses configured count", settings: map[string]string{"databases": "32"}, expected: 32},
		{name: "falls back to 16 when missing", settings: map[string]string{}, expected: 16},
		{name: "falls back to 16 when unparsable", settings: map[string]string{"databases": "abc"}, expected: 16},
		{name: "falls back to 16 when zero", settings: map[string]string{"databases": "0"}, expected: 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if count := databaseCount(tt.settings); count != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, count)
			}
		})
	}
}

func TestFilterDatabases(t *testing.T) {
	result := filterDatabases(16, "1")
	expected := []string{"1", "10", "11", "12", "13", "14", "15"}
	if !slices.Equal(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}

	all := filterDatabases(3, "")
	if !slices.Equal(all, []string{"0", "1", "2"}) {
		t.Errorf("expected all 3 databases, got %v", all)
	}
}

func TestSchemaNotSupportedOperations(t *testing.T) {
	adapter := NewAdapter()
	ctx := database.Context{}

	type op struct {
		name string
		call func() error
	}
	ops := []op{
		{"ListSchemas", func() error { _, e := adapter.ListSchemas(ctx, ""); return e }},
		{"ListTables", func() error { _, e := adapter.ListTables(ctx, "", ""); return e }},
	}

	for _, o := range ops {
		t.Run(o.name, func(t *testing.T) {
			if err := o.call(); !errors.Is(err, ErrNotSupported) {
				t.Errorf("expected ErrNotSupported, got %v", err)
			}
		})
	}
}
