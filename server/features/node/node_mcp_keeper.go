package node

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (r *McpRouter) KeeperNodeList(ctx context.Context, req *mcp.CallToolRequest, in KeeperOneRequest) (*mcp.CallToolResult, []KeeperOneResponse, error) {
	list, _, err := r.service.KeeperNodeList(in)
	if err != nil {
		return nil, nil, err
	}
	return nil, list, nil
}

func (r *McpRouter) KeeperConfigGet(ctx context.Context, req *mcp.CallToolRequest, in KeeperOneRequest) (*mcp.CallToolResult, any, error) {
	config, _, err := r.service.KeeperConfigGet(in)
	if err != nil {
		return nil, nil, err
	}
	return nil, config, nil
}

func (r *McpRouter) KeeperConfigUpdate(ctx context.Context, req *mcp.CallToolRequest, in KeeperOneRequest) (*mcp.CallToolResult, any, error) {
	config, _, err := r.service.KeeperConfigUpdate(in)
	if err != nil {
		return nil, nil, err
	}
	return nil, config, nil
}

func (r *McpRouter) KeeperSwitchover(ctx context.Context, req *mcp.CallToolRequest, in KeeperOneRequest) (*mcp.CallToolResult, *string, error) {
	return r.keeperCall(in, r.service.KeeperSwitchover)
}

func (r *McpRouter) KeeperSwitchoverCancel(ctx context.Context, req *mcp.CallToolRequest, in KeeperOneRequest) (*mcp.CallToolResult, *string, error) {
	return r.keeperCall(in, r.service.KeeperSwitchoverDelete)
}

func (r *McpRouter) KeeperReinitialize(ctx context.Context, req *mcp.CallToolRequest, in KeeperOneRequest) (*mcp.CallToolResult, *string, error) {
	return r.keeperCall(in, r.service.KeeperReinitialize)
}

func (r *McpRouter) KeeperRestart(ctx context.Context, req *mcp.CallToolRequest, in KeeperOneRequest) (*mcp.CallToolResult, *string, error) {
	return r.keeperCall(in, r.service.KeeperRestart)
}

func (r *McpRouter) KeeperRestartCancel(ctx context.Context, req *mcp.CallToolRequest, in KeeperOneRequest) (*mcp.CallToolResult, *string, error) {
	return r.keeperCall(in, r.service.KeeperRestartDelete)
}

func (r *McpRouter) KeeperReload(ctx context.Context, req *mcp.CallToolRequest, in KeeperOneRequest) (*mcp.CallToolResult, *string, error) {
	return r.keeperCall(in, r.service.KeeperReload)
}

func (r *McpRouter) KeeperFailover(ctx context.Context, req *mcp.CallToolRequest, in KeeperOneRequest) (*mcp.CallToolResult, *string, error) {
	return r.keeperCall(in, r.service.KeeperFailover)
}

func (r *McpRouter) KeeperActivate(ctx context.Context, req *mcp.CallToolRequest, in KeeperOneRequest) (*mcp.CallToolResult, *string, error) {
	return r.keeperCall(in, r.service.KeeperActivate)
}

func (r *McpRouter) KeeperPause(ctx context.Context, req *mcp.CallToolRequest, in KeeperOneRequest) (*mcp.CallToolResult, *string, error) {
	return r.keeperCall(in, r.service.KeeperPause)
}

// keeperCall is the shape every keeper operation shares: run it, and report
// what the keeper itself said - which for some operations is nothing.
func (r *McpRouter) keeperCall(
	request KeeperOneRequest,
	action func(KeeperOneRequest) (*string, int, error),
) (*mcp.CallToolResult, *string, error) {
	response, _, err := action(request)
	if err != nil {
		return nil, nil, err
	}
	return nil, response, nil
}
