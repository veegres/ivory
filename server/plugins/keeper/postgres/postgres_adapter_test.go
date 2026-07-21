package postgres

import (
	"errors"
	"ivory/plugins/keeper"
	"net/http"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestMapNode(t *testing.T) {
	tests := []struct {
		name         string
		host         string
		port         int
		inRecovery   bool
		lag          int64
		expectedRole keeper.Role
		expectedLag  int64
		expectedKey  string
	}{
		{
			name: "primary is leader with zero lag",
			host: "db1", port: 5432, inRecovery: false, lag: 123,
			expectedRole: keeper.Leader, expectedLag: 0, expectedKey: "db1:5432",
		},
		{
			name: "replica keeps lag",
			host: "db2", port: 5433, inRecovery: true, lag: 456,
			expectedRole: keeper.Replica, expectedLag: 456, expectedKey: "db2:5433",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := mapNode(tt.host, tt.port, tt.inRecovery, tt.lag)
			if response.Role != tt.expectedRole {
				t.Errorf("expected role %v, got %v", tt.expectedRole, response.Role)
			}
			if response.Lag != tt.expectedLag {
				t.Errorf("expected lag %d, got %d", tt.expectedLag, response.Lag)
			}
			if response.Key == nil || *response.Key != tt.expectedKey {
				t.Errorf("expected key %q, got %v", tt.expectedKey, response.Key)
			}
			if response.State != keeper.StateRunning {
				t.Errorf("expected state running, got %q", response.State)
			}
			if response.Status == nil || *response.Status != keeper.Active {
				t.Errorf("expected active status, got %v", response.Status)
			}
			if response.DiscoveredHost == nil || *response.DiscoveredHost != tt.host {
				t.Errorf("expected discovered host %q, got %v", tt.host, response.DiscoveredHost)
			}
			if response.DiscoveredKeeperPort == nil || *response.DiscoveredKeeperPort != tt.port {
				t.Errorf("expected discovered keeper port %d, got %v", tt.port, response.DiscoveredKeeperPort)
			}
			if response.DiscoveredDbPort == nil || *response.DiscoveredDbPort != tt.port {
				t.Errorf("expected discovered db port %d, got %v", tt.port, response.DiscoveredDbPort)
			}
		})
	}
}

func TestMapSyncStandby(t *testing.T) {
	tests := []struct {
		name         string
		syncState    string
		expectedSync bool
	}{
		{"sync standby", "sync", true},
		{"quorum standby", "quorum", true},
		{"async standby", "async", false},
		{"potential standby", "potential", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := mapSyncStandby("10.0.0.2", 5432, tt.syncState)
			if response.Role != keeper.Replica {
				t.Errorf("expected role replica, got %v", response.Role)
			}
			if response.Sync != tt.expectedSync {
				t.Errorf("expected sync %v, got %v", tt.expectedSync, response.Sync)
			}
			if response.State != keeper.StateRunning {
				t.Errorf("expected state running, got %q", response.State)
			}
			if response.DiscoveredHost == nil || *response.DiscoveredHost != "10.0.0.2" {
				t.Errorf("expected discovered host 10.0.0.2, got %v", response.DiscoveredHost)
			}
			if response.DiscoveredKeeperPort == nil || *response.DiscoveredKeeperPort != 5432 {
				t.Errorf("expected discovered keeper port 5432 (reused from the primary's own connection port), got %v", response.DiscoveredKeeperPort)
			}
			if response.DiscoveredDbPort == nil || *response.DiscoveredDbPort != 5432 {
				t.Errorf("expected discovered db port 5432, got %v", response.DiscoveredDbPort)
			}
		})
	}
}

func TestMapUnavailableNode(t *testing.T) {
	response := mapUnavailableNode("db1", 5432, keeper.StateStarting)

	if response.State != keeper.StateStarting {
		t.Errorf("expected state %q, got %q", keeper.StateStarting, response.State)
	}
	if response.Role != keeper.Unknown {
		t.Errorf("expected role unknown, got %v", response.Role)
	}
	if response.Status == nil || *response.Status != keeper.Active {
		t.Errorf("expected active status, got %v", response.Status)
	}
	if response.DiscoveredHost == nil || *response.DiscoveredHost != "db1" {
		t.Errorf("expected discovered host db1, got %v", response.DiscoveredHost)
	}
	if response.DiscoveredKeeperPort == nil || *response.DiscoveredKeeperPort != 5432 {
		t.Errorf("expected discovered keeper port 5432, got %v", response.DiscoveredKeeperPort)
	}
}

func TestMapUnavailableState(t *testing.T) {
	tests := []struct {
		name          string
		err           error
		expectedState keeper.State
		expectedOk    bool
	}{
		{"nil error", nil, "", false},
		{"unrelated error", errors.New("boom"), "", false},
		{"wrong sqlstate", &pgconn.PgError{Code: "28000"}, "", false},
		{
			"starting up sqlstate",
			&pgconn.PgError{Code: sqlStateCannotConnectNow, Message: "the database system is starting up"},
			keeper.StateStarting, true,
		},
		{
			"shutting down sqlstate",
			&pgconn.PgError{Code: sqlStateCannotConnectNow, Message: "the database system is shutting down"},
			keeper.StateStopping, true,
		},
		{
			"recovery mode sqlstate defaults to starting",
			&pgconn.PgError{Code: sqlStateCannotConnectNow, Message: "the database system is in recovery mode"},
			keeper.StateStarting, true,
		},
		{
			"wrapped starting up sqlstate",
			errors.Join(errors.New("failed to connect"), &pgconn.PgError{Code: sqlStateCannotConnectNow, Message: "the database system is starting up"}),
			keeper.StateStarting, true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state, ok := mapUnavailableState(tt.err)
			if ok != tt.expectedOk {
				t.Errorf("expected ok %v, got %v", tt.expectedOk, ok)
			}
			if state != tt.expectedState {
				t.Errorf("expected state %q, got %q", tt.expectedState, state)
			}
		})
	}
}

func TestListRequiresCredentials(t *testing.T) {
	adapter := NewAdapter()
	request := keeper.Request{Host: "localhost", Port: 5432}

	_, status, err := adapter.List(request)

	if status != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", status)
	}
	if !errors.Is(err, ErrCredentialsRequired) {
		t.Errorf("expected ErrCredentialsRequired, got %v", err)
	}
}

func TestReloadRequiresCredentials(t *testing.T) {
	adapter := NewAdapter()
	request := keeper.Request{Host: "localhost", Port: 5432}

	_, status, err := adapter.Reload(request)

	if status != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", status)
	}
	if !errors.Is(err, ErrCredentialsRequired) {
		t.Errorf("expected ErrCredentialsRequired, got %v", err)
	}
}

func TestFailoverRequiresCredentials(t *testing.T) {
	adapter := NewAdapter()
	request := keeper.Request{Host: "localhost", Port: 5432}

	_, status, err := adapter.Failover(request)

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
