package node

import (
	"io"
	"ivory/clients/console"
	"ivory/core/config"
	"ivory/core/service/vault"
	"ivory/plugins/keeper"
	"ivory/plugins/platform"
	"testing"
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
	calls []string
}

func (r *recordingPlatform) SupportedFeatures() map[config.Feature]bool { return nil }

func (r *recordingPlatform) StopContainer(platform.Connection, string) console.Command {
	return recordingCommand{name: "stop", calls: &r.calls}
}

func (r *recordingPlatform) DownContainer(platform.Connection, string) console.Command {
	return recordingCommand{name: "rm", calls: &r.calls}
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
