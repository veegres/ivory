package linux

import (
	"errors"
	"ivory/src/clients/ssh"
	"ivory/src/plugins/os"
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
		if !errors.Is(err, os.ErrInvalidCpuMetrics) {
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
		if !errors.Is(err, os.ErrInvalidMemoryMetrics) {
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
		if !errors.Is(err, os.ErrInvalidNetworkMetrics) {
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
