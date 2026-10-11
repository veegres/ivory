package deployment

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type McpRouter struct {
	service *Service
}

func NewMcpRouter(service *Service) *McpRouter {
	return &McpRouter{service: service}
}

// TemplateList returns the shipped templates alongside the stored ones, the
// way the template screen lists them: the shipped half is computed from the
// keeper plugins on every request rather than stored.
func (r *McpRouter) TemplateList(ctx context.Context, req *mcp.CallToolRequest, in ListRequest) (*mcp.CallToolResult, []Template, error) {
	stored, err := r.service.List(in)
	if err != nil {
		return nil, nil, err
	}
	return nil, append(r.service.Defaults(in), stored...), nil
}

func (r *McpRouter) TemplateCreate(ctx context.Context, req *mcp.CallToolRequest, in TemplateRequest) (*mcp.CallToolResult, *Template, error) {
	created, err := r.service.Create(in)
	if err != nil {
		return nil, nil, err
	}
	return nil, created, nil
}

func (r *McpRouter) TemplateUpdate(ctx context.Context, req *mcp.CallToolRequest, in RouterUpdateRequest) (*mcp.CallToolResult, *Template, error) {
	id, err := parseMcpTemplateId(in.Id)
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
	id, err := parseMcpTemplateId(in.Id)
	if err != nil {
		return nil, "", err
	}
	if err := r.service.Delete(id); err != nil {
		return nil, "", err
	}
	return nil, "deleted", nil
}

// parseMcpTemplateId is parseTemplateId for a tool argument rather than a url
// parameter: a shipped template's synthetic id says why it is refused instead
// of failing as an unparseable uuid.
func parseMcpTemplateId(param string) (uuid.UUID, error) {
	if isDefaultId(param) {
		return uuid.Nil, ErrTemplateReadOnly
	}
	id, err := uuid.Parse(param)
	if err != nil {
		return uuid.Nil, errors.New("invalid template id")
	}
	return id, nil
}
