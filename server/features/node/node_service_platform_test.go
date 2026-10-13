package node

import (
	"errors"
	"io"
	"ivory/clients/console"
	"ivory/core/config"
	"ivory/core/service/vault"
	"ivory/core/utils"
	"ivory/plugins/keeper"
	"ivory/plugins/platform"
	"testing"
	"time"

	"github.com/google/uuid"
)

// TestResolveCommand covers what replaced the escaper: the command is split
// into arguments first and the values go into arguments that are already
// separated, so a credential is content within its argument and nothing about
// it can be read as syntax. The platform is handed the finished arguments and
// never sees a placeholder.
func TestResolveCommand(t *testing.T) {
	values := keeper.Values{
		Name:   "etcd1",
		DbUser: "root",
		DbPass: `it's a "test" \ $HOME ` + "`id`",
		DbPort: "2379",
	}

	tests := []struct {
		name     string
		template string
		expected []string
	}{
		{
			name:     "readability newlines collapse between arguments",
			template: "run -d\n  --name {{name}}\n  img",
			expected: []string{"run", "-d", "--name", "etcd1", "img"},
		},
		{
			name:     "a hostile value stays one argument",
			template: `etcdctl --endpoints=http://localhost:{{dbPort}} user add {{dbUser}}:{{dbPass}}`,
			expected: []string{"etcdctl", "--endpoints=http://localhost:2379", "user", "add", `root:it's a "test" \ $HOME ` + "`id`"},
		},
		{
			name:     "a value carrying a space does not become two arguments",
			template: `-e P={{dbPass}} img`,
			expected: []string{"-e", `P=it's a "test" \ $HOME ` + "`id`", "img"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := resolveCommand(test.template, values)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if len(got) != len(test.expected) {
				t.Fatalf("expected %q, got %q", test.expected, got)
			}
			for i := range got {
				if got[i] != test.expected[i] {
					t.Fatalf("expected %q, got %q", test.expected, got)
				}
			}
		})
	}

	t.Run("a variable with no value is reported, not run", func(t *testing.T) {
		if _, err := resolveCommand("run --name {{name}} --cluster {{cluster}}", values); err == nil {
			t.Fatal("expected an unresolved placeholder to be refused")
		}
	})
}

// recordingCommand records that it ran, so a service test can assert which
// platform primitives were issued and in what order.
type recordingCommand struct {
	name  string
	calls *[]string
}

func (c recordingCommand) Id() string      { return c.name }
func (c recordingCommand) KeepAlive() bool { return false }
func (c recordingCommand) Persist() bool   { return false }
func (c recordingCommand) Abort() error    { return nil }
func (c recordingCommand) Wait() error     { return nil }

func (c recordingCommand) Start() (io.Reader, error) { return nil, nil }

func (c recordingCommand) Execute() ([]string, error) {
	*c.calls = append(*c.calls, c.name)
	return []string{c.name}, nil
}

type recordingPlatform struct {
	platform.Plugin
	calls   []string
	infoErr error
}

func (r *recordingPlatform) Info(platform.Connection) ([]platform.InfoItem, error) {
	r.calls = append(r.calls, "info")
	return nil, r.infoErr
}

func (r *recordingPlatform) SupportedFeatures() map[config.Feature]bool { return nil }

func (r *recordingPlatform) CopyId(platform.Connection, string) error {
	r.calls = append(r.calls, "copy-id")
	return nil
}

func (r *recordingPlatform) StopContainer(platform.Connection, string) console.Command {
	return recordingCommand{name: "stop", calls: &r.calls}
}

func (r *recordingPlatform) DownContainer(platform.Connection, string) console.Command {
	return recordingCommand{name: "rm", calls: &r.calls}
}

func (r *recordingPlatform) InspectContainer(_ platform.Connection, name string) (*platform.ContainerInspect, error) {
	r.calls = append(r.calls, "inspect "+name)
	return &platform.ContainerInspect{
		Status: "exited", ExitCode: 137, OomKilled: true, Error: "oom",
		StartedAt: time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC), FinishedAt: time.Date(2026, 10, 1, 10, 5, 0, 0, time.UTC),
		RestartCount: 2, RestartPolicy: "on-failure:5", NetworkMode: "host", Image: "etcd:3.6",
	}, nil
}

type stubPlatformRegistry struct{ plugin platform.Plugin }

func (s stubPlatformRegistry) Get(platform.PluginType) (platform.Plugin, error) {
	return s.plugin, nil
}

// TestPlatformContainerDownStopsBeforeRemoving pins that Down is usable on a
// running deployment: the platform refuses to remove one, so removing without
// stopping first failed on exactly the deployments a user asks to remove.
func TestPlatformContainerDownStopsBeforeRemoving(t *testing.T) {
	s, vaultService := createTestNodeService(t)

	key, _, errCreate := vaultService.Create(vault.Vault{Type: vault.SSH_KEY, Username: "root"})
	if errCreate != nil {
		t.Fatalf("failed to seed ssh key vault: %v", errCreate)
	}

	recorder := &recordingPlatform{}
	s.platformRegistry = stubPlatformRegistry{plugin: recorder}

	logs, err := s.PlatformContainerDown(PlatformActionRequest{
		Name:       "db1",
		Connection: PlatformVaultConnection{Host: "host1", Port: 22, VaultId: *key},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(recorder.calls) != 2 || recorder.calls[0] != "stop" || recorder.calls[1] != "rm" {
		t.Fatalf("expected a stop followed by a remove, got %v", recorder.calls)
	}
	if len(logs) != 2 {
		t.Fatalf("expected both commands' output, got %v", logs)
	}
}

// TestPlatformSystemCopyIdRequiresPlatform pins that the adapter is the
// caller's answer rather than a hardcoded default: a PlatformCredConnection is
// only ever built from a request, so there is no cluster stored before
// platforms were selectable to fall back for.
func TestPlatformSystemCopyIdRequiresPlatform(t *testing.T) {
	s, _ := createTestNodeService(t)
	s.platformRegistry = stubPlatformRegistry{plugin: &recordingPlatform{}}

	t.Run("a request naming no platform is rejected", func(t *testing.T) {
		if _, err := s.PlatformSystemCopyId(PlatformCopyIdRequest{
			PlatformCredConnection: PlatformCredConnection{Host: "host1", Port: 22, Username: "root", Password: "secret"},
			PublicKey:              "ssh-rsa key",
		}); !errors.Is(err, ErrPlatformRequired) {
			t.Fatalf("expected ErrPlatformRequired, got %v", err)
		}
	})

	t.Run("a request naming its platform runs against it", func(t *testing.T) {
		recorder := &recordingPlatform{}
		s.platformRegistry = stubPlatformRegistry{plugin: recorder}

		if _, err := s.PlatformSystemCopyId(PlatformCopyIdRequest{
			PlatformCredConnection: PlatformCredConnection{Host: "host1", Port: 22, Username: "root", Password: "secret", Platform: platform.Docker},
			PublicKey:              "ssh-rsa key",
		}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(recorder.calls) != 1 || recorder.calls[0] != "copy-id" {
			t.Fatalf("expected the key to be installed once, got %v", recorder.calls)
		}
	})
}

func TestPlatformContainerInspectReportsEveryField(t *testing.T) {
	s, vaultService := createTestNodeService(t)

	key, _, errCreate := vaultService.Create(vault.Vault{Type: vault.SSH_KEY, Username: "root"})
	if errCreate != nil {
		t.Fatalf("failed to seed ssh key vault: %v", errCreate)
	}

	recorder := &recordingPlatform{}
	s.platformRegistry = stubPlatformRegistry{plugin: recorder}

	inspect, err := s.PlatformContainerInspect(PlatformActionRequest{
		Name:       "etcd2",
		Connection: PlatformVaultConnection{Host: "host1", Port: 22, VaultId: *key},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(recorder.calls) != 1 || recorder.calls[0] != "inspect etcd2" {
		t.Fatalf("expected one inspect of the named container, got %v", recorder.calls)
	}
	expected := PlatformInspectResponse{
		Status: "exited", ExitCode: 137, OomKilled: true, Error: "oom",
		StartedAt: time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC), FinishedAt: time.Date(2026, 10, 1, 10, 5, 0, 0, time.UTC),
		RestartCount: 2, RestartPolicy: "on-failure:5", NetworkMode: "host", Image: "etcd:3.6",
	}
	if *inspect != expected {
		t.Fatalf("expected %+v, got %+v", expected, *inspect)
	}
}

func TestPlatformCheck(t *testing.T) {
	s, vaultService := createTestNodeService(t)
	key, _, errCreate := vaultService.Create(vault.Vault{Type: vault.SSH_KEY, Username: "root"})
	if errCreate != nil {
		t.Fatalf("failed to seed ssh key vault: %v", errCreate)
	}
	unknown := *key
	unknown[0]++

	tests := []struct {
		name     string
		vaultId  uuid.UUID
		infoErr  error
		expected []utils.CheckStep
	}{
		{
			name:    "a readable vault reaches the platform",
			vaultId: *key,
			expected: []utils.CheckStep{
				{Stage: utils.CheckStageVault, Result: utils.CheckResultDone},
				{Stage: utils.CheckStageConnect, Result: utils.CheckResultDone},
			},
		},
		{
			name:    "a vault that cannot be read stops before the platform",
			vaultId: unknown,
			expected: []utils.CheckStep{
				{Stage: utils.CheckStageVault, Result: utils.CheckResultFailed},
				{Stage: utils.CheckStageConnect, Result: utils.CheckResultSkipped, Message: "vault failed"},
			},
		},
		{
			name:    "the platform's own refusal is its step's message",
			vaultId: *key,
			infoErr: errors.New("ssh: handshake failed"),
			expected: []utils.CheckStep{
				{Stage: utils.CheckStageVault, Result: utils.CheckResultDone},
				{Stage: utils.CheckStageConnect, Result: utils.CheckResultFailed, Message: "ssh: handshake failed"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := &recordingPlatform{infoErr: tt.infoErr}
			s.platformRegistry = stubPlatformRegistry{plugin: recorder}
			steps, err := s.PlatformCheck(PlatformVaultConnection{Host: "host1", Port: 22, VaultId: tt.vaultId})
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			assertCheckSteps(t, steps, tt.expected)
			reached := tt.expected[1].Result != utils.CheckResultSkipped
			if reached != (len(recorder.calls) == 1) {
				t.Fatalf("expected the platform to be reached: %v, got calls %v", reached, recorder.calls)
			}
		})
	}
}
