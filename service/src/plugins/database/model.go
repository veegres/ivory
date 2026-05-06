package database

import (
	"crypto/tls"
)

// COMMON (WEB AND SERVER)

type Plugin string

const (
	POSTGRES Plugin = "postgres"
	ETCD     Plugin = "etcd"
)

type Config struct {
	Plugin Plugin  `json:"plugin"`
	Host   string  `json:"host"`
	Port   int     `json:"port"`
	Name   *string `json:"name"`
	Schema *string `json:"schema"`
}

// SPECIFIC (SERVER)

type Connection struct {
	Config      Config       `json:"config" form:"config"`
	Credentials *Credentials `json:"credentials" form:"credentials"`
	TlsConfig   *tls.Config  `json:"tlsConfig" form:"tlsConfig"`
}

type Context struct {
	Connection  *Connection `json:"connection"`
	Application string      `json:"application"`
}

type Credentials struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

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

type QueryOptions struct {
	Params []any   `json:"params"`
	Limit  *string `json:"limit"`
	Trim   *bool   `json:"trim"`
}

type QueryField struct {
	Name        string `json:"name"`
	DataType    string `json:"dataType"`
	DataTypeOID uint32 `json:"dataTypeOID"`
}

type QueryFields struct {
	Fields    []QueryField  `json:"fields"`
	Rows      [][]any       `json:"rows"`
	StartTime int64         `json:"startTime"`
	EndTime   int64         `json:"endTime"`
	Url       string        `json:"url"`
	Options   *QueryOptions `json:"options"`
}

type SystemRequestCategory int8

const (
	BLOAT SystemRequestCategory = iota
	ACTIVITY
	REPLICATION
	STATISTIC
	OTHER
)

type SystemRequestVariety int8

const (
	DatabaseSensitive SystemRequestVariety = iota
	MasterOnly
	ReplicaRecommended
)

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

type SystemRequest struct {
	Name        string
	Type        SystemRequestCategory
	Description string
	Query       string
	Varieties   []SystemRequestVariety
	Params      []string
}
