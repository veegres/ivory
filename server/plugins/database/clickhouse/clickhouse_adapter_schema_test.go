package clickhouse

import (
	"errors"
	"ivory/plugins/database"
	"testing"
)

func TestListSchemasNotSupported(t *testing.T) {
	adapter := NewAdapter()
	if _, err := adapter.ListSchemas(database.Context{}, ""); !errors.Is(err, ErrNotSupported) {
		t.Errorf("expected ErrNotSupported, got %v", err)
	}
}
