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
	Description string         `json:"description"`
	Default     string         `json:"default"`
	Custom      string         `json:"custom"`
	CreatedAt   int64          `json:"createdAt"`
}

type QueryRequest struct {
	Name        *string        `json:"name"`
	Type        *QueryType     `json:"type"`
	Description *string        `json:"description"`
	Query       string         `json:"query"`
	Varieties   []QueryVariety `json:"varieties"`
	Params      []string       `json:"params"`
}

type QueryPostgresRequest struct {
	CredentialId uuid.UUID `json:"credentialId"`
	Db           Database  `json:"db"`
}

type QueryRunRequest struct {
	QueryPostgresRequest
	QueryUuid   uuid.UUID `json:"queryUuid"`
	QueryParams []any     `json:"queryParams"`
}

type QueryKillRequest struct {
	QueryPostgresRequest
	Pid int `json:"pid"`
}

type QueryChartRequest struct {
	QueryPostgresRequest
}

type QueryDatabasesRequest struct {
	QueryPostgresRequest
	Name string `json:"name"`
}

type QuerySchemasRequest struct {
	QueryPostgresRequest
	Name string `json:"name"`
}

type QueryTablesRequest struct {
	QueryPostgresRequest
	Schema string `json:"schema"`
	Name   string `json:"name"`
}

type QueryField struct {
	Name        string `json:"name"`
	DataType    string `json:"dataType"`
	DataTypeOID uint32 `json:"dataTypeOID"`
}

type QueryChart struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}

type QueryRunResponse struct {
	Fields []QueryField `json:"fields"`
	Rows   [][]any      `json:"rows"`
	Url    string       `json:"url"`
}

// SPECIFIC (SERVER)
