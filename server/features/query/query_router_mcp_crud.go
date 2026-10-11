package query

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (r *McpRouter) QueryList(ctx context.Context, req *mcp.CallToolRequest, in McpQueryListRequest) (*mcp.CallToolResult, []Response, error) {
	list, err := r.service.GetList(in.Type, in.Plugin)
	if err != nil {
		return nil, nil, err
	}
	return nil, list, nil
}

func (r *McpRouter) QueryCreate(ctx context.Context, req *mcp.CallToolRequest, in Request) (*mcp.CallToolResult, *Response, error) {
	_, created, err := r.service.Create(Manual, in)
	if err != nil {
		return nil, nil, err
	}
	return nil, created, nil
}

func (r *McpRouter) QueryUpdate(ctx context.Context, req *mcp.CallToolRequest, in McpQueryUpdateRequest) (*mcp.CallToolResult, *Response, error) {
	id, err := r.queryId(in.Id)
	if err != nil {
		return nil, nil, err
	}
	_, updated, err := r.service.Update(id, in.Query)
	if err != nil {
		return nil, nil, err
	}
	return nil, updated, nil
}

func (r *McpRouter) QueryDelete(ctx context.Context, req *mcp.CallToolRequest, in McpQueryIdRequest) (*mcp.CallToolResult, string, error) {
	id, err := r.queryId(in.Id)
	if err != nil {
		return nil, "", err
	}
	if err := r.service.Delete(id); err != nil {
		return nil, "", err
	}
	return nil, "deleted", nil
}
