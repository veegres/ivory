package clickhouse

import (
	"errors"
	"ivory/plugins/keeper"
	"net/http"
	"testing"
)

func TestMapNode(t *testing.T) {
	tests := []struct {
		name          string
		readonly      bool
		absoluteDelay uint64
		expectedState keeper.State
		expectedRole  keeper.Role
		expectedLag   int64
	}{
		{name: "healthy node is a running replica", readonly: false, absoluteDelay: 0, expectedState: keeper.StateRunning, expectedRole: keeper.Replica, expectedLag: 0},
		// NOTE: still a replica - the node answered, so what it is stays
		// known; only its state changes
		{name: "readonly node stays a replica", readonly: true, absoluteDelay: 12, expectedState: keeper.StateStopping, expectedRole: keeper.Replica, expectedLag: 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := mapNode("ch1", 9000, tt.readonly, tt.absoluteDelay)
			if response.State != tt.expectedState {
				t.Errorf("expected state %q, got %q", tt.expectedState, response.State)
			}
			if response.Lag != tt.expectedLag {
				t.Errorf("expected lag %d, got %d", tt.expectedLag, response.Lag)
			}
			if response.Role != tt.expectedRole {
				t.Errorf("expected role %q, got %q", tt.expectedRole, response.Role)
			}
			if response.Sync {
				t.Error("clickhouse has no synchronous-replica set, so Sync must stay false")
			}
			if response.Key == nil || *response.Key != "ch1:9000" {
				t.Errorf("expected key ch1:9000, got %v", response.Key)
			}
			if response.Status == nil || *response.Status != keeper.Active {
				t.Errorf("expected active status, got %v", response.Status)
			}
			if response.DiscoveredHost == nil || *response.DiscoveredHost != "ch1" {
				t.Errorf("expected discovered host ch1, got %v", response.DiscoveredHost)
			}
			if response.DiscoveredKeeperPort == nil || *response.DiscoveredKeeperPort != 9000 {
				t.Errorf("expected discovered keeper port 9000, got %v", response.DiscoveredKeeperPort)
			}
			if response.DiscoveredDbPort == nil || *response.DiscoveredDbPort != 9000 {
				t.Errorf("expected discovered db port 9000, got %v", response.DiscoveredDbPort)
			}
		})
	}
}

func TestListRequiresCredentials(t *testing.T) {
	adapter := NewAdapter()
	request := keeper.Request{Host: "localhost", Port: 9000}

	_, status, err := adapter.List(request)

	if status != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", status)
	}
	if !errors.Is(err, ErrCredentialsRequired) {
		t.Errorf("expected ErrCredentialsRequired, got %v", err)
	}
}

func TestConfigRequiresCredentials(t *testing.T) {
	adapter := NewAdapter()
	request := keeper.Request{Host: "localhost", Port: 9000}

	_, status, err := adapter.Config(request)

	if status != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", status)
	}
	if !errors.Is(err, ErrCredentialsRequired) {
		t.Errorf("expected ErrCredentialsRequired, got %v", err)
	}
}

func TestReloadRequiresCredentials(t *testing.T) {
	adapter := NewAdapter()
	request := keeper.Request{Host: "localhost", Port: 9000}

	_, status, err := adapter.Reload(request)

	if status != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", status)
	}
	if !errors.Is(err, ErrCredentialsRequired) {
		t.Errorf("expected ErrCredentialsRequired, got %v", err)
	}
}

func TestUnsupportedOperations(t *testing.T) {
	adapter := NewAdapter()
	request := keeper.Request{}

	type op struct {
		name string
		call func() (int, error)
	}
	ops := []op{
		{"ConfigUpdate", func() (int, error) { _, s, e := adapter.ConfigUpdate(request); return s, e }},
		{"Switchover", func() (int, error) { _, s, e := adapter.Switchover(request); return s, e }},
		{"DeleteSwitchover", func() (int, error) { _, s, e := adapter.DeleteSwitchover(request); return s, e }},
		{"Reinitialize", func() (int, error) { _, s, e := adapter.Reinitialize(request); return s, e }},
		{"Restart", func() (int, error) { _, s, e := adapter.Restart(request); return s, e }},
		{"DeleteRestart", func() (int, error) { _, s, e := adapter.DeleteRestart(request); return s, e }},
		{"Failover", func() (int, error) { _, s, e := adapter.Failover(request); return s, e }},
		{"Activate", func() (int, error) { _, s, e := adapter.Activate(request); return s, e }},
		{"Pause", func() (int, error) { _, s, e := adapter.Pause(request); return s, e }},
	}

	for _, o := range ops {
		t.Run(o.name, func(t *testing.T) {
			status, err := o.call()
			if status != http.StatusNotImplemented {
				t.Errorf("expected status 501, got %d", status)
			}
			if !errors.Is(err, keeper.ErrNotSupported) {
				t.Errorf("expected ErrNotSupported, got %v", err)
			}
		})
	}
}
