package platform

import (
	"errors"
	"ivory/clients/console"
)

var ErrInvalidCpuMetrics = errors.New("invalid cpu metrics output")
var ErrInvalidMemoryMetrics = errors.New("invalid memory metrics output")
var ErrInvalidNetworkMetrics = errors.New("invalid network metrics output")
var ErrInvalidContainerMetrics = errors.New("invalid container metrics output")

// Connection contains the bare minimum details to execute commands on a platform node.
type Connection struct {
	Host       string
	Port       int
	Username   string
	Password   string
	PrivateKey []byte
}

type Adapter interface {
	VmManager
	ContainerManager
}

type VmManager interface {
	Metrics(connection Connection) (*Metrics, error)
	CopyId(connection Connection, publicKey string) error
	Logs(connection Connection, path string, tail int, follow bool) console.Command
}

type ContainerManager interface {
	ListContainer(connection Connection) console.Command
	UpContainer(connection Connection, options, image string) console.Command
	DownContainer(connection Connection, name string) console.Command
	StartContainer(connection Connection, name string) console.Command
	StopContainer(connection Connection, name string) console.Command
	RestartContainer(connection Connection, name string) console.Command
	LogsContainer(connection Connection, name string, tail int, follow bool) console.Command
	MetricsContainer(connection Connection, name string) (*Metrics, error)
}
