package query

import (
	"ivory/core/service/cert"
	"ivory/plugins/database"
	"slices"

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
	Name        string        `json:"name" binding:"required,notblank"`
	Type        *Type         `json:"type" binding:"required,valid"`
	Plugin      DbPlugin      `json:"plugin" binding:"omitempty,valid"`
	Description *string       `json:"description"`
	Query       string        `json:"query" binding:"required,notblank"`
	Varieties   []VarietyType `json:"varieties" binding:"dive,valid"`
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

type DbPlugin = database.PluginType

type DbConfig struct {
	Plugin DbPlugin `json:"plugin" binding:"required,valid"`
	Host   string   `json:"host" binding:"required,notblank"`
	Port   int      `json:"port" binding:"required,min=1,max=65535"`
	Name   *string  `json:"name"`
	Schema *string  `json:"schema"`
}

type Connection struct {
	Db      DbConfig    `json:"db" binding:"required"`
	Certs   *cert.Certs `json:"certs"`
	VaultId *uuid.UUID  `json:"vaultId"`
}

type RouterTemplateRequest struct {
	Connection Connection `json:"connection" binding:"required"`
	Options    *DbOptions `json:"options"`
	QueryUuid  uuid.UUID  `json:"queryUuid" binding:"required"`
}

type RouterConsoleRequest struct {
	Connection Connection `json:"connection" binding:"required"`
	Options    *DbOptions `json:"options"`
	Query      string     `json:"query" binding:"required,notblank"`
}

type RouterKillRequest struct {
	Connection Connection `json:"connection" binding:"required"`
	Pid        int        `json:"pid" binding:"required,min=1"`
}

type RouterChartRequest struct {
	Connection Connection `json:"connection" binding:"required"`
	Type       ChartType  `json:"type" binding:"required,valid"`
}

type RouterDatabasesRequest struct {
	Connection Connection `json:"connection" binding:"required"`
	Name       string     `json:"name"`
}

type RouterSchemasRequest struct {
	Connection Connection `json:"connection" binding:"required"`
	Name       string     `json:"name"`
}

type RouterTablesRequest struct {
	Connection Connection `json:"connection" binding:"required"`
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

type RouterListRequest struct {
	Type   *Type     `json:"type,omitempty" form:"type" binding:"omitempty,valid" jsonschema:"only queries of this category"`
	Plugin *DbPlugin `json:"plugin,omitempty" form:"plugin" binding:"omitempty,valid" jsonschema:"only queries for this database engine"`
}
type RouterIdRequest struct {
	Id uuid.UUID `json:"id" binding:"required" jsonschema:"uuid of the saved query"`
}
type RouterUpdateRequest struct {
	Id    uuid.UUID `json:"id" binding:"required" jsonschema:"uuid of the saved query to overwrite"`
	Query Request   `json:"query" binding:"required" jsonschema:"the new query definition"`
}

// SPECIFIC (SERVER)

type Context struct {
	Connection Connection
	Session    string
}

func mapSystemRequest(plugin DbPlugin, req database.SystemRequest) Request {
	t := req.Type
	v := make([]VarietyType, len(req.Varieties))
	copy(v, req.Varieties)
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

// withShipped brings a stored system query up to the version a plugin ships now.
// Description and type are not the user's to change, so they always follow;
// Default is the shipped text, so a reset restores the current one. The text,
// params and varieties follow only while the user has not edited the query.
func (r Response) withShipped(shipped Request) (Response, bool) {
	updated := r
	updated.Type = *shipped.Type
	updated.Description = shipped.Description
	updated.Default = shipped.Query
	if r.Custom == r.Default {
		updated.Custom = shipped.Query
		updated.Params = shipped.Params
		updated.Varieties = shipped.Varieties
	}
	changed := updated.Type != r.Type ||
		updated.Default != r.Default ||
		updated.Custom != r.Custom ||
		!slices.Equal(updated.Params, r.Params) ||
		!slices.Equal(updated.Varieties, r.Varieties) ||
		(updated.Description == nil) != (r.Description == nil) ||
		(updated.Description != nil && r.Description != nil && *updated.Description != *r.Description)
	return updated, changed
}

func (c DbConfig) toDatabaseConfig() database.Config {
	return database.Config{
		Plugin: c.Plugin,
		Host:   c.Host,
		Port:   c.Port,
		Name:   c.Name,
		Schema: c.Schema,
	}
}

func (o *DbOptions) toQueryOptions() *database.QueryOptions {
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
	rows := r.Rows
	if rows == nil {
		rows = make([][]any, 0)
	}
	var opts *DbOptions
	if r.Options != nil {
		opts = &DbOptions{Params: r.Options.Params, Limit: r.Options.Limit, Trim: r.Options.Trim}
	}
	return &DbResponse{
		Fields:    fields,
		Rows:      rows,
		StartTime: r.StartTime,
		EndTime:   r.EndTime,
		Url:       r.Url,
		Options:   opts,
	}
}
