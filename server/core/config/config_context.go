package config

import "context"

// AuthContextKey are the gin.Context keys the auth middlewares populate and other features read.
var AuthContextKey = struct {
	Enabled    string
	Authorised string
	Username   string
	Type       string
	Superuser  string
	Error      string
	Session    string
}{
	Enabled:    "auth",
	Authorised: "authorised",
	Username:   "username",
	Type:       "authType",
	Superuser:  "superuser",
	Error:      "authError",
	Session:    "session",
}

// PermissionContextKey is the gin.Context key permission.ValidateMiddleware
// populates with the caller's PermissionMap. It is read by the per-route
// permission check and by the MCP router, which gates tools per call rather
// than per route.
const PermissionContextKey = "permissions"

// McpCaller is who an MCP request runs as, for the tools that have to describe
// the caller back to them. It carries no permissions: which tools a caller is
// offered is decided by the MCP engine, not by a feature.
type McpCaller struct {
	Username    string
	Authorised  bool
	AuthEnabled bool
}

type mcpCallerContextKey struct{}

// WithMcpCaller marks a request's context with the caller the MCP engine
// authenticated. A tool handler has no gin.Context to read it from, so this is
// how the one identity a request belongs to reaches it.
func WithMcpCaller(ctx context.Context, caller McpCaller) context.Context {
	return context.WithValue(ctx, mcpCallerContextKey{}, caller)
}

// McpCallerFrom is the caller this request runs as, zero-valued when the
// engine set none.
func McpCallerFrom(ctx context.Context) McpCaller {
	if caller, ok := ctx.Value(mcpCallerContextKey{}).(McpCaller); ok {
		return caller
	}
	return McpCaller{}
}
