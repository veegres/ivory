package database

import (
	"crypto/tls"
)

// COMMON (WEB AND SERVER)

// PluginType identifies which database engine a Config connects to.
type PluginType string

const (
	POSTGRES   PluginType = "postgres"
	ETCD       PluginType = "etcd"
	REDIS      PluginType = "redis"
	CLICKHOUSE PluginType = "clickhouse"
	ZOOKEEPER  PluginType = "zookeeper"
	MONGO      PluginType = "mongo"
)

func (p PluginType) String() string {
	return string(p)
}

// Config is the persisted, user-facing database connection configuration.
type Config struct {
	Plugin PluginType `json:"plugin"`
	Host   string     `json:"host"`
	Port   int        `json:"port"`
	Name   *string    `json:"name"`
	Schema *string    `json:"schema"`
}

// SPECIFIC (SERVER)

// Connection is Config plus the runtime credentials and TLS settings needed
// to actually open a connection.
type Connection struct {
	Config      Config       `json:"config" form:"config"`
	Credentials *Credentials `json:"credentials" form:"credentials"`
	TlsConfig   *tls.Config  `json:"tlsConfig" form:"tlsConfig"`
}

// Context carries a Connection and the calling application's name through an
// Adapter call, e.g. so ActiveQueries can filter by application_name.
type Context struct {
	Connection  *Connection `json:"connection"`
	Application string      `json:"application"`
}

// Credentials authenticates a database Connection.
type Credentials struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

// QueryAnalysis is a shallow keyword-count analysis of a query, used to decide
// whether it is safe to append a LIMIT clause (see postgres addLimitToQuery).
type QueryAnalysis struct {
	LIMIT     int
	UPDATE    int
	DELETE    int
	INSERT    int
	SELECT    int
	FROM      int
	EXPLAIN   int
	Semicolon bool
}

// QueryOptions controls how GetFields runs and post-processes a query.
type QueryOptions struct {
	Params []any   `json:"params"`
	Limit  *string `json:"limit"`
	Trim   *bool   `json:"trim"`
}

// QueryField describes a single result column.
type QueryField struct {
	Name        string `json:"name"`
	DataType    string `json:"dataType"`
	DataTypeOID uint32 `json:"dataTypeOID"`
}

// QueryFields is the full result of a GetFields/ActiveQueries call: the
// columns, the rows, and metadata about how/where the query ran.
type QueryFields struct {
	Fields    []QueryField  `json:"fields"`
	Rows      [][]any       `json:"rows"`
	StartTime int64         `json:"startTime"`
	EndTime   int64         `json:"endTime"`
	Url       string        `json:"url"`
	Options   *QueryOptions `json:"options"`
}

// SystemRequestCategory groups a SystemRequest for display purposes.
type SystemRequestCategory int8

const (
	BLOAT SystemRequestCategory = iota
	ACTIVITY
	REPLICATION
	STATISTIC
	OTHER
)

// SystemRequestVariety flags a SystemRequest's applicability, e.g. whether it
// only makes sense on the primary or needs a specific database selected.
type SystemRequestVariety int8

const (
	DatabaseSensitive SystemRequestVariety = iota
	MasterOnly
	ReplicaRecommended
)

// SystemChartType identifies one of the built-in dashboard charts.
type SystemChartType string

const (
	Databases      SystemChartType = "Databases"
	Connections                    = "Connections"
	DatabaseSize                   = "Database Size"
	DatabaseUptime                 = "Database Uptime"
	Schemas                        = "Schemas"
	TablesSize                     = "Tables Size"
	IndexesSize                    = "Indexes Size"
	TotalSize                      = "Total Size"
)

// SystemRequest is one built-in, ready-made query a database plugin ships,
// shown to the user as a shortcut instead of writing the query by hand.
type SystemRequest struct {
	Name        string
	Type        SystemRequestCategory
	Description string
	Query       string
	Varieties   []SystemRequestVariety
	Params      []string
}
