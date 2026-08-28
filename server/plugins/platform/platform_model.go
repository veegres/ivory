package platform

// Plugin identifies which platform adapter a deployment target uses.
type Plugin string

const (
	Linux Plugin = "linux"
)

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
