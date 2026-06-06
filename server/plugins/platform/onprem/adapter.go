package onprem

import (
	"crypto/ed25519"
	"fmt"
	"ivory/clients/console"
	"ivory/clients/console/ssh"
	"ivory/plugins/platform"
	"strconv"
	"strings"
)

// NOTE: validate that is matches interface in compile-time
var _ platform.Adapter = (*Adapter)(nil)

type Adapter struct {
	sshClient *ssh.Client
}

func NewAdapter(sshClient *ssh.Client) *Adapter {
	return &Adapter{sshClient}
}

func (a *Adapter) Metrics(connection platform.Connection) (*platform.Metrics, error) {
	result, err := a.execute(connection, MetricsCommand)
	if err != nil {
		return nil, err
	}
	return a.parseMetrics(strings.Join(result, "\n"))
}

func (a *Adapter) CopyId(connection platform.Connection, publicKey string) error {
	command := fmt.Sprintf(`mkdir -p ~/.ssh && chmod 700 ~/.ssh && (grep -qF "%s" ~/.ssh/authorized_keys 2>/dev/null || echo "%s" >> ~/.ssh/authorized_keys) && chmod 600 ~/.ssh/authorized_keys`, publicKey, publicKey)
	_, err := a.execute(connection, command)
	return err
}

func (a *Adapter) execute(connection platform.Connection, command string) ([]string, error) {
	cmd := a.sshClient.Command(a.mapToSshCommand(connection), command)
	return cmd.Execute()
}

func (a *Adapter) ListCommand(connection platform.Connection) console.Command {
	return a.sshClient.Command(a.mapToSshCommand(connection), a.normalizeDockerCommand("ps -a"))
}

func (a *Adapter) DeployCommand(connection platform.Connection, options, image string) console.Command {
	return a.sshClient.Command(a.mapToSshCommand(connection), a.normalizeDockerCommand(fmt.Sprintf("run -d %s %s", options, image)))
}

func (a *Adapter) StopCommand(connection platform.Connection, name string) console.Command {
	return a.sshClient.Command(a.mapToSshCommand(connection), a.normalizeDockerCommand("stop "+name))
}

func (a *Adapter) DeleteCommand(connection platform.Connection, name string) console.Command {
	return a.sshClient.Command(a.mapToSshCommand(connection), a.normalizeDockerCommand("rm "+name))
}

func (a *Adapter) LogsCommand(connection platform.Connection, name string, tail int) console.Command {
	return a.sshClient.Command(a.mapToSshCommand(connection), a.normalizeDockerCommand(a.getLogsCommand(name, tail)))
}

func (a *Adapter) getLogsCommand(name string, tail int) string {
	command := "logs "
	if tail > 0 {
		command += "--tail " + strconv.Itoa(tail) + " "
	}
	command += name
	return command
}

func (a *Adapter) mapToSshCommand(conn platform.Connection) ssh.Command {
	var prvKey *ed25519.PrivateKey
	if len(conn.PrivateKey) > 0 {
		pk := ed25519.PrivateKey(conn.PrivateKey)
		prvKey = &pk
	}
	return ssh.Command{
		Host:       conn.Host,
		Port:       conn.Port,
		Username:   conn.Username,
		Password:   conn.Password,
		PrivateKey: prvKey,
	}
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

func (a *Adapter) parseMetrics(output string) (*platform.Metrics, error) {
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

	return &platform.Metrics{
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

func (a *Adapter) parseCpuMetrics(lines []string) (platform.CpuMetrics, error) {
	if len(lines) == 0 {
		return platform.CpuMetrics{}, platform.ErrInvalidCpuMetrics
	}

	total, idle, err := a.parseCpuLine(lines[0])
	if err != nil {
		return platform.CpuMetrics{}, err
	}

	return platform.CpuMetrics{TotalTicks: total, IdleTicks: idle}, nil
}

func (a *Adapter) parseCpuLine(line string) (uint64, uint64, error) {
	fields := strings.Fields(line)
	if len(fields) < 5 || fields[0] != "cpu" {
		return 0, 0, platform.ErrInvalidCpuMetrics
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

func (a *Adapter) parseMemoryMetrics(lines []string) (platform.MemoryMetrics, error) {
	if len(lines) < 2 {
		return platform.MemoryMetrics{}, platform.ErrInvalidMemoryMetrics
	}

	values := make(map[string]uint64)
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return platform.MemoryMetrics{}, platform.ErrInvalidMemoryMetrics
		}
		value, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return platform.MemoryMetrics{}, err
		}
		values[strings.TrimSuffix(fields[0], ":")] = value * 1024
	}

	total := values["MemTotal"]
	available := values["MemAvailable"]
	if total == 0 {
		return platform.MemoryMetrics{}, platform.ErrInvalidMemoryMetrics
	}

	return platform.MemoryMetrics{
		TotalBytes:     total,
		AvailableBytes: available,
	}, nil
}

func (a *Adapter) parseNetworkMetrics(lines []string) (platform.NetworkMetrics, error) {
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
			return platform.NetworkMetrics{}, platform.ErrInvalidNetworkMetrics
		}

		rx, err := strconv.ParseUint(fields[0], 10, 64)
		if err != nil {
			return platform.NetworkMetrics{}, err
		}
		tx, err := strconv.ParseUint(fields[8], 10, 64)
		if err != nil {
			return platform.NetworkMetrics{}, err
		}

		received += rx
		transmitted += tx
	}

	return platform.NetworkMetrics{
		ReceivedBytes:    received,
		TransmittedBytes: transmitted,
	}, nil
}
