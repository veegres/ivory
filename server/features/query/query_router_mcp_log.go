package query

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (r *McpRouter) QueryLogGet(ctx context.Context, req *mcp.CallToolRequest, in RouterIdRequest) (*mcp.CallToolResult, []DbResponse, error) {
	results, err := r.service.GetLog(in.Id)
	if err != nil {
		return nil, nil, err
	}
	return nil, results, nil
}

func (r *McpRouter) QueryLogDelete(ctx context.Context, req *mcp.CallToolRequest, in RouterIdRequest) (*mcp.CallToolResult, string, error) {
	if err := r.service.DeleteLog(in.Id); err != nil {
		return nil, "", err
	}
	return nil, "deleted", nil
}
