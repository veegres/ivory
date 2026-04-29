package database

import (
	"errors"
	"ivory/src/features"
)

var ErrCannotLimitWithoutTrim = errors.New("cannot limit query without trimming it")
var ErrDatabaseHostOrPortNotSpecified = errors.New("database host or port are not specified")
var ErrPasswordNotSet = errors.New("password is not set")
var ErrClientNotImplemented = errors.New("client is not implemented")

type QueryExecutor interface {
	GetMany(ctx Context, query string, queryParams []any) ([]string, error)
	GetOne(ctx Context, query string) (any, error)
	GetFields(ctx Context, query string, options *QueryOptions) (*QueryFields, error)
}

type SchemaInquirer interface {
	ListDatabases(ctx Context, name string) ([]string, error)
	ListSchemas(ctx Context, name string) ([]string, error)
	ListTables(ctx Context, schema string, name string) ([]string, error)
}

type SessionManager interface {
	Cancel(ctx Context, pid int) error
	Terminate(ctx Context, pid int) error
	ActiveQueries(ctx Context, options *QueryOptions) (*QueryFields, error)
}

type MetadataProvider interface {
	SupportedFeatures() []features.Feature
	SystemRequests() []SystemRequest
	SystemCharts() map[SystemChartType]string
}

type Adapter interface {
	QueryExecutor
	SchemaInquirer
	SessionManager
	MetadataProvider
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

func (r *Registry) All() map[Type]Adapter {
	return r.clients
}
