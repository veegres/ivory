package docker

import (
	"errors"
	"ivory/clients/console/ssh"
	"ivory/plugins/platform"
	"math"
	"strings"
	"testing"
)

// TestUpContainerRunsTheCommandAsWritten covers what the deployment model
// rests on: the command the user reads is the command that runs. Readability
// whitespace between flags collapses, but a quoted multi-line startup script -
// whose newlines are real statement separators - survives byte for byte.
func TestUpContainerRunsTheCommandAsWritten(t *testing.T) {
	adapter := NewPlugin(ssh.NewClient())
	connection := platform.Connection{}

	tests := []struct {
		name     string
		command  string
		expected string
	}{
		{
			name: "newlines between flags collapse to spaces",
			command: `docker run -d
  --name etcd-1
  --hostname 10.0.0.1
  -p 2379:2379
  quay.io/coreos/etcd:v3.6.5`,
			expected: `docker run -d --name etcd-1 --hostname 10.0.0.1 -p 2379:2379 quay.io/coreos/etcd:v3.6.5`,
		},
		{
			name: "a quoted multi-line entry script stays one argument",
			command: `docker run -d --name pg-2 postgres:18
  sh -c '
until pg_isready -h pg-1; do sleep 1; done
exec docker-entrypoint.sh postgres
'`,
			expected: `docker run -d --name pg-2 postgres:18 sh -c '
until pg_isready -h pg-1; do sleep 1; done
exec docker-entrypoint.sh postgres
'`,
		},
		{
			// NOTE: the verb is Ivory's, so the user's own is dropped rather
			// than duplicated - a command is usually copied in whole
			name:     "a restated docker run is not duplicated",
			command:  `docker run -d --name n1 redis:7`,
			expected: `docker run -d --name n1 redis:7`,
		},
		{
			name:     "options alone still become a docker run",
			command:  `-d --name n1 redis:7`,
			expected: `docker run -d --name n1 redis:7`,
		},
		{
			// NOTE: UpContainer starts a container and nothing else - a command
			// text naming its own executable would make it a remote shell
			name:     "another executable cannot replace the verb",
			command:  `rm -rf /`,
			expected: `docker run rm -rf /`,
		},
		{
			name:     "a quoted value with a space stays one argument",
			command:  `docker run -d -e SCOPE="my cluster" postgres:18`,
			expected: `docker run -d -e 'SCOPE=my cluster' postgres:18`,
		},
		{
			name:     "tabs and carriage returns collapse like any other whitespace",
			command:  "docker run -d --name test \n\t --restart always  \r\n -p 80:80 redis:7",
			expected: `docker run -d --name test --restart always -p 80:80 redis:7`,
		},
		{
			name:     "the user's own --name is honoured, not overridden",
			command:  `docker run -d --name my-own-name postgres:18`,
			expected: `docker run -d --name my-own-name postgres:18`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := adapter.UpContainer(connection, platform.SplitCommand(test.command)).(*ssh.Command).Command
			if got != test.expected {
				t.Fatalf("expected %q, got %q", test.expected, got)
			}
		})
	}
}

func TestContainerCommandsQuoteShellArguments(t *testing.T) {
	adapter := NewPlugin(ssh.NewClient())
	connection := platform.Connection{}
	name := "foo; rm -rf /"

	tests := []struct {
		name     string
		command  string
		expected string
	}{
		{
			name:     "up re-quotes an argument holding shell metacharacters",
			command:  adapter.UpContainer(connection, platform.SplitCommand(`docker run -d --name "foo; rm -rf /" postgres:16`)).(*ssh.Command).Command,
			expected: `docker run -d --name 'foo; rm -rf /' postgres:16`,
		},
		{
			name:     "down quotes name",
			command:  adapter.DownContainer(connection, name).(*ssh.Command).Command,
			expected: `docker rm -- 'foo; rm -rf /'`,
		},
		{
			name:     "start quotes name",
			command:  adapter.StartContainer(connection, name).(*ssh.Command).Command,
			expected: `docker start -- 'foo; rm -rf /'`,
		},
		{
			name:     "stop quotes name",
			command:  adapter.StopContainer(connection, name).(*ssh.Command).Command,
			expected: `docker stop -- 'foo; rm -rf /'`,
		},
		{
			name:     "restart quotes name",
			command:  adapter.RestartContainer(connection, name).(*ssh.Command).Command,
			expected: `docker restart -- 'foo; rm -rf /'`,
		},
		{
			name:     "container logs quotes name",
			command:  adapter.LogsContainer(connection, name, 50, true).(*ssh.Command).Command,
			expected: `docker logs --tail 50 --follow -- 'foo; rm -rf /'`,
		},
		{
			name:     "exec quotes name and command fields",
			command:  adapter.ExecContainer(connection, name, platform.SplitCommand(`etcdctl user add 'root:se cret'`)).(*ssh.Command).Command,
			expected: `docker exec -- 'foo; rm -rf /' etcdctl user add 'root:se cret'`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.command != test.expected {
				t.Fatalf("expected %q, got %q", test.expected, test.command)
			}
		})
	}
}

// TestSplitCommandKeepsAValueInsideItsArgument is why the split happens before
// interpolation and lives outside the adapter: a value filled into an argument
// that is already separated is content, not syntax. Before, the value went into
// the command text and this tokenizer parsed it - a quote closed the template
// author's own span and the value came apart into several arguments.
func TestSplitCommandKeepsAValueInsideItsArgument(t *testing.T) {
	adapter := NewPlugin(ssh.NewClient())
	connection := platform.Connection{}

	hostile := `it's a "test" \ $HOME ` + "`id`"
	fill := func(command []string) []string {
		for i, argument := range command {
			command[i] = strings.ReplaceAll(argument, "{{dbPass}}", hostile)
		}
		return command
	}

	tests := []struct {
		name     string
		run      func() string
		expected string
	}{
		{
			name: "run keeps the value in one argument",
			run: func() string {
				return adapter.UpContainer(connection, fill(platform.SplitCommand(`docker run -d -e P="{{dbPass}}" img`))).(*ssh.Command).Command
			},
			expected: `docker run -d -e 'P=it'\''s a "test" \ $HOME ` + "`id`" + `' img`,
		},
		{
			name: "exec keeps the value in one argument",
			run: func() string {
				return adapter.ExecContainer(connection, "n1", fill(platform.SplitCommand(`etcdctl user add root:{{dbPass}}`))).(*ssh.Command).Command
			},
			expected: `docker exec -- 'n1' etcdctl user add 'root:it'\''s a "test" \ $HOME ` + "`id`" + `'`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.run(); got != test.expected {
				t.Fatalf("expected %q, got %q", test.expected, got)
			}
		})
	}
}

// TestExecContainerKeepsAuthorBackslashInSingleQuotedSpan pins the other half
// of that contract: inside a single-quoted span a backslash belongs to the
// template author and is a literal character, so a post script's own \" has
// to reach the inner `sh -c` still escaped. Consuming it here is what turned
// mongo's rs.initiate into `{_id: mongo-cluster}` - bare identifiers where
// strings belonged - once the remote shell parsed the script a second time.
func TestExecContainerKeepsAuthorBackslashInSingleQuotedSpan(t *testing.T) {
	adapter := NewPlugin(ssh.NewClient())
	connection := platform.Connection{}

	command := `sh -c 'mongosh --eval "rs.initiate({_id: \"rs0\"})"'`
	got := adapter.ExecContainer(connection, "n1", platform.SplitCommand(command)).(*ssh.Command).Command
	expected := `docker exec -- 'n1' sh -c 'mongosh --eval "rs.initiate({_id: \"rs0\"})"'`
	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

// TestSplitShellFieldsAppliesShellBackslashRules covers the quote states a
// backslash means different things in, so the tokenizer keeps matching a real
// shell rather than the template author having to guess.
func TestSplitCommandAppliesShellBackslashRules(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single-quoted span keeps the backslash literal",
			input:    `sh -c 'echo \"x\"'`,
			expected: []string{"sh", "-c", `echo \"x\"`},
		},
		{
			name:     "double-quoted span consumes an escape it interprets",
			input:    `-e A="say \"hi\""`,
			expected: []string{"-e", `A=say "hi"`},
		},
		{
			name:     "double-quoted span keeps a backslash before an ordinary rune",
			input:    `-e A="C:\path"`,
			expected: []string{"-e", `A=C:\path`},
		},
		{
			name:     "unquoted backslash escapes the next rune",
			input:    `-e A=one\ two`,
			expected: []string{"-e", "A=one two"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := platform.SplitCommand(test.input)
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
}

func TestMetricsContainerQuotesShellArguments(t *testing.T) {
	name := "foo; rm -rf /"

	command := "docker stats --no-stream --format " + shellQuote("{{json .}}") + " -- " + shellQuote(name)
	expected := `docker stats --no-stream --format '{{json .}}' -- 'foo; rm -rf /'`
	if command != expected {
		t.Fatalf("expected %q, got %q", expected, command)
	}
}

func TestParseContainerMetrics(t *testing.T) {
	output := `{"CPUPerc":"12.34%","MemUsage":"20.5MiB / 1.952GiB","NetIO":"1.2kB / 648B"}`

	adapter := NewPlugin(ssh.NewClient())
	metrics, err := adapter.parseContainerMetrics(output)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if metrics.Cpu.TotalTicks != containerCpuScale {
		t.Fatalf("expected total ticks %d, got %v", containerCpuScale, metrics.Cpu.TotalTicks)
	}
	if metrics.Cpu.IdleTicks != containerCpuScale-1234 {
		t.Fatalf("expected idle ticks %d, got %v", containerCpuScale-1234, metrics.Cpu.IdleTicks)
	}

	expectedLimit := uint64(math.Round(1.952 * 1024 * 1024 * 1024))
	if metrics.Memory.TotalBytes != expectedLimit {
		t.Fatalf("expected total bytes %d, got %v", expectedLimit, metrics.Memory.TotalBytes)
	}
	expectedUsed := uint64(math.Round(20.5 * 1024 * 1024))
	if metrics.Memory.AvailableBytes != expectedLimit-expectedUsed {
		t.Fatalf("expected available bytes %d, got %v", expectedLimit-expectedUsed, metrics.Memory.AvailableBytes)
	}

	if metrics.Network.ReceivedBytes != 1200 {
		t.Fatalf("expected received bytes 1200, got %v", metrics.Network.ReceivedBytes)
	}
	if metrics.Network.TransmittedBytes != 648 {
		t.Fatalf("expected transmitted bytes 648, got %v", metrics.Network.TransmittedBytes)
	}
}

func TestParseContainerMetricsCpuClampsAboveHundredPercent(t *testing.T) {
	output := `{"CPUPerc":"230.50%","MemUsage":"1B / 2B","NetIO":"0B / 0B"}`

	adapter := NewPlugin(ssh.NewClient())
	metrics, err := adapter.parseContainerMetrics(output)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if metrics.Cpu.IdleTicks != 0 {
		t.Fatalf("expected idle ticks 0, got %v", metrics.Cpu.IdleTicks)
	}
}

func TestParseContainerMetricsErrorsAndEdgeCases(t *testing.T) {
	adapter := NewPlugin(ssh.NewClient())

	t.Run("empty output", func(t *testing.T) {
		_, err := adapter.parseContainerMetrics("   ")
		if !errors.Is(err, platform.ErrContainerNotRunning) {
			t.Errorf("expected ErrContainerNotRunning, got %v", err)
		}
	})

	t.Run("malformed json", func(t *testing.T) {
		_, err := adapter.parseContainerMetrics("not json")
		if !errors.Is(err, platform.ErrInvalidContainerMetrics) {
			t.Errorf("expected ErrInvalidContainerMetrics, got %v", err)
		}
	})

	t.Run("malformed cpu", func(t *testing.T) {
		_, err := adapter.parseContainerMetrics(`{"CPUPerc":"abc","MemUsage":"1B / 2B","NetIO":"0B / 0B"}`)
		if !errors.Is(err, platform.ErrInvalidContainerMetrics) {
			t.Errorf("expected ErrInvalidContainerMetrics, got %v", err)
		}
	})

	t.Run("malformed memory - missing slash", func(t *testing.T) {
		_, err := adapter.parseContainerMetrics(`{"CPUPerc":"1%","MemUsage":"1B","NetIO":"0B / 0B"}`)
		if !errors.Is(err, platform.ErrInvalidContainerMetrics) {
			t.Errorf("expected ErrInvalidContainerMetrics, got %v", err)
		}
	})

	t.Run("malformed network - bad unit", func(t *testing.T) {
		_, err := adapter.parseContainerMetrics(`{"CPUPerc":"1%","MemUsage":"1B / 2B","NetIO":"1XB / 0B"}`)
		if !errors.Is(err, platform.ErrInvalidContainerMetrics) {
			t.Errorf("expected ErrInvalidContainerMetrics, got %v", err)
		}
	})

	t.Run("zero memory limit", func(t *testing.T) {
		_, err := adapter.parseContainerMetrics(`{"CPUPerc":"1%","MemUsage":"0B / 0B","NetIO":"0B / 0B"}`)
		if !errors.Is(err, platform.ErrInvalidContainerMetrics) {
			t.Errorf("expected ErrInvalidContainerMetrics, got %v", err)
		}
	})
}
