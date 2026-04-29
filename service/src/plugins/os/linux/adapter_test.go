package linux

import (
	"ivory/src/clients/ssh"
	"testing"
)

func TestParseMetrics(t *testing.T) {
	output := `__IVORY_CPU_1__
cpu  100 0 100 800 0 0 0 0 0 0
__IVORY_CPU_2__
cpu  150 0 150 900 0 0 0 0 0 0
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

	if metrics.Cpu.UsagePercent != 50 {
		t.Fatalf("expected cpu usage 50, got %v", metrics.Cpu.UsagePercent)
	}
	if metrics.Memory.TotalBytes != 1024*1024 {
		t.Fatalf("unexpected total memory: %d", metrics.Memory.TotalBytes)
	}
	if metrics.Memory.UsedBytes != 768*1024 {
		t.Fatalf("unexpected used memory: %d", metrics.Memory.UsedBytes)
	}
	if metrics.Network.ReceivedBytes != 2048 {
		t.Fatalf("unexpected received bytes: %d", metrics.Network.ReceivedBytes)
	}
	if metrics.Network.TransmittedBytes != 4096 {
		t.Fatalf("unexpected transmitted bytes: %d", metrics.Network.TransmittedBytes)
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
