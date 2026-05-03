package os

type Plugin string

const (
	Linux Plugin = "linux"
)

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

type DockerResult struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exitCode"`
}
