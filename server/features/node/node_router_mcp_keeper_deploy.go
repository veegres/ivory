package node

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (r *McpRouter) KeeperDeploy(ctx context.Context, req *mcp.CallToolRequest, in KeeperDeployRequest) (*mcp.CallToolResult, RouterPlatformResponse[KeeperDeployResponse], error) {
	return handlePlatformRequest(in, r.service.KeeperDeploy)
}
