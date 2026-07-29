// Package database defines the plugin boundary for database engines (postgres,
// etcd, ...). A database.Adapter lets the rest of Ivory run queries, inspect
// schema, manage sessions and describe plugin capabilities without knowing
// which concrete engine it is talking to.
package database

import (
	"errors"
	"ivory/core/config"
)

var ErrCannotLimitWithoutTrim = errors.New("cannot limit query without trimming it")
var ErrDatabaseHostOrPortNotSpecified = errors.New("database host or port are not specified")
var ErrPasswordNotSet = errors.New("password is not set")

// Adapter is implemented by every database plugin (postgres, etcd, ...). It
// combines query execution, schema inspection, session management and
// self-description into the single interface the rest of Ivory depends on.
type Adapter interface {
	QueryExecutor
	SchemaInquirer
	SessionManager
	MetadataProvider
}

// QueryExecutor runs queries against the database and reports their results.
type QueryExecutor interface {
	// GetMany runs query and returns the first column of every row as a string.
	GetMany(ctx Context, query string, queryParams []any) ([]string, error)
	// GetOne runs query and returns the first column of the first row.
	GetOne(ctx Context, query string) (any, error)
	// GetFields runs query and returns its full result set (fields and rows),
	// honoring options for parameters, trimming and limiting.
	GetFields(ctx Context, query string, options *QueryOptions) (*QueryFields, error)
}

// SchemaInquirer lists the database's own structural objects.
type SchemaInquirer interface {
	// ListDatabases returns database names matching name.
	ListDatabases(ctx Context, name string) ([]string, error)
	// ListSchemas returns schema names matching name.
	ListSchemas(ctx Context, name string) ([]string, error)
	// ListTables returns table names within schema matching name.
	ListTables(ctx Context, schema string, name string) ([]string, error)
}

// SessionManager inspects and controls in-flight backend sessions.
type SessionManager interface {
	// Cancel asks the backend running pid to cancel its current query.
	Cancel(ctx Context, pid int) error
	// Terminate forcibly ends the backend running pid.
	Terminate(ctx Context, pid int) error
	// ActiveQueries returns the currently running queries for this connection's
	// application, in the same shape as GetFields.
	ActiveQueries(ctx Context, options *QueryOptions) (*QueryFields, error)
}

// MetadataProvider lets a database plugin describe itself: which optional
// features it supports and which built-in queries/charts it offers.
type MetadataProvider interface {
	// SupportedFeatures reports, for every db-related config.Feature this plugin
	// knows about, whether it supports it. A feature absent from the map is
	// not a database capability at all and is left unrestricted.
	SupportedFeatures() map[config.Feature]bool
	// SystemRequests returns the catalog of built-in queries this plugin ships
	// (activity, replication, bloat, ...), shown to the user as ready-made queries.
	SystemRequests() []SystemRequest
	// SystemCharts returns the queries backing the dashboard charts this
	// plugin supports, keyed by chart type.
	SystemCharts() map[SystemChartType]string
}
