package linux

import (
	"fmt"
	"ivory/src/clients/ssh"
	"ivory/src/plugins/os"
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

func (a *Adapter) DockerList(connection ssh.Connection) (*os.DockerResult, error) {
	return a.executeDocker(connection, "ps -a")
}

func (a *Adapter) DockerDeploy(connection ssh.Connection, image, options string) (*os.DockerResult, error) {
	command := fmt.Sprintf("pull %s && run -d %s %s", image, options, image)
	return a.executeDocker(connection, command)
}

func (a *Adapter) DockerRun(connection ssh.Connection, options, image string) (*os.DockerResult, error) {
	return a.executeDocker(connection, fmt.Sprintf("run -d %s %s", options, image))
}

func (a *Adapter) DockerStop(connection ssh.Connection, container string) (*os.DockerResult, error) {
	return a.executeDocker(connection, "stop "+container)
}

func (a *Adapter) DockerDelete(connection ssh.Connection, container string) (*os.DockerResult, error) {
	return a.executeDocker(connection, "rm "+container)
}

func (a *Adapter) DockerLogs(connection ssh.Connection, container string, tail int) (*os.DockerResult, error) {
	command := "logs "
	if tail > 0 {
		command += "--tail " + strconv.Itoa(tail) + " "
	}
	command += container
	return a.executeDocker(connection, command)
}

func (a *Adapter) executeDocker(connection ssh.Connection, command string) (*os.DockerResult, error) {
	res, err := a.sshClient.Execute(connection, a.normalizeDockerCommand(command))
	if err != nil {
		return nil, err
	}
	return &os.DockerResult{Stdout: res.Stdout, Stderr: res.Stderr, ExitCode: res.ExitCode}, nil
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

func (a *Adapter) parseMetrics(output string) (*os.Metrics, error) {
	sections := a.splitMetricsOutput(output)

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

func (a *Adapter) parseCpuMetrics(lines []string) (os.CpuMetrics, error) {
	if len(lines) == 0 {
		return os.CpuMetrics{}, os.ErrInvalidCpuMetrics
	}

	total, idle, err := a.parseCpuLine(lines[0])
	if err != nil {
		return os.CpuMetrics{}, err
	}

	return os.CpuMetrics{TotalTicks: total, IdleTicks: idle}, nil
}

func (a *Adapter) parseCpuLine(line string) (uint64, uint64, error) {
	fields := strings.Fields(line)
	if len(fields) < 5 || fields[0] != "cpu" {
		return 0, 0, os.ErrInvalidCpuMetrics
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
