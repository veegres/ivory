package model

import (
	"github.com/google/uuid"
)

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

type Query struct {
	Name        string        `json:"name"`
	Type        QueryType     `json:"type"`
	Creation    QueryCreation `json:"creation"`
	Description string        `json:"description"`
	Default     string        `json:"default"`
	Custom      string        `json:"custom"`
}

type QueryRequest struct {
	Name        *string    `json:"name"`
	Type        *QueryType `json:"type"`
	Description *string    `json:"description"`
	Query       string     `json:"query"`
}

type QueryRunRequest struct {
	QueryUuid   uuid.UUID `json:"queryUuid"`
	ClusterName string    `json:"clusterName"`
	Db          Database  `json:"db"`
}

type QueryField struct {
	Name        string `json:"name"`
	DataType    string `json:"dataType"`
	DataTypeOID uint32 `json:"dataTypeOID"`
}

type QueryRunResponse struct {
	Fields []QueryField `json:"fields"`
	Rows   [][]any      `json:"rows"`
}
