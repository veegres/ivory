package tag

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// McpRouter is this feature's half of the MCP engine, the way Router is its
// half of the http one: it binds a tool's arguments, calls the service and
// hands back what the service returned. Which tools exist and what each one
// needs is declared in core/engine, beside the route table.
//
// A tool's structured result may be any JSON value, not only an object, so a
// handler returns the service's own type. There is no transport shape to wrap
// it in.
type McpRouter struct {
	tagService *Service
}

func NewMcpRouter(tagService *Service) *McpRouter {
	return &McpRouter{tagService: tagService}
}

func (r *McpRouter) TagList(ctx context.Context, req *mcp.CallToolRequest, in struct{}) (*mcp.CallToolResult, map[string][]string, error) {
	tags, err := r.tagService.GetMap()
	if err != nil {
		return nil, nil, err
	}
	return nil, tags, nil
}
