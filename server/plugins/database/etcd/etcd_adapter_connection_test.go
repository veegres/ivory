package etcd

import (
	"errors"
	"ivory/plugins/database"
	"testing"
)

func TestSessionNotSupportedOperations(t *testing.T) {
	adapter := NewAdapter()
	ctx := database.Context{}

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
			if err := o.call(); !errors.Is(err, ErrNotSupported) {
				t.Errorf("expected ErrNotSupported, got %v", err)
			}
		})
	}
}

func TestAdapterPingRequiresAnEndpoint(t *testing.T) {
	adapter := NewAdapter()
	tests := []struct {
		name string
		ctx  database.Context
	}{
		{name: "missing host", ctx: database.Context{Connection: &database.Connection{Config: database.Config{Port: 2379}}}},
		{name: "missing port", ctx: database.Context{Connection: &database.Connection{Config: database.Config{Host: "localhost"}}}},
		{name: "placeholder host", ctx: database.Context{Connection: &database.Connection{Config: database.Config{Host: "-", Port: 2379}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := adapter.Ping(tt.ctx); !errors.Is(err, database.ErrDatabaseHostOrPortNotSpecified) {
				t.Fatalf("expected ErrDatabaseHostOrPortNotSpecified, got %v", err)
			}
		})
	}
}
