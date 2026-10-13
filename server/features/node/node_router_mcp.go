package node

import (
	"context"
	"errors"
	"ivory/core/service/job"
	"strings"

	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// McpRouter is this feature's half of the MCP engine, the way Router is its
// half of the http one. Which tools exist and what each one needs is declared
// in core/engine, beside the route table.
//
// A tool's structured result may be any JSON value, not only an object, so a
// handler hands back the service's own type.
type McpRouter struct {
	service *Service
}

func NewMcpRouter(service *Service) *McpRouter {
	return &McpRouter{service: service}
}

// handleKeeperRequest and handlePlatformRequest are the tool twins of the http
// handlers of the same name: what an action produced is sent beside its error,
// never instead of it. A failure is still flagged as one, through IsError.
// They are free functions because Go does not allow a method to introduce type
// parameters.
func handleKeeperRequest[T any](in KeeperOneRequest, action func(KeeperOneRequest) (T, int, error)) (*mcp.CallToolResult, RouterKeeperResponse[T], error) {
	body, status, err := action(in)
	response := RouterKeeperResponse[T]{Status: status, Response: body}
	if err != nil {
		response.Error = err.Error()
		return &mcp.CallToolResult{IsError: true}, response, nil
	}
	return nil, response, nil
}

func handlePlatformRequest[R any, T any](in R, action func(R) (T, error)) (*mcp.CallToolResult, RouterPlatformResponse[T], error) {
	body, err := action(in)
	response := RouterPlatformResponse[T]{Response: body}
	if err != nil {
		response.Error = err.Error()
		return &mcp.CallToolResult{IsError: true}, response, nil
	}
	return nil, response, nil
}

// collectLogs drains a streaming service method into a bounded result, the way
// handleStreamRequest drains it into an SSE response. The services report a
// failure down the stream rather than returning it: a message before the
// stream opens is a request refused outright, and the one following a FAILED
// status is why the job died, so either becomes the tool's error. The request
// context is the backstop for a job that never ends.
func (r *McpRouter) collectLogs(
	ctx context.Context,
	stream func(job.SubscriberID, <-chan struct{}, func(job.Message)),
) ([]string, error) {
	lines := make([]string, 0)
	opened, failed := false, false
	reason := ""
	stream(job.SubscriberID(uuid.New().String()), ctx.Done(), func(event job.Message) {
		switch event.Type {
		case job.STREAM:
			opened = true
		case job.STATUS:
			failed = failed || event.Message == job.FAILED.String()
		case job.SERVER:
			if reason == "" && (!opened || failed) {
				reason = event.Message
			}
		case job.LOG:
			lines = append(lines, event.Message)
		}
	})
	if !opened || failed {
		return nil, errors.New(strings.Join(append([]string{reason}, lines...), "\n"))
	}
	return lines, nil
}

// boundedLogs is what a tool reads: a tail rather than a live stream, because
// a tool call returns once. Zero means the agent said nothing, and an
// unbounded tail of a busy log is not a useful result.
func boundedLogs(request PlatformLogsRequest) PlatformLogsRequest {
	if request.Tail <= 0 {
		request.Tail = 100
	}
	request.Follow = false
	return request
}
