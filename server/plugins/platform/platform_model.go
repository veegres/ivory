package platform

// PluginType identifies which platform adapter a deployment target uses.
type PluginType string

const (
	Docker PluginType = "docker"
)

// ValueEscape marks the single character after it as literal content the
// server injected into a command (a vault credential), whatever quoting
// surrounds it. It is deliberately not a backslash: a backslash in a command
// belongs to the template author and has to keep its ordinary shell meaning,
// so escaping both with the same character made a plugin's own `\"` and an
// escaped credential indistinguishable - and the tokenizer ate both.
const ValueEscape = '\x00'

// Metrics is a snapshot of resource usage for a node or a single container.
type Metrics struct {
	Cpu     CpuMetrics     `json:"cpu"`
	Memory  MemoryMetrics  `json:"memory"`
	Network NetworkMetrics `json:"network"`
}

// CpuMetrics reports cumulative CPU tick counters; callers derive a
// percentage from the delta between two samples over time.
type CpuMetrics struct {
	TotalTicks uint64 `json:"totalTicks"`
	IdleTicks  uint64 `json:"idleTicks"`
}

// MemoryMetrics reports memory usage in bytes.
type MemoryMetrics struct {
	TotalBytes     uint64 `json:"totalBytes"`
	AvailableBytes uint64 `json:"availableBytes"`
}

// NetworkMetrics reports cumulative network byte counters; callers derive a
// throughput from the delta between two samples over time.
type NetworkMetrics struct {
	ReceivedBytes    uint64 `json:"receivedBytes"`
	TransmittedBytes uint64 `json:"transmittedBytes"`
}

// InfoItem is a single arbitrary platform-reported detail about a node (e.g.
// "OS": "Ubuntu 22.04"). Content and order are entirely up to each adapter.
type InfoItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Process is a single row from a host-level process listing, analogous to
// what tools like btop/top display: identity, resource usage and the command.
type Process struct {
	Pid         int     `json:"pid"`
	Program     string  `json:"program"`
	Command     string  `json:"command"`
	Threads     int     `json:"threads"`
	User        string  `json:"user"`
	MemoryBytes uint64  `json:"memoryBytes"`
	MemPercent  float64 `json:"memPercent"`
	CpuPercent  float64 `json:"cpuPercent"`
}
