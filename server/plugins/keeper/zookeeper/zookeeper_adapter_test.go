package zookeeper

import (
	"errors"
	"ivory/plugins/keeper"
	"net/http"
	"testing"
)

func TestParseLines(t *testing.T) {
	t.Run("tab separated mntr output", func(t *testing.T) {
		output := "zk_version\t3.9.2\r\nzk_server_state\tleader\r\nzk_znode_count\t5\r\n"
		fields := parseLines(output, "\t")

		if fields["zk_version"] != "3.9.2" {
			t.Errorf("expected zk_version 3.9.2, got %q", fields["zk_version"])
		}
		if fields["zk_server_state"] != "leader" {
			t.Errorf("expected zk_server_state leader, got %q", fields["zk_server_state"])
		}
	})

	t.Run("equals separated conf output skips headerless lines", func(t *testing.T) {
		output := "clientPort=2181\ndataDir=/data/version-2\nmembership: \nserver.1=zoo1:2888:3888:participant\n"
		fields := parseLines(output, "=")

		if fields["clientPort"] != "2181" {
			t.Errorf("expected clientPort 2181, got %q", fields["clientPort"])
		}
		if fields["dataDir"] != "/data/version-2" {
			t.Errorf("expected dataDir /data/version-2, got %q", fields["dataDir"])
		}
		if fields["server.1"] != "zoo1:2888:3888:participant" {
			t.Errorf("expected server.1 entry, got %q", fields["server.1"])
		}
		if _, ok := fields["membership: "]; ok {
			t.Errorf("expected the headerless membership line to be skipped, got %v", fields)
		}
	})

	t.Run("empty input", func(t *testing.T) {
		if fields := parseLines("", "="); len(fields) != 0 {
			t.Errorf("expected no fields, got %v", fields)
		}
	})
}

func TestMapNode(t *testing.T) {
	tests := []struct {
		name         string
		state        string
		expectedRole keeper.Role
	}{
		{"leader", "leader", keeper.Leader},
		{"standalone counts as leader", "standalone", keeper.Leader},
		{"follower", "follower", keeper.Replica},
		{"observer counts as replica", "observer", keeper.Replica},
		{"unknown state", "electing", keeper.Unknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := mapNode("zk1", 2181, tt.state)
			if response.Role != tt.expectedRole {
				t.Errorf("expected role %v, got %v", tt.expectedRole, response.Role)
			}
			if response.Lag != 0 {
				t.Errorf("expected lag 0, got %d", response.Lag)
			}
			if response.Key == nil || *response.Key != "zk1:2181" {
				t.Errorf("expected key zk1:2181, got %v", response.Key)
			}
			if response.State != keeper.StateRunning {
				t.Errorf("expected state running, got %q", response.State)
			}
			if response.Status == nil || *response.Status != keeper.Active {
				t.Errorf("expected active status, got %v", response.Status)
			}
			if response.DiscoveredHost == nil || *response.DiscoveredHost != "zk1" {
				t.Errorf("expected discovered host zk1, got %v", response.DiscoveredHost)
			}
			if response.DiscoveredKeeperPort == nil || *response.DiscoveredKeeperPort != 2181 {
				t.Errorf("expected discovered keeper port 2181, got %v", response.DiscoveredKeeperPort)
			}
			if response.DiscoveredDbPort == nil || *response.DiscoveredDbPort != 2181 {
				t.Errorf("expected discovered db port 2181, got %v", response.DiscoveredDbPort)
			}
		})
	}
}

func TestListRequiresReachableHost(t *testing.T) {
	adapter := NewPlugin()
	request := keeper.Request{Host: "127.0.0.1", Port: 1}

	_, status, err := adapter.List(request)

	if status != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", status)
	}
	if err == nil {
		t.Error("expected a connection error for an unreachable port")
	}
}

func TestUnsupportedOperations(t *testing.T) {
	adapter := NewPlugin()
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
