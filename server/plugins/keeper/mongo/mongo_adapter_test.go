package mongo

import (
	"errors"
	"ivory/plugins/keeper"
	"net/http"
	"testing"
	"time"
)

func TestMapRoleState(t *testing.T) {
	tests := []struct {
		name          string
		stateStr      string
		health        float64
		expectedRole  keeper.Role
		expectedState keeper.State
	}{
		{"primary", "PRIMARY", 1, keeper.Leader, keeper.StateRunning},
		{"secondary", "SECONDARY", 1, keeper.Replica, keeper.StateRunning},
		{"recovering", "RECOVERING", 1, keeper.Replica, keeper.StateRestarting},
		{"rollback", "ROLLBACK", 1, keeper.Replica, keeper.StateRestarting},
		{"startup", "STARTUP", 1, keeper.Unknown, keeper.StateStarting},
		{"startup2", "STARTUP2", 1, keeper.Unknown, keeper.StateStarting},
		{"arbiter", "ARBITER", 1, keeper.Unknown, keeper.StateRunning},
		{"down", "DOWN", 1, keeper.Unknown, keeper.StateUnreachable},
		{"removed", "REMOVED", 1, keeper.Unknown, keeper.StateStopped},
		{"unknown state", "UNKNOWN", 1, keeper.Unknown, keeper.StateUnknown},
		{"unhealthy overrides primary", "PRIMARY", 0, keeper.Unknown, keeper.StateUnreachable},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			role, state := mapRoleState(tt.stateStr, tt.health)
			if role != tt.expectedRole {
				t.Errorf("expected role %v, got %v", tt.expectedRole, role)
			}
			if state != tt.expectedState {
				t.Errorf("expected state %v, got %v", tt.expectedState, state)
			}
		})
	}
}

func TestMapStatus(t *testing.T) {
	now := time.Now()
	status := &replSetStatus{
		Set: "rs0",
		Members: []replSetMember{
			{Name: "h1:27017", StateStr: "PRIMARY", Health: 1, OptimeDate: now},
			{Name: "h2:27017", StateStr: "SECONDARY", Health: 1, OptimeDate: now.Add(-5 * time.Second)},
			{Name: "h3:27017", StateStr: "ARBITER", Health: 1, OptimeDate: time.Time{}},
		},
	}

	responses := mapStatus(status)
	if len(responses) != 3 {
		t.Fatalf("expected 3 responses, got %d", len(responses))
	}

	primary := responses[0]
	if primary.Role != keeper.Leader || primary.Lag != 0 {
		t.Errorf("expected leader with zero lag, got role %v lag %d", primary.Role, primary.Lag)
	}
	if primary.Key == nil || *primary.Key != "h1:27017" {
		t.Errorf("expected key h1:27017, got %v", primary.Key)
	}
	if primary.DiscoveredHost == nil || *primary.DiscoveredHost != "h1" {
		t.Errorf("expected discovered host h1, got %v", primary.DiscoveredHost)
	}
	if primary.DiscoveredKeeperPort == nil || *primary.DiscoveredKeeperPort != 27017 {
		t.Errorf("expected discovered keeper port 27017, got %v", primary.DiscoveredKeeperPort)
	}
	if primary.DiscoveredDbPort == nil || *primary.DiscoveredDbPort != 27017 {
		t.Errorf("expected discovered db port 27017, got %v", primary.DiscoveredDbPort)
	}

	secondary := responses[1]
	if secondary.Role != keeper.Replica || secondary.Lag != 5 {
		t.Errorf("expected replica with 5s lag, got role %v lag %d", secondary.Role, secondary.Lag)
	}

	arbiter := responses[2]
	if arbiter.Role != keeper.Unknown || arbiter.Lag != 0 {
		t.Errorf("expected arbiter unknown role with zero lag, got role %v lag %d", arbiter.Role, arbiter.Lag)
	}
}

func TestMapStatusNoPrimary(t *testing.T) {
	status := &replSetStatus{
		Members: []replSetMember{
			{Name: "h1:27017", StateStr: "SECONDARY", Health: 1, OptimeDate: time.Now()},
		},
	}
	responses := mapStatus(status)
	if responses[0].Lag != 0 {
		t.Errorf("expected zero lag with no primary in view, got %d", responses[0].Lag)
	}
}

func TestSelfMember(t *testing.T) {
	status := &replSetStatus{
		Members: []replSetMember{
			{Name: "h1:27017", Self: false},
			{Name: "h2:27017", Self: true},
		},
	}
	self, ok := selfMember(status)
	if !ok || self.Name != "h2:27017" {
		t.Errorf("expected self h2:27017, got %+v ok=%v", self, ok)
	}

	_, ok = selfMember(&replSetStatus{})
	if ok {
		t.Error("expected no self member for an empty status")
	}
}

func TestConfigUpdateRequiresBody(t *testing.T) {
	adapter := NewPlugin()
	request := keeper.Request{Host: "localhost", Port: 27017}

	_, status, err := adapter.ConfigUpdate(request)

	if status != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", status)
	}
	if !errors.Is(err, ErrConfigBodyRequired) {
		t.Errorf("expected ErrConfigBodyRequired, got %v", err)
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

func TestListConnectFailure(t *testing.T) {
	adapter := NewPlugin()
	request := keeper.Request{Host: "", Port: 0}

	_, status, err := adapter.List(request)

	if status != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", status)
	}
	if err == nil {
		t.Error("expected an error for an unspecified host/port")
	}
}
