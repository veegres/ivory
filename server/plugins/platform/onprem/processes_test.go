package onprem

import (
	"errors"
	"ivory/clients/console/ssh"
	"ivory/plugins/platform"
	"testing"
)

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
