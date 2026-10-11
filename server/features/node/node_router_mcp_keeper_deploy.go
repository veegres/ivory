package node

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (r *McpRouter) KeeperDeploy(ctx context.Context, req *mcp.CallToolRequest, in KeeperDeployRequest) (*mcp.CallToolResult, PlatformResponse, error) {
	output, err := r.service.KeeperDeploy(in)
	if err != nil {
		return nil, nil, err
	}
	return nil, output, nil
}
