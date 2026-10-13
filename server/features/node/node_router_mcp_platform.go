package node

import (
	"context"
	"ivory/core/service/job"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// SYSTEM

func (r *McpRouter) SystemMetrics(ctx context.Context, req *mcp.CallToolRequest, in PlatformMetricsRequest) (*mcp.CallToolResult, RouterPlatformResponse[*PlatformMetricsResponse], error) {
	return handlePlatformRequest(in, r.service.PlatformSystemMetrics)
}

func (r *McpRouter) SystemProcesses(ctx context.Context, req *mcp.CallToolRequest, in PlatformProcessesRequest) (*mcp.CallToolResult, RouterPlatformResponse[PlatformProcessesResponse], error) {
	return handlePlatformRequest(in, r.service.PlatformSystemProcesses)
}

func (r *McpRouter) SystemInfo(ctx context.Context, req *mcp.CallToolRequest, in PlatformInfoRequest) (*mcp.CallToolResult, RouterPlatformResponse[PlatformInfoResponse], error) {
	return handlePlatformRequest(in, r.service.PlatformSystemInfo)
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

func (r *McpRouter) SystemCopyId(ctx context.Context, req *mcp.CallToolRequest, in PlatformCopyIdRequest) (*mcp.CallToolResult, RouterPlatformResponse[string], error) {
	return handlePlatformRequest(in, r.service.PlatformSystemCopyId)
}

// CONTAINER

func (r *McpRouter) ContainerList(ctx context.Context, req *mcp.CallToolRequest, in PlatformVaultConnection) (*mcp.CallToolResult, RouterPlatformResponse[[]string], error) {
	return handlePlatformRequest(in, r.service.PlatformContainerList)
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

func (r *McpRouter) ContainerMetrics(ctx context.Context, req *mcp.CallToolRequest, in PlatformActionRequest) (*mcp.CallToolResult, RouterPlatformResponse[*PlatformMetricsResponse], error) {
	return handlePlatformRequest(in, r.service.PlatformContainerMetrics)
}

func (r *McpRouter) ContainerStart(ctx context.Context, req *mcp.CallToolRequest, in PlatformActionRequest) (*mcp.CallToolResult, RouterPlatformResponse[[]string], error) {
	return handlePlatformRequest(in, r.service.PlatformContainerStart)
}

func (r *McpRouter) ContainerStop(ctx context.Context, req *mcp.CallToolRequest, in PlatformActionRequest) (*mcp.CallToolResult, RouterPlatformResponse[[]string], error) {
	return handlePlatformRequest(in, r.service.PlatformContainerStop)
}

func (r *McpRouter) ContainerRestart(ctx context.Context, req *mcp.CallToolRequest, in PlatformActionRequest) (*mcp.CallToolResult, RouterPlatformResponse[[]string], error) {
	return handlePlatformRequest(in, r.service.PlatformContainerRestart)
}

func (r *McpRouter) ContainerDown(ctx context.Context, req *mcp.CallToolRequest, in PlatformActionRequest) (*mcp.CallToolResult, RouterPlatformResponse[[]string], error) {
	return handlePlatformRequest(in, r.service.PlatformContainerDown)
}

func (r *McpRouter) ContainerUp(ctx context.Context, req *mcp.CallToolRequest, in PlatformUpRequest) (*mcp.CallToolResult, RouterPlatformResponse[[]string], error) {
	return handlePlatformRequest(in, r.service.PlatformContainerUp)
}

func (r *McpRouter) ContainerExec(ctx context.Context, req *mcp.CallToolRequest, in PlatformExecRequest) (*mcp.CallToolResult, RouterPlatformResponse[[]string], error) {
	return handlePlatformRequest(in, r.service.PlatformContainerExec)
}
