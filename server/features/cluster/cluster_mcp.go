package cluster

import (
	"context"
	"ivory/features/node"

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

// McpClusterNameRequest and McpDeployRequest carry what the http route takes
// from the url path. They exist because a tool's arguments are one typed
// object and there is no path to read a name out of - not to restate a model.
type McpClusterNameRequest struct {
	Name string `json:"name" binding:"required" jsonschema:"name of the cluster as registered in Ivory"`
}

// McpDeployRequest is DeployRequest with CommonConfig's inline credentials left
// out: a deploy names vault entries instead, so no password passes through the
// agent's context. Deploy refuses a vault id beside an inline pair anyway.
type McpDeployRequest struct {
	Cluster  string              `json:"cluster" binding:"required" jsonschema:"name to register the new cluster under"`
	Dcs      string              `json:"dcs,omitempty" jsonschema:"address of the coordination store the whole deployment registers with, substituted for {{dcs}}"`
	Parallel bool                `json:"parallel,omitempty" jsonschema:"deploy every node at once instead of in order"`
	Platform node.PlatformPlugin `json:"platform" binding:"required" jsonschema:"deployment target, e.g. docker"`
	Nodes    []DeployNode        `json:"nodes" binding:"required" jsonschema:"one entry per node: its name, host, ssh/keeper/database ports, the command that deploys it and that command's post-scripts"`
	Options  Options             `json:"options" jsonschema:"keeper and database plugins, TLS flags, certificates, tags, and the vault ids holding the ssh key and the keeper/database credentials"`
}

func (in McpDeployRequest) toDeployRequest() DeployRequest {
	return DeployRequest{
		Parallel:       in.Parallel,
		Nodes:          in.Nodes,
		Platform:       in.Platform,
		CommonConfig:   CommonConfig{Cluster: in.Cluster, Dcs: in.Dcs},
		ClusterOptions: in.Options,
	}
}

// McpDeploy is its own type because Deploy answers two things at once: what
// happened, and whether all of it worked. Complete is every node up and every
// post-script run; a false here with no error is a partial deploy, and the
// cluster is registered with all of its nodes either way.
type McpDeploy struct {
	Logs     []string `json:"logs"`
	Complete bool     `json:"complete"`
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

func (r *McpRouter) ClusterGet(ctx context.Context, req *mcp.CallToolRequest, in McpClusterNameRequest) (*mcp.CallToolResult, *Response, error) {
	c, err := r.clusterService.Get(in.Name)
	if err != nil {
		return nil, nil, err
	}
	return nil, &c, nil
}

func (r *McpRouter) ClusterOverview(ctx context.Context, req *mcp.CallToolRequest, in McpClusterNameRequest) (*mcp.CallToolResult, *Overview, error) {
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

func (r *McpRouter) ClusterDelete(ctx context.Context, req *mcp.CallToolRequest, in McpClusterNameRequest) (*mcp.CallToolResult, string, error) {
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

func (r *McpRouter) ClusterFix(ctx context.Context, req *mcp.CallToolRequest, in McpClusterNameRequest) (*mcp.CallToolResult, *Response, error) {
	fixed, err := r.clusterService.Fix(in.Name)
	if err != nil {
		return nil, nil, err
	}
	return nil, fixed, nil
}

func (r *McpRouter) ClusterDeploy(ctx context.Context, req *mcp.CallToolRequest, in McpDeployRequest) (*mcp.CallToolResult, McpDeploy, error) {
	logs, complete, err := r.clusterService.Deploy(in.toDeployRequest())
	if err != nil {
		return nil, McpDeploy{}, err
	}
	return nil, McpDeploy{Logs: logs, Complete: complete}, nil
}
