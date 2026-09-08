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
			response := mapNode("zk1", 2181, tt.state, map[string]string{})
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

func TestMapTags(t *testing.T) {
	tests := []struct {
		name     string
		state    string
		fields   map[string]string
		expected map[string]any
	}{
		{
			name:  "leader reports how much of its ensemble is synced",
			state: "leader",
			fields: map[string]string{
				"zk_version": "3.9.1-abc123, built on 2023-10-01", "zk_znode_count": "1204",
				"zk_num_alive_connections": "6", "zk_synced_followers": "2", "zk_followers": "2",
			},
			expected: map[string]any{"version": "3.9.1", "znodes": "1204", "connections": "6", "syncedFollowers": "2/2"},
		},
		{
			name:     "follower states no server state of its own",
			state:    "follower",
			fields:   map[string]string{"zk_version": "3.9.1", "zk_znode_count": "1204"},
			expected: map[string]any{"version": "3.9.1", "znodes": "1204"},
		},
		{
			name:     "observer is told apart from a voting follower",
			state:    "observer",
			fields:   map[string]string{"zk_znode_count": "1204"},
			expected: map[string]any{"serverState": "observer", "znodes": "1204"},
		},
		{
			name:     "standalone is told apart from an elected leader",
			state:    "standalone",
			fields:   map[string]string{},
			expected: map[string]any{"serverState": "standalone"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tags := mapTags(tt.state, tt.fields)
			if tags == nil {
				t.Fatal("expected tags, got nil")
			}
			if len(*tags) != len(tt.expected) {
				t.Errorf("expected %d tags, got %v", len(tt.expected), *tags)
			}
			for key, value := range tt.expected {
				if (*tags)[key] != value {
					t.Errorf("expected tag %q to be %v, got %v", key, value, (*tags)[key])
				}
			}
		})
	}

	t.Run("nothing to report stays nil", func(t *testing.T) {
		if tags := mapTags("leader", map[string]string{}); tags != nil {
			t.Errorf("expected no tags, got %v", *tags)
		}
	})
}

func TestShortVersion(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		expected string
	}{
		{"build hash and date are dropped", "3.9.1-abc123, built on 2023-10-01", "3.9.1"},
		{"comma alone is dropped", "3.8.0, built on 2022-02-25", "3.8.0"},
		{"a bare version survives", "3.9.1", "3.9.1"},
		{"empty stays empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if version := shortVersion(tt.version); version != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, version)
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

// TestEnsembleIdentity covers what a zookeeper node can say about its cluster
// at all. mntr describes nothing but the node answering, so the ensemble it was
// configured with is the only thing that can contradict a node pasted in from
// another one - and it has to agree between members regardless of the parts
// each node writes only for itself.
func TestEnsembleIdentity(t *testing.T) {
	tests := []struct {
		name     string
		settings map[string]string
		expected string
	}{
		{
			name: "every member of one ensemble reports the same list",
			settings: map[string]string{
				"clientPort": "2181",
				"server.1":   "10.0.0.1:2888:3888;2181",
				"server.2":   "10.0.0.2:2888:3888;2181",
				"server.3":   "10.0.0.3:2888:3888;2181",
			},
			expected: "10.0.0.1:2888,10.0.0.2:2888,10.0.0.3:2888",
		},
		{
			name: "the order the lines are read in cannot change the identity",
			settings: map[string]string{
				"server.3": "10.0.0.3:2888:3888;2181",
				"server.1": "10.0.0.1:2888:3888;2181",
				"server.2": "10.0.0.2:2888:3888;2181",
			},
			expected: "10.0.0.1:2888,10.0.0.2:2888,10.0.0.3:2888",
		},
		{
			name: "a single-host ensemble states no client ports and is still one ensemble",
			settings: map[string]string{
				"server.1": "10.0.0.1:2888:3888",
				"server.2": "10.0.0.1:2890:3890",
				"server.3": "10.0.0.1:2892:3892",
			},
			expected: "10.0.0.1:2888,10.0.0.1:2890,10.0.0.1:2892",
		},
		{
			name: "the role and the address a node writes for itself are dropped",
			settings: map[string]string{
				"server.1": "10.0.0.1:2888:3888:participant;0.0.0.0:2181",
				"server.2": "10.0.0.2:2888:3888:participant;2181",
			},
			expected: "10.0.0.1:2888,10.0.0.2:2888",
		},
		{
			name:     "a conf without a membership section reports no ensemble",
			settings: map[string]string{"clientPort": "2181", "dataDir": "/data"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			identity := ensembleIdentity(tt.settings)
			if tt.expected == "" {
				if identity != nil {
					t.Fatalf("expected no ensemble identity, got %q", *identity)
				}
				return
			}
			if identity == nil {
				t.Fatalf("expected ensemble identity %q, got nil", tt.expected)
			}
			if *identity != tt.expected {
				t.Errorf("expected ensemble identity %q, got %q", tt.expected, *identity)
			}
		})
	}
}

// TestEnsembleIdentityTellsEnsemblesApart is the case the whole field exists
// for: two real ensembles must never produce one string, or a node from the
// wrong one goes unnoticed.
func TestEnsembleIdentityTellsEnsemblesApart(t *testing.T) {
	ours := ensembleIdentity(map[string]string{
		"server.1": "10.0.0.1:2888:3888;2181",
		"server.2": "10.0.0.2:2888:3888;2181",
	})
	theirs := ensembleIdentity(map[string]string{
		"server.1": "10.9.9.1:2888:3888;2181",
		"server.2": "10.9.9.2:2888:3888;2181",
	})

	if ours == nil || theirs == nil {
		t.Fatal("expected both ensembles to report an identity")
	}
	if *ours == *theirs {
		t.Errorf("expected two ensembles to differ, both reported %q", *ours)
	}
}

// TestMemberEndpoint covers the one part of a server.N line Ivory can actually
// connect to. The quorum and election ports are how the ensemble talks to
// itself; the client port after the ";" is the only address a keeper request
// can use, and a line without it names no reachable node.
func TestMemberEndpoint(t *testing.T) {
	tests := []struct {
		name   string
		server string
		host   string
		port   int
		ok     bool
	}{
		{name: "bare client port", server: "10.0.0.1:2888:3888;2181", host: "10.0.0.1", port: 2181, ok: true},
		{name: "client address and port", server: "10.0.0.2:2888:3888:participant;0.0.0.0:2181", host: "10.0.0.2", port: 2181, ok: true},
		{name: "a single-host line states no client port and is not addressable", server: "10.0.0.1:2888:3888"},
		{name: "an unparseable client port is dropped rather than defaulted", server: "10.0.0.1:2888:3888;zk"},
		{name: "empty", server: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			host, port, ok := memberEndpoint(tt.server)
			if ok != tt.ok {
				t.Fatalf("expected ok %v, got %v", tt.ok, ok)
			}
			if !ok {
				return
			}
			if host != tt.host || port != tt.port {
				t.Errorf("expected %s:%d, got %s:%d", tt.host, tt.port, host, port)
			}
		})
	}
}

// TestMapEnsembleMembers covers what a single zookeeper can say about the rest
// of its ensemble. It reports the members it can address and skips itself,
// which serverId identifies - a self entry would be a second, stateless copy of
// the node that is already answering for itself.
func TestMapEnsembleMembers(t *testing.T) {
	tests := []struct {
		name     string
		settings map[string]string
		expected []string
	}{
		{
			name: "a multi-host ensemble reports its peers and skips itself",
			settings: map[string]string{
				"serverId": "1",
				"server.1": "10.0.0.1:2888:3888;2181",
				"server.2": "10.0.0.2:2888:3888;2181",
				"server.3": "10.0.0.3:2888:3888;2181",
			},
			expected: []string{"10.0.0.2:2181", "10.0.0.3:2181"},
		},
		{
			name: "a single-host ensemble states no client ports, so it reports no members at all",
			settings: map[string]string{
				"serverId": "1",
				"server.1": "10.0.0.1:2888:3888",
				"server.2": "10.0.0.1:2890:3890",
				"server.3": "10.0.0.1:2892:3892",
			},
		},
		{
			name:     "a conf with no membership section reports nothing",
			settings: map[string]string{"clientPort": "2181"},
		},
		{
			name: "a conf that does not name the node still reports every member it can address",
			settings: map[string]string{
				"server.1": "10.0.0.1:2888:3888;2181",
				"server.2": "10.0.0.2:2888:3888;2181",
			},
			expected: []string{"10.0.0.1:2181", "10.0.0.2:2181"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			members := mapEnsembleMembers(tt.settings)
			if len(members) != len(tt.expected) {
				t.Fatalf("expected %d members, got %d: %v", len(tt.expected), len(members), members)
			}
			for i, key := range tt.expected {
				if members[i].Key == nil || *members[i].Key != key {
					t.Errorf("expected member %d to be %q, got %v", i, key, members[i].Key)
				}
				if members[i].State != keeper.StateUnknown {
					t.Errorf("nothing contacted %q, expected unknown state, got %q", key, members[i].State)
				}
			}
		})
	}
}
