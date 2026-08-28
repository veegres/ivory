package features

import (
	"ivory/clients/auth/basic"
	"ivory/clients/auth/ldap"
	"ivory/clients/auth/oidc"
	"ivory/clients/storage"
	"ivory/core"
	coreConfig "ivory/core/config"
	"ivory/core/service/encryption"
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
	"ivory/plugins/database"
	"ivory/plugins/keeper"
	"ivory/plugins/platform"
	"ivory/tools"
)

type Router struct {
	Auth       *auth.Router
	Cluster    *cluster.Router
	Permission *permission.Router
	Tag        *tag.Router
	Node       *node.Router
	Query      *query.Router
	Deployment *deployment.Router
	Management *management.Router
	Config     *config.Router
}

type Context struct {
	Router *Router
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

	// FILES
	configFiles := storage.NewFileStorage("config", ".json")
	queryLogFiles := storage.NewFileStorage("query", ".jsonl")

	// REPOS
	clusterRepo := cluster.NewRepository(clusterBucket)
	tagRepo := tag.NewRepository(tagBucket)
	permissionRepo := permission.NewRepository(permissionBucket)
	queryRepo := query.NewRepository(queryBucket, queryLogFiles)
	deploymentRepo := deployment.NewRepository(deploymentBucket)

	// AUTH PROVIDER
	basicProvider := basic.NewProvider()
	ldapProvider := ldap.NewProvider()
	oidcProvider := oidc.NewProvider()

	// CORE SERVICES
	vaultService := coreService.Vault
	jobService := coreService.Job
	certService := coreService.Cert

	encryptionService := encryption.NewService()
	secretService := vaultService.SecretService()
	permissionService := permission.NewService(permissionRepo)

	nodeService := node.NewService(platformRegistry, keeperRegistry, vaultService, certService, jobService)
	tagService := tag.NewService(tagRepo)
	queryService := query.NewService(queryRepo, databaseRegistry, vaultService, certService, env.Version.Label)
	deploymentService := deployment.NewService(deploymentRepo, keeperRegistry, platformRegistry)
	clusterService := cluster.NewService(clusterRepo, nodeService, tagService, queryService, vaultService, toolRegistry)
	authService := auth.NewService(secretService, basicProvider, ldapProvider, oidcProvider, permissionService)
	configService := config.NewService(configFiles, encryptionService, secretService, authService, permissionService, basicProvider, ldapProvider, oidcProvider)
	backupService := backup.NewService(clusterService, queryService, permissionService, deploymentService)
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
		toolRegistry,
	)

	return &Context{
		Router: &Router{
			Auth:       auth.NewRouter(authService, env.Config.UrlPath, env.Config.TlsEnabled),
			Cluster:    cluster.NewRouter(clusterService),
			Permission: permission.NewRouter(permissionService),
			Tag:        tag.NewRouter(tagService),
			Node:       node.NewRouter(nodeService),
			Query:      query.NewRouter(queryService, configService),
			Deployment: deployment.NewRouter(deploymentService),
			Management: management.NewRouter(managementService),
			Config:     config.NewRouter(configService),
		},
	}
}
