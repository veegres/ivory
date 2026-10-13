package engine

import (
	"context"
	"encoding/json"
	"ivory/clients/console/ssh"
	"ivory/clients/storage"
	"ivory/core"
	coreConfig "ivory/core/config"
	"ivory/core/service/cert"
	"ivory/core/service/encryption"
	"ivory/core/service/secret"
	"ivory/core/service/vault"
	"ivory/core/utils"
	"ivory/features"
	"ivory/features/cluster"
	"ivory/features/deployment"
	"ivory/features/node"
	"ivory/features/permission"
	"ivory/features/query"
	"ivory/features/tag"
	"ivory/plugins/database"
	"ivory/plugins/keeper"
	"ivory/plugins/platform"
	"ivory/tools"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// createTestMcpRouter wires each feature's MCP router against a temporary
// database, the way management_service_test.go wires its services. The plugin
// registries are empty, so nothing here can reach a host.
// mcpTestRouters is the pair NewMcpServer takes, carried as one value so a
// test passes it around the way it used to pass the features half alone.
type mcpTestRouters struct {
	core     *core.McpRouter
	features *features.McpRouter
}

func createTestMcpRouter(t *testing.T) mcpTestRouters {
	t.Helper()

	tmpDir, errDir := os.MkdirTemp("", "engine-mcp-test-*")
	if errDir != nil {
		t.Fatalf("failed to create temp dir: %v", errDir)
	}
	t.Cleanup(func() { os.RemoveAll(tmpDir) })

	db, errOpen := bolt.Open(filepath.Join(tmpDir, "test.db"), 0600, nil)
	if errOpen != nil {
		t.Fatalf("failed to open test database: %v", errOpen)
	}
	t.Cleanup(func() { db.Close() })

	oldWd, errWd := os.Getwd()
	if errWd != nil {
		t.Fatalf("failed to get working dir: %v", errWd)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	t.Cleanup(func() { os.Chdir(oldWd) })

	secretService := secret.NewService(
		secret.NewRepository(storage.NewDbBucket[string](db, "Secret")),
		encryption.NewService(),
	)
	if err := secretService.SetDefault(); err != nil {
		t.Fatalf("failed to set default secret: %v", err)
	}
	vaultService := vault.NewService(
		vault.NewRepository(storage.NewDbBucket[vault.Vault](db, "Vault")),
		ssh.NewClient(),
		secretService,
		encryption.NewService(),
	)
	certService := cert.NewService(
		cert.NewRepository(storage.NewDbBucket[cert.Cert](db, "Cert"), storage.NewFileStorage("cert", "")),
	)
	tagService := tag.NewService(tag.NewRepository(storage.NewDbBucket[[]string](db, "Tag")))
	queryService := query.NewService(
		query.NewRepository(storage.NewDbBucket[query.Response](db, "Query"), storage.NewFileStorage("query-logs", ".log")),
		utils.NewRegistry[database.PluginType, database.Adapter](),
		vaultService,
		certService,
		"ivory",
	)
	deploymentService := deployment.NewService(
		deployment.NewRepository(storage.NewDbBucket[deployment.Template](db, "DeploymentTemplate")),
		utils.NewRegistry[keeper.PluginType, keeper.Plugin](),
		utils.NewRegistry[platform.PluginType, platform.Plugin](),
	)
	nodeService := node.NewService(
		utils.NewRegistry[platform.PluginType, platform.Plugin](),
		utils.NewRegistry[keeper.PluginType, keeper.Plugin](),
		vaultService,
		certService,
		nil,
	)
	clusterService := cluster.NewService(
		cluster.NewRepository(storage.NewDbBucket[cluster.Response](db, "Cluster")),
		nodeService,
		tagService,
		queryService,
		vaultService,
		utils.NewRegistry[tools.Tool, tools.Adapter](),
	)

	sshVault := uuid.New()
	port := func(value int) *int { return &value }
	if _, err := clusterService.Create(cluster.Request{
		Name:  "prod",
		Nodes: []cluster.NodeConfig{{Name: "pg1", Host: "10.0.0.1", SshPort: port(22), KeeperPort: port(8008), DbPort: port(5432)}},
		Options: cluster.Options{
			Plugins: cluster.Plugins{Keeper: "patroni_postgres", Database: "postgres"},
			Vaults:  cluster.Vaults{SshKeyId: &sshVault},
		},
	}); err != nil {
		t.Fatalf("failed to store the test cluster: %v", err)
	}

	cr := &core.McpRouter{
		Vault: vault.NewMcpRouter(vaultService),
		Cert:  cert.NewMcpRouter(certService),
	}
	return mcpTestRouters{core: cr, features: &features.McpRouter{
		Cluster:    cluster.NewMcpRouter(clusterService),
		Node:       node.NewMcpRouter(nodeService),
		Tag:        tag.NewMcpRouter(tagService),
		Query:      query.NewMcpRouter(queryService),
		Deployment: deployment.NewMcpRouter(deploymentService),
		Management: nil,
	}}
}

func mcpContextWith(features ...coreConfig.Feature) context.Context {
	permissions := make(permission.PermissionMap, len(features))
	for _, feature := range features {
		permissions[feature] = permission.GRANTED
	}
	return withMcpPermissions(context.Background(), permissions)
}

// session connects a client to the server this caller is given.
func session(t *testing.T, mr mcpTestRouters, ctx context.Context) *mcpsdk.ClientSession {
	t.Helper()
	serverTransport, clientTransport := mcpsdk.NewInMemoryTransports()
	server := newMcpServer(ctx, mr.core, mr.features, "test")
	go func() { _ = server.Run(context.Background(), serverTransport) }()

	clientSession, err := mcpsdk.NewClient(&mcpsdk.Implementation{Name: "test", Version: "v0"}, nil).
		Connect(context.Background(), clientTransport, nil)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(func() { clientSession.Close() })
	return clientSession
}

func offeredTools(t *testing.T, mr mcpTestRouters, ctx context.Context) []*mcpsdk.Tool {
	t.Helper()
	result, err := session(t, mr, ctx).ListTools(context.Background(), nil)
	if err != nil {
		t.Fatalf("list tools: %v", err)
	}
	return result.Tools
}

func listTools(t *testing.T, mr mcpTestRouters, ctx context.Context) []string {
	t.Helper()
	tools := offeredTools(t, mr, ctx)
	names := make([]string, 0, len(tools))
	for _, tool := range tools {
		names = append(names, tool.Name)
	}
	slices.Sort(names)
	return names
}

func viewFeatures() []coreConfig.Feature {
	features := make([]coreConfig.Feature, 0)
	for _, feature := range coreConfig.All {
		if strings.HasPrefix(string(feature), "view.") {
			features = append(features, feature)
		}
	}
	return features
}

func TestMcpToolsAreGatedByPermission(t *testing.T) {
	mr := createTestMcpRouter(t)

	t.Run("a caller holding nothing is offered only the ungated tool", func(t *testing.T) {
		names := listTools(t, mr, mcpContextWith())
		if !slices.Equal(names, []string{"ivory_app_info"}) {
			t.Errorf("tools = %v, want only ivory_app_info", names)
		}
	})

	t.Run("one feature offers only the tools it covers", func(t *testing.T) {
		names := listTools(t, mr, mcpContextWith(coreConfig.ViewClusterList))
		if !slices.Equal(names, []string{"ivory_app_info", "ivory_cluster_list"}) {
			t.Errorf("tools = %v, want app_info and cluster_list", names)
		}
	})

	t.Run("every tool a view-only caller is offered declares itself read-only", func(t *testing.T) {
		tools := offeredTools(t, mr, mcpContextWith(viewFeatures()...))
		if len(tools) == 0 {
			t.Fatal("a caller holding every view feature was offered nothing")
		}
		for _, tool := range tools {
			if tool.Annotations == nil || !tool.Annotations.ReadOnlyHint {
				t.Errorf("tool %q is reachable with only view permissions but is not annotated read-only", tool.Name)
			}
		}
	})

	t.Run("every mutating tool needs a manage feature", func(t *testing.T) {
		view := listTools(t, mr, mcpContextWith(viewFeatures()...))
		for _, tool := range offeredTools(t, mr, mcpContextWith(coreConfig.All...)) {
			if tool.Annotations != nil && tool.Annotations.ReadOnlyHint {
				continue
			}
			if slices.Contains(view, tool.Name) {
				t.Errorf("mutating tool %q is reachable with view permissions alone", tool.Name)
			}
		}
	})

	t.Run("granting everything offers every tool, each one gated by a real feature", func(t *testing.T) {
		all := listTools(t, mr, mcpContextWith(coreConfig.All...))
		if len(all) < 40 {
			t.Fatalf("only %d tools registered, expected the full surface", len(all))
		}

		// Granting one feature at a time and unioning the result proves every
		// tool hangs off a feature that is actually in config.All: a tool gated
		// by anything else could never appear in a single-feature server.
		union := make(map[string]bool)
		for _, feature := range coreConfig.All {
			for _, name := range listTools(t, mr, mcpContextWith(feature)) {
				union[name] = true
			}
		}
		for _, name := range all {
			if !union[name] {
				t.Errorf("tool %q is gated by a feature that is not in config.All", name)
			}
		}
	})
}

// TestMcpSchemasAcceptIvorysOwnTypes guards the two corrections mcpSchemaFor
// makes. Both were real failures: a uuid.UUID is a [16]byte, so an inferred
// schema calls it an array and rejects the vault ids Ivory writes itself; and
// Ivory's models carry no omitempty, so every nullable pointer came back
// required and an agent had to spell out explicit nulls to say "none".
func TestMcpSchemasAcceptIvorysOwnTypes(t *testing.T) {
	mr := createTestMcpRouter(t)
	ctx := context.Background()

	t.Run("a vault id passes as a uuid string, and nullable fields may be omitted", func(t *testing.T) {
		clientSession := session(t, mr, mcpContextWith(coreConfig.ManageClusterCreate))
		result, err := clientSession.CallTool(ctx, &mcpsdk.CallToolParams{
			Name: "ivory_cluster_create",
			Arguments: map[string]any{
				"name": "staging",
				"nodes": []any{
					map[string]any{"name": "pg1", "host": "10.0.1.1", "sshPort": 22, "keeperPort": 8008, "dbPort": 5432},
				},
				"plugins": map[string]any{"keeper": "patroni_postgres", "database": "postgres"},
				"tls":     map[string]any{"keeper": false, "database": false},
				"certs":   map[string]any{},
				"vaults":  map[string]any{"sshKeyId": "b6686f05-5beb-40c2-9d25-600fc81ea5d7"},
			},
		})
		if err != nil {
			t.Fatalf("call: %v", err)
		}
		if result.IsError {
			for _, content := range result.Content {
				if text, ok := content.(*mcpsdk.TextContent); ok {
					t.Fatalf("the schema refused Ivory's own shapes: %s", text.Text)
				}
			}
			t.Fatal("the schema refused Ivory's own shapes")
		}
	})

	t.Run("a uuid comes back out as a string too", func(t *testing.T) {
		clientSession := session(t, mr, mcpContextWith(coreConfig.ViewClusterItem))
		result, err := clientSession.CallTool(ctx, &mcpsdk.CallToolParams{
			Name:      "ivory_cluster_get",
			Arguments: map[string]any{"name": "prod"},
		})
		if err != nil {
			t.Fatalf("call: %v", err)
		}
		if result.IsError {
			t.Fatalf("output validation refused the stored cluster: %v", result.Content)
		}

		// The result is the feature's own cluster.Response, not a wrapper
		// around it: a tool's structured result may be any JSON value.
		encoded, _ := json.Marshal(result.StructuredContent)
		var out struct {
			Name   string `json:"name"`
			Vaults struct {
				SshKeyId string `json:"sshKeyId"`
			} `json:"vaults"`
		}
		if err := json.Unmarshal(encoded, &out); err != nil {
			t.Fatalf("a vault id did not come back as a string: %v", err)
		}
		if out.Name != "prod" || out.Vaults.SshKeyId == "" {
			t.Errorf("cluster = %+v", out)
		}
	})

	t.Run("a rich type lists its values, and an integer one binds", func(t *testing.T) {
		var schema any
		for _, tool := range offeredTools(t, mr, mcpContextWith(coreConfig.ManageQueryCrudCreate)) {
			if tool.Name == "ivory_query_create" {
				schema = tool.InputSchema
			}
		}
		encoded, _ := json.Marshal(schema)
		for _, want := range []string{`"postgres"`, `"clickhouse"`, `"enum":[0,1,2,3,4]`} {
			if !strings.Contains(string(encoded), want) {
				t.Errorf("ivory_query_create schema does not list %s: %s", want, encoded)
			}
		}

		clientSession := session(t, mr, mcpContextWith(coreConfig.ManageQueryCrudCreate))
		result, err := clientSession.CallTool(ctx, &mcpsdk.CallToolParams{
			Name:      "ivory_query_create",
			Arguments: map[string]any{"name": "bloat", "type": 1, "plugin": "postgres", "query": "select 1", "description": nil, "varieties": []any{0}, "params": nil},
		})
		if err != nil {
			t.Fatalf("call: %v", err)
		}
		if result.IsError {
			t.Fatalf("a legal category and plugin were refused: %v", result.Content)
		}

		refused, err := clientSession.CallTool(ctx, &mcpsdk.CallToolParams{
			Name:      "ivory_query_create",
			Arguments: map[string]any{"name": "bloat", "type": 1, "plugin": "oracle", "query": "select 1"},
		})
		if err != nil {
			t.Fatalf("call: %v", err)
		}
		if !refused.IsError {
			t.Error("expected an unknown plugin to be refused")
		}
	})

	t.Run("every output schema marshals as an object, never the boolean an any output reflects to", func(t *testing.T) {
		for _, tool := range offeredTools(t, mr, mcpContextWith(coreConfig.All...)) {
			if tool.OutputSchema == nil {
				continue
			}
			encoded, _ := json.Marshal(tool.OutputSchema)
			var schema map[string]any
			if err := json.Unmarshal(encoded, &schema); err != nil {
				t.Errorf("tool %q output schema is not a JSON object: %s", tool.Name, encoded)
			}
		}
	})

	t.Run("a required argument the request states is still refused", func(t *testing.T) {
		clientSession := session(t, mr, mcpContextWith(coreConfig.ViewClusterItem))
		result, err := clientSession.CallTool(ctx, &mcpsdk.CallToolParams{Name: "ivory_cluster_get", Arguments: map[string]any{}})
		if err != nil {
			t.Fatalf("call: %v", err)
		}
		if !result.IsError {
			t.Error("expected a schema error when the cluster name is omitted")
		}
	})
}

// TestEveryFeatureIsWiredIn calls one representative tool per feature through
// a real session. The plugin registries are empty, so none of these can reach
// a host - what the test proves is that each feature's McpRouter is bound to
// the tool table and that its arguments bind, which is the part engine_mcp.go
// owns. A mis-wired feature fails here rather than in front of an agent.
func TestEveryFeatureIsWiredIn(t *testing.T) {
	mr := createTestMcpRouter(t)
	ctx := context.Background()
	clientSession := session(t, mr, mcpContextWith(coreConfig.All...))

	connection := map[string]any{"host": "10.0.0.1", "port": 8008, "vaultId": uuid.New().String(), "platform": "docker"}

	calls := []struct {
		feature   string
		tool      string
		arguments map[string]any
	}{
		{"cluster", "ivory_cluster_list", map[string]any{}},
		{"cluster/overview", "ivory_cluster_overview", map[string]any{"name": "prod", "host": "10.0.0.1", "port": 8008}},
		{"cluster", "ivory_cluster_connections", map[string]any{"name": "prod"}},
		{"cluster", "ivory_cluster_check", map[string]any{"name": "prod"}},
		{"vault", "ivory_vault_list", map[string]any{"type": 3}},
		{"cert", "ivory_cert_list", map[string]any{}},
		{"tag", "ivory_tag_list", map[string]any{}},
		{"node/keeper", "ivory_keeper_node_list", map[string]any{
			"host": "10.0.0.1", "port": 8008, "cluster": "prod", "plugin": "patroni_postgres",
		}},
		{"node/keeper", "ivory_node_keeper_check", map[string]any{
			"host": "10.0.0.1", "port": 8008, "cluster": "prod", "plugin": "patroni_postgres",
		}},
		{"node/platform", "ivory_node_platform_check", connection},
		{"node/system", "ivory_system_info", connection},
		{"node/container", "ivory_container_list", connection},
		{"node/container", "ivory_container_inspect", map[string]any{"connection": connection, "name": "pg1"}},
		{"query/crud", "ivory_query_list", map[string]any{}},
		{"query/db", "ivory_query_check", map[string]any{
			"db": map[string]any{"plugin": "postgres", "host": "10.0.0.1", "port": 5432},
		}},
		{"query/db", "ivory_query_databases", map[string]any{
			"name": "",
			"connection": map[string]any{
				"db": map[string]any{"plugin": "postgres", "host": "10.0.0.1", "port": 5432},
			},
		}},
		{"deployment", "ivory_deployment_template_list", map[string]any{}},
	}

	for _, call := range calls {
		t.Run(call.feature+" answers "+call.tool, func(t *testing.T) {
			result, err := clientSession.CallTool(ctx, &mcpsdk.CallToolParams{Name: call.tool, Arguments: call.arguments})
			if err != nil {
				t.Fatalf("the tool could not be called at all: %v", err)
			}
			// An error from the plugin boundary is fine; a schema error is not,
			// because it means the arguments never reached the feature.
			for _, content := range result.Content {
				text, ok := content.(*mcpsdk.TextContent)
				if ok && strings.Contains(text.Text, `validating "arguments"`) {
					t.Errorf("arguments were refused before reaching the feature: %s", text.Text)
				}
			}
		})
	}
}

// TestMcpFailureKeepsTheEnvelope pins the tool twin of the http envelope: a
// failed keeper or platform call is flagged as an error and still answers with
// {response, error}, so whatever the action produced is not dropped with it.
func TestMcpFailureKeepsTheEnvelope(t *testing.T) {
	mr := createTestMcpRouter(t)
	clientSession := session(t, mr, mcpContextWith(coreConfig.All...))

	calls := map[string]map[string]any{
		"ivory_keeper_node_list": {"host": "10.0.0.1", "port": 8008, "cluster": "prod", "plugin": "patroni_postgres"},
		"ivory_container_list":   {"host": "10.0.0.1", "port": 22, "vaultId": uuid.New().String(), "platform": "docker"},
	}
	for tool, arguments := range calls {
		t.Run(tool, func(t *testing.T) {
			result, err := clientSession.CallTool(context.Background(), &mcpsdk.CallToolParams{Name: tool, Arguments: arguments})
			if err != nil {
				t.Fatalf("the tool could not be called at all: %v", err)
			}
			if !result.IsError {
				t.Fatalf("expected the failure to be flagged as one, got %+v", result)
			}
			raw, errMarshal := json.Marshal(result.StructuredContent)
			if errMarshal != nil {
				t.Fatalf("marshal structured content: %v", errMarshal)
			}
			var envelope map[string]any
			if err := json.Unmarshal(raw, &envelope); err != nil {
				t.Fatalf("structured content is not an object: %s", raw)
			}
			if _, ok := envelope["response"]; !ok {
				t.Errorf("expected a response field beside the error, got %s", raw)
			}
			if message, _ := envelope["error"].(string); message == "" {
				t.Errorf("expected the error to be named, got %s", raw)
			}
		})
	}
}

// TestMcpCrossOriginProtection covers the one thing the in-memory transport
// cannot reach: the endpoint is reached over http, and Ivory's session cookie
// is a valid credential on it, so a browser page must not be able to drive
// tool calls from another origin. Non-browser MCP clients send neither header
// and have to keep working.
func TestMcpCrossOriginProtection(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mr := createTestMcpRouter(t)
	engine := gin.New()
	engine.POST("/mcp", NewMcpServer(
		&coreConfig.Environment{Version: coreConfig.Version{Label: "test"}},
		mr.core,
		mr.features,
	))

	body := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2026-07-28","capabilities":{},"clientInfo":{"name":"t","version":"0"}}}`

	tests := []struct {
		name      string
		headers   map[string]string
		forbidden bool
	}{
		{
			name:      "a cross-site browser request is refused",
			headers:   map[string]string{"Sec-Fetch-Site": "cross-site", "Origin": "https://evil.example"},
			forbidden: true,
		},
		{
			name:    "a non-browser client sends neither header and is allowed",
			headers: map[string]string{},
		},
		{
			name:    "a same-origin browser request is allowed",
			headers: map[string]string{"Sec-Fetch-Site": "same-origin"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(body))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json, text/event-stream")
			for name, value := range test.headers {
				request.Header.Set(name, value)
			}
			recorder := httptest.NewRecorder()
			engine.ServeHTTP(recorder, request)

			if test.forbidden && recorder.Code != http.StatusForbidden {
				t.Errorf("status = %d, want 403 for a cross-origin request", recorder.Code)
			}
			if !test.forbidden && recorder.Code == http.StatusForbidden {
				t.Errorf("status = 403, want the request to reach the mcp handler")
			}
		})
	}
}

// TestMcpSession pins what an MCP caller's session is for: the database keeps
// seven characters of it as application_name, so two names sharing a prefix
// must still differ there, and the same name must always land on the same one.
func TestMcpSession(t *testing.T) {
	if got := mcpSession(""); got != "mcp" {
		t.Errorf("a caller with no name: got %q, want mcp", got)
	}
	if mcpSession("alice") != mcpSession("alice") {
		t.Errorf("expected one caller to keep one session")
	}
	if mcpSession("alice.smith")[:7] == mcpSession("alice.jones")[:7] {
		t.Errorf("expected names sharing a prefix to differ within the first seven characters")
	}
}
