package ssh

import (
	"bytes"
	"errors"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

var ErrCommandEmpty = errors.New("command cannot be empty")
var ErrHostEmpty = errors.New("vm host cannot be empty")
var ErrInvalidCpuMetrics = errors.New("invalid cpu metrics output")
var ErrInvalidMemoryMetrics = errors.New("invalid memory metrics output")
var ErrInvalidNetworkMetrics = errors.New("invalid network metrics output")

type Client interface {
	Execute(vm VM, command string) (*CommandResult, error)
	ExecuteDocker(vm VM, command string) (*CommandResult, error)
	Metrics(vm VM) (*Metrics, error)
}

type client struct {
	timeout time.Duration
}

func NewClient() Client {
	return &client{timeout: 10 * time.Second}
}

func (c *client) Execute(vm VM, command string) (*CommandResult, error) {
	trimmed := strings.TrimSpace(command)
	if trimmed == "" {
		return nil, ErrCommandEmpty
	}

	signer, err := ssh.ParsePrivateKey([]byte(vm.SshKey))
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User:            vm.Username,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         c.timeout,
	}

	target, err := getDialAddress(vm)
	if err != nil {
		return nil, err
	}

	conn, err := ssh.Dial("tcp", target, config)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	result := &CommandResult{}
	err = session.Run(trimmed)
	result.Stdout = stdout.String()
	result.Stderr = stderr.String()

	var exitErr *ssh.ExitError
	if errors.As(err, &exitErr) {
		result.ExitCode = exitErr.ExitStatus()
		return result, nil
	}
	if err != nil {
		return nil, err
	}

	result.ExitCode = 0
	return result, nil
}

func (c *client) ExecuteDocker(vm VM, command string) (*CommandResult, error) {
	return c.Execute(vm, normalizeDockerCommand(command))
}

func (c *client) Metrics(vm VM) (*Metrics, error) {
	result, err := c.Execute(vm, MetricsCommand)
	if err != nil {
		return nil, err
	}
	return parseMetrics(result.Stdout)
}

func normalizeDockerCommand(command string) string {
	trimmed := strings.TrimSpace(command)
	if trimmed == "" {
		return ""
	}
	if strings.HasPrefix(trimmed, "docker ") || trimmed == "docker" || strings.HasPrefix(trimmed, "sudo docker ") {
		return trimmed
	}
	return "docker " + trimmed
}

func getDialAddress(vm VM) (string, error) {
	host := strings.TrimSpace(vm.Host)
	if host == "" {
		return "", ErrHostEmpty
	}

	if strings.Contains(host, "://") {
		parsed, err := url.Parse(host)
		if err != nil {
			return "", err
		}
		if parsed.Hostname() != "" {
			port := vm.Port
			if parsed.Port() != "" {
				parsedPort, err := strconv.Atoi(parsed.Port())
				if err != nil {
					return "", err
				}
				port = parsedPort
			}
			return net.JoinHostPort(parsed.Hostname(), strconv.Itoa(port)), nil
		}
	}

	if _, _, err := net.SplitHostPort(host); err == nil {
		return host, nil
	}

	return net.JoinHostPort(host, strconv.Itoa(vm.Port)), nil
}

func parseMetrics(output string) (*Metrics, error) {
	sections := splitMetricsOutput(output)

	cpu, err := parseCpuMetrics(sections["__IVORY_CPU_1__"], sections["__IVORY_CPU_2__"])
	if err != nil {
		return nil, err
	}
	memory, err := parseMemoryMetrics(sections["__IVORY_MEM__"])
	if err != nil {
		return nil, err
	}
	network, err := parseNetworkMetrics(sections["__IVORY_NET__"])
	if err != nil {
		return nil, err
	}

	return &Metrics{
		Cpu:     cpu,
		Memory:  memory,
		Network: network,
	}, nil
}

func splitMetricsOutput(output string) map[string][]string {
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

func parseCpuMetrics(first []string, second []string) (CpuMetrics, error) {
	if len(first) == 0 || len(second) == 0 {
		return CpuMetrics{}, ErrInvalidCpuMetrics
	}

	prev, prevIdle, err := parseCpuLine(first[0])
	if err != nil {
		return CpuMetrics{}, err
	}
	next, nextIdle, err := parseCpuLine(second[0])
	if err != nil {
		return CpuMetrics{}, err
	}

	totalDiff := next - prev
	idleDiff := nextIdle - prevIdle
	if totalDiff <= 0 {
		return CpuMetrics{}, ErrInvalidCpuMetrics
	}

	usagePercent := (float64(totalDiff-idleDiff) / float64(totalDiff)) * 100
	return CpuMetrics{UsagePercent: usagePercent}, nil
}

func parseCpuLine(line string) (uint64, uint64, error) {
	fields := strings.Fields(line)
	if len(fields) < 5 || fields[0] != "cpu" {
		return 0, 0, ErrInvalidCpuMetrics
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

func parseMemoryMetrics(lines []string) (MemoryMetrics, error) {
	if len(lines) < 2 {
		return MemoryMetrics{}, ErrInvalidMemoryMetrics
	}

	values := make(map[string]uint64)
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return MemoryMetrics{}, ErrInvalidMemoryMetrics
		}
		value, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return MemoryMetrics{}, err
		}
		values[strings.TrimSuffix(fields[0], ":")] = value * 1024
	}

	total := values["MemTotal"]
	available := values["MemAvailable"]
	if total == 0 {
		return MemoryMetrics{}, ErrInvalidMemoryMetrics
	}
	used := total - available

	return MemoryMetrics{
		TotalBytes:     total,
		UsedBytes:      used,
		AvailableBytes: available,
		UsagePercent:   (float64(used) / float64(total)) * 100,
	}, nil
}

func parseNetworkMetrics(lines []string) (NetworkMetrics, error) {
	if len(lines) < 3 {
		return NetworkMetrics{}, ErrInvalidNetworkMetrics
	}

	var received uint64
	var transmitted uint64
	for _, line := range lines[2:] {
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			return NetworkMetrics{}, ErrInvalidNetworkMetrics
		}

		iface := strings.TrimSpace(parts[0])
		if iface == "lo" {
			continue
		}

		fields := strings.Fields(parts[1])
		if len(fields) < 9 {
			return NetworkMetrics{}, ErrInvalidNetworkMetrics
		}

		rx, err := strconv.ParseUint(fields[0], 10, 64)
		if err != nil {
			return NetworkMetrics{}, err
		}
		tx, err := strconv.ParseUint(fields[8], 10, 64)
		if err != nil {
			return NetworkMetrics{}, err
		}

		received += rx
		transmitted += tx
	}

	return NetworkMetrics{
		ReceivedBytes:    received,
		TransmittedBytes: transmitted,
	}, nil
}
