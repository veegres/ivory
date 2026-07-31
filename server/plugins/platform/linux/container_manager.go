package linux

import (
	"encoding/json"
	"ivory/clients/console"
	"ivory/plugins/platform"
	"math"
	"strconv"
	"strings"
	"unicode"
)

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

// RenderOptions renders the deploy spec into docker run flags, one per line,
// keeping {{placeholder}} templates intact for later interpolation.
func (a *Adapter) RenderOptions(spec platform.DeploySpec) string {
	lines := make([]string, 0)
	if spec.Name != "" {
		lines = append(lines, "--name "+spec.Name)
	}
	if spec.Hostname != "" {
		lines = append(lines, "--hostname "+spec.Hostname)
	}
	if spec.HostNetwork {
		lines = append(lines, "--network host")
	}
	if spec.RestartPolicy != "" {
		lines = append(lines, "--restart "+spec.RestartPolicy)
	}
	for _, port := range spec.Ports {
		lines = append(lines, "-p "+port+":"+port)
	}
	for _, volume := range spec.Volumes {
		lines = append(lines, "-v "+volume.HostPath+":"+volume.ContainerPath)
	}
	for _, env := range spec.Env {
		lines = append(lines, "-e "+env.Name+"="+env.Value)
	}
	return strings.Join(lines, "\n")
}

func (a *Adapter) ListContainer(connection platform.Connection) console.Command {
	return a.sshClient.Command(a.mapToSshCommand(connection), a.normalizeDockerCommand("ps -a"))
}

func (a *Adapter) UpContainer(connection platform.Connection, name, options, image, entryScript string) console.Command {
	// NOTE: name is enforced here rather than trusted from the (user-editable)
	// options text, so a --name line removed or altered during editing can
	// never desync the real container name from what Ivory expects to manage
	// afterwards (see removeFlag).
	parts := []string{"run", "-d", "--name", shellQuote(name)}
	parts = append(parts, quoteFields(removeFlag(splitShellFields(options), "--name"))...)
	parts = append(parts, "--")
	parts = append(parts, shellQuote(image))
	if entryScript != "" {
		parts = append(parts, shellQuoteFields(entryScript)...)
	}
	return a.sshClient.Command(a.mapToSshCommand(connection), a.normalizeDockerCommand(strings.Join(parts, " ")))
}

func (a *Adapter) DownContainer(connection platform.Connection, name string) console.Command {
	return a.sshClient.Command(a.mapToSshCommand(connection), a.normalizeDockerCommand("rm -- "+shellQuote(name)))
}

func (a *Adapter) StartContainer(connection platform.Connection, name string) console.Command {
	return a.sshClient.Command(a.mapToSshCommand(connection), a.normalizeDockerCommand("start -- "+shellQuote(name)))
}

func (a *Adapter) StopContainer(connection platform.Connection, name string) console.Command {
	return a.sshClient.Command(a.mapToSshCommand(connection), a.normalizeDockerCommand("stop -- "+shellQuote(name)))
}

func (a *Adapter) RestartContainer(connection platform.Connection, name string) console.Command {
	return a.sshClient.Command(a.mapToSshCommand(connection), a.normalizeDockerCommand("restart -- "+shellQuote(name)))
}

func (a *Adapter) ExecContainer(connection platform.Connection, name string, command string) console.Command {
	parts := []string{"exec", "--", shellQuote(name)}
	parts = append(parts, shellQuoteFields(command)...)
	return a.sshClient.Command(a.mapToSshCommand(connection), a.normalizeDockerCommand(strings.Join(parts, " ")))
}

func (a *Adapter) LogsContainer(connection platform.Connection, name string, tail int, follow bool) console.Command {
	commandStr := "logs "
	if tail > 0 {
		commandStr += "--tail " + strconv.Itoa(tail) + " "
	}
	if follow {
		commandStr += "--follow "
	}
	commandStr += "-- " + shellQuote(name)
	command := a.sshClient.Command(a.mapToSshCommand(connection), a.normalizeDockerCommand(commandStr))
	command.JobKeepAlive = false
	return command
}

func (a *Adapter) MetricsContainer(connection platform.Connection, name string) (*platform.Metrics, error) {
	command := a.normalizeDockerCommand("stats --no-stream --format " + shellQuote("{{json .}}") + " -- " + shellQuote(name))
	result, err := a.execute(connection, command)
	if err != nil {
		return nil, err
	}
	return a.parseContainerMetrics(strings.Join(result, "\n"))
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

// removeFlag drops flag and the value token right after it from an already
// tokenized options list, operating on tokens (not the raw string) so a
// value elsewhere in the string can't be mistaken for the flag itself and
// a token containing a literal space (from a quoted span) never gets
// re-split by a later round of tokenizing.
func removeFlag(fields []string, flag string) []string {
	kept := make([]string, 0, len(fields))
	for i := 0; i < len(fields); i++ {
		if fields[i] == flag {
			i++ // also skip the flag's value
			continue
		}
		kept = append(kept, fields[i])
	}
	return kept
}

func quoteFields(fields []string) []string {
	quoted := make([]string, 0, len(fields))
	for _, field := range fields {
		quoted = append(quoted, shellQuote(field))
	}
	return quoted
}

func shellQuoteFields(value string) []string {
	return quoteFields(splitShellFields(value))
}

// splitShellFields splits a docker options string into individual arguments,
// honoring single/double-quoted spans the way a real shell would, so a flag
// value like `-e SCOPE="my cluster"` (as produced by RenderOptions for an env
// value) stays one argument instead of being broken apart at the space inside
// the quotes. A backslash escapes the very next rune literally regardless of
// quote state, so a value inserted into a plugin's own hand-written quoted
// span (e.g. a PostScript wrapping {{dbUser}}:{{dbPass}} in literal quotes)
// can contain that same quote character - escaped by the caller before
// interpolation - without prematurely closing the span.
func splitShellFields(value string) []string {
	fields := make([]string, 0)
	var current strings.Builder
	hasToken := false
	var quote rune
	runes := []rune(value)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		switch {
		case r == '\\' && i+1 < len(runes):
			i++
			current.WriteRune(runes[i])
			hasToken = true
		case quote != 0:
			if r == quote {
				quote = 0
			} else {
				current.WriteRune(r)
			}
		case r == '\'' || r == '"':
			quote = r
			hasToken = true
		case unicode.IsSpace(r):
			if hasToken {
				fields = append(fields, current.String())
				current.Reset()
				hasToken = false
			}
		default:
			current.WriteRune(r)
			hasToken = true
		}
	}
	if hasToken {
		fields = append(fields, current.String())
	}
	return fields
}

func (a *Adapter) parseContainerMetrics(output string) (*platform.Metrics, error) {
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
