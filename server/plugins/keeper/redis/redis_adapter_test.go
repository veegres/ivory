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

func TestMapTags(t *testing.T) {
	tests := []struct {
		name     string
		fields   map[string]string
		role     keeper.Role
		expected map[string]any
	}{
		{
			name:     "master reports its replica count",
			fields:   map[string]string{"redis_version": "7.2.4", "used_memory_human": "1.51M", "connected_slaves": "2"},
			role:     keeper.Leader,
			expected: map[string]any{"version": "7.2.4", "memory": "1.51M", "replicas": "2"},
		},
		{
			name: "replica reports the master it follows and the link to it",
			fields: map[string]string{
				"redis_version": "7.2.4", "used_memory_human": "1.51M", "connected_slaves": "0",
				"master_link_status": "down", "master_host": "10.0.0.1", "master_port": "6379",
			},
			role:     keeper.Replica,
			expected: map[string]any{"version": "7.2.4", "memory": "1.51M", "link": "down", "master": "10.0.0.1:6379"},
		},
		{
			name:     "nothing to report stays nil",
			fields:   map[string]string{},
			role:     keeper.Leader,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tags := mapTags(tt.fields, tt.role)
			if tt.expected == nil {
				if tags != nil {
					t.Fatalf("expected no tags, got %v", *tags)
				}
				return
			}
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
}

func TestListRequiresCredentials(t *testing.T) {
	adapter := NewPlugin()
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
	adapter := NewPlugin()
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
	adapter := NewPlugin()
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

// TestMapMembers covers the one direction a redis node can be believed about:
// the master a replica names comes out of its own replicaof, so it is an
// address somebody configured. It is what stops a node borrowed from another
// redis being accepted in silence - the master it names is one nobody
// configured, so it surfaces as a node found in the keeper response but
// missing from the cluster.
func TestMapMembers(t *testing.T) {
	tests := []struct {
		name     string
		fields   map[string]string
		expected []string
	}{
		{
			name:     "a replica names the master it follows",
			fields:   map[string]string{"role": "slave", "master_host": "10.0.0.1", "master_port": "6379"},
			expected: []string{"10.0.0.1:6379"},
		},
		{
			name:   "a master with no replicas names nobody",
			fields: map[string]string{"role": "master", "connected_slaves": "0"},
		},
		{
			name:   "a replica that has not been told a master yet names nobody",
			fields: map[string]string{"role": "slave"},
		},
		{
			name:   "an unparseable master port is dropped rather than defaulted",
			fields: map[string]string{"role": "slave", "master_host": "10.0.0.1", "master_port": "six"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			members := mapMembers(tt.fields)
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
				if members[i].DiscoveredKeeperPort == nil || members[i].DiscoveredDbPort == nil {
					t.Errorf("expected %q to carry both ports so it can be matched to a configured node", key)
				}
			}
		})
	}
}

// TestMapMembersClaimsNoRoleForTheMaster pins why the one direction reported
// still claims nothing: a replica usually follows the master, but redis allows
// a replica of a replica, so Leader would be a guess and a wrong one puts a
// second leader on the overview.
func TestMapMembersClaimsNoRoleForTheMaster(t *testing.T) {
	master := mapMembers(map[string]string{"role": "slave", "master_host": "10.0.0.1", "master_port": "6379"})
	if master[0].Role != keeper.Unknown {
		t.Errorf("a chained replica makes the master a guess, expected unknown role, got %q", master[0].Role)
	}
}

// TestMapMembersDoesNotEnumerateReplicas pins the deliberate gap: the ip on a
// slaveN: line is the source address the master observed, so under the shipped
// single-host template it is ::1 and under bridge networking the docker
// gateway. Reporting those invented three nodes nobody configured and warned
// about each one, which is worse than reporting no member list at all.
func TestMapMembersDoesNotEnumerateReplicas(t *testing.T) {
	members := mapMembers(map[string]string{
		"role":             "master",
		"connected_slaves": "2",
		"slave0":           "ip=::1,port=6380,state=online,offset=1,lag=0",
		"slave1":           "ip=::1,port=6381,state=online,offset=1,lag=0",
	})
	if len(members) != 0 {
		t.Errorf("a master must name no replicas, got %v", members)
	}
}

// TestMapTagsReportsReplicaCount pins what replaced the replica list: the count
// is the only honest thing a master can say about its replicas, and it is the
// one place a replica missing from the cluster shows up at all.
func TestMapTagsReportsReplicaCount(t *testing.T) {
	tags := mapTags(map[string]string{"connected_slaves": "2"}, keeper.Leader)
	if tags == nil {
		t.Fatal("expected tags, got nil")
	}
	if (*tags)["replicas"] != "2" {
		t.Errorf("expected replicas tag %q, got %v", "2", (*tags)["replicas"])
	}
}
