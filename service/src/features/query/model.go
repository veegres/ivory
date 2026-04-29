package query

import (
	"ivory/src/clients/database"
	"ivory/src/features/cert"

	"github.com/google/uuid"
)

// COMMON (WEB AND SERVER)

type Type int8

const (
	BLOAT Type = iota
	ACTIVITY
	REPLICATION
	STATISTIC
	OTHER
)

type VarietyType int8

const (
	DatabaseSensitive VarietyType = iota
	MasterOnly
	ReplicaRecommended
)

type CreationType string

const (
	Manual CreationType = "manual"
	System              = "system"
)

type ChartType string

const (
	Databases      ChartType = "Databases"
	Connections              = "Connections"
	DatabaseSize             = "Database Size"
	DatabaseUptime           = "Database Uptime"
	Schemas                  = "Schemas"
	TablesSize               = "Tables Size"
	IndexesSize              = "Indexes Size"
	TotalSize                = "Total Size"
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
	Type       *ChartType `json:"type"`
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
