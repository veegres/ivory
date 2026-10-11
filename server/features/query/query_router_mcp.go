package query

import (
	"context"
	"errors"
	"ivory/core/config"

	"github.com/google/uuid"
)

// ErrQueryUuidRequired is the rule PostExecuteTemplateQuery states inline: a
// template query names the saved query it runs.
var ErrQueryUuidRequired = errors.New("queryUuid is required")

// McpRouter is this feature's half of the MCP engine, the way Router is its
// half of the http one. Which tools exist and what each one needs is declared
// in core/engine, beside the route table.
//
// A tool's structured result may be any JSON value, not only an object, so a
// handler hands back the service's own type rather than a transport shape.
type McpRouter struct {
	service *Service
}

func NewMcpRouter(service *Service) *McpRouter {
	return &McpRouter{service: service}
}

// McpQueryListRequest, McpQueryIdRequest and McpQueryUpdateRequest carry what
// the http routes take from the url path and query string. They exist because
// a tool has arguments and no url, not to restate a model.
type McpQueryListRequest struct {
	Type   *Type     `json:"type,omitempty" jsonschema:"only queries of this category"`
	Plugin *DbPlugin `json:"plugin,omitempty" jsonschema:"only queries for this database engine"`
}

type McpQueryIdRequest struct {
	Id string `json:"id" binding:"required" jsonschema:"uuid of the saved query"`
}

type McpQueryUpdateRequest struct {
	Id    string  `json:"id" binding:"required" jsonschema:"uuid of the saved query to overwrite"`
	Query Request `json:"query" binding:"required" jsonschema:"the new query definition"`
}

// queryContext is withQueryContext for the MCP engine. The session only ever
// becomes the connection's application_name, and an MCP client carries no
// browser session, so the caller's own name is what shows up in the engine's
// activity view.
func (r *McpRouter) queryContext(ctx context.Context, connection Connection) Context {
	session := config.McpCallerFrom(ctx).Username
	if session == "" {
		session = "mcp"
	}
	return Context{Connection: connection, Session: session}
}

func (r *McpRouter) queryId(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}
