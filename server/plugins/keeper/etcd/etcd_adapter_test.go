package etcd

import (
	"errors"
	"ivory/plugins/keeper"
	"net/http"
	"testing"
)

func TestMapMembers(t *testing.T) {
	members := []member{
		{ID: 1, Name: "etcd1", ClientURLs: []string{"http://etcd1:2379"}},
		{ID: 2, Name: "etcd2", ClientURLs: []string{"http://etcd2:2379"}},
		{ID: 3, Name: "etcd3", ClientURLs: []string{"http://etcd3:2379"}},
	}
	statuses := map[uint64]endpointStatus{
		1: {Leader: 1, RaftIndex: 100},
		2: {Leader: 1, RaftIndex: 95},
		3: {Err: errors.New("connection refused")},
	}

	responses := mapMembers(members, statuses)
	if len(responses) != 3 {
		t.Fatalf("expected 3 responses, got %d", len(responses))
	}

	leader := responses[0]
	if leader.Role != keeper.Leader || leader.Lag != 0 {
		t.Errorf("expected leader with lag 0, got role %v lag %d", leader.Role, leader.Lag)
	}
	if leader.Key == nil || *leader.Key != "etcd1" {
		t.Errorf("expected key etcd1, got %v", leader.Key)
	}
	if leader.DiscoveredName == nil || *leader.DiscoveredName != "etcd1" {
		t.Errorf("expected discovered name etcd1, got %v", leader.DiscoveredName)
	}
	if leader.DiscoveredHost == nil || *leader.DiscoveredHost != "etcd1" {
		t.Errorf("expected discovered host etcd1, got %v", leader.DiscoveredHost)
	}
	if leader.DiscoveredKeeperPort == nil || *leader.DiscoveredKeeperPort != 2379 {
		t.Errorf("expected discovered keeper port 2379, got %v", leader.DiscoveredKeeperPort)
	}
	if leader.DiscoveredDbPort == nil || *leader.DiscoveredDbPort != 2379 {
		t.Errorf("expected discovered db port 2379, got %v", leader.DiscoveredDbPort)
	}

	replica := responses[1]
	if replica.Role != keeper.Replica {
		t.Errorf("expected replica role, got %v", replica.Role)
	}
	if replica.Lag != 5 {
		t.Errorf("expected raft index lag 5, got %d", replica.Lag)
	}

	unreachable := responses[2]
	if unreachable.Role != keeper.Unknown {
		t.Errorf("expected unknown role, got %v", unreachable.Role)
	}
	if unreachable.State != keeper.StateUnreachable {
		t.Errorf("expected unreachable state, got %q", unreachable.State)
	}
	if unreachable.Lag != -1 {
		t.Errorf("expected lag -1, got %d", unreachable.Lag)
	}
}

func TestMapMembersEdgeCases(t *testing.T) {
	t.Run("member without client urls is starting", func(t *testing.T) {
		members := []member{{ID: 1, Name: "etcd1"}}
		responses := mapMembers(members, map[uint64]endpointStatus{})
		if responses[0].State != keeper.StateStarting {
			t.Errorf("expected starting state, got %q", responses[0].State)
		}
		if responses[0].DiscoveredHost != nil {
			t.Errorf("expected nil discovered host, got %v", *responses[0].DiscoveredHost)
		}
	})

	t.Run("learner gets tag", func(t *testing.T) {
		members := []member{{ID: 2, Name: "learner1", ClientURLs: []string{"http://l1:2379"}, IsLearner: true}}
		statuses := map[uint64]endpointStatus{2: {Leader: 1, RaftIndex: 90}}
		responses := mapMembers(members, statuses)
		if responses[0].Tags == nil || (*responses[0].Tags)["learner"] != true {
			t.Errorf("expected learner tag, got %v", responses[0].Tags)
		}
	})

	t.Run("replica ahead of stale leader index does not underflow", func(t *testing.T) {
		members := []member{
			{ID: 1, Name: "etcd1", ClientURLs: []string{"http://e1:2379"}},
			{ID: 2, Name: "etcd2", ClientURLs: []string{"http://e2:2379"}},
		}
		statuses := map[uint64]endpointStatus{
			1: {Leader: 1, RaftIndex: 90},
			2: {Leader: 1, RaftIndex: 95},
		}
		responses := mapMembers(members, statuses)
		if responses[1].Lag != -1 {
			t.Errorf("expected lag -1 when replica index is ahead, got %d", responses[1].Lag)
		}
	})
}

func TestFindLeader(t *testing.T) {
	t.Run("majority agreement designates the leader", func(t *testing.T) {
		statuses := map[uint64]endpointStatus{
			1: {Leader: 1, RaftIndex: 100},
			2: {Leader: 1, RaftIndex: 95},
			3: {Err: errors.New("connection refused")},
		}
		leaderID, raftIndex := findLeader(statuses)
		if leaderID != 1 {
			t.Errorf("expected leader 1, got %d", leaderID)
		}
		if raftIndex != 100 {
			t.Errorf("expected raft index 100, got %d", raftIndex)
		}
	})

	t.Run("split members disagreeing on the leader report no leader", func(t *testing.T) {
		statuses := map[uint64]endpointStatus{
			1: {Leader: 1, RaftIndex: 100},
			2: {Leader: 2, RaftIndex: 95},
		}
		leaderID, raftIndex := findLeader(statuses)
		if leaderID != 0 {
			t.Errorf("expected no leader (0), got %d", leaderID)
		}
		if raftIndex != 0 {
			t.Errorf("expected raft index 0, got %d", raftIndex)
		}
	})

	t.Run("minority report is not enough to designate a leader", func(t *testing.T) {
		statuses := map[uint64]endpointStatus{
			1: {Leader: 1, RaftIndex: 100},
			2: {Err: errors.New("connection refused")},
			3: {Err: errors.New("connection refused")},
		}
		leaderID, _ := findLeader(statuses)
		if leaderID != 0 {
			t.Errorf("expected no leader (0), got %d", leaderID)
		}
	})
}

func TestParseSwitchoverBody(t *testing.T) {
	candidate := "etcd2"
	schedule := "2026-07-04T10:00:00"
	empty := ""

	tests := []struct {
		name        string
		body        any
		expectedErr error
	}{
		{
			name: "valid candidate",
			body: map[string]any{"leader": "etcd1", "candidate": candidate},
		},
		{
			name:        "schedule rejected",
			body:        map[string]any{"candidate": candidate, "scheduled_at": schedule},
			expectedErr: ErrScheduleNotSupported,
		},
		{
			name:        "missing candidate",
			body:        map[string]any{"leader": "etcd1"},
			expectedErr: ErrCandidateRequired,
		},
		{
			name:        "empty candidate",
			body:        map[string]any{"candidate": empty},
			expectedErr: ErrCandidateRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := parseSwitchoverBody(tt.body)
			if tt.expectedErr != nil {
				if !errors.Is(err, tt.expectedErr) {
					t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if parsed.Candidate == nil || *parsed.Candidate != candidate {
				t.Errorf("expected candidate %q, got %v", candidate, parsed.Candidate)
			}
		})
	}
}

func TestResolveCandidate(t *testing.T) {
	members := []member{
		{ID: 1, Name: "etcd1", ClientURLs: []string{"http://10.0.0.1:2379"}},
		{ID: 2, Name: "etcd2", ClientURLs: []string{"http://10.0.0.2:2379"}},
	}

	tests := []struct {
		name        string
		candidate   string
		expectedID  uint64
		expectedErr error
	}{
		{name: "resolve by member name", candidate: "etcd2", expectedID: 2},
		{name: "resolve by client url host", candidate: "10.0.0.1", expectedID: 1},
		{name: "unknown candidate", candidate: "missing", expectedErr: ErrCandidateNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := resolveCandidate(members, tt.candidate)
			if tt.expectedErr != nil {
				if !errors.Is(err, tt.expectedErr) {
					t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if id != tt.expectedID {
				t.Errorf("expected member id %d, got %d", tt.expectedID, id)
			}
		})
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
		{"Config", func() (int, error) { _, s, e := adapter.Config(request); return s, e }},
		{"ConfigUpdate", func() (int, error) { _, s, e := adapter.ConfigUpdate(request); return s, e }},
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
