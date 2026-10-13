package etcd

import (
	"errors"
	"ivory/plugins/keeper"
	"net/http"
	"slices"
	"strings"
	"testing"
)

func TestMapTags(t *testing.T) {
	tests := []struct {
		name     string
		member   member
		status   endpointStatus
		expected map[string]any
	}{
		{
			name:     "reachable voting member",
			member:   member{ID: 1, Name: "etcd1"},
			status:   endpointStatus{Version: "3.5.9", RaftTerm: 7},
			expected: map[string]any{"version": "3.5.9", "raftTerm": uint64(7)},
		},
		{
			name:     "learner is flagged",
			member:   member{ID: 2, Name: "etcd2", IsLearner: true},
			status:   endpointStatus{Version: "3.5.9", RaftTerm: 7},
			expected: map[string]any{"learner": true, "version": "3.5.9", "raftTerm": uint64(7)},
		},
		{
			// Ivory does not know the backend quota, so a size on its own says
			// nothing about the read-only alarm it would warn of.
			name:     "the backend size is not a tag",
			member:   member{ID: 4, Name: "etcd4"},
			status:   endpointStatus{Version: "3.5.9"},
			expected: map[string]any{"version": "3.5.9"},
		},
		{
			name:     "unreachable member reports nothing it cannot know",
			member:   member{ID: 3, Name: "etcd3"},
			status:   endpointStatus{Err: errors.New("unreachable")},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertTags(t, mapTags(tt.member, tt.status), tt.expected)
		})
	}
}

func assertTags(t *testing.T, tags *map[string]any, expected map[string]any) {
	t.Helper()
	if expected == nil {
		if tags != nil {
			t.Fatalf("expected no tags, got %v", *tags)
		}
		return
	}
	if tags == nil {
		t.Fatal("expected tags, got nil")
	}
	if len(*tags) != len(expected) {
		t.Errorf("expected %d tags, got %v", len(expected), *tags)
	}
	for key, value := range expected {
		if (*tags)[key] != value {
			t.Errorf("expected tag %q to be %v, got %v", key, value, (*tags)[key])
		}
	}
}

func TestMapMembers(t *testing.T) {
	members := []member{
		{ID: 1, Name: "etcd1", ClientURLs: []string{"http://etcd1:2379"}},
		{ID: 2, Name: "etcd2", ClientURLs: []string{"http://etcd2:2379"}},
		{ID: 3, Name: "etcd3", ClientURLs: []string{"http://etcd3:2379"}},
	}

	t.Run("asked member is a follower", func(t *testing.T) {
		responses := mapMembers(members, 2, endpointStatus{Leader: 1, RaftIndex: 95}, leaderStatus{ID: 1, RaftIndex: 100})
		if len(responses) != 3 {
			t.Fatalf("expected 3 responses, got %d", len(responses))
		}

		self := responses[1]
		if self.Role != keeper.Replica || self.State != keeper.StateRunning {
			t.Errorf("expected running replica, got role %v state %q", self.Role, self.State)
		}
		if self.Lag != 5 {
			t.Errorf("expected raft index lag 5, got %d", self.Lag)
		}

		leader := responses[0]
		if leader.Role != keeper.Leader {
			t.Errorf("expected the leader the asked member sees, got %v", leader.Role)
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

		// NOTE: a peer nobody contacted claims topology and no liveness, which
		// is what lets its own answer replace it in the merge.
		for _, peer := range []keeper.Response{responses[0], responses[2]} {
			if peer.State != keeper.StateUnknown || peer.Lag != -1 {
				t.Errorf("expected peer %s as hearsay, got state %q lag %d", *peer.Key, peer.State, peer.Lag)
			}
			if peer.Warnings != nil {
				t.Errorf("expected no warnings on peer %s, got %v", *peer.Key, peer.Warnings)
			}
		}
		if responses[2].Role != keeper.Replica {
			t.Errorf("expected replica role, got %v", responses[2].Role)
		}
	})

	t.Run("asked member is the leader", func(t *testing.T) {
		responses := mapMembers(members, 1, endpointStatus{Leader: 1, RaftIndex: 100}, leaderStatus{ID: 1, RaftIndex: 100})
		if responses[0].Role != keeper.Leader || responses[0].Lag != 0 || responses[0].State != keeper.StateRunning {
			t.Errorf("expected running leader with lag 0, got role %v lag %d state %q", responses[0].Role, responses[0].Lag, responses[0].State)
		}
	})

	t.Run("unreachable leader leaves lag unknown and says why", func(t *testing.T) {
		leader := leaderStatus{ID: 1, Err: errors.New("leader is unreachable at http://etcd1:2379: connection refused")}
		responses := mapMembers(members, 2, endpointStatus{Leader: 1, RaftIndex: 95}, leader)
		self := responses[1]
		if self.State != keeper.StateRunning {
			t.Errorf("the asked member answered, expected running, got %q", self.State)
		}
		if self.Lag != -1 {
			t.Errorf("expected lag -1, got %d", self.Lag)
		}
		if !slices.Equal(self.Warnings, []string{leader.Err.Error()}) {
			t.Errorf("expected the leader probe failure as a warning, got %v", self.Warnings)
		}
	})

	t.Run("no leader leaves every role unknown", func(t *testing.T) {
		responses := mapMembers(members, 2, endpointStatus{RaftIndex: 95}, leaderStatus{})
		if responses[0].Role != keeper.Unknown || responses[2].Role != keeper.Unknown {
			t.Errorf("expected unknown peer roles, got %v and %v", responses[0].Role, responses[2].Role)
		}
		if responses[1].Lag != -1 {
			t.Errorf("expected lag -1 without a leader, got %d", responses[1].Lag)
		}
	})
}

func TestMapMembersEdgeCases(t *testing.T) {
	t.Run("member without client urls is starting", func(t *testing.T) {
		members := []member{{ID: 1, Name: "etcd1"}}
		responses := mapMembers(members, 2, endpointStatus{}, leaderStatus{})
		if responses[0].State != keeper.StateStarting {
			t.Errorf("expected starting state, got %q", responses[0].State)
		}
		if responses[0].DiscoveredHost != nil {
			t.Errorf("expected nil discovered host, got %v", *responses[0].DiscoveredHost)
		}
	})

	t.Run("learner gets tag", func(t *testing.T) {
		members := []member{{ID: 2, Name: "learner1", ClientURLs: []string{"http://l1:2379"}, IsLearner: true}}
		responses := mapMembers(members, 2, endpointStatus{Leader: 1, RaftIndex: 90}, leaderStatus{ID: 1, RaftIndex: 100})
		if responses[0].Tags == nil || (*responses[0].Tags)["learner"] != true {
			t.Errorf("expected learner tag, got %v", responses[0].Tags)
		}
	})

	t.Run("learner peer gets tag", func(t *testing.T) {
		members := []member{{ID: 2, Name: "learner1", ClientURLs: []string{"http://l1:2379"}, IsLearner: true}}
		responses := mapMembers(members, 1, endpointStatus{Leader: 1}, leaderStatus{ID: 1})
		if responses[0].Tags == nil || (*responses[0].Tags)["learner"] != true {
			t.Errorf("expected learner tag, got %v", responses[0].Tags)
		}
	})

	t.Run("replica ahead of stale leader index does not underflow", func(t *testing.T) {
		members := []member{
			{ID: 1, Name: "etcd1", ClientURLs: []string{"http://e1:2379"}},
			{ID: 2, Name: "etcd2", ClientURLs: []string{"http://e2:2379"}},
		}
		responses := mapMembers(members, 2, endpointStatus{Leader: 1, RaftIndex: 95}, leaderStatus{ID: 1, RaftIndex: 90})
		if responses[1].Lag != -1 {
			t.Errorf("expected lag -1 when replica index is ahead, got %d", responses[1].Lag)
		}
	})

	t.Run("asked member missing from the list reports only peers", func(t *testing.T) {
		members := []member{{ID: 1, Name: "etcd1", ClientURLs: []string{"http://e1:2379"}}}
		responses := mapMembers(members, 9, endpointStatus{Leader: 1}, leaderStatus{ID: 1})
		if len(responses) != 1 || responses[0].State != keeper.StateUnknown {
			t.Errorf("expected one hearsay peer, got %v", responses)
		}
	})
}

func TestFindLeader(t *testing.T) {
	members := []member{
		{ID: 1, Name: "etcd1", ClientURLs: []string{"http://etcd1:2379", "http://10.0.0.1:2379"}},
		{ID: 2, Name: "etcd2", ClientURLs: []string{"http://etcd2:2379"}},
		{ID: 3, Name: "etcd3"},
	}
	tests := []struct {
		name   string
		selfID uint64
		self   endpointStatus
		leader leaderStatus
		url    string
	}{
		{name: "asked member is the leader, no probe", selfID: 1, self: endpointStatus{Leader: 1, RaftIndex: 100}, leader: leaderStatus{ID: 1, RaftIndex: 100}},
		{name: "follower probes the leader it sees", selfID: 2, self: endpointStatus{Leader: 1, RaftIndex: 95}, leader: leaderStatus{ID: 1}, url: "http://etcd1:2379"},
		{name: "no leader, no probe", selfID: 2, self: endpointStatus{RaftIndex: 95}, leader: leaderStatus{}},
		{name: "leader without client urls cannot be probed", selfID: 2, self: endpointStatus{Leader: 3}, leader: leaderStatus{ID: 3}},
		{name: "leader missing from the list cannot be probed", selfID: 2, self: endpointStatus{Leader: 7}, leader: leaderStatus{ID: 7}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leader, url := findLeader(members, tt.selfID, tt.self)
			if leader != tt.leader {
				t.Errorf("expected leader %+v, got %+v", tt.leader, leader)
			}
			if url != tt.url {
				t.Errorf("expected url %q, got %q", tt.url, url)
			}
		})
	}
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
	adapter := NewPlugin()
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

// TestMapMemberCarriesItsOwnFailure pins that a member Ivory could not reach
// explains itself on its own row. The reason used to be returned as an error
// for the whole List, which made every healthy node Ivory asked carry a
// "failed to get Keeper response" naming somebody else's outage.
func TestMapMemberCarriesItsOwnFailure(t *testing.T) {
	down := member{ID: 3, Name: "etcd3", ClientURLs: []string{"http://etcd3:2379"}}
	response := mapMember(down, endpointStatus{Err: errors.New("connection refused")}, 1, 100)

	if response.State != keeper.StateUnreachable {
		t.Errorf("expected unreachable, got %q", response.State)
	}
	if len(response.Warnings) != 1 || !strings.Contains(response.Warnings[0], "connection refused") || !strings.Contains(response.Warnings[0], "http://etcd3:2379") {
		t.Errorf("expected the member to carry its own reason, got %v", response.Warnings)
	}

	up := member{ID: 1, Name: "etcd1", ClientURLs: []string{"http://etcd1:2379"}}
	if warnings := mapMember(up, endpointStatus{RaftIndex: 100}, 1, 100).Warnings; warnings != nil {
		t.Errorf("a member that answered warns about nothing, got %v", warnings)
	}
}

func TestMapWarnings(t *testing.T) {
	m := member{ID: 1, Name: "etcd1", ClientURLs: []string{"http://etcd1:2379"}}
	tests := []struct {
		name     string
		member   member
		status   endpointStatus
		warnings []string
	}{
		{name: "healthy", member: m, status: endpointStatus{RaftIndex: 100, RaftAppliedIndex: 100}, warnings: nil},
		{name: "alarm", member: m, status: endpointStatus{RaftIndex: 100, RaftAppliedIndex: 100, Errors: []string{"memberID:1 alarm:NOSPACE "}}, warnings: []string{"member reports: memberID:1 alarm:NOSPACE "}},
		{name: "apply gap within the limit", member: m, status: endpointStatus{RaftIndex: 5100, RaftAppliedIndex: 100}, warnings: nil},
		{name: "apply gap past the limit", member: m, status: endpointStatus{RaftIndex: 5101, RaftAppliedIndex: 100}, warnings: []string{"5001 committed entries are not applied yet, past 5000 etcd refuses new requests"}},
		{name: "unreachable names the url", member: m, status: endpointStatus{Err: errors.New("timeout")}, warnings: []string{"member is unreachable at http://etcd1:2379: timeout"}},
		{name: "unreachable without url", member: member{ID: 2}, status: endpointStatus{Err: errors.New("timeout")}, warnings: []string{"member is unreachable: timeout"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapWarnings(tt.member, tt.status)
			if !slices.Equal(got, tt.warnings) {
				t.Errorf("expected %v, got %v", tt.warnings, got)
			}
		})
	}
}
