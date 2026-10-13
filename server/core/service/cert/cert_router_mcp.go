package cert

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// McpRouter is this service's half of the MCP engine, the way HttpRouter is
// its half of the http one.
type McpRouter struct {
	certService *Service
}

func NewMcpRouter(certService *Service) *McpRouter {
	return &McpRouter{certService: certService}
}

func (r *McpRouter) CertList(ctx context.Context, req *mcp.CallToolRequest, in RouterListRequest) (*mcp.CallToolResult, CertMap, error) {
	if in.Type != nil {
		list, err := r.certService.ListByType(*in.Type)
		if err != nil {
			return nil, nil, err
		}
		return nil, list, nil
	}
	list, err := r.certService.List()
	if err != nil {
		return nil, nil, err
	}
	return nil, list, nil
}
