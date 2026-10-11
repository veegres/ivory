package query

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"ivory/core/config"

	"github.com/google/uuid"
)

// ErrQueryUuidRequired is the rule PostExecuteTemplateQuery states inline: a
// template query names the saved query it runs.
var ErrQueryUuidRequired = errors.New("queryUuid is required")

// McpRouter is this feature's half of the MCP engine, the way Router is its
// half of the http one. Which tools exist and what each one needs is declared
// in core/engine, beside the route table.
//
// A tool's structured result may be any JSON value, not only an object, so a
// handler hands back the service's own type rather than a transport shape.
type McpRouter struct {
	service *Service
}

func NewMcpRouter(service *Service) *McpRouter {
	return &McpRouter{service: service}
}

// queryContext is withQueryContext for the MCP engine. The session only ever
// becomes the connection's application_name, which keeps seven characters of
// it, so it is a hash of the caller's name: a prefix of the name itself would
// merge everyone whose names start alike into one activity view.
func (r *McpRouter) queryContext(ctx context.Context, connection Connection) Context {
	session := "mcp"
	if username := config.McpCallerFrom(ctx).Username; username != "" {
		sum := sha256.Sum256([]byte(username))
		session = hex.EncodeToString(sum[:])
	}
	return Context{Connection: connection, Session: session}
}

func (r *McpRouter) queryId(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}
