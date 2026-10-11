package node

import (
	"context"
	"ivory/core/service/job"

	"github.com/google/uuid"
)

// McpRouter is this feature's half of the MCP engine, the way Router is its
// half of the http one. Which tools exist and what each one needs is declared
// in core/engine, beside the route table.
//
// A tool's structured result may be any JSON value, not only an object, so a
// handler hands back the service's own type. The one exception is McpLogs,
// which carries something the service does not return.
type McpRouter struct {
	service *Service
}

func NewMcpRouter(service *Service) *McpRouter {
	return &McpRouter{service: service}
}

// McpLogs is a bounded log read. It is its own type because the log services
// report a failure by pushing it down the stream rather than returning an
// error, so the lines and that commentary arrive together and both have to
// reach the caller.
type McpLogs struct {
	Lines []string `json:"lines"`
	Notes []string `json:"notes,omitempty"`
}

// collectLogs drains a streaming service method into a bounded result, the way
// handleStreamRequest drains it into an SSE response. The request context is
// the backstop: a job that never ends is cut off when the client goes away
// rather than holding the tool open forever.
func (r *McpRouter) collectLogs(
	ctx context.Context,
	stream func(job.SubscriberID, <-chan struct{}, func(job.Message)),
) McpLogs {
	lines := make([]string, 0)
	notes := make([]string, 0)
	stream(job.SubscriberID(uuid.New().String()), ctx.Done(), func(event job.Message) {
		switch event.Type {
		case job.LOG:
			lines = append(lines, event.Message)
		case job.SERVER:
			notes = append(notes, event.Message)
		}
	})
	return McpLogs{Lines: lines, Notes: notes}
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
