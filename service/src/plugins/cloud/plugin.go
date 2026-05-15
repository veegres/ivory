package cloud

import (
	"errors"
	"ivory/src/clients/ssh"
)

var ErrInvalidCpuMetrics = errors.New("invalid cpu metrics output")
var ErrInvalidMemoryMetrics = errors.New("invalid memory metrics output")
var ErrInvalidNetworkMetrics = errors.New("invalid network metrics output")
var ErrClientNotImplemented = errors.New("client is not implemented")

type Adapter interface {
	CloudManager
	ContainerManager
}

type CloudManager interface {
	Metrics(connection ssh.Connection) (*Metrics, error)
	CopyId(connection ssh.Connection, publicKey string) error
}

type ContainerManager interface {
	ContainerList(connection ssh.Connection) (*Container, error)
	ContainerRun(connection ssh.Connection, options, image string) (*Container, error)
	ContainerStop(connection ssh.Connection, container string) (*Container, error)
	ContainerDelete(connection ssh.Connection, container string) (*Container, error)
	ContainerLogs(connection ssh.Connection, container string, tail int) (*Container, error)
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
