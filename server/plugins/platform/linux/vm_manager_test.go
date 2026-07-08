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

func TestParseInfo(t *testing.T) {
	output := `__IVORY_HOST__
product_name=21KDA04PCD
product_family=ThinkPad X1 Carbon Gen 12
hostname=myhost
__IVORY_OS__
NAME="Ubuntu"
VERSION="22.04.3 LTS (Jammy Jellyfish)"
PRETTY_NAME="Ubuntu 22.04.3 LTS"
__IVORY_KERNEL__
Linux 5.15.0-91-generic
__IVORY_UPTIME__
90061.35 12345.6
__IVORY_CPU__
model name	: Intel(R) Xeon(R) CPU E5-2670 v3 @ 2.30GHz
8
__IVORY_GPU__
name=0000:00:02.0 VGA compatible controller: Intel Corporation Meteor Lake-P [Intel Arc Graphics] (rev 08)
addr=0000:00:02.0
freq_mhz=2250
__IVORY_MEM__
MemTotal:       16777216 kB
__IVORY_SWAP__
SwapTotal:       2097152 kB
__IVORY_DISK__
/dev/sda1 20971520 8388608 12147712 41% /
__IVORY_IP__
192.168.1.10 172.17.0.1
__IVORY_LOCALE__
LANG=en_US.UTF-8
`
	adapter := NewAdapter(ssh.NewClient())
	items := adapter.parseInfo(output)

	got := map[string]string{}
	for _, item := range items {
		got[item.Key] = item.Value
	}

	expected := map[string]string{
		"Host":      "21KDA04PCD (ThinkPad X1 Carbon Gen 12)",
		"OS":        "Ubuntu 22.04.3 LTS",
		"Kernel":    "Linux 5.15.0-91-generic",
		"Uptime":    "1d 1h 1m",
		"CPU":       "Intel(R) Xeon(R) CPU E5-2670 v3 @ 2.30GHz",
		"CPU Cores": "8 cores",
		"GPU":       "Intel Arc Graphics @ 2.25 GHz [Integrated]",
		"Memory":    "16.0 GiB",
		"Swap":      "2.0 GiB",
		"Disk":      "20.0 GiB",
		"Local IP":  "192.168.1.10, 172.17.0.1",
		"Locale":    "en_US.UTF-8",
	}

	if len(items) != len(expected) {
		t.Fatalf("expected %d items, got %d: %+v", len(expected), len(items), items)
	}
	for key, value := range expected {
		if got[key] != value {
			t.Errorf("expected %q=%q, got %q", key, value, got[key])
		}
	}
}

// TestParseInfoOmitsMissingSections verifies parseInfo is best-effort: a
// platform-info dump is display-only, so absent/empty sections (e.g. no GPU,
// no swap configured) are simply left out rather than erroring.
func TestParseInfoOmitsMissingSections(t *testing.T) {
	output := `__IVORY_HOST__
product_name=
product_family=
hostname=myhost
__IVORY_OS__
PRETTY_NAME="Ubuntu 22.04.3 LTS"
__IVORY_KERNEL__
Linux 5.15.0-91-generic
__IVORY_UPTIME__
60 0
__IVORY_CPU__
model name	: Intel(R) Xeon(R) CPU E5-2670 v3 @ 2.30GHz
1
__IVORY_GPU__
name=
addr=
freq_mhz=
__IVORY_MEM__
MemTotal:       1048576 kB
__IVORY_SWAP__
SwapTotal:       0 kB
__IVORY_DISK__
/dev/sda1 1048576 524288 524288 50% /
__IVORY_IP__
__IVORY_LOCALE__
`
	adapter := NewAdapter(ssh.NewClient())
	items := adapter.parseInfo(output)

	got := map[string]string{}
	for _, item := range items {
		got[item.Key] = item.Value
	}

	for _, missing := range []string{"GPU", "Swap", "Local IP", "Locale"} {
		if _, ok := got[missing]; ok {
			t.Errorf("expected %q to be omitted, got %q", missing, got[missing])
		}
	}
	if got["Host"] != "myhost" {
		t.Errorf("expected Host=myhost, got %q", got["Host"])
	}
	if got["Uptime"] != "1m" {
		t.Errorf("expected Uptime=1m, got %q", got["Uptime"])
	}
}

func TestParseGpuModel(t *testing.T) {
	adapter := NewAdapter(ssh.NewClient())

	tests := map[string]struct {
		line     string
		expected string
	}{
		"intel bracket already has vendor": {
			line:     "00:02.0 VGA compatible controller: Intel Corporation Meteor Lake-P [Intel Arc Graphics] (rev 08)",
			expected: "Intel Arc Graphics",
		},
		"nvidia bracket missing vendor": {
			line:     "01:00.0 VGA compatible controller: NVIDIA Corporation GP104 [GeForce GTX 1080] (rev a1)",
			expected: "NVIDIA GeForce GTX 1080",
		},
		"no bracket falls back to raw description": {
			line:     "00:02.0 VGA compatible controller: Cirrus Logic GD 5446 (rev a1)",
			expected: "Cirrus Logic GD 5446",
		},
		"empty": {
			line:     "",
			expected: "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := adapter.parseGpuModel(tt.line)
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestParseGpu(t *testing.T) {
	adapter := NewAdapter(ssh.NewClient())

	tests := map[string]struct {
		lines    []string
		expected string
	}{
		"intel integrated at reserved pci slot": {
			lines: []string{
				"name=0000:00:02.0 VGA compatible controller: Intel Corporation Meteor Lake-P [Intel Arc Graphics] (rev 08)",
				"addr=0000:00:02.0",
				"freq_mhz=2250",
			},
			expected: "Intel Arc Graphics @ 2.25 GHz [Integrated]",
		},
		"nvidia is always discrete": {
			lines: []string{
				"name=0000:01:00.0 VGA compatible controller: NVIDIA Corporation GP104 [GeForce GTX 1080] (rev a1)",
				"addr=0000:01:00.0",
				"freq_mhz=",
			},
			expected: "NVIDIA GeForce GTX 1080 [Discrete]",
		},
		"unknown vendor at non-reserved slot has no type tag": {
			lines: []string{
				"name=0000:00:02.0 VGA compatible controller: Cirrus Logic GD 5446 (rev a1)",
				"addr=0000:03:00.0",
				"freq_mhz=",
			},
			expected: "Cirrus Logic GD 5446",
		},
		"no frequency when unavailable": {
			lines: []string{
				"name=0000:00:02.0 VGA compatible controller: Intel Corporation Meteor Lake-P [Intel Arc Graphics] (rev 08)",
				"addr=0000:00:02.0",
				"freq_mhz=",
			},
			expected: "Intel Arc Graphics [Integrated]",
		},
		"empty section": {
			lines:    nil,
			expected: "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := adapter.parseGpu(tt.lines)
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}
