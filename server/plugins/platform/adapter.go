package platform

import (
	"errors"
	"ivory/clients/console"
)

var ErrInvalidCpuMetrics = errors.New("invalid cpu metrics output")
var ErrInvalidMemoryMetrics = errors.New("invalid memory metrics output")
var ErrInvalidNetworkMetrics = errors.New("invalid network metrics output")

// Connection contains the bare minimum details to execute commands on a platform node.
type Connection struct {
	Host       string
	Port       int
	Username   string
	Password   string
	PrivateKey []byte
}

type Adapter interface {
	PlatformManager
	DeploymentManager
}

type PlatformManager interface {
	Metrics(connection Connection) (*Metrics, error)
	CopyId(connection Connection, publicKey string) error
}

type DeploymentManager interface {
	ListCommand(connection Connection) console.Command
	DeployCommand(connection Connection, options, image string) console.Command
	StopCommand(connection Connection, name string) console.Command
	DeleteCommand(connection Connection, name string) console.Command
	LogsCommand(connection Connection, name string, tail int) console.Command
}
