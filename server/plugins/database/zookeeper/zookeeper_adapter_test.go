package zookeeper

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
		{name: "missing host", ctx: database.Context{Connection: &database.Connection{Config: database.Config{Port: 2181}}}},
		{name: "missing port", ctx: database.Context{Connection: &database.Connection{Config: database.Config{Host: "localhost"}}}},
		{name: "placeholder host", ctx: database.Context{Connection: &database.Connection{Config: database.Config{Host: "-", Port: 2181}}}},
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
