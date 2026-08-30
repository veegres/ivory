package docker

import (
	"crypto/ed25519"
	"errors"
	"fmt"
	"ivory/clients/console"
	"ivory/clients/console/ssh"
	"ivory/plugins/platform"
	"path"
	"strconv"
	"strings"
	"time"
)

const MetricsCommand = `sh -c '
echo __IVORY_CPU__; head -n 1 /proc/stat;
echo __IVORY_MEM__; grep -E "MemTotal|MemAvailable" /proc/meminfo;
echo __IVORY_NET__; cat /proc/net/dev'`

const ProcessesCommand = `ps -eo pid,user,pcpu,pmem,rss,nlwp,comm,args --no-headers --sort=-pcpu | head -n 100`

const InfoCommand = `sh -c '
echo __IVORY_HOST__;
printf "product_name=%s\n" "$(cat /sys/devices/virtual/dmi/id/product_name 2>/dev/null)";
printf "product_family=%s\n" "$(cat /sys/devices/virtual/dmi/id/product_family 2>/dev/null)";
printf "hostname=%s\n" "$(hostname)";
echo __IVORY_OS__; cat /etc/os-release 2>/dev/null;
echo __IVORY_KERNEL__; uname -sr;
echo __IVORY_UPTIME__; cat /proc/uptime;
echo __IVORY_CPU__; grep -m1 "model name" /proc/cpuinfo; grep -c ^processor /proc/cpuinfo;
echo __IVORY_GPU__;
gpu_line="$(lspci -D 2>/dev/null | grep -i vga | head -n1)";
printf "name=%s\n" "$gpu_line";
gpu_addr="${gpu_line%% *}";
printf "addr=%s\n" "$gpu_addr";
gpu_dir="/sys/bus/pci/devices/$gpu_addr";
freq="$(cat "$gpu_dir/tile0/gt0/freq0/max_freq" 2>/dev/null)";
if [ -z "$freq" ]; then drm_card="$(ls "$gpu_dir/drm" 2>/dev/null | grep "^card" | head -n1)"; freq="$(cat "$gpu_dir/drm/$drm_card/gt_max_freq_mhz" 2>/dev/null)"; fi;
printf "freq_mhz=%s\n" "$freq";
echo __IVORY_MEM__; grep MemTotal /proc/meminfo;
echo __IVORY_SWAP__; grep SwapTotal /proc/meminfo;
echo __IVORY_DISK__; df -Pk / | tail -n 1;
echo __IVORY_IP__; hostname -I 2>/dev/null;
echo __IVORY_LOCALE__; locale 2>/dev/null | grep ^LANG='`

var metricsSectionKeys = []string{"__IVORY_CPU__", "__IVORY_MEM__", "__IVORY_NET__"}

var infoSectionKeys = []string{
	"__IVORY_HOST__", "__IVORY_OS__", "__IVORY_KERNEL__", "__IVORY_UPTIME__",
	"__IVORY_CPU__", "__IVORY_GPU__", "__IVORY_MEM__", "__IVORY_SWAP__",
	"__IVORY_DISK__", "__IVORY_IP__", "__IVORY_LOCALE__",
}

var ErrInvalidPublicKey = errors.New("public key cannot be empty or contain newlines")

// NOTE: validate that is matches interface in compile-time
var _ platform.System = (*Plugin)(nil)

type Plugin struct {
	sshClient *ssh.Client
}

func NewPlugin(sshClient *ssh.Client) *Plugin {
	return &Plugin{sshClient}
}

func (p *Plugin) Metrics(connection platform.Connection) (*platform.Metrics, error) {
	result, err := p.execute(connection, MetricsCommand)
	if err != nil {
		return nil, err
	}
	return p.parseMetrics(strings.Join(result, "\n"))
}

func (p *Plugin) CopyId(connection platform.Connection, publicKey string) error {
	if strings.TrimSpace(publicKey) == "" || strings.ContainsAny(publicKey, "\r\n") {
		return ErrInvalidPublicKey
	}
	key := shellQuote(publicKey)
	command := fmt.Sprintf(`umask 077 && mkdir -p ~/.ssh && touch ~/.ssh/authorized_keys && chmod 700 ~/.ssh && chmod 600 ~/.ssh/authorized_keys && (grep -qxF -- %s ~/.ssh/authorized_keys || printf '%%s\n' %s >> ~/.ssh/authorized_keys)`, key, key)
	_, err := p.execute(connection, command)
	return err
}

func (p *Plugin) Logs(connection platform.Connection, filePath string, tail int, follow bool) console.Command {
	commandStr := "tail "
	if tail > 0 {
		commandStr += "-n " + strconv.Itoa(tail) + " "
	}
	if follow {
		commandStr += "-f "
	}
	commandStr += "-- " + shellQuote(filePath)
	command := p.sshClient.Command(p.mapToSshCommand(connection), commandStr)
	command.JobKeepAlive = false
	return command
}

func (p *Plugin) Processes(connection platform.Connection) ([]platform.Process, error) {
	result, err := p.execute(connection, ProcessesCommand)
	if err != nil {
		return nil, err
	}
	return p.parseProcesses(result)
}

func (p *Plugin) Info(connection platform.Connection) ([]platform.InfoItem, error) {
	result, err := p.execute(connection, InfoCommand)
	if err != nil {
		return nil, err
	}
	return p.parseInfo(strings.Join(result, "\n")), nil
}

func (p *Plugin) execute(connection platform.Connection, command string) ([]string, error) {
	cmd := p.sshClient.Command(p.mapToSshCommand(connection), command)
	return cmd.Execute()
}

func (p *Plugin) mapToSshCommand(conn platform.Connection) ssh.Connection {
	var prvKey *ed25519.PrivateKey
	if len(conn.PrivateKey) > 0 {
		pk := ed25519.PrivateKey(conn.PrivateKey)
		prvKey = &pk
	}
	return ssh.Connection{
		Host:       conn.Host,
		Port:       conn.Port,
		Username:   conn.Username,
		Password:   conn.Password,
		PrivateKey: prvKey,
	}
}

func shellQuote(value string) string {
	if value == "" {
		return "''"
	}
	return "'" + strings.ReplaceAll(value, "'", `'\''`) + "'"
}

func (p *Plugin) parseMetrics(output string) (*platform.Metrics, error) {
	sections := p.splitSections(output, metricsSectionKeys)

	for _, key := range metricsSectionKeys {
		if _, ok := sections[key]; !ok {
			return nil, fmt.Errorf("metrics output missing section %q", key)
		}
	}

	cpu, err := p.parseCpuMetrics(sections["__IVORY_CPU__"])
	if err != nil {
		return nil, err
	}
	memory, err := p.parseMemoryMetrics(sections["__IVORY_MEM__"])
	if err != nil {
		return nil, err
	}
	network, err := p.parseNetworkMetrics(sections["__IVORY_NET__"])
	if err != nil {
		return nil, err
	}

	return &platform.Metrics{
		Cpu:     cpu,
		Memory:  memory,
		Network: network,
	}, nil
}

// splitSections splits sentinel-marked command output (each section preceded
// by one of keys on its own line) into per-section lines.
func (p *Plugin) splitSections(output string, keys []string) map[string][]string {
	markers := make(map[string]bool, len(keys))
	for _, key := range keys {
		markers[key] = true
	}

	sections := map[string][]string{}
	current := ""

	for line := range strings.SplitSeq(output, "\n") {
		trimmed := strings.TrimSpace(line)
		if markers[trimmed] {
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

func (p *Plugin) parseCpuMetrics(lines []string) (platform.CpuMetrics, error) {
	if len(lines) == 0 {
		return platform.CpuMetrics{}, platform.ErrInvalidCpuMetrics
	}

	total, idle, err := p.parseCpuLine(lines[0])
	if err != nil {
		return platform.CpuMetrics{}, err
	}

	return platform.CpuMetrics{TotalTicks: total, IdleTicks: idle}, nil
}

func (p *Plugin) parseCpuLine(line string) (uint64, uint64, error) {
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

func (p *Plugin) parseMemoryMetrics(lines []string) (platform.MemoryMetrics, error) {
	if len(lines) < 2 {
		return platform.MemoryMetrics{}, platform.ErrInvalidMemoryMetrics
	}

	values := p.parseMeminfoFields(lines)
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

// parseMeminfoFields extracts numeric fields (in bytes) from /proc/meminfo-style
// "Key: value kB" lines. Unparsable lines are skipped rather than failing,
// since callers apply their own strictness on top (e.g. requiring MemTotal).
func (p *Plugin) parseMeminfoFields(lines []string) map[string]uint64 {
	values := make(map[string]uint64, len(lines))
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		value, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}
		values[strings.TrimSuffix(fields[0], ":")] = value * 1024
	}
	return values
}

func (p *Plugin) parseNetworkMetrics(lines []string) (platform.NetworkMetrics, error) {
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

// parseProcesses is lenient: a single malformed row (e.g. a command name with
// an unexpected shape) is skipped rather than failing the whole listing.
func (p *Plugin) parseProcesses(lines []string) ([]platform.Process, error) {
	processes := make([]platform.Process, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		fields := strings.Fields(trimmed)
		if len(fields) < 8 {
			continue
		}

		pid, err := strconv.Atoi(fields[0])
		if err != nil {
			continue
		}
		cpuPercent, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			continue
		}
		memPercent, err := strconv.ParseFloat(fields[3], 64)
		if err != nil {
			continue
		}
		rssKb, err := strconv.ParseUint(fields[4], 10, 64)
		if err != nil {
			continue
		}
		threads, err := strconv.Atoi(fields[5])
		if err != nil {
			continue
		}

		processes = append(processes, platform.Process{
			Pid:         pid,
			Program:     path.Base(fields[6]),
			Command:     strings.Join(fields[7:], " "),
			Threads:     threads,
			User:        fields[1],
			MemoryBytes: rssKb * 1024,
			MemPercent:  memPercent,
			CpuPercent:  cpuPercent,
		})
	}

	if len(processes) == 0 && len(lines) > 0 {
		return nil, platform.ErrInvalidProcesses
	}

	return processes, nil
}

// parseInfo is best-effort: it is a display-only inventory, so a missing or
// unparsable section is simply omitted rather than failing the whole call.
func (p *Plugin) parseInfo(output string) []platform.InfoItem {
	sections := p.splitSections(output, infoSectionKeys)

	items := make([]platform.InfoItem, 0, len(infoSectionKeys))
	add := func(key, value string) {
		if value == "" {
			return
		}
		items = append(items, platform.InfoItem{Key: key, Value: value})
	}

	add("Host", p.parseHost(sections["__IVORY_HOST__"]))
	add("OS", p.parseOsRelease(sections["__IVORY_OS__"]))
	add("Kernel", p.firstLine(sections["__IVORY_KERNEL__"]))
	add("Uptime", p.parseUptime(sections["__IVORY_UPTIME__"]))
	cpu, cores := p.parseCpuInfo(sections["__IVORY_CPU__"])
	add("CPU", cpu)
	add("CPU Cores", cores)
	add("GPU", p.parseGpu(sections["__IVORY_GPU__"]))
	add("Memory", p.parseMeminfoTotal(sections["__IVORY_MEM__"], "MemTotal"))
	add("Swap", p.parseMeminfoTotal(sections["__IVORY_SWAP__"], "SwapTotal"))
	add("Disk", p.parseDiskSummary(sections["__IVORY_DISK__"]))
	add("Local IP", p.parseLocalIp(sections["__IVORY_IP__"]))
	add("Locale", p.parseLocale(sections["__IVORY_LOCALE__"]))

	return items
}

func (p *Plugin) firstLine(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	return lines[0]
}

// parseHost prefers the DMI-reported hardware model ("<product_name>
// (<product_family>)", e.g. "21KDA04PCD (ThinkPad X1 Carbon Gen 12)") and
// falls back to the network hostname when DMI data isn't available (e.g.
// inside a VM/container without access to /sys/devices/virtual/dmi).
func (p *Plugin) parseHost(lines []string) string {
	var productName, productFamily, hostname string
	for _, line := range lines {
		if value, ok := strings.CutPrefix(line, "product_name="); ok {
			productName = strings.TrimSpace(value)
			continue
		}
		if value, ok := strings.CutPrefix(line, "product_family="); ok {
			productFamily = strings.TrimSpace(value)
			continue
		}
		if value, ok := strings.CutPrefix(line, "hostname="); ok {
			hostname = strings.TrimSpace(value)
		}
	}

	if productName == "" {
		return hostname
	}
	if productFamily == "" || productFamily == productName {
		return productName
	}
	return fmt.Sprintf("%s (%s)", productName, productFamily)
}

func (p *Plugin) parseOsRelease(lines []string) string {
	for _, line := range lines {
		if name, ok := strings.CutPrefix(line, "PRETTY_NAME="); ok {
			return strings.Trim(name, `"`)
		}
	}
	return ""
}

func (p *Plugin) parseUptime(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	fields := strings.Fields(lines[0])
	if len(fields) == 0 {
		return ""
	}
	seconds, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return ""
	}

	d := time.Duration(seconds) * time.Second
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	switch {
	case days > 0:
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	case hours > 0:
		return fmt.Sprintf("%dh %dm", hours, minutes)
	default:
		return fmt.Sprintf("%dm", minutes)
	}
}

// parseCpuInfo reads the two lines produced by `grep -m1 "model name"` and
// `grep -c ^processor` against /proc/cpuinfo.
func (p *Plugin) parseCpuInfo(lines []string) (model string, cores string) {
	for _, line := range lines {
		if _, value, found := strings.Cut(line, ":"); found && strings.HasPrefix(line, "model name") {
			model = strings.TrimSpace(value)
			continue
		}
		if _, err := strconv.Atoi(line); err == nil {
			cores = line + " cores"
		}
	}
	return model, cores
}

// intelIntegratedPciAddr is the PCI bus address conventionally reserved for
// the CPU's integrated graphics on x86 platforms - the same heuristic
// fastfetch uses (comparing against domain 0, bus 0, device 2, function 0)
// to classify a GPU as integrated vs discrete.
const intelIntegratedPciAddr = "0000:00:02.0"

// parseGpu reports "<model> @ <freq> GHz [Integrated|Discrete]", omitting
// whichever pieces aren't available (a discrete/multi-GPU or headless setup
// may not expose the max-frequency sysfs files this reads, and the
// Integrated/Discrete tag is only added when it can be determined with
// confidence - see the addr/vendor checks below).
func (p *Plugin) parseGpu(lines []string) string {
	var name, addr, freqMhz string
	for _, line := range lines {
		if value, ok := strings.CutPrefix(line, "name="); ok {
			name = strings.TrimSpace(value)
			continue
		}
		if value, ok := strings.CutPrefix(line, "addr="); ok {
			addr = strings.TrimSpace(value)
			continue
		}
		if value, ok := strings.CutPrefix(line, "freq_mhz="); ok {
			freqMhz = strings.TrimSpace(value)
		}
	}

	model := p.parseGpuModel(name)
	if model == "" {
		return ""
	}

	result := model
	if freq, err := strconv.ParseFloat(freqMhz, 64); err == nil && freq > 0 {
		result += fmt.Sprintf(" @ %.2f GHz", freq/1000)
	}
	switch {
	case addr == intelIntegratedPciAddr:
		result += " [Integrated]"
	case strings.Contains(strings.ToUpper(model), "NVIDIA"):
		// NVIDIA has never shipped an integrated GPU, so this is always safe.
		result += " [Discrete]"
	}
	return result
}

// parseGpuModel cleans up a `lspci | grep -i vga` line (e.g. "00:02.0 VGA
// compatible controller: Intel Corporation Meteor Lake-P [Intel Arc
// Graphics] (rev 08)") into a short vendor+model string ("Intel Arc
// Graphics"): it prefers the bracketed marketing name lspci reports over the
// raw PCI device string, and re-adds the vendor prefix if the bracketed name
// dropped it (e.g. NVIDIA's bracket omits "NVIDIA").
func (p *Plugin) parseGpuModel(line string) string {
	if line == "" {
		return ""
	}
	description := line
	if _, value, found := strings.Cut(line, ": "); found {
		description = value
	}
	description = strings.TrimSpace(description)

	vendor := ""
	if fields := strings.Fields(description); len(fields) > 0 {
		vendor = fields[0]
	}

	model := description
	if start := strings.LastIndex(description, "["); start != -1 {
		if end := strings.Index(description[start:], "]"); end != -1 {
			model = description[start+1 : start+end]
		}
	} else {
		if idx := strings.Index(description, " (rev "); idx != -1 {
			model = description[:idx]
		}
		model = strings.NewReplacer(" Corporation", "", " Inc.", "", " Co., Ltd.", "").Replace(model)
	}
	model = strings.TrimSpace(model)

	if vendor != "" && !strings.HasPrefix(model, vendor) {
		model = vendor + " " + model
	}
	return model
}

// parseMeminfoTotal reports the total capacity for a /proc/meminfo key (e.g.
// MemTotal, SwapTotal).
func (p *Plugin) parseMeminfoTotal(lines []string, totalKey string) string {
	total := p.parseMeminfoFields(lines)[totalKey]
	if total == 0 {
		return ""
	}
	return formatBytes(total)
}

// parseDiskSummary reports the total capacity of the root filesystem.
func (p *Plugin) parseDiskSummary(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	fields := strings.Fields(lines[0])
	if len(fields) < 2 {
		return ""
	}
	totalKb, err := strconv.ParseUint(fields[1], 10, 64)
	if err != nil || totalKb == 0 {
		return ""
	}
	return formatBytes(totalKb * 1024)
}

// parseLocalIp reports every address `hostname -I` returns (a node may have
// several, e.g. a LAN address plus a docker bridge address).
func (p *Plugin) parseLocalIp(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	fields := strings.Fields(lines[0])
	if len(fields) == 0 {
		return ""
	}
	return strings.Join(fields, ", ")
}

func (p *Plugin) parseLocale(lines []string) string {
	for _, line := range lines {
		if value, ok := strings.CutPrefix(line, "LANG="); ok {
			return value
		}
	}
	return ""
}

// formatBytes renders a byte count as a human-readable binary-unit string
// (e.g. "1.5 GiB").
func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
