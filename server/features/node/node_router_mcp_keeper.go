package node

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (r *McpRouter) KeeperNodeList(ctx context.Context, req *mcp.CallToolRequest, in KeeperOneRequest) (*mcp.CallToolResult, RouterKeeperResponse[[]KeeperOneResponse], error) {
	return handleKeeperRequest(in, r.service.KeeperNodeList)
}

func (r *McpRouter) KeeperConfigGet(ctx context.Context, req *mcp.CallToolRequest, in KeeperOneRequest) (*mcp.CallToolResult, RouterKeeperResponse[any], error) {
	return handleKeeperRequest(in, r.service.KeeperConfigGet)
}

func (r *McpRouter) KeeperConfigUpdate(ctx context.Context, req *mcp.CallToolRequest, in KeeperOneRequest) (*mcp.CallToolResult, RouterKeeperResponse[any], error) {
	return handleKeeperRequest(in, r.service.KeeperConfigUpdate)
}

func (r *McpRouter) KeeperSwitchover(ctx context.Context, req *mcp.CallToolRequest, in KeeperOneRequest) (*mcp.CallToolResult, RouterKeeperResponse[*string], error) {
	return handleKeeperRequest(in, r.service.KeeperSwitchover)
}

func (r *McpRouter) KeeperSwitchoverCancel(ctx context.Context, req *mcp.CallToolRequest, in KeeperOneRequest) (*mcp.CallToolResult, RouterKeeperResponse[*string], error) {
	return handleKeeperRequest(in, r.service.KeeperSwitchoverDelete)
}

func (r *McpRouter) KeeperReinitialize(ctx context.Context, req *mcp.CallToolRequest, in KeeperOneRequest) (*mcp.CallToolResult, RouterKeeperResponse[*string], error) {
	return handleKeeperRequest(in, r.service.KeeperReinitialize)
}

func (r *McpRouter) KeeperRestart(ctx context.Context, req *mcp.CallToolRequest, in KeeperOneRequest) (*mcp.CallToolResult, RouterKeeperResponse[*string], error) {
	return handleKeeperRequest(in, r.service.KeeperRestart)
}

func (r *McpRouter) KeeperRestartCancel(ctx context.Context, req *mcp.CallToolRequest, in KeeperOneRequest) (*mcp.CallToolResult, RouterKeeperResponse[*string], error) {
	return handleKeeperRequest(in, r.service.KeeperRestartDelete)
}

func (r *McpRouter) KeeperReload(ctx context.Context, req *mcp.CallToolRequest, in KeeperOneRequest) (*mcp.CallToolResult, RouterKeeperResponse[*string], error) {
	return handleKeeperRequest(in, r.service.KeeperReload)
}

func (r *McpRouter) KeeperFailover(ctx context.Context, req *mcp.CallToolRequest, in KeeperOneRequest) (*mcp.CallToolResult, RouterKeeperResponse[*string], error) {
	return handleKeeperRequest(in, r.service.KeeperFailover)
}

func (r *McpRouter) KeeperActivate(ctx context.Context, req *mcp.CallToolRequest, in KeeperOneRequest) (*mcp.CallToolResult, RouterKeeperResponse[*string], error) {
	return handleKeeperRequest(in, r.service.KeeperActivate)
}

func (r *McpRouter) KeeperPause(ctx context.Context, req *mcp.CallToolRequest, in KeeperOneRequest) (*mcp.CallToolResult, RouterKeeperResponse[*string], error) {
	return handleKeeperRequest(in, r.service.KeeperPause)
}
