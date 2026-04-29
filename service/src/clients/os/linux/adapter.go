package linux

import (
	"ivory/src/clients/os"
	"ivory/src/clients/ssh"
	"strconv"
	"strings"
)

// NOTE: validate that is matches interface in compile-time
var _ os.Adapter = (*Adapter)(nil)

type Adapter struct {
	sshClient *ssh.Client
}

func NewAdapter(sshClient *ssh.Client) *Adapter {
	return &Adapter{sshClient}
}

func (a *Adapter) Metrics(connection ssh.Connection) (*os.Metrics, error) {
	result, err := a.sshClient.Execute(connection, MetricsCommand)
	if err != nil {
		return nil, err
	}
	return a.parseMetrics(result.Stdout)
}

func (a *Adapter) parseMetrics(output string) (*os.Metrics, error) {
	sections := a.splitMetricsOutput(output)

	cpu, err := a.parseCpuMetrics(sections["__IVORY_CPU_1__"], sections["__IVORY_CPU_2__"])
	if err != nil {
		return nil, err
	}
	memory, err := a.parseMemoryMetrics(sections["__IVORY_MEM__"])
	if err != nil {
		return nil, err
	}
	network, err := a.parseNetworkMetrics(sections["__IVORY_NET__"])
	if err != nil {
		return nil, err
	}

	return &os.Metrics{
		Cpu:     cpu,
		Memory:  memory,
		Network: network,
	}, nil
}

func (a *Adapter) splitMetricsOutput(output string) map[string][]string {
	sections := map[string][]string{}
	current := ""

	for _, line := range strings.Split(output, "\n") {
		trimmed := strings.TrimSpace(line)
		switch trimmed {
		case "__IVORY_CPU_1__", "__IVORY_CPU_2__", "__IVORY_MEM__", "__IVORY_NET__":
			current = trimmed
			continue
		}
		if current == "" || trimmed == "" {
			continue
		}
		sections[current] = append(sections[current], trimmed)
	}

	return sections
}

func (a *Adapter) parseCpuMetrics(first []string, second []string) (os.CpuMetrics, error) {
	if len(first) == 0 || len(second) == 0 {
		return os.CpuMetrics{}, os.ErrInvalidCpuMetrics
	}

	prev, prevIdle, err := a.parseCpuLine(first[0])
	if err != nil {
		return os.CpuMetrics{}, err
	}
	next, nextIdle, err := a.parseCpuLine(second[0])
	if err != nil {
		return os.CpuMetrics{}, err
	}

	totalDiff := next - prev
	idleDiff := nextIdle - prevIdle
	if totalDiff <= 0 {
		return os.CpuMetrics{}, os.ErrInvalidCpuMetrics
	}

	usagePercent := (float64(totalDiff-idleDiff) / float64(totalDiff)) * 100
	return os.CpuMetrics{UsagePercent: usagePercent}, nil
}

func (a *Adapter) parseCpuLine(line string) (uint64, uint64, error) {
	fields := strings.Fields(line)
	if len(fields) < 5 || fields[0] != "cpu" {
		return 0, 0, os.ErrInvalidCpuMetrics
	}

	var total uint64
	for _, field := range fields[1:] {
		value, err := strconv.ParseUint(field, 10, 64)
		if err != nil {
			return 0, 0, err
		}
		total += value
	}

	idle, err := strconv.ParseUint(fields[4], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	return total, idle, nil
}

func (a *Adapter) parseMemoryMetrics(lines []string) (os.MemoryMetrics, error) {
	if len(lines) < 2 {
		return os.MemoryMetrics{}, os.ErrInvalidMemoryMetrics
	}

	values := make(map[string]uint64)
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return os.MemoryMetrics{}, os.ErrInvalidMemoryMetrics
		}
		value, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return os.MemoryMetrics{}, err
		}
		values[strings.TrimSuffix(fields[0], ":")] = value * 1024
	}

	total := values["MemTotal"]
	available := values["MemAvailable"]
	if total == 0 {
		return os.MemoryMetrics{}, os.ErrInvalidMemoryMetrics
	}
	used := total - available

	return os.MemoryMetrics{
		TotalBytes:     total,
		UsedBytes:      used,
		AvailableBytes: available,
		UsagePercent:   (float64(used) / float64(total)) * 100,
	}, nil
}

func (a *Adapter) parseNetworkMetrics(lines []string) (os.NetworkMetrics, error) {
	if len(lines) < 3 {
		return os.NetworkMetrics{}, os.ErrInvalidNetworkMetrics
	}

	var received uint64
	var transmitted uint64
	for _, line := range lines[2:] {
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			return os.NetworkMetrics{}, os.ErrInvalidNetworkMetrics
		}

		iface := strings.TrimSpace(parts[0])
		if iface == "lo" {
			continue
		}

		fields := strings.Fields(parts[1])
		if len(fields) < 9 {
			return os.NetworkMetrics{}, os.ErrInvalidNetworkMetrics
		}

		rx, err := strconv.ParseUint(fields[0], 10, 64)
		if err != nil {
			return os.NetworkMetrics{}, err
		}
		tx, err := strconv.ParseUint(fields[8], 10, 64)
		if err != nil {
			return os.NetworkMetrics{}, err
		}

		received += rx
		transmitted += tx
	}

	return os.NetworkMetrics{
		ReceivedBytes:    received,
		TransmittedBytes: transmitted,
	}, nil
}
