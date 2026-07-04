package onprem

import (
	"errors"
	"ivory/clients/console/ssh"
	"ivory/plugins/platform"
	"math"
	"testing"
)

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
		if !errors.Is(err, platform.ErrInvalidContainerMetrics) {
			t.Errorf("expected ErrInvalidContainerMetrics, got %v", err)
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

func TestMetricsContainerQuotesShellArguments(t *testing.T) {
	adapter := NewAdapter(ssh.NewClient())
	name := "foo; rm -rf /"

	command := adapter.normalizeDockerCommand("stats --no-stream --format " + shellQuote("{{json .}}") + " -- " + shellQuote(name))
	expected := `docker stats --no-stream --format '{{json .}}' -- 'foo; rm -rf /'`
	if command != expected {
		t.Fatalf("expected %q, got %q", expected, command)
	}
}
