package onprem

import (
	"fmt"
	ssh2 "ivory/clients/ssh"
	platform2 "ivory/plugins/platform"
	"strconv"
	"strings"
)

// NOTE: validate that is matches interface in compile-time
var _ platform2.Adapter = (*Adapter)(nil)

type Adapter struct {
	sshClient *ssh2.Client
}

func NewAdapter(sshClient *ssh2.Client) *Adapter {
	return &Adapter{sshClient}
}

func (a *Adapter) Metrics(connection ssh2.Connection) (*platform2.Metrics, error) {
	result, err := a.sshClient.Execute(connection, MetricsCommand)
	if err != nil {
		return nil, err
	}
	return a.parseMetrics(result.Stdout)
}

func (a *Adapter) CopyId(connection ssh2.Connection, publicKey string) error {
	return a.sshClient.CopyId(connection, publicKey)
}

func (a *Adapter) List(connection ssh2.Connection) (*platform2.OperationResult, error) {
	return a.executeDocker(connection, "ps -a")
}

func (a *Adapter) Deploy(connection ssh2.Connection, options, image string) (*platform2.OperationResult, error) {
	return a.executeDocker(connection, fmt.Sprintf("run -d %s %s", options, image))
}

func (a *Adapter) Stop(connection ssh2.Connection, name string) (*platform2.OperationResult, error) {
	return a.executeDocker(connection, "stop "+name)
}

func (a *Adapter) Delete(connection ssh2.Connection, name string) (*platform2.OperationResult, error) {
	return a.executeDocker(connection, "rm "+name)
}

func (a *Adapter) Logs(connection ssh2.Connection, name string, tail int) (*platform2.OperationResult, error) {
	command := "logs "
	if tail > 0 {
		command += "--tail " + strconv.Itoa(tail) + " "
	}
	command += name
	return a.executeDocker(connection, command)
}

func (a *Adapter) executeDocker(connection ssh2.Connection, command string) (*platform2.OperationResult, error) {
	res, err := a.sshClient.Execute(connection, a.normalizeDockerCommand(command))
	if err != nil {
		return nil, err
	}
	return &platform2.OperationResult{Stdout: res.Stdout, Stderr: res.Stderr, ExitCode: res.ExitCode}, nil
}

func (a *Adapter) normalizeDockerCommand(command string) string {
	trimmed := strings.TrimSpace(command)
	if trimmed == "" {
		return ""
	}
	if strings.HasPrefix(trimmed, "docker ") || trimmed == "docker" || strings.HasPrefix(trimmed, "sudo docker ") {
		return trimmed
	}
	return "docker " + trimmed
}

func (a *Adapter) parseMetrics(output string) (*platform2.Metrics, error) {
	sections := a.splitMetricsOutput(output)

	for _, key := range []string{"__IVORY_CPU__", "__IVORY_MEM__", "__IVORY_NET__"} {
		if _, ok := sections[key]; !ok {
			return nil, fmt.Errorf("metrics output missing section %q", key)
		}
	}

	cpu, err := a.parseCpuMetrics(sections["__IVORY_CPU__"])
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

	return &platform2.Metrics{
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
		case "__IVORY_CPU__", "__IVORY_MEM__", "__IVORY_NET__":
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

func (a *Adapter) parseCpuMetrics(lines []string) (platform2.CpuMetrics, error) {
	if len(lines) == 0 {
		return platform2.CpuMetrics{}, platform2.ErrInvalidCpuMetrics
	}

	total, idle, err := a.parseCpuLine(lines[0])
	if err != nil {
		return platform2.CpuMetrics{}, err
	}

	return platform2.CpuMetrics{TotalTicks: total, IdleTicks: idle}, nil
}

func (a *Adapter) parseCpuLine(line string) (uint64, uint64, error) {
	fields := strings.Fields(line)
	if len(fields) < 5 || fields[0] != "cpu" {
		return 0, 0, platform2.ErrInvalidCpuMetrics
	}

	var total uint64
	var idle uint64
	for i, field := range fields {
		if i == 0 {
			continue
		}
		value, err := strconv.ParseUint(field, 10, 64)
		if err != nil {
			return 0, 0, err
		}
		// Sum all fields up to 'steal' (index 8) to avoid double counting guest time
		// guest (9) and guest_nice (10) are already included in user (1) and nice (2)
		if i <= 8 {
			total += value
		}
		// idle (4) and iowait (5) are both considered "not working" time
		if i == 4 || i == 5 {
			idle += value
		}
	}

	return total, idle, nil
}

func (a *Adapter) parseMemoryMetrics(lines []string) (platform2.MemoryMetrics, error) {
	if len(lines) < 2 {
		return platform2.MemoryMetrics{}, platform2.ErrInvalidMemoryMetrics
	}

	values := make(map[string]uint64)
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return platform2.MemoryMetrics{}, platform2.ErrInvalidMemoryMetrics
		}
		value, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return platform2.MemoryMetrics{}, err
		}
		values[strings.TrimSuffix(fields[0], ":")] = value * 1024
	}

	total := values["MemTotal"]
	available := values["MemAvailable"]
	if total == 0 {
		return platform2.MemoryMetrics{}, platform2.ErrInvalidMemoryMetrics
	}

	return platform2.MemoryMetrics{
		TotalBytes:     total,
		AvailableBytes: available,
	}, nil
}

func (a *Adapter) parseNetworkMetrics(lines []string) (platform2.NetworkMetrics, error) {
	var received uint64
	var transmitted uint64
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}

		iface := strings.TrimSpace(parts[0])
		if strings.Contains(iface, "|") || iface == "lo" {
			continue
		}

		fields := strings.Fields(parts[1])
		if len(fields) < 16 {
			return platform2.NetworkMetrics{}, platform2.ErrInvalidNetworkMetrics
		}

		rx, err := strconv.ParseUint(fields[0], 10, 64)
		if err != nil {
			return platform2.NetworkMetrics{}, err
		}
		tx, err := strconv.ParseUint(fields[8], 10, 64)
		if err != nil {
			return platform2.NetworkMetrics{}, err
		}

		received += rx
		transmitted += tx
	}

	return platform2.NetworkMetrics{
		ReceivedBytes:    received,
		TransmittedBytes: transmitted,
	}, nil
}
