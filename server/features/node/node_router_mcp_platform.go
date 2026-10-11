package node

import (
	"context"
	"ivory/core/service/job"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// SYSTEM

func (r *McpRouter) SystemMetrics(ctx context.Context, req *mcp.CallToolRequest, in PlatformMetricsRequest) (*mcp.CallToolResult, *PlatformMetricsResponse, error) {
	metrics, err := r.service.PlatformSystemMetrics(in)
	if err != nil {
		return nil, nil, err
	}
	return nil, metrics, nil
}

func (r *McpRouter) SystemProcesses(ctx context.Context, req *mcp.CallToolRequest, in PlatformProcessesRequest) (*mcp.CallToolResult, PlatformProcessesResponse, error) {
	processes, err := r.service.PlatformSystemProcesses(in)
	if err != nil {
		return nil, nil, err
	}
	return nil, processes, nil
}

func (r *McpRouter) SystemInfo(ctx context.Context, req *mcp.CallToolRequest, in PlatformInfoRequest) (*mcp.CallToolResult, PlatformInfoResponse, error) {
	info, err := r.service.PlatformSystemInfo(in)
	if err != nil {
		return nil, nil, err
	}
	return nil, info, nil
}

func (r *McpRouter) SystemLogs(ctx context.Context, req *mcp.CallToolRequest, in PlatformLogsRequest) (*mcp.CallToolResult, []string, error) {
	request := boundedLogs(in)
	lines, err := r.collectLogs(ctx, func(id job.SubscriberID, close <-chan struct{}, send func(job.Message)) {
		r.service.PlatformSystemLogs(request, id, close, send)
	})
	if err != nil {
		return nil, nil, err
	}
	return nil, lines, nil
}

func (r *McpRouter) SystemCopyId(ctx context.Context, req *mcp.CallToolRequest, in PlatformCopyIdRequest) (*mcp.CallToolResult, string, error) {
	response, err := r.service.PlatformSystemCopyId(in)
	if err != nil {
		return nil, "", err
	}
	return nil, response, nil
}

// CONTAINER

func (r *McpRouter) ContainerList(ctx context.Context, req *mcp.CallToolRequest, in PlatformVaultConnection) (*mcp.CallToolResult, []string, error) {
	names, err := r.service.PlatformContainerList(in)
	if err != nil {
		return nil, nil, err
	}
	return nil, names, nil
}

func (r *McpRouter) ContainerLogs(ctx context.Context, req *mcp.CallToolRequest, in PlatformLogsRequest) (*mcp.CallToolResult, []string, error) {
	request := boundedLogs(in)
	lines, err := r.collectLogs(ctx, func(id job.SubscriberID, close <-chan struct{}, send func(job.Message)) {
		r.service.PlatformContainerLogs(request, id, close, send)
	})
	if err != nil {
		return nil, nil, err
	}
	return nil, lines, nil
}

func (r *McpRouter) ContainerMetrics(ctx context.Context, req *mcp.CallToolRequest, in PlatformActionRequest) (*mcp.CallToolResult, *PlatformMetricsResponse, error) {
	metrics, err := r.service.PlatformContainerMetrics(in)
	if err != nil {
		return nil, nil, err
	}
	return nil, metrics, nil
}

func (r *McpRouter) ContainerStart(ctx context.Context, req *mcp.CallToolRequest, in PlatformActionRequest) (*mcp.CallToolResult, RouterPlatformResponse, error) {
	return r.containerCall(r.service.PlatformContainerStart, in)
}

func (r *McpRouter) ContainerStop(ctx context.Context, req *mcp.CallToolRequest, in PlatformActionRequest) (*mcp.CallToolResult, RouterPlatformResponse, error) {
	return r.containerCall(r.service.PlatformContainerStop, in)
}

func (r *McpRouter) ContainerRestart(ctx context.Context, req *mcp.CallToolRequest, in PlatformActionRequest) (*mcp.CallToolResult, RouterPlatformResponse, error) {
	return r.containerCall(r.service.PlatformContainerRestart, in)
}

func (r *McpRouter) ContainerDown(ctx context.Context, req *mcp.CallToolRequest, in PlatformActionRequest) (*mcp.CallToolResult, RouterPlatformResponse, error) {
	return r.containerCall(r.service.PlatformContainerDown, in)
}

func (r *McpRouter) containerCall(
	action func(PlatformActionRequest) ([]string, error),
	in PlatformActionRequest,
) (*mcp.CallToolResult, RouterPlatformResponse, error) {
	output, err := action(in)
	if err != nil {
		return nil, nil, err
	}
	return nil, output, nil
}

func (r *McpRouter) ContainerUp(ctx context.Context, req *mcp.CallToolRequest, in PlatformUpRequest) (*mcp.CallToolResult, RouterPlatformResponse, error) {
	output, err := r.service.PlatformContainerUp(in)
	if err != nil {
		return nil, nil, err
	}
	return nil, output, nil
}

func (r *McpRouter) ContainerExec(ctx context.Context, req *mcp.CallToolRequest, in PlatformExecRequest) (*mcp.CallToolResult, RouterPlatformResponse, error) {
	output, err := r.service.PlatformContainerExec(in)
	if err != nil {
		return nil, nil, err
	}
	return nil, output, nil
}
