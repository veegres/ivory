package query

import (
	"ivory/core/service/cert"
	"ivory/plugins/database"

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
	Plugin      DbPlugin      `json:"plugin"`
	Description *string       `json:"description"`
	Query       string        `json:"query"`
	Varieties   []VarietyType `json:"varieties"`
	Params      []string      `json:"params"`
}

type Response struct {
	Id          uuid.UUID     `json:"id"`
	Name        string        `json:"name"`
	Type        Type          `json:"type"`
	Plugin      DbPlugin      `json:"plugin"`
	Creation    CreationType  `json:"creation"`
	Varieties   []VarietyType `json:"varieties"`
	Params      []string      `json:"params"`
	Description *string       `json:"description"`
	Default     string        `json:"default"`
	Custom      string        `json:"custom"`
	CreatedAt   int64         `json:"createdAt"`
}

type DbPlugin = database.Plugin

type DbConfig struct {
	Plugin DbPlugin `json:"plugin"`
	Host   string   `json:"host"`
	Port   int      `json:"port"`
	Name   *string  `json:"name"`
	Schema *string  `json:"schema"`
}

type Connection struct {
	Db      DbConfig    `json:"db"`
	Certs   *cert.Certs `json:"certs"`
	VaultId *uuid.UUID  `json:"vaultId"`
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

type DbOptions struct {
	Params []any   `json:"params"`
	Limit  *string `json:"limit"`
	Trim   *bool   `json:"trim"`
}

type DbRow struct {
	Name        string `json:"name"`
	DataType    string `json:"dataType"`
	DataTypeOID uint32 `json:"dataTypeOID"`
}

type DbResponse struct {
	Fields    []DbRow    `json:"fields"`
	Rows      [][]any    `json:"rows"`
	StartTime int64      `json:"startTime"`
	EndTime   int64      `json:"endTime"`
	Url       string     `json:"url"`
	Options   *DbOptions `json:"options"`
}

// SPECIFIC (SERVER)

type Context struct {
	Connection Connection
	Session    string
}

func mapSystemRequest(plugin DbPlugin, req database.SystemRequest) Request {
	t := req.Type
	v := make([]VarietyType, len(req.Varieties))
	for i, val := range req.Varieties {
		v[i] = val
	}
	return Request{
		Name:        req.Name,
		Type:        &t,
		Plugin:      plugin,
		Description: &req.Description,
		Query:       req.Query,
		Varieties:   v,
		Params:      req.Params,
	}
}

func mapDbConfig(c DbConfig) database.Config {
	return database.Config{
		Plugin: c.Plugin,
		Host:   c.Host,
		Port:   c.Port,
		Name:   c.Name,
		Schema: c.Schema,
	}
}

func mapDbOptions(o *DbOptions) *database.QueryOptions {
	if o == nil {
		return nil
	}
	return &database.QueryOptions{
		Params: o.Params,
		Limit:  o.Limit,
		Trim:   o.Trim,
	}
}

func mapDbResponse(r *database.QueryFields) *DbResponse {
	if r == nil {
		return nil
	}
	fields := make([]DbRow, len(r.Fields))
	for i, f := range r.Fields {
		fields[i] = DbRow{Name: f.Name, DataType: f.DataType, DataTypeOID: f.DataTypeOID}
	}
	var opts *DbOptions
	if r.Options != nil {
		opts = &DbOptions{Params: r.Options.Params, Limit: r.Options.Limit, Trim: r.Options.Trim}
	}
	return &DbResponse{
		Fields:    fields,
		Rows:      r.Rows,
		StartTime: r.StartTime,
		EndTime:   r.EndTime,
		Url:       r.Url,
		Options:   opts,
	}
}
