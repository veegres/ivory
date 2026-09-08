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
		expectedState keeper.State
		expectedRole  keeper.Role
		expectedLag   int64
	}{
		{name: "healthy node is a running replica", absoluteDelay: 0, expectedState: keeper.StateRunning, expectedRole: keeper.Replica, expectedLag: 0},
		// NOTE: a node that answered is running whatever its replication is
		// doing - it serves reads - so a delay changes neither state nor role.
		{name: "a lagging node is still a running replica", absoluteDelay: 12, expectedState: keeper.StateRunning, expectedRole: keeper.Replica, expectedLag: 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := mapNode("ch1", 9000, tt.absoluteDelay)
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
		name                          string
		readonly                      bool
		activeReplicas, totalReplicas uint32
		stuckCount                    uint64
		sampleError                   string
		expectedCount                 int
	}{
		{name: "no replicated tables yet", activeReplicas: 0, totalReplicas: 0, stuckCount: 0, expectedCount: 0},
		{name: "every replica has a session and nothing is stuck", activeReplicas: 3, totalReplicas: 3, stuckCount: 0, expectedCount: 0},
		{name: "a replica dropped its coordination-store session", activeReplicas: 2, totalReplicas: 3, stuckCount: 0, expectedCount: 1},
		{name: "sessions are fine but a fetch is stuck", activeReplicas: 3, totalReplicas: 3, stuckCount: 1, sampleError: "Connection refused", expectedCount: 1},
		{name: "both at once", activeReplicas: 1, totalReplicas: 3, stuckCount: 2, sampleError: "DNS_ERROR", expectedCount: 2},
		// A read-only node is running and serving reads, so it is a warning
		// rather than a state - the same call redis makes for a dead link.
		{name: "read-only tables warn without changing the state", readonly: true, activeReplicas: 3, totalReplicas: 3, stuckCount: 0, expectedCount: 1},
		{name: "read-only on top of a lost session is two warnings", readonly: true, activeReplicas: 2, totalReplicas: 3, stuckCount: 0, expectedCount: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warnings := replicaWarnings(tt.readonly, tt.activeReplicas, tt.totalReplicas, tt.stuckCount, tt.sampleError)
			if len(warnings) != tt.expectedCount {
				t.Fatalf("expected %d warning(s), got %v", tt.expectedCount, warnings)
			}
			if tt.stuckCount > 0 && !containsSubstring(warnings, tt.sampleError) {
				t.Errorf("expected the stuck-queue warning to surface the real error %q, got %v", tt.sampleError, warnings)
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
		declared bool
		err      error
		expected string
	}{
		{
			name:     "a node that declares the cluster says nothing",
			cluster:  "prod",
			declared: true,
		},
		{
			name:     "a node that has never heard of the cluster says so",
			cluster:  "prod",
			expected: `this node is not in the cluster "prod"`,
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
