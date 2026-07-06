package etcd

import (
	"errors"
	"ivory/plugins/database"
	"testing"
)

func TestSchemaNotSupportedOperations(t *testing.T) {
	adapter := NewAdapter()
	ctx := database.Context{}

	type op struct {
		name string
		call func() error
	}
	ops := []op{
		{"ListDatabases", func() error { _, e := adapter.ListDatabases(ctx, ""); return e }},
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
