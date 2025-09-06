package query

import (
	"ivory/src/clients/database"
	"ivory/src/features/cert"

	"github.com/google/uuid"
)

// COMMON (WEB AND SERVER)

type QueryCreation string

const (
	Manual QueryCreation = "manual"
	System               = "system"
)

type Query struct {
	Id          uuid.UUID               `json:"id"`
	Name        string                  `json:"name"`
	Type        database.QueryType      `json:"type"`
	Creation    QueryCreation           `json:"creation"`
	Varieties   []database.QueryVariety `json:"varieties"`
	Params      []string                `json:"params"`
	Description *string                 `json:"description"`
	Default     string                  `json:"default"`
	Custom      string                  `json:"custom"`
	CreatedAt   int64                   `json:"createdAt"`
}

type QueryConnection struct {
	Db           database.Database `json:"db"`
	Certs        *cert.Certs       `json:"certs"`
	CredentialId *uuid.UUID        `json:"credentialId"`
}

type QueryRunRequest struct {
	Connection   QueryConnection        `json:"connection"`
	QueryUuid    *uuid.UUID             `json:"queryUuid"`
	Query        *string                `json:"query"`
	QueryOptions *database.QueryOptions `json:"queryOptions"`
}

type QueryKillRequest struct {
	Connection QueryConnection `json:"connection"`
	Pid        int             `json:"pid"`
}

type QueryChartRequest struct {
	Connection QueryConnection          `json:"connection"`
	Type       *database.QueryChartType `json:"type"`
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

// SPECIFIC (SERVER)

type QueryContext struct {
	Connection QueryConnection
	Session    string
}
