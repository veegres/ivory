package mongo

import (
	"errors"
	"ivory/plugins/database"
	"testing"
)

func TestFilterNames(t *testing.T) {
	names := []string{"users", "orders", "user_sessions"}

	tests := []struct {
		name     string
		filter   string
		expected []string
	}{
		{"empty filter returns all", "", names},
		{"filter matches substring", "user", []string{"users", "user_sessions"}},
		{"filter matches nothing", "zzz", []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterNames(names, tt.filter)
			if len(result) != len(tt.expected) {
				t.Fatalf("expected %v, got %v", tt.expected, result)
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestListSchemasNotSupported(t *testing.T) {
	adapter := NewAdapter()
	if _, err := adapter.ListSchemas(database.Context{}, ""); !errors.Is(err, ErrNotSupported) {
		t.Errorf("expected ErrNotSupported, got %v", err)
	}
}

func TestListDatabasesConnectFailure(t *testing.T) {
	adapter := NewAdapter()
	ctx := database.Context{Connection: &database.Connection{Config: database.Config{}}}

	if _, err := adapter.ListDatabases(ctx, ""); !errors.Is(err, database.ErrDatabaseHostOrPortNotSpecified) {
		t.Errorf("expected ErrDatabaseHostOrPortNotSpecified, got %v", err)
	}
}
