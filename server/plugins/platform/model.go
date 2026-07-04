package platform

type Plugin string

const (
	Linux Plugin = "linux"
)

// DeploySpec describes a single deployment in platform-neutral terms. Field
// values may contain {{placeholder}} templates; adapters render them as-is
// and interpolation happens later, right before deployment.
type DeploySpec struct {
	Name          string
	Hostname      string
	RestartPolicy string
	HostNetwork   bool
	Ports         []string
	Volumes       []VolumeMount
	Env           []EnvVar
}

type EnvVar struct {
	Name  string
	Value string
}

type VolumeMount struct {
	HostPath      string
	ContainerPath string
}

type Metrics struct {
	Cpu     CpuMetrics     `json:"cpu"`
	Memory  MemoryMetrics  `json:"memory"`
	Network NetworkMetrics `json:"network"`
}

type CpuMetrics struct {
	TotalTicks uint64 `json:"totalTicks"`
	IdleTicks  uint64 `json:"idleTicks"`
}

type MemoryMetrics struct {
	TotalBytes     uint64 `json:"totalBytes"`
	AvailableBytes uint64 `json:"availableBytes"`
}

type NetworkMetrics struct {
	ReceivedBytes    uint64 `json:"receivedBytes"`
	TransmittedBytes uint64 `json:"transmittedBytes"`
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
