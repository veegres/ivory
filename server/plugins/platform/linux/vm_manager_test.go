package linux

import (
	"errors"
	"ivory/clients/console/ssh"
	"ivory/plugins/platform"
	"strings"
	"testing"
)

func TestParseMetrics(t *testing.T) {
	output := `__IVORY_CPU__
cpu  100 0 100 800 0 0 0 0 0 0
__IVORY_MEM__
MemTotal:       1024 kB
MemAvailable:    256 kB
__IVORY_NET__
Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
  lo: 100 1 0 0 0 0 0 0 100 1 0 0 0 0 0 0
eth0: 2048 2 0 0 0 0 0 0 4096 4 0 0 0 0 0 0
`
	sshClient := ssh.NewClient()
	adapter := NewAdapter(sshClient)
	metrics, err := adapter.parseMetrics(output)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if metrics.Cpu.TotalTicks != 1000 {
		t.Fatalf("expected total ticks 1000, got %v", metrics.Cpu.TotalTicks)
	}
	if metrics.Cpu.IdleTicks != 800 {
		t.Fatalf("expected idle ticks 800, got %v", metrics.Cpu.IdleTicks)
	}
	if metrics.Memory.TotalBytes != 1024*1024 {
		t.Fatalf("unexpected total memory: %d", metrics.Memory.TotalBytes)
	}
	if metrics.Network.ReceivedBytes != 2048 {
		t.Fatalf("unexpected received bytes: %d", metrics.Network.ReceivedBytes)
	}
	if metrics.Network.TransmittedBytes != 4096 {
		t.Fatalf("unexpected transmitted bytes: %d", metrics.Network.TransmittedBytes)
	}
}

func TestParseMetricsErrorsAndEdgeCases(t *testing.T) {
	adapter := NewAdapter(ssh.NewClient())

	t.Run("multi-interface", func(t *testing.T) {
		output := `__IVORY_CPU__
cpu  100 0 100 800 0 0 0 0 0 0
__IVORY_MEM__
MemTotal:       1024 kB
MemAvailable:    256 kB
__IVORY_NET__
Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
eth0: 1000 1 0 0 0 0 0 0 2000 1 0 0 0 0 0 0
eth1: 500 1 0 0 0 0 0 0 300 1 0 0 0 0 0 0
`
		m, err := adapter.parseMetrics(output)
		if err != nil {
			t.Fatal(err)
		}
		if m.Network.ReceivedBytes != 1500 {
			t.Errorf("expected 1500 rx, got %d", m.Network.ReceivedBytes)
		}
		if m.Network.TransmittedBytes != 2300 {
			t.Errorf("expected 2300 tx, got %d", m.Network.TransmittedBytes)
		}
	})

	t.Run("malformed-cpu", func(t *testing.T) {
		output := `__IVORY_CPU__
cpu 1 2 3
__IVORY_MEM__
MemTotal: 1024 kB
MemAvailable: 256 kB
__IVORY_NET__
eth0: 1 1 0 0 0 0 0 0 1 1 0 0 0 0 0 0
`
		_, err := adapter.parseMetrics(output)
		if !errors.Is(err, platform.ErrInvalidCpuMetrics) {
			t.Errorf("expected ErrInvalidCpuMetrics, got %v", err)
		}
	})

	t.Run("malformed-memory", func(t *testing.T) {
		output := `__IVORY_CPU__
cpu  100 0 100 800 0 0 0 0 0 0
__IVORY_MEM__
MemTotal: 1024 kB
__IVORY_NET__
eth0: 1 1 0 0 0 0 0 0 1 1 0 0 0 0 0 0
`
		_, err := adapter.parseMetrics(output)
		if !errors.Is(err, platform.ErrInvalidMemoryMetrics) {
			t.Errorf("expected ErrInvalidMemoryMetrics, got %v", err)
		}
	})

	t.Run("malformed-network", func(t *testing.T) {
		output := `__IVORY_CPU__
cpu  100 0 100 800 0 0 0 0 0 0
__IVORY_MEM__
MemTotal: 1024 kB
MemAvailable: 256 kB
__IVORY_NET__
eth0: 1 2 3 4 5 6 7 8 9
`
		_, err := adapter.parseMetrics(output)
		if !errors.Is(err, platform.ErrInvalidNetworkMetrics) {
			t.Errorf("expected ErrInvalidNetworkMetrics, got %v", err)
		}
	})

	t.Run("only-loopback", func(t *testing.T) {
		output := `__IVORY_CPU__
cpu  100 0 100 800 0 0 0 0 0 0
__IVORY_MEM__
MemTotal: 1024 kB
MemAvailable: 256 kB
__IVORY_NET__
lo: 100 1 0 0 0 0 0 0 100 1 0 0 0 0 0 0
`
		m, err := adapter.parseMetrics(output)
		if err != nil {
			t.Fatal(err)
		}
		if m.Network.ReceivedBytes != 0 || m.Network.TransmittedBytes != 0 {
			t.Errorf("expected 0/0, got %d/%d", m.Network.ReceivedBytes, m.Network.TransmittedBytes)
		}
	})
}

func TestCopyIdRejectsNewlinePublicKey(t *testing.T) {
	adapter := NewAdapter(ssh.NewClient())

	err := adapter.CopyId(platform.Connection{}, "ssh-ed25519 AAAA attacker\nssh-ed25519 BBBB")
	if !errors.Is(err, ErrInvalidPublicKey) {
		t.Fatalf("expected ErrInvalidPublicKey, got %v", err)
	}
}

func TestShellQuote(t *testing.T) {
	got := shellQuote("foo'$(touch /tmp/pwned);bar")
	if got != `'foo'\''$(touch /tmp/pwned);bar'` {
		t.Fatalf("unexpected quote result: %q", got)
	}
	if strings.Contains(got[1:len(got)-1], "';") {
		t.Fatalf("expected semicolon to remain inside quoted argument: %q", got)
	}
}

func TestLogsQuotesShellArguments(t *testing.T) {
	adapter := NewAdapter(ssh.NewClient())
	connection := platform.Connection{}

	command := adapter.Logs(connection, `/tmp/app; rm -rf /`, 10, false).(*ssh.Command).Command
	expected := `tail -n 10 -- '/tmp/app; rm -rf /'`
	if command != expected {
		t.Fatalf("expected %q, got %q", expected, command)
	}
}

func TestParseProcesses(t *testing.T) {
	lines := []string{
		"1234 postgres 12.3 4.5 102400 4 /usr/lib/postgresql/16/bin/postgres postgres: main process",
		"   1 root      0.0  0.1   1024 1 init init",
	}

	adapter := NewAdapter(ssh.NewClient())
	processes, err := adapter.parseProcesses(lines)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(processes) != 2 {
		t.Fatalf("expected 2 processes, got %d", len(processes))
	}

	first := processes[0]
	if first.Pid != 1234 || first.User != "postgres" || first.CpuPercent != 12.3 || first.MemPercent != 4.5 {
		t.Fatalf("unexpected first process: %+v", first)
	}
	if first.MemoryBytes != 102400*1024 {
		t.Fatalf("expected memory bytes %d, got %d", 102400*1024, first.MemoryBytes)
	}
	if first.Threads != 4 {
		t.Fatalf("expected 4 threads, got %d", first.Threads)
	}
	if first.Program != "postgres" {
		t.Fatalf("expected program basename %q, got %q", "postgres", first.Program)
	}
	if first.Command != "postgres: main process" {
		t.Fatalf("unexpected command: %q", first.Command)
	}
}

func TestParseProcessesProgramStripsPath(t *testing.T) {
	lines := []string{"1 root 0.0 0.1 1024 1 /usr/lib/postgresql/16/bin/postgres /usr/lib/postgresql/16/bin/postgres -D /var/lib/postgresql/data"}

	adapter := NewAdapter(ssh.NewClient())
	processes, err := adapter.parseProcesses(lines)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if processes[0].Program != "postgres" {
		t.Fatalf("expected program basename %q, got %q", "postgres", processes[0].Program)
	}
	if processes[0].Command != "/usr/lib/postgresql/16/bin/postgres -D /var/lib/postgresql/data" {
		t.Fatalf("expected full command line, got %q", processes[0].Command)
	}
}

func TestParseProcessesSkipsMalformedLines(t *testing.T) {
	lines := []string{
		"1234 postgres 12.3 4.5 102400 4 postgres postgres",
		"garbage line that is too short",
		"",
		"   ",
	}

	adapter := NewAdapter(ssh.NewClient())
	processes, err := adapter.parseProcesses(lines)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(processes) != 1 {
		t.Fatalf("expected 1 process after skipping malformed lines, got %d", len(processes))
	}
}

func TestParseProcessesAllMalformedReturnsError(t *testing.T) {
	lines := []string{"not a valid ps line"}

	adapter := NewAdapter(ssh.NewClient())
	_, err := adapter.parseProcesses(lines)
	if !errors.Is(err, platform.ErrInvalidProcesses) {
		t.Fatalf("expected ErrInvalidProcesses, got %v", err)
	}
}

func TestParseProcessesEmptyInput(t *testing.T) {
	adapter := NewAdapter(ssh.NewClient())
	processes, err := adapter.parseProcesses([]string{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(processes) != 0 {
		t.Fatalf("expected 0 processes, got %d", len(processes))
	}
}
