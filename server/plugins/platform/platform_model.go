package platform

// PluginType identifies which platform adapter a deployment target uses.
type PluginType string

const (
	Docker PluginType = "docker"
)

// renamedPlugins maps a platform key that used to be stored under a different
// name to the one it goes by now. Deployment templates and backups persist the
// key they were written with, so a rename has to stay readable or every stored
// template reads back as a platform no adapter answers to.
var renamedPlugins = map[PluginType]PluginType{
	"linux": Docker,
}

// Current resolves a stored platform key to the name it goes by now, leaving
// anything already current untouched.
func (p PluginType) Current() PluginType {
	if current, ok := renamedPlugins[p]; ok {
		return current
	}
	return p
}

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
