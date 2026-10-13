package node

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (r *McpRouter) KeeperDeploy(ctx context.Context, req *mcp.CallToolRequest, in KeeperDeployRequest) (*mcp.CallToolResult, RouterKeeperDeployResponse, error) {
	logs, complete, err := r.service.KeeperDeploy(in)
	if err != nil {
		return nil, RouterKeeperDeployResponse{}, err
	}
	return nil, RouterKeeperDeployResponse{Logs: logs, Complete: complete}, nil
}
