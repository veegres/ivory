package node

import (
	"io"
	"ivory/clients/console"
	"ivory/core/config"
	"ivory/core/service/vault"
	"ivory/plugins/platform"
	"testing"
)

func TestEscapeInterpolatedValue(t *testing.T) {
	esc := string(platform.ValueEscape)
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "plain value is unchanged", input: "postgres", expected: "postgres"},
		{name: "embedded single quote is escaped", input: "it's a test", expected: "it" + esc + "'s a test"},
		{name: "embedded double quote is escaped", input: `pa"ss`, expected: "pa" + esc + `"ss`},
		{name: "embedded backslash is escaped", input: `pa\ss`, expected: "pa" + esc + `\ss`},
		{name: "backslash and quote together", input: `pa\'ss`, expected: "pa" + esc + `\` + esc + `'ss`},
		{name: "a marker in the value is dropped, never trusted", input: "pa" + esc + "ss", expected: "pass"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := escapeInterpolatedValue(tt.input); got != tt.expected {
				t.Errorf("escapeInterpolatedValue(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
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
	platform.Adapter
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
