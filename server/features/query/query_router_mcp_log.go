package query

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (r *McpRouter) QueryLogGet(ctx context.Context, req *mcp.CallToolRequest, in McpQueryIdRequest) (*mcp.CallToolResult, []DbResponse, error) {
	id, err := r.queryId(in.Id)
	if err != nil {
		return nil, nil, err
	}
	results, err := r.service.GetLog(id)
	if err != nil {
		return nil, nil, err
	}
	return nil, results, nil
}

func (r *McpRouter) QueryLogDelete(ctx context.Context, req *mcp.CallToolRequest, in McpQueryIdRequest) (*mcp.CallToolResult, string, error) {
	id, err := r.queryId(in.Id)
	if err != nil {
		return nil, "", err
	}
	if err := r.service.DeleteLog(id); err != nil {
		return nil, "", err
	}
	return nil, "deleted", nil
}
