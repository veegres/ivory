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
}

type Registry struct {
	clients map[Type]Adapter
}

func NewRegistry() *Registry {
	return &Registry{
		clients: make(map[Type]Adapter),
	}
}

func (r *Registry) Register(clientType Type, client Adapter) {
	r.clients[clientType] = client
}

func (r *Registry) Get(t Type) (Adapter, error) {
	if client, ok := r.clients[t]; ok {
		return client, nil
	}
	return nil, ErrClientNotImplemented
}
