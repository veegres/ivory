package database

import (
	"errors"
	"ivory/src/features"
)

var ErrClientNotImplemented = errors.New("client is not implemented")

type Client interface {
	GetApplicationName(session string) string
	GetMany(ctx Context, query string, queryParams []any) ([]string, error)
	GetOne(ctx Context, query string) (any, error)
	GetFields(ctx Context, query string, options *QueryOptions) (*QueryFields, error)
	Cancel(ctx Context, pid int) error
	Terminate(ctx Context, pid int) error
	SupportedFeatures() []features.Feature
}

type Registry struct {
	clients map[Type]Client
}

func NewRegistry() *Registry {
	return &Registry{
		clients: make(map[Type]Client),
	}
}

func (f *Registry) Register(clientType Type, client Client) {
	f.clients[clientType] = client
}

func (f *Registry) Get(t Type) (Client, error) {
	if client, ok := f.clients[t]; ok {
		return client, nil
	}
	return nil, ErrClientNotImplemented
}
