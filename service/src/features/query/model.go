package query

import (
	"ivory/src/clients/database"
	"ivory/src/features/cert"

	"github.com/google/uuid"
)

// COMMON (WEB AND SERVER)

type Type = database.SystemRequestCategory

const (
	BLOAT       = database.BLOAT
	ACTIVITY    = database.ACTIVITY
	REPLICATION = database.REPLICATION
	STATISTIC   = database.STATISTIC
	OTHER       = database.OTHER
)

type VarietyType = database.SystemRequestVariety

const (
	DatabaseSensitive  = database.DatabaseSensitive
	MasterOnly         = database.MasterOnly
	ReplicaRecommended = database.ReplicaRecommended
)

type CreationType string

const (
	Manual CreationType = "manual"
	System              = "system"
)

type ChartType = database.SystemChartType

const (
	Databases      = database.Databases
	Connections    = database.Connections
	DatabaseSize   = database.DatabaseSize
	DatabaseUptime = database.DatabaseUptime
	Schemas        = database.Schemas
	TablesSize     = database.TablesSize
	IndexesSize    = database.IndexesSize
	TotalSize      = database.TotalSize
)

type Request struct {
	Name        string        `json:"name"`
	Type        *Type         `json:"type"`
	Description *string       `json:"description"`
	Query       string        `json:"query"`
	Varieties   []VarietyType `json:"varieties"`
	Params      []string      `json:"params"`
}

type Response struct {
	Id          uuid.UUID     `json:"id"`
	Name        string        `json:"name"`
	Type        Type          `json:"type"`
	Creation    CreationType  `json:"creation"`
	Varieties   []VarietyType `json:"varieties"`
	Params      []string      `json:"params"`
	Description *string       `json:"description"`
	Default     string        `json:"default"`
	Custom      string        `json:"custom"`
	CreatedAt   int64         `json:"createdAt"`
}

type Connection struct {
	Db      database.Config `json:"db"`
	Certs   *cert.Certs     `json:"certs"`
	VaultId *uuid.UUID      `json:"vaultId"`
}

type TemplateRequest struct {
	Connection Connection `json:"connection"`
	Options    *DbOptions `json:"options"`
	QueryUuid  *uuid.UUID `json:"queryUuid"`
}

type ConsoleRequest struct {
	Connection Connection `json:"connection"`
	Options    *DbOptions `json:"options"`
	Query      string     `json:"query"`
}

type KillRequest struct {
	Connection Connection `json:"connection"`
	Pid        int        `json:"pid"`
}

type ChartRequest struct {
	Connection Connection `json:"connection"`
	Type       ChartType  `json:"type"`
}

type DatabasesRequest struct {
	Connection Connection `json:"connection"`
	Name       string     `json:"name"`
}

type SchemasRequest struct {
	Connection Connection `json:"connection"`
	Name       string     `json:"name"`
}

type TablesRequest struct {
	Connection Connection `json:"connection"`
	Schema     string     `json:"schema"`
	Name       string     `json:"name"`
}

type Chart struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}

type DbOptions = database.QueryOptions
type DbRow = database.QueryField
type DbResponse = database.QueryFields

// SPECIFIC (SERVER)

type Context struct {
	Connection Connection
	Session    string
}
