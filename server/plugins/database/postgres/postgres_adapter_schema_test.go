package postgres

import (
	"errors"
	"ivory/plugins/database"
	"testing"
)

func TestSchemaPropagatesConnectionErrors(t *testing.T) {
	adapter := NewAdapter()
	ctx := noConnectionContext()

	type op struct {
		name string
		call func() error
	}
	ops := []op{
		{"ListDatabases", func() error { _, e := adapter.ListDatabases(ctx, ""); return e }},
		{"ListSchemas", func() error { _, e := adapter.ListSchemas(ctx, ""); return e }},
		{"ListTables", func() error { _, e := adapter.ListTables(ctx, "public", ""); return e }},
	}

	for _, o := range ops {
		t.Run(o.name, func(t *testing.T) {
			if err := o.call(); !errors.Is(err, database.ErrDatabaseHostOrPortNotSpecified) {
				t.Errorf("expected ErrDatabaseHostOrPortNotSpecified, got %v", err)
			}
		})
	}
}
