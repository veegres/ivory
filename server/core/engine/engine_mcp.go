package engine

import (
	"context"
	"fmt"
	coreConfig "ivory/core/config"
	"ivory/core/utils"
	"ivory/features"
	"ivory/features/permission"
	"ivory/plugins/database"
	"ivory/plugins/keeper"
	"ivory/plugins/platform"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/google/uuid"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// engine_mcp.go is the Model Context Protocol engine, the sibling of
// engine_http.go. It does for tools exactly what that file does for routes:
// declares each one next to the permission that gates it and the feature
// router that serves it. Every feature keeps its own MCP implementation in its
// own package, in <feature>_router_mcp.go, beside the router it already has.
//
// The one structural difference is when the table runs. A route is registered
// once at startup; a tool list belongs to one caller, so these functions run
// per request and drop the tools that caller may not call.

// maxMcpRequestBodyBytes bounds one JSON-RPC request. Tool arguments are
// command text and query text, not uploads.
const maxMcpRequestBodyBytes = 4 << 20

// mcpServerName is how Ivory introduces itself to an MCP client.
const mcpServerName = "ivory"

const mcpInstructions = `Ivory manages high-availability database clusters.

Start with ivory_cluster_list to see the registered clusters, then
ivory_cluster_overview for one cluster's health: every node's role, state,
replication lag, pending restarts and the warnings Ivory derived from
comparing the keeper's own answer against the stored configuration. Those
warnings are the fastest route to a diagnosis.

Tools name a cluster and a node by the names in its configuration, never by
host and port. Keeper member keys are different: switchover and failover
identify the leader and the candidate by the keeper's own member key, which
the overview reports as each node's keeper.key.

Only the tools this user holds permission for are listed.`

// NewMcpServer builds the MCP engine, the way NewHttpServer builds the http
// one. It returns the handler engine_http.go mounts, since MCP rides the same
// listener rather than opening one of its own.
func NewMcpServer(env *coreConfig.Environment, mr *features.McpRouter) gin.HandlerFunc {
	handler := mcpsdk.NewStreamableHTTPHandler(
		func(request *http.Request) *mcpsdk.Server {
			return newMcpServer(request.Context(), mr, env.Version.Label)
		},
		&mcpsdk.StreamableHTTPOptions{
			// Stateless rebuilds the server per request, which is what keeps
			// the tool list honest: a stateful session would pin it to the
			// permissions the caller held when they connected.
			Stateless:                    true,
			JSONResponse:                 true,
			MaxRequestBodyBytes:          maxMcpRequestBodyBytes,
			PropagateRequestCancellation: true,
		},
	)

	// The session cookie is a valid credential here, so a browser page could
	// otherwise drive tool calls cross-origin. Non-browser clients send no
	// Origin and are unaffected. This wraps the handler rather than setting the
	// SDK's own option, which is deprecated and ignores a deny handler.
	protected := http.NewCrossOriginProtection().Handler(handler)

	return func(context *gin.Context) {
		request := context.Request.WithContext(mcpRequestContext(context))
		protected.ServeHTTP(context.Writer, request)
	}
}

// newMcpServer is the per-caller equivalent of the route table: it declares
// every tool group, and each group drops the tools this caller may not call.
func newMcpServer(ctx context.Context, mr *features.McpRouter, version string) *mcpsdk.Server {
	utils.RegisterValidators()
	r := &mcpRegistrar{
		server: mcpsdk.NewServer(
			&mcpsdk.Implementation{Name: mcpServerName, Version: version},
			&mcpsdk.ServerOptions{Instructions: mcpInstructions, SchemaCache: mcpResolvedSchemas},
		),
		permissions: mcpPermissionsFrom(ctx),
	}

	managementTools(r, mr)
	clusterTools(r, mr)
	nodeTools(r, mr)
	tagTools(r, mr)
	queryTools(r, mr)
	deploymentTools(r, mr)

	return r.server
}

// managementTools holds the one tool behind no permission, mirroring the
// unguarded /info route.
func managementTools(r *mcpRegistrar, mr *features.McpRouter) {
	addToolUngated(r, &mcpsdk.Tool{
		Name:        "ivory_app_info",
		Description: "Ivory's own version, configuration state and which authentication methods it accepts.",
		Annotations: readOnly,
	}, mr.Management.AppInfo)
}

func clusterTools(r *mcpRegistrar, mr *features.McpRouter) {
	addTool(r, coreConfig.ViewClusterList, &mcpsdk.Tool{
		Name:        "ivory_cluster_list",
		Description: "List the clusters registered in Ivory, with each one's nodes, keeper and database plugins, TLS flags and tags. Start here.",
		Annotations: readOnly,
	}, mr.Cluster.ClusterList)

	addTool(r, coreConfig.ViewClusterItem, &mcpsdk.Tool{
		Name:        "ivory_cluster_get",
		Description: "One cluster's stored configuration: every node's name, host and ssh/keeper/database ports. Node names from here are what the other tools take.",
		Annotations: readOnly,
	}, mr.Cluster.ClusterGet)

	addTool(r, coreConfig.ViewClusterOverview, &mcpsdk.Tool{
		Name: "ivory_cluster_overview",
		Description: "Live health of every node in a cluster: role, state, replication lag, sync status, pending restarts, " +
			"keeper member key, plus the warnings Ivory derives by comparing what the keeper reports against the stored " +
			"configuration (a missing leader, two leaders, a node the keeper does not know, a node belonging to a different " +
			"cluster). Read this before diagnosing anything.",
		Annotations: readOnly,
	}, mr.Cluster.ClusterOverview)

	addTool(r, coreConfig.ManageClusterCreate, &mcpsdk.Tool{
		Name:        "ivory_cluster_create",
		Description: "Register a cluster from a full node list. Registers it in Ivory only; it deploys nothing.",
		Annotations: mutating,
	}, mr.Cluster.ClusterCreate)

	addTool(r, coreConfig.ManageClusterUpdate, &mcpsdk.Tool{
		Name:        "ivory_cluster_update",
		Description: "Overwrite a registered cluster's configuration. Every node must state a name unique within the cluster.",
		Annotations: mutating,
	}, mr.Cluster.ClusterUpdate)

	addTool(r, coreConfig.ManageClusterDelete, &mcpsdk.Tool{
		Name:        "ivory_cluster_delete",
		Description: "Remove a cluster from Ivory. Deletes the registration only - no container is stopped and no data is touched.",
		Annotations: destructive,
	}, mr.Cluster.ClusterDelete)

	addTool(r, coreConfig.ManageClusterCreate, &mcpsdk.Tool{
		Name:        "ivory_cluster_detect",
		Description: "Ask one running keeper endpoint which members it knows, then register the whole cluster from that answer.",
		Annotations: mutating,
	}, mr.Cluster.ClusterDetect)

	addTool(r, coreConfig.ManageClusterUpdate, &mcpsdk.Tool{
		Name:        "ivory_cluster_fix",
		Description: "Replace a registered cluster's node list with the members its keeper currently reports. Use when the overview warns that the configuration and the keeper disagree.",
		Annotations: mutating,
	}, mr.Cluster.ClusterFix)

	addTool(r, coreConfig.ManageClusterCreate, &mcpsdk.Tool{
		Name: "ivory_cluster_deploy",
		Description: "Deploy a whole cluster: run each node's command on its host over ssh, then its post-scripts. The cluster " +
			"record is written before anything starts and is kept whether or not every node came up, so a failed container stays " +
			"on its host with its logs intact - that is deliberate, it is the only evidence of what went wrong. Credentials are " +
			"named as vault ids, never typed here. Returns the logs and whether every node and post-script succeeded.",
		Annotations: destructive,
	}, mr.Cluster.ClusterDeploy)
}

// nodeTools covers all three halves of a node, the way nodeRouter does:
// the keeper managing it, the host it runs on, and the containers on that host.
func nodeTools(r *mcpRegistrar, mr *features.McpRouter) {
	// keeper
	addTool(r, coreConfig.ViewNodeKeeperOverview, &mcpsdk.Tool{
		Name:        "ivory_keeper_node_list",
		Description: "Ask one node's keeper what it knows about the cluster's members. ivory_cluster_overview is usually the better read: it unions every node's answer and adds Ivory's own warnings.",
		Annotations: readOnly,
	}, mr.Node.KeeperNodeList)

	addTool(r, coreConfig.ViewNodeKeeperConfig, &mcpsdk.Tool{
		Name:        "ivory_keeper_config_get",
		Description: "The keeper's current configuration as that keeper reports it.",
		Annotations: readOnly,
	}, mr.Node.KeeperConfigGet)

	addTool(r, coreConfig.ManageNodeKeeperConfigUpdate, &mcpsdk.Tool{
		Name:        "ivory_keeper_config_update",
		Description: "Patch the keeper's configuration. The patch shape is the keeper's own; read ivory_keeper_config_get first. Some settings take effect only after a restart.",
		Annotations: mutating,
	}, mr.Node.KeeperConfigUpdate)

	addTool(r, coreConfig.ManageNodeKeeperSwitchover, &mcpsdk.Tool{
		Name: "ivory_keeper_switchover",
		Description: "Hand leadership to another member. Causes a brief write outage. Send this to the current leader and name " +
			"leader and candidate by their keeper member keys (the overview's keeper.key), not by host. Omitting the candidate " +
			"lets the keeper choose.",
		Annotations: destructive,
	}, mr.Node.KeeperSwitchover)

	addTool(r, coreConfig.ManageNodeKeeperSwitchover, &mcpsdk.Tool{
		Name:        "ivory_keeper_switchover_cancel",
		Description: "Cancel a scheduled switchover.",
		Annotations: mutating,
	}, mr.Node.KeeperSwitchoverCancel)

	addTool(r, coreConfig.ManageNodeKeeperReinitialize, &mcpsdk.Tool{
		Name:        "ivory_keeper_reinitialize",
		Description: "DESTRUCTIVE: wipe this member's data directory and re-clone it from the leader. The node is unavailable for the whole rebuild, which on a large database is slow. Only ever use on a replica.",
		Annotations: destructive,
	}, mr.Node.KeeperReinitialize)

	addTool(r, coreConfig.ManageNodeKeeperRestart, &mcpsdk.Tool{
		Name:        "ivory_keeper_restart",
		Description: "Restart the database process this keeper manages. Drops every connection to that node. Pass restartPending to act only on a node that reports one.",
		Annotations: destructive,
	}, mr.Node.KeeperRestart)

	addTool(r, coreConfig.ManageNodeKeeperRestart, &mcpsdk.Tool{
		Name:        "ivory_keeper_restart_cancel",
		Description: "Cancel a scheduled restart.",
		Annotations: mutating,
	}, mr.Node.KeeperRestartCancel)

	addTool(r, coreConfig.ManageNodeKeeperReload, &mcpsdk.Tool{
		Name:        "ivory_keeper_reload",
		Description: "Ask the database process to re-read its configuration, without a restart.",
		Annotations: mutating,
	}, mr.Node.KeeperReload)

	addTool(r, coreConfig.ManageNodeKeeperFailover, &mcpsdk.Tool{
		Name: "ivory_keeper_failover",
		Description: "DESTRUCTIVE: force an immediate leader change, bypassing the checks a switchover makes. CAN LOSE " +
			"COMMITTED DATA if the candidate is behind. Use only when there is no healthy leader to switch over from. Name the " +
			"candidate by its keeper member key.",
		Annotations: destructive,
	}, mr.Node.KeeperFailover)

	addTool(r, coreConfig.ManageNodeKeeperActivation, &mcpsdk.Tool{
		Name:        "ivory_keeper_activate",
		Description: "Resume automatic failover after a pause.",
		Annotations: mutating,
	}, mr.Node.KeeperActivate)

	addTool(r, coreConfig.ManageNodeKeeperActivation, &mcpsdk.Tool{
		Name:        "ivory_keeper_pause",
		Description: "Suspend automatic failover. The cluster stays up but stops protecting itself, and nothing reminds you later - always pair this with ivory_keeper_activate.",
		Annotations: destructive,
	}, mr.Node.KeeperPause)

	// system - the node itself
	addTool(r, coreConfig.ViewNodeSystem, &mcpsdk.Tool{
		Name:        "ivory_system_metrics",
		Description: "CPU ticks, memory bytes and network bytes for the node's host. Counters, not percentages: call twice and take the difference.",
		Annotations: readOnly,
	}, mr.Node.SystemMetrics)

	addTool(r, coreConfig.ViewNodeSystem, &mcpsdk.Tool{
		Name:        "ivory_system_processes",
		Description: "Processes on the node's host, with cpu and memory use.",
		Annotations: readOnly,
	}, mr.Node.SystemProcesses)

	addTool(r, coreConfig.ViewNodeSystem, &mcpsdk.Tool{
		Name:        "ivory_system_info",
		Description: "Host facts for the node: kernel, distribution, cpu count, memory and disk.",
		Annotations: readOnly,
	}, mr.Node.SystemInfo)

	addTool(r, coreConfig.ViewNodeSystem, &mcpsdk.Tool{
		Name:        "ivory_system_logs",
		Description: "Read the tail of a file on the node's host over ssh.",
		Annotations: readOnly,
	}, mr.Node.SystemLogs)

	addTool(r, coreConfig.ManageNodeSystem, &mcpsdk.Tool{
		Name:        "ivory_system_copy_id",
		Description: "Append a public key to the host's authorized_keys, authenticating once with an ssh password. The password is used only to install the key.",
		Annotations: mutating,
	}, mr.Node.SystemCopyId)

	// container - what is deployed onto it
	addTool(r, coreConfig.ViewNodePlatformContainer, &mcpsdk.Tool{
		Name:        "ivory_container_list",
		Description: "Containers on the node's host, as the platform reports them. Names from here are what the other container tools take.",
		Annotations: readOnly,
	}, mr.Node.ContainerList)

	addTool(r, coreConfig.ViewNodePlatformContainer, &mcpsdk.Tool{
		Name:        "ivory_container_logs",
		Description: "The tail of one container's output. This is where a container that failed to start says why.",
		Annotations: readOnly,
	}, mr.Node.ContainerLogs)

	addTool(r, coreConfig.ViewNodePlatformContainer, &mcpsdk.Tool{
		Name:        "ivory_container_metrics",
		Description: "CPU, memory and network counters for one container.",
		Annotations: readOnly,
	}, mr.Node.ContainerMetrics)

	addTool(r, coreConfig.ManageNodePlatformContainer, &mcpsdk.Tool{
		Name:        "ivory_container_start",
		Description: "Start a stopped container.",
		Annotations: mutating,
	}, mr.Node.ContainerStart)

	addTool(r, coreConfig.ManageNodePlatformContainer, &mcpsdk.Tool{
		Name:        "ivory_container_stop",
		Description: "Stop a running container. On a database node this takes that node out of the cluster.",
		Annotations: destructive,
	}, mr.Node.ContainerStop)

	addTool(r, coreConfig.ManageNodePlatformContainer, &mcpsdk.Tool{
		Name:        "ivory_container_restart",
		Description: "Restart a container, dropping every connection to it.",
		Annotations: destructive,
	}, mr.Node.ContainerRestart)

	addTool(r, coreConfig.ManageNodePlatformContainer, &mcpsdk.Tool{
		Name: "ivory_container_up",
		Description: "Run a deployment command on the node's host. Ivory writes the `docker run` verb itself, so the command " +
			"supplies options only; {{placeholder}} variables are filled in from the cluster and its vaults, and any left " +
			"unresolved is an error rather than text reaching the shell.",
		Annotations: destructive,
	}, mr.Node.ContainerUp)

	addTool(r, coreConfig.ManageNodePlatformContainer, &mcpsdk.Tool{
		Name:        "ivory_container_down",
		Description: "DESTRUCTIVE: stop the container AND remove it. Anything not on a mounted volume is gone, including the logs that would explain a failure - read ivory_container_logs first.",
		Annotations: destructive,
	}, mr.Node.ContainerDown)

	addTool(r, coreConfig.ManageNodePlatformContainer, &mcpsdk.Tool{
		Name: "ivory_container_exec",
		Description: "Run a command inside a container. The executable is yours to choose (etcdctl, mongosh, psql), and the " +
			"command runs with whatever privileges that container has. Arbitrary execution: there is no allow-list.",
		Annotations: destructive,
	}, mr.Node.ContainerExec)

	addTool(r, coreConfig.ManageNodePlatformContainer, &mcpsdk.Tool{
		Name:        "ivory_container_keeper_deploy",
		Description: "Deploy one node of a cluster: run its command on its host, then its post-scripts. Use ivory_cluster_deploy for a whole new cluster.",
		Annotations: destructive,
	}, mr.Node.KeeperDeploy)
}

func tagTools(r *mcpRegistrar, mr *features.McpRouter) {
	addTool(r, coreConfig.ViewTagList, &mcpsdk.Tool{
		Name:        "ivory_tag_list",
		Description: "Every cluster tag, with the clusters carrying it.",
		Annotations: readOnly,
	}, mr.Tag.TagList)
}

func queryTools(r *mcpRegistrar, mr *features.McpRouter) {
	addTool(r, coreConfig.ViewQueryCrudList, &mcpsdk.Tool{
		Name:        "ivory_query_list",
		Description: "The saved queries, both the ones Ivory ships for monitoring and diagnostics and the user's own. Run one with ivory_query_template.",
		Annotations: readOnly,
	}, mr.Query.QueryList)

	addTool(r, coreConfig.ViewQueryLogList, &mcpsdk.Tool{
		Name:        "ivory_query_log_get",
		Description: "Past result sets recorded for one saved query.",
		Annotations: readOnly,
	}, mr.Query.QueryLogGet)

	addTool(r, coreConfig.ViewQueryDbInfo, &mcpsdk.Tool{
		Name:        "ivory_query_activity",
		Description: "Queries currently running on the node, as Ivory's own connections see them. Use the pids here with ivory_query_cancel or ivory_query_terminate.",
		Annotations: readOnly,
	}, mr.Query.QueryActivity)

	addTool(r, coreConfig.ViewQueryDbInfo, &mcpsdk.Tool{
		Name:        "ivory_query_databases",
		Description: "Databases on the node.",
		Annotations: readOnly,
	}, mr.Query.QueryDatabases)

	addTool(r, coreConfig.ViewQueryDbInfo, &mcpsdk.Tool{
		Name:        "ivory_query_schemas",
		Description: "Schemas in the selected database.",
		Annotations: readOnly,
	}, mr.Query.QuerySchemas)

	addTool(r, coreConfig.ViewQueryDbInfo, &mcpsdk.Tool{
		Name:        "ivory_query_tables",
		Description: "Tables in the selected schema.",
		Annotations: readOnly,
	}, mr.Query.QueryTables)

	addTool(r, coreConfig.ViewQueryDbChart, &mcpsdk.Tool{
		Name:        "ivory_query_chart",
		Description: "One precomputed summary figure about the node: database count, connection count, sizes, uptime.",
		Annotations: readOnly,
	}, mr.Query.QueryChart)

	addTool(r, coreConfig.ManageQueryDbTemplate, &mcpsdk.Tool{
		Name:        "ivory_query_template",
		Description: "Run a saved query against a node. Prefer this over writing SQL by hand: the shipped queries are the maintained, engine-correct versions.",
		Annotations: mutating,
	}, mr.Query.QueryTemplate)

	addTool(r, coreConfig.ManageQueryDbConsole, &mcpsdk.Tool{
		Name: "ivory_query_console",
		Description: "Run a statement against a node's database, verbatim. Ivory does not inspect or classify it: DROP, DELETE " +
			"and ALTER are as reachable as SELECT, and nothing here is undone afterwards. Prefer ivory_query_template for " +
			"anything it already covers, and read before you write.",
		Annotations: destructive,
	}, mr.Query.QueryConsole)

	addTool(r, coreConfig.ManageQueryDbCancel, &mcpsdk.Tool{
		Name:        "ivory_query_cancel",
		Description: "Ask one backend to abandon the statement it is running, leaving the connection open.",
		Annotations: mutating,
	}, mr.Query.QueryCancel)

	addTool(r, coreConfig.ManageQueryDbTerminate, &mcpsdk.Tool{
		Name:        "ivory_query_terminate",
		Description: "Kill a backend connection outright. The client sees its connection drop and any open transaction rolls back.",
		Annotations: destructive,
	}, mr.Query.QueryTerminate)

	addTool(r, coreConfig.ManageQueryCrudCreate, &mcpsdk.Tool{
		Name:        "ivory_query_create",
		Description: "Save a new query for reuse. Does not run it.",
		Annotations: mutating,
	}, mr.Query.QueryCreate)

	addTool(r, coreConfig.ManageQueryCrudUpdate, &mcpsdk.Tool{
		Name:        "ivory_query_update",
		Description: "Overwrite a saved query. A shipped query keeps its name, type and description; only its text can change.",
		Annotations: mutating,
	}, mr.Query.QueryUpdate)

	addTool(r, coreConfig.ManageQueryCrudDelete, &mcpsdk.Tool{
		Name:        "ivory_query_delete",
		Description: "Delete a saved query. Shipped queries cannot be deleted.",
		Annotations: destructive,
	}, mr.Query.QueryDelete)

	addTool(r, coreConfig.ManageQueryLogDelete, &mcpsdk.Tool{
		Name:        "ivory_query_log_delete",
		Description: "Discard the recorded results of one saved query.",
		Annotations: destructive,
	}, mr.Query.QueryLogDelete)
}

func deploymentTools(r *mcpRegistrar, mr *features.McpRouter) {
	addTool(r, coreConfig.ViewDeploymentTemplateList, &mcpsdk.Tool{
		Name: "ivory_deployment_template_list",
		Description: "Deployment templates: the ones Ivory ships for each keeper plugin plus the user's own. A template is an " +
			"ordered list of commands, one per node, with the {{placeholder}} variables a deploy fills in. Shipped templates " +
			"are computed rather than stored, so they are read-only - copy one to adjust it.",
		Annotations: readOnly,
	}, mr.Deployment.TemplateList)

	addTool(r, coreConfig.ManageDeploymentTemplateCreate, &mcpsdk.Tool{
		Name:        "ivory_deployment_template_create",
		Description: "Save a new deployment template. Every {{placeholder}} must come from Ivory's closed variable list or the save is refused.",
		Annotations: mutating,
	}, mr.Deployment.TemplateCreate)

	addTool(r, coreConfig.ManageDeploymentTemplateUpdate, &mcpsdk.Tool{
		Name:        "ivory_deployment_template_update",
		Description: "Overwrite one of the user's own templates. Its keeper and platform cannot change, and shipped templates cannot be edited.",
		Annotations: mutating,
	}, mr.Deployment.TemplateUpdate)

	addTool(r, coreConfig.ManageDeploymentTemplateDelete, &mcpsdk.Tool{
		Name:        "ivory_deployment_template_delete",
		Description: "Delete one of the user's own templates. Shipped templates cannot be deleted.",
		Annotations: destructive,
	}, mr.Deployment.TemplateDelete)
}

// SUPPORT
//
// Everything below serves the table above: how a caller reaches it, how a tool
// is gated, and how a handler's argument and result types become schemas.

// mcpRequestContext carries what the middlewares already worked out down to
// the SDK, which hands the server builder an *http.Request and nothing else.
func mcpRequestContext(context *gin.Context) context.Context {
	ctx := coreConfig.WithMcpCaller(context.Request.Context(), coreConfig.McpCaller{
		Username:    context.GetString(coreConfig.AuthContextKey.Username),
		Authorised:  context.GetBool(coreConfig.AuthContextKey.Authorised),
		AuthEnabled: context.GetBool(coreConfig.AuthContextKey.Enabled),
	})
	if value, ok := context.Get(coreConfig.PermissionContextKey); ok {
		if permissions, ok := value.(permission.PermissionMap); ok {
			ctx = withMcpPermissions(ctx, permissions)
		}
	}
	return ctx
}

type mcpPermissionsContextKey struct{}

func withMcpPermissions(ctx context.Context, permissions permission.PermissionMap) context.Context {
	return context.WithValue(ctx, mcpPermissionsContextKey{}, permissions)
}

func mcpPermissionsFrom(ctx context.Context) permission.PermissionMap {
	if permissions, ok := ctx.Value(mcpPermissionsContextKey{}).(permission.PermissionMap); ok {
		return permissions
	}
	return permission.PermissionMap{}
}

// mcpRegistrar builds one caller's server.
type mcpRegistrar struct {
	server      *mcpsdk.Server
	permissions permission.PermissionMap
}

func (r *mcpRegistrar) granted(feature coreConfig.Feature) bool {
	return r.permissions[feature] == permission.GRANTED
}

// addTool registers a tool only when the caller's map grants its feature, so a
// tool the user may not call is never offered - the equivalent of
// ValidateMethodMiddleware on a route, decided per tool instead. It is a free
// function because Go does not allow a method to introduce type parameters.
func addTool[In, Out any](r *mcpRegistrar, feature coreConfig.Feature, tool *mcpsdk.Tool, handler mcpsdk.ToolHandlerFor[In, Out]) {
	if !r.granted(feature) {
		return
	}
	addToolUngated(r, tool, handler)
}

// addToolUngated fills in both schemas before handing the tool over. The SDK
// infers them itself only when they are nil, and its inference has no hook for
// mcpTypeSchemas, so they are built here instead. The input then goes through
// gin's own validator, the one ShouldBindJSON runs, because a schema can say a
// field must be present but not that it must be non-zero.
func addToolUngated[In, Out any](r *mcpRegistrar, tool *mcpsdk.Tool, handler mcpsdk.ToolHandlerFor[In, Out]) {
	if tool.InputSchema == nil {
		tool.InputSchema = mcpSchemaFor[In](tool.Name)
	}
	if tool.OutputSchema == nil {
		tool.OutputSchema = mcpSchemaFor[Out](tool.Name)
	}
	mcpsdk.AddTool(r.server, tool, func(ctx context.Context, req *mcpsdk.CallToolRequest, in In) (*mcpsdk.CallToolResult, Out, error) {
		if err := binding.Validator.ValidateStruct(in); err != nil {
			var zero Out
			return nil, zero, err
		}
		return handler(ctx, req, in)
	})
}

// hintFalse and hintTrue back the annotation pointers below; the SDK only
// reads them.
var (
	hintFalse = false
	hintTrue  = true

	readOnly    = &mcpsdk.ToolAnnotations{ReadOnlyHint: true}
	mutating    = &mcpsdk.ToolAnnotations{DestructiveHint: &hintFalse}
	destructive = &mcpsdk.ToolAnnotations{DestructiveHint: &hintTrue}
)

// mcpTypeSchemas corrects the types reflection cannot read off a Go type
// alone. uuid.UUID is a [16]byte, so inference calls it an array, while it
// marshals as a string - the SDK then rejects Ivory's own vault ids against
// the schema it generated for them. The rich types list their legal values,
// the same ones their Valid method accepts, so an agent reads them instead of
// guessing.
var mcpTypeSchemas = map[reflect.Type]*jsonschema.Schema{
	reflect.TypeFor[uuid.UUID]():                      {Type: "string", Format: "uuid"},
	reflect.TypeFor[keeper.PluginType]():              mcpEnum("string", keeper.PluginTypes),
	reflect.TypeFor[database.PluginType]():            mcpEnum("string", database.PluginTypes),
	reflect.TypeFor[platform.PluginType]():            mcpEnum("string", []platform.PluginType{platform.Docker}),
	reflect.TypeFor[database.SystemChartType]():       mcpEnum("string", database.SystemChartTypes),
	reflect.TypeFor[database.SystemRequestCategory](): mcpEnum("integer", []database.SystemRequestCategory{database.BLOAT, database.ACTIVITY, database.REPLICATION, database.STATISTIC, database.OTHER}),
	reflect.TypeFor[database.SystemRequestVariety]():  mcpEnum("integer", []database.SystemRequestVariety{database.DatabaseSensitive, database.MasterOnly, database.ReplicaRecommended}),
}

func mcpEnum[T any](schemaType string, values []T) *jsonschema.Schema {
	enum := make([]any, 0, len(values))
	for _, value := range values {
		enum = append(enum, value)
	}
	return &jsonschema.Schema{Type: schemaType, Enum: enum}
}

// mcpSchemas and mcpResolvedSchemas exist because a stateless server is built
// per request: a schema depends only on its Go type, so it is reflected once
// and then shared, and the SDK resolves each shared pointer once too.
var (
	mcpSchemas         sync.Map
	mcpResolvedSchemas = mcpsdk.NewSchemaCache()
)

func mcpSchemaFor[T any](name string) *jsonschema.Schema {
	goType := reflect.TypeFor[T]()
	if cached, ok := mcpSchemas.Load(goType); ok {
		return cached.(*jsonschema.Schema)
	}
	schema, err := jsonschema.For[T](&jsonschema.ForOptions{TypeSchemas: mcpTypeSchemas})
	if err != nil {
		panic(fmt.Sprintf("mcp: tool %q: %v", name, err))
	}
	mcpRequired(goType, schema)
	cached, _ := mcpSchemas.LoadOrStore(goType, schema)
	return cached.(*jsonschema.Schema)
}

// mcpRequired rewrites an inferred schema's required list to match what the
// http engine actually requires. Inference marks every field required because
// Ivory's models spell them all out without omitempty, but gin binds an absent
// field to its zero value unless the struct says binding:"required" - so the
// MCP engine was refusing requests the REST API accepts, and an agent had to
// send "body": null to read a keeper config.
func mcpRequired(goType reflect.Type, schema *jsonschema.Schema) {
	if schema == nil {
		return
	}
	for goType != nil && goType.Kind() == reflect.Pointer {
		goType = goType.Elem()
	}
	if goType == nil {
		return
	}

	switch goType.Kind() {
	case reflect.Struct:
		if schema.Properties == nil {
			return
		}
		required := make([]string, 0, len(schema.Required))
		collectRequired(goType, schema, &required)
		schema.Required = required
	case reflect.Slice, reflect.Array:
		mcpRequired(goType.Elem(), schema.Items)
	case reflect.Map:
		mcpRequired(goType.Elem(), schema.AdditionalProperties)
	}
}

// collectRequired walks one struct level, following embedded fields into the
// same schema because their json fields are flattened there.
func collectRequired(goType reflect.Type, schema *jsonschema.Schema, required *[]string) {
	for i := range goType.NumField() {
		field := goType.Field(i)
		if !field.IsExported() {
			continue
		}
		if field.Anonymous && fieldKind(field.Type) == reflect.Struct {
			collectRequired(structType(field.Type), schema, required)
			continue
		}
		name := jsonName(field)
		if name == "" {
			continue
		}
		property, ok := schema.Properties[name]
		if !ok {
			continue
		}
		if strings.Contains(field.Tag.Get("binding"), "required") {
			*required = append(*required, name)
		}
		mcpRequired(field.Type, property)
	}
}

func fieldKind(goType reflect.Type) reflect.Kind {
	for goType.Kind() == reflect.Pointer {
		goType = goType.Elem()
	}
	return goType.Kind()
}

func structType(goType reflect.Type) reflect.Type {
	for goType.Kind() == reflect.Pointer {
		goType = goType.Elem()
	}
	return goType
}

func jsonName(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	if tag == "-" {
		return ""
	}
	if name, _, _ := strings.Cut(tag, ","); name != "" {
		return name
	}
	return field.Name
}
