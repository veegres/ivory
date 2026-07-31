package linux

import (
	"errors"
	"ivory/clients/console/ssh"
	"ivory/plugins/platform"
	"math"
	"testing"
)

func TestRenderOptions(t *testing.T) {
	adapter := &Adapter{}

	t.Run("full spec renders docker flags in order", func(t *testing.T) {
		spec := platform.DeploySpec{
			Name:          "{{host}}",
			Hostname:      "{{host}}",
			RestartPolicy: "unless-stopped",
			Ports:         []string{"{{keeperPort}}", "{{dbPort}}"},
			Volumes:       []platform.VolumeMount{{HostPath: "/data/postgres", ContainerPath: "/home/postgres/pgdata"}},
			Env: []platform.EnvVar{
				{Name: "SCOPE", Value: `"{{cluster}}"`},
				{Name: "PGPORT", Value: `{{dbPort}}`},
			},
		}
		expected := "--name {{host}}\n" +
			"--hostname {{host}}\n" +
			"--restart unless-stopped\n" +
			"-p {{keeperPort}}:{{keeperPort}}\n" +
			"-p {{dbPort}}:{{dbPort}}\n" +
			"-v /data/postgres:/home/postgres/pgdata\n" +
			`-e SCOPE="{{cluster}}"` + "\n" +
			"-e PGPORT={{dbPort}}"
		if got := adapter.RenderOptions(spec); got != expected {
			t.Fatalf("expected:\n%s\ngot:\n%s", expected, got)
		}
	})

	t.Run("host network renders without ports, volumes and restart policy", func(t *testing.T) {
		spec := platform.DeploySpec{
			Name:        "{{host}}",
			Hostname:    "{{host}}",
			HostNetwork: true,
			Env:         []platform.EnvVar{{Name: "SCOPE", Value: `"{{cluster}}"`}},
		}
		expected := "--name {{host}}\n" +
			"--hostname {{host}}\n" +
			"--network host\n" +
			`-e SCOPE="{{cluster}}"`
		if got := adapter.RenderOptions(spec); got != expected {
			t.Fatalf("expected:\n%s\ngot:\n%s", expected, got)
		}
	})

	t.Run("empty spec renders empty string", func(t *testing.T) {
		if got := adapter.RenderOptions(platform.DeploySpec{}); got != "" {
			t.Fatalf("expected empty string, got %q", got)
		}
	})
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
			name:     "up quotes options and image",
			command:  adapter.UpContainer(connection, name, `--name foo;rm -rf / -e POSTGRES_PASSWORD=pass`, `postgres:16; reboot`, "").(*ssh.Command).Command,
			expected: `docker run -d --name 'foo; rm -rf /' '-rf' '/' '-e' 'POSTGRES_PASSWORD=pass' -- 'postgres:16; reboot'`,
		},
		{
			name:     "up keeps a quoted value with a space as a single argument",
			command:  adapter.UpContainer(connection, name, `-e SCOPE="my cluster" -e NAME=test`, `postgres:16`, "").(*ssh.Command).Command,
			expected: `docker run -d --name 'foo; rm -rf /' '-e' 'SCOPE=my cluster' '-e' 'NAME=test' -- 'postgres:16'`,
		},
		{
			name:     "up always enforces the given name, discarding any --name in options",
			command:  adapter.UpContainer(connection, "mycontainer", `--name foo`, `postgres:18`, `sh -c 'echo "hi"'`).(*ssh.Command).Command,
			expected: `docker run -d --name 'mycontainer' -- 'postgres:18' 'sh' '-c' 'echo "hi"'`,
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
			expected: `docker exec -- 'foo; rm -rf /' 'etcdctl' 'user' 'add' 'root:se cret'`,
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
	expected := `docker exec -- 'n1' 'etcdctl' 'user' 'add' 'root:it'\''s a test'`
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
