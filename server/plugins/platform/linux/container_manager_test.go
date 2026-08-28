package linux

import (
	"errors"
	"ivory/clients/console/ssh"
	"ivory/plugins/platform"
	"math"
	"testing"
)

// TestUpContainerRunsTheCommandAsWritten covers what the deployment model
// rests on: the command the user reads is the command that runs. Readability
// whitespace between flags collapses, but a quoted multi-line startup script -
// whose newlines are real statement separators - survives byte for byte.
func TestUpContainerRunsTheCommandAsWritten(t *testing.T) {
	adapter := NewAdapter(ssh.NewClient())
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
			name:     "a leading docker prefix is not duplicated",
			command:  `docker run -d --name n1 redis:7`,
			expected: `docker run -d --name n1 redis:7`,
		},
		{
			name:     "a sudo docker prefix is normalized away",
			command:  `sudo docker run -d --name n1 redis:7`,
			expected: `docker run -d --name n1 redis:7`,
		},
		{
			name:     "a command without the prefix gets one",
			command:  `run -d --name n1 redis:7`,
			expected: `docker run -d --name n1 redis:7`,
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
			got := adapter.UpContainer(connection, test.command).(*ssh.Command).Command
			if got != test.expected {
				t.Fatalf("expected %q, got %q", test.expected, got)
			}
		})
	}
}

func TestNormalizeDockerCommand(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected string
	}{
		{name: "prefix plain command", command: "ps", expected: "docker ps"},
		{name: "keep docker command", command: "docker ps", expected: "docker ps"},
		{name: "keep sudo docker command", command: "sudo docker ps", expected: "sudo docker ps"},
		{name: "trim spaces", command: "  images  ", expected: "docker images"},
	}

	adapter := &Adapter{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := adapter.normalizeDockerCommand(test.command)
			if actual != test.expected {
				t.Fatalf("expected %q, got %q", test.expected, actual)
			}
		})
	}
}

func TestContainerCommandsQuoteShellArguments(t *testing.T) {
	adapter := NewAdapter(ssh.NewClient())
	connection := platform.Connection{}
	name := "foo; rm -rf /"

	tests := []struct {
		name     string
		command  string
		expected string
	}{
		{
			name:     "up re-quotes an argument holding shell metacharacters",
			command:  adapter.UpContainer(connection, `docker run -d --name "foo; rm -rf /" postgres:16`).(*ssh.Command).Command,
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
			command:  adapter.ExecContainer(connection, name, `etcdctl user add 'root:se cret'`).(*ssh.Command).Command,
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

// TestExecContainerRecoversEscapedQuoteInHandWrittenSpan reproduces a
// plugin's own hand-written single-quoted span wrapped around an
// interpolated value (e.g. `user add '{{dbUser}}:{{dbPass}}'`): a password
// containing a quote and a space must survive tokenizing intact once the
// caller (node.getExecutionValues) has backslash-escaped it, instead of the
// naive quote-toggle parser closing the span early and mangling the value.
func TestExecContainerRecoversEscapedQuoteInHandWrittenSpan(t *testing.T) {
	adapter := NewAdapter(ssh.NewClient())
	connection := platform.Connection{}

	// value as produced by escapeInterpolatedValue("it's a test") interpolated
	// into etcd's own "'root:" + dbUser + ":" + dbPass + "'" template.
	command := `etcdctl user add 'root:it\'s a test'`
	got := adapter.ExecContainer(connection, "n1", command).(*ssh.Command).Command
	expected := `docker exec -- 'n1' etcdctl user add 'root:it'\''s a test'`
	// NOTE: the last field is shellQuote's own re-escaping of the recovered
	// literal value "root:it's a test" - what matters is that it was
	// recovered as one field, not split apart at the embedded quote.
	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestMetricsContainerQuotesShellArguments(t *testing.T) {
	adapter := NewAdapter(ssh.NewClient())
	name := "foo; rm -rf /"

	command := adapter.normalizeDockerCommand("stats --no-stream --format " + shellQuote("{{json .}}") + " -- " + shellQuote(name))
	expected := `docker stats --no-stream --format '{{json .}}' -- 'foo; rm -rf /'`
	if command != expected {
		t.Fatalf("expected %q, got %q", expected, command)
	}
}

func TestParseContainerMetrics(t *testing.T) {
	output := `{"CPUPerc":"12.34%","MemUsage":"20.5MiB / 1.952GiB","NetIO":"1.2kB / 648B"}`

	adapter := NewAdapter(ssh.NewClient())
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

	adapter := NewAdapter(ssh.NewClient())
	metrics, err := adapter.parseContainerMetrics(output)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if metrics.Cpu.IdleTicks != 0 {
		t.Fatalf("expected idle ticks 0, got %v", metrics.Cpu.IdleTicks)
	}
}

func TestParseContainerMetricsErrorsAndEdgeCases(t *testing.T) {
	adapter := NewAdapter(ssh.NewClient())

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
