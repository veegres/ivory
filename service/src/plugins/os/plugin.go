package os

import (
	"errors"
	"ivory/src/clients/ssh"
)

var ErrInvalidCpuMetrics = errors.New("invalid cpu metrics output")
var ErrInvalidMemoryMetrics = errors.New("invalid memory metrics output")
var ErrInvalidNetworkMetrics = errors.New("invalid network metrics output")
var ErrClientNotImplemented = errors.New("client is not implemented")

type Adapter interface {
	Metrics(connection ssh.Connection) (*Metrics, error)
	DockerManager
}

type DockerManager interface {
	DockerList(connection ssh.Connection) (*Docker, error)
	DockerDeploy(connection ssh.Connection, image, options string) (*Docker, error)
	DockerRun(connection ssh.Connection, options, image string) (*Docker, error)
	DockerStop(connection ssh.Connection, container string) (*Docker, error)
	DockerDelete(connection ssh.Connection, container string) (*Docker, error)
	DockerLogs(connection ssh.Connection, container string, tail int) (*Docker, error)
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
