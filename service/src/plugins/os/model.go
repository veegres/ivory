package os

type Type int8

const (
	Linux Type = iota
)

type Metrics struct {
	Cpu     CpuMetrics     `json:"cpu"`
	Memory  MemoryMetrics  `json:"memory"`
	Network NetworkMetrics `json:"network"`
}

type CpuMetrics struct {
	UsagePercent float64 `json:"usagePercent"`
}

type MemoryMetrics struct {
	TotalBytes     uint64  `json:"totalBytes"`
	UsedBytes      uint64  `json:"usedBytes"`
	AvailableBytes uint64  `json:"availableBytes"`
	UsagePercent   float64 `json:"usagePercent"`
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
