package query

import (
	"context"
	"ivory/core/config"
)

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

// queryContext is withQueryContext for the MCP engine.
func (r *McpRouter) queryContext(ctx context.Context, connection Connection) Context {
	return Context{Connection: connection, Session: config.McpCallerFrom(ctx).Session}
}
