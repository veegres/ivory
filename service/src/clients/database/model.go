package database

import (
	"crypto/tls"
	"ivory/src/features/cert"

	"github.com/google/uuid"
)

// COMMON (WEB AND SERVER)

type Database struct {
	Host   string  `json:"host"`
	Port   int     `json:"port"`
	Name   *string `json:"name"`
	Schema *string `json:"schema"`
}

type ConnectionRequest struct {
	Db           Database    `json:"db"`
	Certs        *cert.Certs `json:"certs"`
	CredentialId *uuid.UUID  `json:"credentialId"`
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

type QueryChartType string

const (
	Databases      QueryChartType = "Databases"
	Connections                   = "Connections"
	DatabaseSize                  = "Database Size"
	DatabaseUptime                = "Database Uptime"
	Schemas                       = "Schemas"
	TablesSize                    = "Tables Size"
	IndexesSize                   = "Indexes Size"
	TotalSize                     = "Total Size"
)

type QueryType int8

const (
	BLOAT QueryType = iota
	ACTIVITY
	REPLICATION
	STATISTIC
	OTHER
)

type QueryVariety int8

const (
	DatabaseSensitive QueryVariety = iota
	MasterOnly
	ReplicaRecommended
)

type Query struct {
	Name        string         `json:"name"`
	Type        *QueryType     `json:"type"`
	Description *string        `json:"description"`
	Query       string         `json:"query"`
	Varieties   []QueryVariety `json:"varieties"`
	Params      []string       `json:"params"`
}

// SPECIFIC (SERVER)

type Connection struct {
	Database    Database     `json:"database" form:"database"`
	Credentials *Credentials `json:"credentials" form:"credentials"`
	TlsConfig   *tls.Config  `json:"tlsConfig" form:"tlsConfig"`
}

type Context struct {
	Connection *Connection `json:"connection"`
	Session    string      `json:"session"`
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
