package mongo

import (
	"errors"
	"ivory/plugins/database"
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestCancelNotSupported(t *testing.T) {
	adapter := NewAdapter()
	if err := adapter.Cancel(database.Context{}, 1); !errors.Is(err, ErrNotSupported) {
		t.Errorf("expected ErrNotSupported, got %v", err)
	}
}

func TestStringField(t *testing.T) {
	if got := stringField("hello"); got != "hello" {
		t.Errorf("expected hello, got %q", got)
	}
	if got := stringField(int32(1)); got != "" {
		t.Errorf("expected empty string for non-string value, got %q", got)
	}
	if got := stringField(nil); got != "" {
		t.Errorf("expected empty string for nil, got %q", got)
	}
}

func TestInt64Field(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		expected int64
	}{
		{"int32", int32(5), 5},
		{"int64", int64(9), 9},
		{"float64", float64(3), 3},
		{"string opid not supported", "shard1:12345", 0},
		{"nil", nil, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := int64Field(tt.value); got != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, got)
			}
		})
	}
}

func TestParseCurrentOp(t *testing.T) {
	inProg := []bson.M{
		{"opid": int32(42), "ns": "test.users", "op": "query", "secs_running": int32(3), "client": "127.0.0.1:1234", "desc": "conn5"},
	}
	fields, rows := parseCurrentOp(inProg)

	if len(fields) != 6 || fields[0].Name != "pid" {
		t.Fatalf("unexpected fields: %+v", fields)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	expected := []any{int64(42), "test.users", "query", int64(3), "127.0.0.1:1234", "conn5"}
	for i := range expected {
		if rows[0][i] != expected[i] {
			t.Errorf("expected row[%d]=%v, got %v", i, expected[i], rows[0][i])
		}
	}
}

func TestTerminateConnectFailure(t *testing.T) {
	adapter := NewAdapter()
	ctx := database.Context{Connection: &database.Connection{Config: database.Config{}}}

	if err := adapter.Terminate(ctx, 1); !errors.Is(err, database.ErrDatabaseHostOrPortNotSpecified) {
		t.Errorf("expected ErrDatabaseHostOrPortNotSpecified, got %v", err)
	}
}
