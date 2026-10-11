package features

import (
	"ivory/clients/auth/basic"
	"ivory/clients/auth/ldap"
	"ivory/clients/auth/oidc"
	"ivory/clients/storage"
	"ivory/core"
	coreConfig "ivory/core/config"
	"ivory/core/service/encryption"
	"ivory/core/service/token"
	"ivory/core/utils"
	"ivory/features/auth"
	"ivory/features/backup"
	"ivory/features/cluster"
	"ivory/features/config"
	"ivory/features/deployment"
	"ivory/features/management"
	"ivory/features/node"
	"ivory/features/permission"
	"ivory/features/query"
	"ivory/features/tag"
	"ivory/features/user"
	"ivory/plugins/database"
	"ivory/plugins/keeper"
	"ivory/plugins/platform"
	"ivory/tools"
)

type HttpRouter struct {
	Auth       *auth.HttpRouter
	Cluster    *cluster.HttpRouter
	Permission *permission.HttpRouter
	Tag        *tag.HttpRouter
	Node       *node.HttpRouter
	Query      *query.HttpRouter
	Deployment *deployment.HttpRouter
	Management *management.HttpRouter
	Config     *config.HttpRouter
	User       *user.HttpRouter
}

// McpRouter is the MCP engine's half of every feature, the way Router is the
// http engine's. Each entry is that feature's own <feature>_mcp.go; which
// tools exist and what gates them is declared in core/engine/engine_mcp.go.
type McpRouter struct {
	Cluster    *cluster.McpRouter
	Node       *node.McpRouter
	Tag        *tag.McpRouter
	Query      *query.McpRouter
	Deployment *deployment.McpRouter
	Management *management.McpRouter
}

type Context struct {
	HttpRouter *HttpRouter
	McpRouter  *McpRouter
}

func NewContext(
	env *coreConfig.Environment,
	databaseRegistry *utils.Registry[database.PluginType, database.Adapter],
	platformRegistry *utils.Registry[platform.PluginType, platform.Plugin],
	keeperRegistry *utils.Registry[keeper.PluginType, keeper.Plugin],
	toolRegistry *utils.Registry[tools.Tool, tools.Adapter],
	coreService *core.Service,
) *Context {
	// DB
	st := storage.NewDbStorage("feature.db")
	clusterBucket := storage.NewDbBucket[cluster.Response](st, "Cluster")
	tagBucket := storage.NewDbBucket[[]string](st, "Tag")
	permissionBucket := storage.NewDbBucket[permission.PermissionMap](st, "Permission")
	queryBucket := storage.NewDbBucket[query.Response](st, "Query")
	deploymentBucket := storage.NewDbBucket[deployment.Template](st, "DeploymentTemplate")
	userBucket := storage.NewDbBucket[user.User](st, "User")

	// FILES
	configFiles := storage.NewFileStorage("config", ".json")
	queryLogFiles := storage.NewFileStorage("query", ".jsonl")

	// REPOS
	clusterRepo := cluster.NewRepository(clusterBucket)
	tagRepo := tag.NewRepository(tagBucket)
	permissionRepo := permission.NewRepository(permissionBucket)
	queryRepo := query.NewRepository(queryBucket, queryLogFiles)
	deploymentRepo := deployment.NewRepository(deploymentBucket)
	userRepo := user.NewRepository(userBucket)

	// CORE SERVICES
	vaultService := coreService.Vault
	jobService := coreService.Job
	certService := coreService.Cert

	encryptionService := encryption.NewService()
	secretService := vaultService.SecretService()
	tokenService := token.NewService(secretService)
	permissionService := permission.NewService(permissionRepo)

	userService := user.NewService(userRepo, encryptionService, secretService, permissionService, tokenService)

	// AUTH PROVIDER
	basicProvider := basic.NewProvider(userService)
	ldapProvider := ldap.NewProvider()
	oidcProvider := oidc.NewProvider()

	nodeService := node.NewService(platformRegistry, keeperRegistry, vaultService, certService, jobService)
	tagService := tag.NewService(tagRepo)
	queryService := query.NewService(queryRepo, databaseRegistry, vaultService, certService, env.Version.Label)
	deploymentService := deployment.NewService(deploymentRepo, keeperRegistry, platformRegistry)
	clusterService := cluster.NewService(clusterRepo, nodeService, tagService, queryService, vaultService, toolRegistry)
	authService := auth.NewService(tokenService, basicProvider, ldapProvider, oidcProvider, userService)
	configService := config.NewService(configFiles, encryptionService, secretService, authService, userService, basicProvider, ldapProvider, oidcProvider)
	backupService := backup.NewService(clusterService, queryService, permissionService, deploymentService, userService)
	managementService := management.NewService(
		env,
		authService,
		vaultService,
		clusterService,
		certService,
		tagService,
		jobService,
		queryService,
		deploymentService,
		nodeService,
		secretService,
		configService,
		permissionService,
		backupService,
		userService,
		toolRegistry,
	)
	return &Context{
		McpRouter: &McpRouter{
			Cluster:    cluster.NewMcpRouter(clusterService),
			Node:       node.NewMcpRouter(nodeService),
			Tag:        tag.NewMcpRouter(tagService),
			Query:      query.NewMcpRouter(queryService),
			Deployment: deployment.NewMcpRouter(deploymentService),
			Management: management.NewMcpRouter(managementService),
		},
		HttpRouter: &HttpRouter{
			Auth:       auth.NewHttpRouter(authService, env.Config.UrlPath, env.Config.TlsEnabled),
			Cluster:    cluster.NewHttpRouter(clusterService),
			Permission: permission.NewHttpRouter(permissionService),
			Tag:        tag.NewHttpRouter(tagService),
			Node:       node.NewHttpRouter(nodeService),
			Query:      query.NewHttpRouter(queryService, configService),
			Deployment: deployment.NewHttpRouter(deploymentService),
			Management: management.NewHttpRouter(managementService),
			Config:     config.NewHttpRouter(configService),
			User:       user.NewHttpRouter(userService),
		},
	}
}
