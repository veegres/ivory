package database

import (
	"crypto/tls"
	"slices"
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

var PluginTypes = []PluginType{POSTGRES, ETCD, REDIS, CLICKHOUSE, ZOOKEEPER, MONGO}

func (p PluginType) String() string {
	return string(p)
}

func (p PluginType) Valid() bool {
	return slices.Contains(PluginTypes, p)
}

// Config is the persisted, user-facing database connection configuration.
type Config struct {
	Plugin PluginType `json:"plugin" binding:"omitempty,valid"`
	Host   string     `json:"host" binding:"required,notblank"`
	Port   int        `json:"port" binding:"required,min=1,max=65535"`
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

func (c SystemRequestCategory) Valid() bool {
	return c >= BLOAT && c <= OTHER
}

// SystemRequestVariety flags a SystemRequest's applicability, e.g. whether it
// only makes sense on the primary or needs a specific database selected.
type SystemRequestVariety int8

const (
	DatabaseSensitive SystemRequestVariety = iota
	MasterOnly
	ReplicaRecommended
)

func (v SystemRequestVariety) Valid() bool {
	return v >= DatabaseSensitive && v <= ReplicaRecommended
}

// SystemChartType identifies one of the built-in dashboard charts.
type SystemChartType string

const (
	Databases      SystemChartType = "Databases"
	Connections    SystemChartType = "Connections"
	DatabaseSize   SystemChartType = "Database Size"
	DatabaseUptime SystemChartType = "Database Uptime"
	Schemas        SystemChartType = "Schemas"
	TablesSize     SystemChartType = "Tables Size"
	IndexesSize    SystemChartType = "Indexes Size"
	TotalSize      SystemChartType = "Total Size"
)

var SystemChartTypes = []SystemChartType{Databases, Connections, DatabaseSize, DatabaseUptime, Schemas, TablesSize, IndexesSize, TotalSize}

func (c SystemChartType) Valid() bool {
	return slices.Contains(SystemChartTypes, c)
}

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
