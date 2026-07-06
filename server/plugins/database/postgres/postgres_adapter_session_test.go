package postgres

import (
	"errors"
	"ivory/plugins/database"
	"testing"
)

func TestSessionPropagatesConnectionErrors(t *testing.T) {
	adapter := NewAdapter()
	ctx := noConnectionContext()

	type op struct {
		name string
		call func() error
	}
	ops := []op{
		{"Cancel", func() error { return adapter.Cancel(ctx, 1) }},
		{"Terminate", func() error { return adapter.Terminate(ctx, 1) }},
		{"ActiveQueries", func() error { _, e := adapter.ActiveQueries(ctx, nil); return e }},
	}

	for _, o := range ops {
		t.Run(o.name, func(t *testing.T) {
			if err := o.call(); !errors.Is(err, database.ErrDatabaseHostOrPortNotSpecified) {
				t.Errorf("expected ErrDatabaseHostOrPortNotSpecified, got %v", err)
			}
		})
	}
}
