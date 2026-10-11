package cluster

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// McpRouter is this feature's half of the MCP engine, the way Router is its
// half of the http one. Which tools exist and what each one needs is declared
// in core/engine, beside the route table.
//
// A tool's structured result may be any JSON value, not only an object, so a
// handler hands back the service's own type rather than a transport shape.
type McpRouter struct {
	clusterService *Service
}

func NewMcpRouter(clusterService *Service) *McpRouter {
	return &McpRouter{clusterService: clusterService}
}

func (r *McpRouter) ClusterList(ctx context.Context, req *mcp.CallToolRequest, in SearchRequest) (*mcp.CallToolResult, []Response, error) {
	if len(in.Tags) == 0 && in.Keeper == nil && in.Database == nil {
		list, err := r.clusterService.List()
		if err != nil {
			return nil, nil, err
		}
		return nil, list, nil
	}
	list, err := r.clusterService.Search(in)
	if err != nil {
		return nil, nil, err
	}
	return nil, list, nil
}

func (r *McpRouter) ClusterGet(ctx context.Context, req *mcp.CallToolRequest, in RouterNameRequest) (*mcp.CallToolResult, *Response, error) {
	c, err := r.clusterService.Get(in.Name)
	if err != nil {
		return nil, nil, err
	}
	return nil, &c, nil
}

func (r *McpRouter) ClusterOverview(ctx context.Context, req *mcp.CallToolRequest, in RouterNameRequest) (*mcp.CallToolResult, *Overview, error) {
	overview, err := r.clusterService.Overview(in.Name, "", 0)
	if err != nil {
		return nil, nil, err
	}
	return nil, overview, nil
}

func (r *McpRouter) ClusterCreate(ctx context.Context, req *mcp.CallToolRequest, in Request) (*mcp.CallToolResult, *Response, error) {
	created, err := r.clusterService.Create(in)
	if err != nil {
		return nil, nil, err
	}
	return nil, created, nil
}

func (r *McpRouter) ClusterUpdate(ctx context.Context, req *mcp.CallToolRequest, in Request) (*mcp.CallToolResult, *Response, error) {
	updated, err := r.clusterService.Update(in)
	if err != nil {
		return nil, nil, err
	}
	return nil, updated, nil
}

func (r *McpRouter) ClusterDelete(ctx context.Context, req *mcp.CallToolRequest, in RouterNameRequest) (*mcp.CallToolResult, string, error) {
	if err := r.clusterService.Delete(in.Name); err != nil {
		return nil, "", err
	}
	return nil, "deleted", nil
}

func (r *McpRouter) ClusterDetect(ctx context.Context, req *mcp.CallToolRequest, in CreateAutoRequest) (*mcp.CallToolResult, *Response, error) {
	detected, err := r.clusterService.Detect(in)
	if err != nil {
		return nil, nil, err
	}
	return nil, detected, nil
}

func (r *McpRouter) ClusterFix(ctx context.Context, req *mcp.CallToolRequest, in RouterNameRequest) (*mcp.CallToolResult, *Response, error) {
	fixed, err := r.clusterService.Fix(in.Name)
	if err != nil {
		return nil, nil, err
	}
	return nil, fixed, nil
}

func (r *McpRouter) ClusterDeploy(ctx context.Context, req *mcp.CallToolRequest, in RouterDeployRequest) (*mcp.CallToolResult, RouterDeployResponse, error) {
	logs, complete, err := r.clusterService.Deploy(in.toDeployRequest())
	if err != nil {
		return nil, RouterDeployResponse{}, err
	}
	return nil, RouterDeployResponse{Logs: logs, Complete: complete}, nil
}
