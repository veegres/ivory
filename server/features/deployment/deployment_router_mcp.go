package deployment

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type McpRouter struct {
	service *Service
}

func NewMcpRouter(service *Service) *McpRouter {
	return &McpRouter{service: service}
}

func (r *McpRouter) TemplateList(ctx context.Context, req *mcp.CallToolRequest, in ListRequest) (*mcp.CallToolResult, []Template, error) {
	list, err := r.service.List(in)
	if err != nil {
		return nil, nil, err
	}
	return nil, list, nil
}

func (r *McpRouter) TemplateCreate(ctx context.Context, req *mcp.CallToolRequest, in TemplateRequest) (*mcp.CallToolResult, *Template, error) {
	created, err := r.service.Create(in)
	if err != nil {
		return nil, nil, err
	}
	return nil, created, nil
}

func (r *McpRouter) TemplateUpdate(ctx context.Context, req *mcp.CallToolRequest, in RouterUpdateRequest) (*mcp.CallToolResult, *Template, error) {
	id, err := parseTemplateId(in.Id)
	if err != nil {
		return nil, nil, err
	}
	updated, err := r.service.Update(id, in.Template)
	if err != nil {
		return nil, nil, err
	}
	return nil, updated, nil
}

func (r *McpRouter) TemplateDelete(ctx context.Context, req *mcp.CallToolRequest, in RouterIdRequest) (*mcp.CallToolResult, string, error) {
	id, err := parseTemplateId(in.Id)
	if err != nil {
		return nil, "", err
	}
	if err := r.service.Delete(id); err != nil {
		return nil, "", err
	}
	return nil, "deleted", nil
}
