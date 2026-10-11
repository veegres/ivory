package management

import (
	"context"
	"ivory/core/config"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type McpRouter struct {
	service *Service
}

func NewMcpRouter(service *Service) *McpRouter {
	return &McpRouter{service: service}
}

// AppInfo mirrors the unguarded /info route. It reads the caller off the
// request context because a tool handler has no gin.Context, and because an
// McpRouter outlives any one request.
func (r *McpRouter) AppInfo(ctx context.Context, req *mcp.CallToolRequest, in struct{}) (*mcp.CallToolResult, *AppInfo, error) {
	caller := config.McpCallerFrom(ctx)
	return nil, r.service.GetAppInfo(caller.Authorised, caller.AuthEnabled, caller.Username, ""), nil
}
