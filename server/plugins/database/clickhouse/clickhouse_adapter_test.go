package clickhouse

import (
	"errors"
	"ivory/plugins/database"
	"testing"
)

func TestConnectRequiresHostAndPort(t *testing.T) {
	tests := []struct {
		name string
		ctx  database.Context
	}{
		{name: "missing host", ctx: database.Context{Connection: &database.Connection{Config: database.Config{Port: 9000}}}},
		{name: "missing port", ctx: database.Context{Connection: &database.Connection{Config: database.Config{Host: "localhost"}}}},
		{name: "placeholder host", ctx: database.Context{Connection: &database.Connection{Config: database.Config{Host: "-", Port: 9000}}}},
	}

	adapter := NewAdapter()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := adapter.connect(tt.ctx)
			if !errors.Is(err, database.ErrDatabaseHostOrPortNotSpecified) {
				t.Errorf("expected ErrDatabaseHostOrPortNotSpecified, got %v", err)
			}
		})
	}
}

func TestConnectRequiresCredentials(t *testing.T) {
	adapter := NewAdapter()
	ctx := database.Context{Connection: &database.Connection{Config: database.Config{Host: "localhost", Port: 9000}}}

	_, _, err := adapter.connect(ctx)
	if !errors.Is(err, database.ErrPasswordNotSet) {
		t.Errorf("expected ErrPasswordNotSet, got %v", err)
	}
}

func TestParseQuery(t *testing.T) {
	adapter := NewAdapter()

	tests := []struct {
		name     string
		query    string
		expected database.QueryAnalysis
	}{
		{
			name:     "simple select",
			query:    "SELECT * FROM system.tables;",
			expected: database.QueryAnalysis{SELECT: 1, FROM: 1, Semicolon: true},
		},
		{
			name:     "select with existing limit",
			query:    "SELECT * FROM system.tables LIMIT 10",
			expected: database.QueryAnalysis{SELECT: 1, FROM: 1, LIMIT: 1},
		},
		{
			name:     "alter counts as update",
			query:    "ALTER TABLE t DELETE WHERE 1=1",
			expected: database.QueryAnalysis{UPDATE: 1, DELETE: 1},
		},
		{
			name:     "truncate counts as delete",
			query:    "TRUNCATE TABLE t",
			expected: database.QueryAnalysis{DELETE: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.parseQuery(tt.query)
			if result != tt.expected {
				t.Errorf("expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}

func TestAddLimitToQuery(t *testing.T) {
	adapter := NewAdapter()

	t.Run("plain select gets a limit appended", func(t *testing.T) {
		analysis := database.QueryAnalysis{SELECT: 1, FROM: 1}
		query, limit := adapter.addLimitToQuery("SELECT * FROM t", analysis, "100")
		if query != "SELECT * FROM t LIMIT 100;" {
			t.Errorf("unexpected query: %q", query)
		}
		if limit == nil || *limit != "100" {
			t.Errorf("expected limit 100, got %v", limit)
		}
	})

	t.Run("select with existing limit is left alone", func(t *testing.T) {
		analysis := database.QueryAnalysis{SELECT: 1, FROM: 1, LIMIT: 1}
		query, limit := adapter.addLimitToQuery("SELECT * FROM t LIMIT 5", analysis, "100")
		if query != "SELECT * FROM t LIMIT 5" {
			t.Errorf("unexpected query: %q", query)
		}
		if limit != nil {
			t.Errorf("expected no limit override, got %v", *limit)
		}
	})

	t.Run("mutation is left alone", func(t *testing.T) {
		analysis := database.QueryAnalysis{INSERT: 1}
		query, limit := adapter.addLimitToQuery("INSERT INTO t VALUES (1)", analysis, "100")
		if query != "INSERT INTO t VALUES (1)" {
			t.Errorf("unexpected query: %q", query)
		}
		if limit != nil {
			t.Errorf("expected no limit override, got %v", *limit)
		}
	})
}

func TestNormalizeQueryRequiresTrimForLimit(t *testing.T) {
	adapter := NewAdapter()
	limit := "10"

	_, _, err := adapter.normalizeQuery("SELECT 1", nil, &limit)
	if !errors.Is(err, database.ErrCannotLimitWithoutTrim) {
		t.Errorf("expected ErrCannotLimitWithoutTrim, got %v", err)
	}
}
