package model

import (
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

type QueryConnection struct {
	Db           Database   `json:"db"`
	Certs        *Certs     `json:"certs"`
	CredentialId *uuid.UUID `json:"credentialId"`
}

type QueryRunRequest struct {
	Connection  QueryConnection `json:"connection"`
	QueryUuid   *uuid.UUID      `json:"queryUuid"`
	Query       *string         `json:"query"`
	QueryParams []any           `json:"queryParams"`
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

type QueryField struct {
	Name        string `json:"name"`
	DataType    string `json:"dataType"`
	DataTypeOID uint32 `json:"dataTypeOID"`
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

type QueryFields struct {
	Fields    []QueryField `json:"fields"`
	Rows      [][]any      `json:"rows"`
	StartTime int64        `json:"startTime"`
	EndTime   int64        `json:"endTime"`
	Url       string       `json:"url"`
}

// SPECIFIC (SERVER)
