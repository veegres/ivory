package platform

import (
	"errors"
	"ivory/features/job"
)

var ErrInvalidCpuMetrics = errors.New("invalid cpu metrics output")
var ErrInvalidMemoryMetrics = errors.New("invalid memory metrics output")
var ErrInvalidNetworkMetrics = errors.New("invalid network metrics output")
var ErrClientNotImplemented = errors.New("client is not implemented")

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
	ListCommand(connection Connection) job.Command
	DeployCommand(connection Connection, options, image string) job.Command
	StopCommand(connection Connection, name string) job.Command
	DeleteCommand(connection Connection, name string) job.Command
	LogsCommand(connection Connection, name string, tail int) job.Command
}

type PluginRegistry struct {
	clients map[Plugin]Adapter
}

func NewPluginRegistry() *PluginRegistry {
	return &PluginRegistry{
		clients: make(map[Plugin]Adapter),
	}
}

func (r *PluginRegistry) Register(name Plugin, client Adapter) {
	r.clients[name] = client
}

func (r *PluginRegistry) Get(t Plugin) (Adapter, error) {
	if client, ok := r.clients[t]; ok {
		return client, nil
	}
	return nil, ErrClientNotImplemented
}
