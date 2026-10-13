package node

import (
	"errors"
	"ivory/core/service/cert"
	"ivory/core/service/vault"
	"ivory/core/utils"
	"ivory/plugins/keeper"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

// fakeKeeperAdapter lets the test control exactly what List() reports, in
// particular returning a usable Response alongside a non-nil error the way
// a real adapter does when it can still describe a node's state despite a
// connection problem (e.g. postgres starting up).
type fakeKeeperAdapter struct {
	keeper.Adapter
	keeper.Metadata
	listResponse []keeper.Response
	listStatus   int
	listErr      error
	lastRequest  *keeper.Request
}

func (f *fakeKeeperAdapter) ReplicationModel() keeper.ReplicationModel { return keeper.SingleLeader }
func (f *fakeKeeperAdapter) Lag() keeper.Lag                           { return keeper.Lag{} }

func (f *fakeKeeperAdapter) List(request keeper.Request) ([]keeper.Response, int, error) {
	f.lastRequest = &request
	return f.listResponse, f.listStatus, f.listErr
}

func TestKeeperNodeList_KeepsResponseAlongsideError(t *testing.T) {
	host := "db1"
	port := 5432
	state := keeper.StateStarting

	keeperRegistry := utils.NewRegistry[keeper.PluginType, keeper.Plugin]()
	keeperRegistry.Register("fake", &fakeKeeperAdapter{
		listResponse: []keeper.Response{{
			State:                state,
			Role:                 keeper.Unknown,
			DiscoveredHost:       &host,
			DiscoveredKeeperPort: &port,
		}},
		listStatus: http.StatusServiceUnavailable,
		listErr:    errors.New("the database system is starting up"),
	})
	s := NewService(nil, keeperRegistry, nil, nil, nil)

	responses, status, err := s.KeeperNodeList(KeeperOneRequest{
		KeeperConnection: KeeperConnection{Host: host, Port: port},
		KeeperOptions:    KeeperOptions{Plugin: "fake"},
	})

	if err == nil {
		t.Fatalf("expected the adapter error to still be returned, got nil")
	}
	if status != http.StatusServiceUnavailable {
		t.Errorf("expected status 503, got %d", status)
	}
	if len(responses) != 1 || responses[0].State != state {
		t.Fatalf("expected the degraded response to still be returned alongside the error, got %v", responses)
	}
}

func TestKeeperCheck(t *testing.T) {
	s, vaultService := createTestNodeService(t)
	adapter := &fakeKeeperAdapter{}
	keeperRegistry := utils.NewRegistry[keeper.PluginType, keeper.Plugin]()
	keeperRegistry.Register("fake", adapter)
	s.keeperRegistry = keeperRegistry

	key, _, errVault := vaultService.Create(vault.Vault{Type: vault.KEEPER_PASSWORD, Username: "root", Secret: "secret"})
	if errVault != nil {
		t.Fatalf("failed to seed vault: %v", errVault)
	}
	unknown := uuid.New()
	request := func(vaultId *uuid.UUID, certs *cert.Certs) KeeperOneRequest {
		return KeeperOneRequest{
			KeeperConnection: KeeperConnection{Host: "10.0.0.1", Port: 2379},
			KeeperOptions:    KeeperOptions{Cluster: "etcd", Plugin: "fake", VaultId: vaultId, Certs: certs},
		}
	}

	tests := []struct {
		name     string
		request  KeeperOneRequest
		listErr  error
		expected []utils.CheckStep
	}{
		{
			name:    "nothing to resolve still asks the keeper",
			request: request(nil, nil),
			expected: []utils.CheckStep{
				{Stage: utils.CheckStageVault, Result: utils.CheckResultSkipped, Message: "no vault is configured"},
				{Stage: utils.CheckStageCert, Result: utils.CheckResultSkipped, Message: "tls is off"},
				{Stage: utils.CheckStageConnect, Result: utils.CheckResultDone},
			},
		},
		{
			name:    "a vault that cannot be read stops before the keeper",
			request: request(&unknown, nil),
			expected: []utils.CheckStep{
				{Stage: utils.CheckStageVault, Result: utils.CheckResultFailed},
				{Stage: utils.CheckStageCert, Result: utils.CheckResultSkipped, Message: "tls is off"},
				{Stage: utils.CheckStageConnect, Result: utils.CheckResultSkipped, Message: "vault failed"},
			},
		},
		{
			name:    "a certificate that cannot be loaded stops before the keeper",
			request: request(key, &cert.Certs{ClientCAId: &unknown}),
			expected: []utils.CheckStep{
				{Stage: utils.CheckStageVault, Result: utils.CheckResultDone},
				{Stage: utils.CheckStageCert, Result: utils.CheckResultFailed},
				{Stage: utils.CheckStageConnect, Result: utils.CheckResultSkipped, Message: "cert failed"},
			},
		},
		{
			name:    "the keeper's own refusal is its step's message",
			request: request(key, nil),
			listErr: errors.New("dial tcp 10.0.0.1:2379: i/o timeout"),
			expected: []utils.CheckStep{
				{Stage: utils.CheckStageVault, Result: utils.CheckResultDone},
				{Stage: utils.CheckStageCert, Result: utils.CheckResultSkipped, Message: "tls is off"},
				{Stage: utils.CheckStageConnect, Result: utils.CheckResultFailed, Message: "dial tcp 10.0.0.1:2379: i/o timeout"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter.lastRequest = nil
			adapter.listErr = tt.listErr
			steps, err := s.KeeperCheck(tt.request)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			assertCheckSteps(t, steps, tt.expected)
			asked := tt.expected[2].Result != utils.CheckResultSkipped
			if asked != (adapter.lastRequest != nil) {
				t.Fatalf("expected the keeper to be asked: %v, got asked: %v", asked, adapter.lastRequest != nil)
			}
			if adapter.lastRequest != nil {
				if adapter.lastRequest.Cluster != "etcd" || adapter.lastRequest.Host != "10.0.0.1" || adapter.lastRequest.Port != 2379 {
					t.Errorf("expected the keeper to be asked about the requested node, got %+v", adapter.lastRequest)
				}
				if tt.request.VaultId != nil && (adapter.lastRequest.Credentials == nil || adapter.lastRequest.Credentials.Username != "root") {
					t.Errorf("expected the vault's credentials to reach the keeper, got %+v", adapter.lastRequest.Credentials)
				}
			}
		})
	}

	t.Run("an unknown plugin is refused before any step", func(t *testing.T) {
		r := request(nil, nil)
		r.Plugin = "unknown"
		if _, err := s.KeeperCheck(r); err == nil {
			t.Fatalf("expected an error for an unregistered plugin")
		}
	})
}

func assertCheckSteps(t *testing.T, steps []utils.CheckStep, expected []utils.CheckStep) {
	t.Helper()
	if len(steps) != len(expected) {
		t.Fatalf("expected %d steps, got %+v", len(expected), steps)
	}
	for i, want := range expected {
		got := steps[i]
		if got.Stage != want.Stage || got.Result != want.Result {
			t.Errorf("step %d = %s/%s, want %s/%s", i, got.Stage, got.Result, want.Stage, want.Result)
		}
		if want.Message != "" && got.Message != want.Message {
			t.Errorf("step %d message = %q, want %q", i, got.Message, want.Message)
		}
		if want.Result == utils.CheckResultFailed && got.Message == "" {
			t.Errorf("step %d failed without saying why", i)
		}
	}
}

func TestGetHostError(t *testing.T) {
	s := &Service{}
	cause := errors.New("connection refused")
	tests := []struct {
		name       string
		statusCode int
		expected   string
	}{
		{name: "the default 400 is left out", statusCode: http.StatusBadRequest, expected: `host "db1": connection refused`},
		{name: "no code is left out", statusCode: 0, expected: `host "db1": connection refused`},
		{name: "a code the keeper chose is kept", statusCode: http.StatusServiceUnavailable, expected: `host "db1" failed with code 503: connection refused`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.getHostError("db1", tt.statusCode, cause)
			if err.Error() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, err.Error())
			}
			if !errors.Is(err, cause) {
				t.Errorf("expected the cause to stay wrapped")
			}
		})
	}
}
