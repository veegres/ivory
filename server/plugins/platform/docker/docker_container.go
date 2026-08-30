package docker

import (
	"encoding/json"
	"ivory/clients/console"
	"ivory/plugins/platform"
	"math"
	"strconv"
	"strings"
)

// NOTE: validate that is matches interface in compile-time
var _ platform.Container = (*Plugin)(nil)

// containerCpuScale is an arbitrary fixed-point scale used to encode docker's
// already-computed CPU percentage into the TotalTicks/IdleTicks counters, so
// container metrics fit the same Metrics shape produced from host /proc stats.
const containerCpuScale = 10000

type dockerStatsEntry struct {
	CPUPerc  string `json:"CPUPerc"`
	MemUsage string `json:"MemUsage"`
	NetIO    string `json:"NetIO"`
}

var byteSizeUnits = map[string]float64{
	"B":   1,
	"KB":  1000,
	"KIB": 1024,
	"MB":  1000 * 1000,
	"MIB": 1024 * 1024,
	"GB":  1000 * 1000 * 1000,
	"GIB": 1024 * 1024 * 1024,
	"TB":  1000 * 1000 * 1000 * 1000,
	"TIB": 1024 * 1024 * 1024 * 1024,
}

func (p *Plugin) ListContainer(connection platform.Connection) console.Command {
	return p.sshClient.Command(p.mapToSshCommand(connection), "docker ps -a")
}

// UpContainer starts a container from the options the user wrote. The verb is
// Ivory's, exactly as it is for every other method here: this deploys a
// container and nothing else, so taking the executable from the command text
// would make it a remote shell that happens to be called UpContainer.
func (p *Plugin) UpContainer(connection platform.Connection, command []string) console.Command {
	return p.sshClient.Command(p.mapToSshCommand(connection), strings.Join(quoteFields(normalizeRun(command)), " "))
}

// normalizeRun puts Ivory's own "docker run" in front of the user's options.
// A "docker" or "run" they wrote themselves is dropped rather than duplicated
// - a command is usually copied in whole - so the executed verb is always
// exactly this one.
func normalizeRun(command []string) []string {
	for _, verb := range []string{"docker", "run"} {
		if len(command) > 0 && command[0] == verb {
			command = command[1:]
		}
	}
	return append([]string{"docker", "run"}, command...)
}

func (p *Plugin) DownContainer(connection platform.Connection, name string) console.Command {
	return p.sshClient.Command(p.mapToSshCommand(connection), "docker rm -- "+shellQuote(name))
}

func (p *Plugin) StartContainer(connection platform.Connection, name string) console.Command {
	return p.sshClient.Command(p.mapToSshCommand(connection), "docker start -- "+shellQuote(name))
}

func (p *Plugin) StopContainer(connection platform.Connection, name string) console.Command {
	return p.sshClient.Command(p.mapToSshCommand(connection), "docker stop -- "+shellQuote(name))
}

func (p *Plugin) RestartContainer(connection platform.Connection, name string) console.Command {
	return p.sshClient.Command(p.mapToSshCommand(connection), "docker restart -- "+shellQuote(name))
}

func (p *Plugin) ExecContainer(connection platform.Connection, name string, command []string) console.Command {
	parts := []string{"docker", "exec", "--", shellQuote(name)}
	parts = append(parts, quoteFields(command)...)
	return p.sshClient.Command(p.mapToSshCommand(connection), strings.Join(parts, " "))
}

func (p *Plugin) LogsContainer(connection platform.Connection, name string, tail int, follow bool) console.Command {
	commandStr := "docker logs "
	if tail > 0 {
		commandStr += "--tail " + strconv.Itoa(tail) + " "
	}
	if follow {
		commandStr += "--follow "
	}
	commandStr += "-- " + shellQuote(name)
	command := p.sshClient.Command(p.mapToSshCommand(connection), commandStr)
	command.JobKeepAlive = false
	return command
}

func (p *Plugin) MetricsContainer(connection platform.Connection, name string) (*platform.Metrics, error) {
	command := "docker stats --no-stream --format " + shellQuote("{{json .}}") + " -- " + shellQuote(name)
	result, err := p.execute(connection, command)
	if err != nil {
		return nil, err
	}
	return p.parseContainerMetrics(strings.Join(result, "\n"))
}

// quoteFields re-quotes tokens for the shell, leaving plain ones alone so the
// executed command still reads like the one the user wrote.
func quoteFields(fields []string) []string {
	quoted := make([]string, 0, len(fields))
	for _, field := range fields {
		if needsShellQuote(field) {
			quoted = append(quoted, shellQuote(field))
		} else {
			quoted = append(quoted, field)
		}
	}
	return quoted
}

func needsShellQuote(field string) bool {
	if field == "" {
		return true
	}
	return strings.ContainsAny(field, " \t\n'\"\\$`&|;<>()*?[]#~!")
}

func (p *Plugin) parseContainerMetrics(output string) (*platform.Metrics, error) {
	line := strings.TrimSpace(output)
	if line == "" {
		return nil, platform.ErrContainerNotRunning
	}

	var entry dockerStatsEntry
	if err := json.Unmarshal([]byte(line), &entry); err != nil {
		return nil, platform.ErrInvalidContainerMetrics
	}

	cpu, err := parseContainerCpuMetrics(entry.CPUPerc)
	if err != nil {
		return nil, err
	}
	memory, err := parseContainerMemoryMetrics(entry.MemUsage)
	if err != nil {
		return nil, err
	}
	network, err := parseContainerNetworkMetrics(entry.NetIO)
	if err != nil {
		return nil, err
	}

	return &platform.Metrics{Cpu: cpu, Memory: memory, Network: network}, nil
}

func parseContainerCpuMetrics(raw string) (platform.CpuMetrics, error) {
	trimmed := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(raw), "%"))
	if trimmed == "" {
		return platform.CpuMetrics{}, platform.ErrInvalidContainerMetrics
	}
	percent, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return platform.CpuMetrics{}, platform.ErrInvalidContainerMetrics
	}

	// Docker reports CPU usage relative to a single core and can exceed 100%
	// for multi-core containers; clamp so it fits the fixed-point scale below.
	if percent > 100 {
		percent = 100
	}
	if percent < 0 {
		percent = 0
	}

	busy := uint64(math.Round(percent / 100 * containerCpuScale))
	return platform.CpuMetrics{TotalTicks: containerCpuScale, IdleTicks: containerCpuScale - busy}, nil
}

func parseContainerMemoryMetrics(raw string) (platform.MemoryMetrics, error) {
	parts := strings.SplitN(raw, "/", 2)
	if len(parts) != 2 {
		return platform.MemoryMetrics{}, platform.ErrInvalidContainerMetrics
	}

	used, err := parseByteSize(parts[0])
	if err != nil {
		return platform.MemoryMetrics{}, err
	}
	limit, err := parseByteSize(parts[1])
	if err != nil {
		return platform.MemoryMetrics{}, err
	}
	if limit == 0 {
		return platform.MemoryMetrics{}, platform.ErrInvalidContainerMetrics
	}

	available := uint64(0)
	if limit > used {
		available = limit - used
	}
	return platform.MemoryMetrics{TotalBytes: limit, AvailableBytes: available}, nil
}

func parseContainerNetworkMetrics(raw string) (platform.NetworkMetrics, error) {
	parts := strings.SplitN(raw, "/", 2)
	if len(parts) != 2 {
		return platform.NetworkMetrics{}, platform.ErrInvalidContainerMetrics
	}

	rx, err := parseByteSize(parts[0])
	if err != nil {
		return platform.NetworkMetrics{}, err
	}
	tx, err := parseByteSize(parts[1])
	if err != nil {
		return platform.NetworkMetrics{}, err
	}
	return platform.NetworkMetrics{ReceivedBytes: rx, TransmittedBytes: tx}, nil
}

func parseByteSize(raw string) (uint64, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return 0, platform.ErrInvalidContainerMetrics
	}

	i := 0
	for i < len(trimmed) && (trimmed[i] == '.' || (trimmed[i] >= '0' && trimmed[i] <= '9')) {
		i++
	}
	numberPart := trimmed[:i]
	unitPart := strings.ToUpper(strings.TrimSpace(trimmed[i:]))
	if numberPart == "" {
		return 0, platform.ErrInvalidContainerMetrics
	}

	number, err := strconv.ParseFloat(numberPart, 64)
	if err != nil {
		return 0, platform.ErrInvalidContainerMetrics
	}
	if unitPart == "" {
		unitPart = "B"
	}

	multiplier, ok := byteSizeUnits[unitPart]
	if !ok {
		return 0, platform.ErrInvalidContainerMetrics
	}
	return uint64(math.Round(number * multiplier)), nil
}
