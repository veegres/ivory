package clickhouse

import (
	"errors"
	"ivory/plugins/keeper"
	"net/http"
	"strings"
	"testing"
)

func TestMapNode(t *testing.T) {
	tests := []struct {
		name          string
		absoluteDelay uint64
		readonly      bool
		expectedState keeper.State
		expectedRole  keeper.Role
		expectedLag   int64
	}{
		{name: "healthy node is a running replica", absoluteDelay: 0, expectedState: keeper.StateRunning, expectedRole: keeper.Replica, expectedLag: 0},
		// NOTE: a node that answered is running whatever its replication is
		// doing - it serves reads - so a delay changes neither state nor role.
		{name: "a lagging node is still a running replica", absoluteDelay: 12, expectedState: keeper.StateRunning, expectedRole: keeper.Replica, expectedLag: 12},
		// NOTE: a read-only replica reports the current unix time as its
		// absolute_delay, which is no lag at all.
		{name: "a read-only node still runs but its lag is unknown", absoluteDelay: 1791797896, readonly: true, expectedState: keeper.StateRunning, expectedRole: keeper.Replica, expectedLag: -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := mapNode("ch1", 9000, tt.absoluteDelay, tt.readonly)
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

// TestReplicaWarnings pins the two independent connectivity signals a real
// cluster showed (see clickhouse_adapter.go's List doc): a dropped
// coordination-store session moves active_replicas below total_replicas even
// though nothing is stuck, and a stuck data fetch leaves every session
// healthy while system.replication_queue accumulates a genuine error -
// either can fire alone, and a healthy cluster with no replicated tables yet
// (total_replicas == 0) reports neither.
func TestReplicaWarnings(t *testing.T) {
	tests := []struct {
		name          string
		status        replicaStatus
		expectedCount int
		contains      string
	}{
		{name: "no replicated tables yet", status: replicaStatus{}, expectedCount: 0},
		{name: "every replica has a session and nothing is stuck", status: replicaStatus{ActiveReplicas: 3, TotalReplicas: 3}, expectedCount: 0},
		{name: "a replica dropped its coordination-store session", status: replicaStatus{ActiveReplicas: 2, TotalReplicas: 3}, expectedCount: 1},
		{name: "sessions are fine but a fetch is stuck", status: replicaStatus{ActiveReplicas: 3, TotalReplicas: 3, StuckCount: 1, StuckError: "Connection refused"}, expectedCount: 1, contains: "Connection refused"},
		{name: "both at once", status: replicaStatus{ActiveReplicas: 1, TotalReplicas: 3, StuckCount: 2, StuckError: "DNS_ERROR"}, expectedCount: 2, contains: "DNS_ERROR"},
		// A read-only node is running and serving reads, so it is a warning
		// rather than a state - the same call redis makes for a dead link.
		{name: "read-only tables warn without changing the state and are named", status: replicaStatus{Readonly: 1, ReadonlyCount: 4, ReadonlyTables: "db.a, db.b, db.c", ActiveReplicas: 3, TotalReplicas: 3}, expectedCount: 1, contains: "4 replicated table(s) are read-only (db.a, db.b, db.c)"},
		{name: "read-only on top of a lost session is two warnings", status: replicaStatus{Readonly: 1, ReadonlyCount: 1, ReadonlyTables: "db.a", ActiveReplicas: 2, TotalReplicas: 3}, expectedCount: 2},
		{name: "delay within the distributed query limit", status: replicaStatus{AbsoluteDelay: 300, MaxDelay: 300}, expectedCount: 0},
		{name: "delay past the distributed query limit", status: replicaStatus{AbsoluteDelay: 301, MaxDelay: 300}, expectedCount: 1, contains: "replication is 301s behind"},
		{name: "a read-only replica's delay is not compared", status: replicaStatus{Readonly: 1, ReadonlyCount: 1, ReadonlyTables: "db.a", AbsoluteDelay: 1700000000, MaxDelay: 300}, expectedCount: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warnings := replicaWarnings(tt.status)
			if len(warnings) != tt.expectedCount {
				t.Fatalf("expected %d warning(s), got %v", tt.expectedCount, warnings)
			}
			if tt.contains != "" && !containsSubstring(warnings, tt.contains) {
				t.Errorf("expected a warning containing %q, got %v", tt.contains, warnings)
			}
		})
	}
}

func containsSubstring(warnings []string, substr string) bool {
	for _, w := range warnings {
		if strings.Contains(w, substr) {
			return true
		}
	}
	return false
}

func TestListRequiresCredentials(t *testing.T) {
	adapter := NewPlugin()
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
	adapter := NewPlugin()
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
	adapter := NewPlugin()
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

// TestMembershipWarnings covers what the node could not establish about the
// cluster, which is a separate failure from anything wrong with the node. The
// case that matters is a node declaring no <remote_servers> entry under the
// name Ivory asked about: it may be a healthy member of three other clusters,
// and rendering it as this one would be the exact silent acceptance the cluster
// name was passed down to prevent.
func TestMembershipWarnings(t *testing.T) {
	tests := []struct {
		name     string
		cluster  string
		declared []string
		err      error
		expected string
	}{
		{
			name:     "a node that declares the cluster says nothing",
			cluster:  "prod",
			declared: []string{"default", "prod"},
		},
		{
			name:     "a node that has never heard of the cluster names the ones it is in",
			cluster:  "prod",
			declared: []string{"default", "test_shard_localhost"},
			expected: `this node is not in the cluster "prod", it declares: default, test_shard_localhost`,
		},
		{
			name:     "a node that declares nothing says so",
			cluster:  "prod",
			expected: `this node is not in the cluster "prod", it declares no clusters`,
		},
		{
			name:     "a cluster that could not be read at all is reported as that, not as absent",
			cluster:  "prod",
			err:      errors.New("Not enough privileges"),
			expected: "cannot read system.clusters: Not enough privileges",
		},
		{
			name: "a single-node action outside any cluster is asked nothing and answers nothing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warnings := membershipWarnings(tt.cluster, tt.declared, tt.err)
			if tt.expected == "" {
				if len(warnings) != 0 {
					t.Fatalf("expected no warnings, got %v", warnings)
				}
				return
			}
			if len(warnings) != 1 || warnings[0] != tt.expected {
				t.Errorf("expected warning %q, got %v", tt.expected, warnings)
			}
		})
	}
}

// TestMapPeer pins what a member read out of <remote_servers> may and may not
// claim. It is the address Ivory would reach the peer on, so it carries a host
// and both ports; nothing contacted it, so it carries no state and no lag. It
// names no cluster either - the name came from Ivory, and handing it back would
// be an echo rather than anything the node discovered.
func TestMapPeer(t *testing.T) {
	peer := mapPeer("10.0.0.2", 9000)

	if peer.State != keeper.StateUnknown {
		t.Errorf("a config file cannot vouch for liveness, expected unknown state, got %q", peer.State)
	}
	if peer.Role != keeper.Replica {
		t.Errorf("every member of a multi-leader engine is a replica, got %q", peer.Role)
	}
	if peer.Lag != -1 {
		t.Errorf("expected unknown lag -1, got %d", peer.Lag)
	}
	if peer.Status != nil {
		t.Errorf("expected no keeper status, got %v", *peer.Status)
	}
	if peer.DiscoveredCluster != nil {
		t.Errorf("expected no discovered cluster, got %q", *peer.DiscoveredCluster)
	}
	if peer.DiscoveredHost == nil || *peer.DiscoveredHost != "10.0.0.2" {
		t.Errorf("expected discovered host 10.0.0.2, got %v", peer.DiscoveredHost)
	}
	if peer.DiscoveredKeeperPort == nil || *peer.DiscoveredKeeperPort != 9000 {
		t.Errorf("expected discovered keeper port 9000, got %v", peer.DiscoveredKeeperPort)
	}
	if peer.DiscoveredDbPort == nil || *peer.DiscoveredDbPort != 9000 {
		t.Errorf("expected discovered db port 9000, got %v", peer.DiscoveredDbPort)
	}
	if peer.Key == nil || *peer.Key != "10.0.0.2:9000" {
		t.Errorf("expected key 10.0.0.2:9000, got %v", peer.Key)
	}
}
