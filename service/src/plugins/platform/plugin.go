package platform

import (
	"errors"
	"ivory/src/clients/ssh"
)

var ErrInvalidCpuMetrics = errors.New("invalid cpu metrics output")
var ErrInvalidMemoryMetrics = errors.New("invalid memory metrics output")
var ErrInvalidNetworkMetrics = errors.New("invalid network metrics output")
var ErrClientNotImplemented = errors.New("client is not implemented")

type Adapter interface {
	PlatformManager
	DeploymentManager
}

type PlatformManager interface {
	Metrics(connection ssh.Connection) (*Metrics, error)
	CopyId(connection ssh.Connection, publicKey string) error
}

type DeploymentManager interface {
	List(connection ssh.Connection) (*OperationResult, error)
	Deploy(connection ssh.Connection, options, image string) (*OperationResult, error)
	Stop(connection ssh.Connection, name string) (*OperationResult, error)
	Delete(connection ssh.Connection, name string) (*OperationResult, error)
	Logs(connection ssh.Connection, name string, tail int) (*OperationResult, error)
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
