package query

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (r *McpRouter) QueryConsole(ctx context.Context, req *mcp.CallToolRequest, in ConsoleRequest) (*mcp.CallToolResult, *DbResponse, error) {
	result, err := r.service.ConsoleQuery(r.queryContext(ctx, in.Connection), in.Query, in.Options)
	if err != nil {
		return nil, nil, err
	}
	return nil, result, nil
}

func (r *McpRouter) QueryTemplate(ctx context.Context, req *mcp.CallToolRequest, in TemplateRequest) (*mcp.CallToolResult, *DbResponse, error) {
	if in.QueryUuid == nil {
		return nil, nil, ErrQueryUuidRequired
	}
	result, err := r.service.TemplateQuery(r.queryContext(ctx, in.Connection), *in.QueryUuid, in.Options)
	if err != nil {
		return nil, nil, err
	}
	return nil, result, nil
}

func (r *McpRouter) QueryActivity(ctx context.Context, req *mcp.CallToolRequest, in Connection) (*mcp.CallToolResult, *DbResponse, error) {
	result, err := r.service.RunningQueriesByApplicationName(r.queryContext(ctx, in))
	if err != nil {
		return nil, nil, err
	}
	return nil, result, nil
}

func (r *McpRouter) QueryDatabases(ctx context.Context, req *mcp.CallToolRequest, in DatabasesRequest) (*mcp.CallToolResult, []string, error) {
	names, err := r.service.DatabasesQuery(r.queryContext(ctx, in.Connection), in.Name)
	if err != nil {
		return nil, nil, err
	}
	return nil, names, nil
}

func (r *McpRouter) QuerySchemas(ctx context.Context, req *mcp.CallToolRequest, in SchemasRequest) (*mcp.CallToolResult, []string, error) {
	names, err := r.service.SchemasQuery(r.queryContext(ctx, in.Connection), in.Name)
	if err != nil {
		return nil, nil, err
	}
	return nil, names, nil
}

func (r *McpRouter) QueryTables(ctx context.Context, req *mcp.CallToolRequest, in TablesRequest) (*mcp.CallToolResult, []string, error) {
	names, err := r.service.TablesQuery(r.queryContext(ctx, in.Connection), in.Schema, in.Name)
	if err != nil {
		return nil, nil, err
	}
	return nil, names, nil
}

func (r *McpRouter) QueryChart(ctx context.Context, req *mcp.CallToolRequest, in ChartRequest) (*mcp.CallToolResult, *Chart, error) {
	chart, err := r.service.ChartQuery(r.queryContext(ctx, in.Connection), in.Type)
	if err != nil {
		return nil, nil, err
	}
	return nil, chart, nil
}

func (r *McpRouter) QueryCancel(ctx context.Context, req *mcp.CallToolRequest, in KillRequest) (*mcp.CallToolResult, string, error) {
	if err := r.service.CancelQuery(r.queryContext(ctx, in.Connection), in.Pid); err != nil {
		return nil, "", err
	}
	return nil, "cancelled", nil
}

func (r *McpRouter) QueryTerminate(ctx context.Context, req *mcp.CallToolRequest, in KillRequest) (*mcp.CallToolResult, string, error) {
	if err := r.service.TerminateQuery(r.queryContext(ctx, in.Connection), in.Pid); err != nil {
		return nil, "", err
	}
	return nil, "terminated", nil
}
