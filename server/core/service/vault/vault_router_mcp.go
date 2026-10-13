package vault

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// McpRouter is this service's half of the MCP engine, the way HttpRouter is
// its half of the http one. It only reads: a secret never passes through it.
type McpRouter struct {
	vaultService *Service
}

func NewMcpRouter(vaultService *Service) *McpRouter {
	return &McpRouter{vaultService: vaultService}
}

func (r *McpRouter) VaultList(ctx context.Context, req *mcp.CallToolRequest, in RouterListRequest) (*mcp.CallToolResult, VaultMap, error) {
	vaults, err := r.vaultService.Map(in.Type)
	if err != nil {
		return nil, nil, err
	}
	return nil, vaults, nil
}
