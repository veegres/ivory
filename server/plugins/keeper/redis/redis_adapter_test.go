package redis

import (
	"errors"
	"ivory/plugins/keeper"
	"net/http"
	"testing"
)

func TestParseInfo(t *testing.T) {
	info := "# Replication\r\nrole:slave\r\nmaster_host:10.0.0.1\r\nmaster_last_io_seconds_ago:3\r\n\r\n# Other\r\nignored\r\n"
	fields := parseInfo(info)

	if fields["role"] != "slave" {
		t.Errorf("expected role slave, got %q", fields["role"])
	}
	if fields["master_host"] != "10.0.0.1" {
		t.Errorf("expected master_host 10.0.0.1, got %q", fields["master_host"])
	}
	if fields["master_last_io_seconds_ago"] != "3" {
		t.Errorf("expected master_last_io_seconds_ago 3, got %q", fields["master_last_io_seconds_ago"])
	}
	if _, ok := fields["ignored"]; ok {
		t.Errorf("expected malformed line without a colon to be skipped, got %v", fields)
	}
}

func TestMapNode(t *testing.T) {
	tests := []struct {
		name         string
		fields       map[string]string
		expectedRole keeper.Role
		expectedLag  int64
	}{
		{
			name:         "master is leader with zero lag",
			fields:       map[string]string{"role": "master"},
			expectedRole: keeper.Leader, expectedLag: 0,
		},
		{
			name:         "slave keeps lag from last io seconds",
			fields:       map[string]string{"role": "slave", "master_last_io_seconds_ago": "7"},
			expectedRole: keeper.Replica, expectedLag: 7,
		},
		{
			name:         "slave with unparsable lag defaults to zero",
			fields:       map[string]string{"role": "slave", "master_last_io_seconds_ago": "n/a"},
			expectedRole: keeper.Replica, expectedLag: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := mapNode("db1", 6379, tt.fields)
			if response.Role != tt.expectedRole {
				t.Errorf("expected role %v, got %v", tt.expectedRole, response.Role)
			}
			if response.Lag != tt.expectedLag {
				t.Errorf("expected lag %d, got %d", tt.expectedLag, response.Lag)
			}
			if response.Key == nil || *response.Key != "db1:6379" {
				t.Errorf("expected key db1:6379, got %v", response.Key)
			}
			if response.State != keeper.StateRunning {
				t.Errorf("expected state running, got %q", response.State)
			}
			if response.Status == nil || *response.Status != keeper.Active {
				t.Errorf("expected active status, got %v", response.Status)
			}
			if response.DiscoveredHost == nil || *response.DiscoveredHost != "db1" {
				t.Errorf("expected discovered host db1, got %v", response.DiscoveredHost)
			}
			if response.DiscoveredKeeperPort == nil || *response.DiscoveredKeeperPort != 6379 {
				t.Errorf("expected discovered keeper port 6379, got %v", response.DiscoveredKeeperPort)
			}
			if response.DiscoveredDbPort == nil || *response.DiscoveredDbPort != 6379 {
				t.Errorf("expected discovered db port 6379, got %v", response.DiscoveredDbPort)
			}
		})
	}
}

func TestListRequiresCredentials(t *testing.T) {
	adapter := NewAdapter()
	request := keeper.Request{Host: "localhost", Port: 6379}

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
	request := keeper.Request{Host: "localhost", Port: 6379}

	_, status, err := adapter.Config(request)

	if status != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", status)
	}
	if !errors.Is(err, ErrCredentialsRequired) {
		t.Errorf("expected ErrCredentialsRequired, got %v", err)
	}
}

func TestFailoverRequiresCredentials(t *testing.T) {
	adapter := NewAdapter()
	request := keeper.Request{Host: "localhost", Port: 6379}

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
		{"Reload", func() (int, error) { _, s, e := adapter.Reload(request); return s, e }},
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
