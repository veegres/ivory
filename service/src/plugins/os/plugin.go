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
	DockerList(connection ssh.Connection) (*DockerResult, error)
	DockerDeploy(connection ssh.Connection, image, options string) (*DockerResult, error)
	DockerRun(connection ssh.Connection, options, image string) (*DockerResult, error)
	DockerStop(connection ssh.Connection, container string) (*DockerResult, error)
	DockerDelete(connection ssh.Connection, container string) (*DockerResult, error)
	DockerLogs(connection ssh.Connection, container string, tail int) (*DockerResult, error)
}

type PluginRegistry struct {
	clients map[Type]Adapter
}

func NewPluginRegistry() *PluginRegistry {
	return &PluginRegistry{
		clients: make(map[Type]Adapter),
	}
}

func (r *PluginRegistry) Register(clientType Type, client Adapter) {
	r.clients[clientType] = client
}

func (r *PluginRegistry) Get(t Type) (Adapter, error) {
	if client, ok := r.clients[t]; ok {
		return client, nil
	}
	return nil, ErrClientNotImplemented
}
