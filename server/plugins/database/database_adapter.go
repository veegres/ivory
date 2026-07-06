package database

import (
	"errors"
	"ivory/core/config"
)

var ErrCannotLimitWithoutTrim = errors.New("cannot limit query without trimming it")
var ErrDatabaseHostOrPortNotSpecified = errors.New("database host or port are not specified")
var ErrPasswordNotSet = errors.New("password is not set")

type Adapter interface {
	QueryExecutor
	SchemaInquirer
	SessionManager
	MetadataProvider
}

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
	// SupportedFeatures reports, for every db-related env.Feature this plugin
	// knows about, whether it supports it. A feature absent from the map is
	// not a database capability at all and is left unrestricted.
	SupportedFeatures() map[env.Feature]bool
	SystemRequests() []SystemRequest
	SystemCharts() map[SystemChartType]string
}
