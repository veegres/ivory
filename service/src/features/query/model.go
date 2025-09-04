package query

import (
	"ivory/src/features/cert"

	"github.com/google/uuid"
)

// COMMON (WEB AND SERVER)

type QueryType int8

const (
	BLOAT QueryType = iota
	ACTIVITY
	REPLICATION
	STATISTIC
	OTHER
)

type QueryCreation string

const (
	Manual QueryCreation = "manual"
	System               = "system"
)

type QueryVariety int8

const (
	DatabaseSensitive QueryVariety = iota
	MasterOnly
	ReplicaRecommended
)

type Query struct {
	Id          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Type        QueryType      `json:"type"`
	Creation    QueryCreation  `json:"creation"`
	Varieties   []QueryVariety `json:"varieties"`
	Params      []string       `json:"params"`
	Description *string        `json:"description"`
	Default     string         `json:"default"`
	Custom      string         `json:"custom"`
	CreatedAt   int64          `json:"createdAt"`
}

type QueryRequest struct {
	Name        string         `json:"name"`
	Type        *QueryType     `json:"type"`
	Description *string        `json:"description"`
	Query       string         `json:"query"`
	Varieties   []QueryVariety `json:"varieties"`
	Params      []string       `json:"params"`
}

// TODO 1. this class is to general, I'm not sure if it should be here
//  2. we have duplicated class in Bloat feature that is called DbConnection (we should change it and always provide clusterId?)
type Database struct {
	Host string  `json:"host"`
	Port int     `json:"port"`
	Name *string `json:"name"`
}

type QueryConnection struct {
	Db           Database    `json:"db"`
	Certs        *cert.Certs `json:"certs"`
	CredentialId *uuid.UUID  `json:"credentialId"`
}

type QueryContext struct {
	Connection QueryConnection
	Session    string
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

type QueryRunRequest struct {
	Connection   QueryConnection `json:"connection"`
	QueryUuid    *uuid.UUID      `json:"queryUuid"`
	Query        *string         `json:"query"`
	QueryOptions *QueryOptions   `json:"queryOptions"`
}

type QueryKillRequest struct {
	Connection QueryConnection `json:"connection"`
	Pid        int             `json:"pid"`
}

type QueryChartRequest struct {
	Connection QueryConnection `json:"connection"`
	Type       *QueryChartType `json:"type"`
}

type QueryDatabasesRequest struct {
	Connection QueryConnection `json:"connection"`
	Name       string          `json:"name"`
}

type QuerySchemasRequest struct {
	Connection QueryConnection `json:"connection"`
	Name       string          `json:"name"`
}

type QueryTablesRequest struct {
	Connection QueryConnection `json:"connection"`
	Schema     string          `json:"schema"`
	Name       string          `json:"name"`
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

type QueryChart struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}

// SPECIFIC (SERVER)
