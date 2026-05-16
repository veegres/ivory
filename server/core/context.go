package core

import (
	"ivory/clients/auth/basic"
	"ivory/clients/auth/ldap"
	"ivory/clients/auth/oidc"
	"ivory/clients/http"
	"ivory/clients/ssh"
	"ivory/features/auth"
	"ivory/features/backup"
	"ivory/features/cert"
	"ivory/features/cluster"
	"ivory/features/config"
	"ivory/features/encryption"
	"ivory/features/management"
	"ivory/features/node"
	"ivory/features/permission"
	"ivory/features/query"
	"ivory/features/secret"
	"ivory/features/tag"
	"ivory/features/tools"
	"ivory/features/vault"
	database2 "ivory/plugins/database"
	"ivory/plugins/database/postgres"
	keeper2 "ivory/plugins/keeper"
	"ivory/plugins/keeper/patroni"
	platform2 "ivory/plugins/platform"
	"ivory/plugins/platform/onprem"
	"ivory/storage/db"
	"ivory/storage/env"
	"ivory/storage/files"
)

type Context struct {
	env              *env.AppEnv
	authRouter       *auth.Router
	clusterRouter    *cluster.Router
	toolsRouter      *tools.Router
	certRouter       *cert.Router
	secretRouter     *secret.Router
	vaultRouter      *vault.Router
	permissionRouter *permission.Router
	tagRouter        *tag.Router
	nodeRouter       *node.Router
	queryRouter      *query.Router
	managementRouter *management.Router
	configRouter     *config.Router
}

func NewContext() *Context {
	appEnv := env.NewAppEnv()

	// DB
	st := db.NewStorage("ivory.db")
	clusterBucket := db.NewBucket[cluster.Response](st, "Cluster")
	certBucket := db.NewBucket[cert.Cert](st, "Cert")
	tagBucket := db.NewBucket[[]string](st, "Tag")
	secretBucket := db.NewBucket[string](st, "Secret")
	vaultBucket := db.NewBucket[vault.Vault](st, "Vault")
	permissionBucket := db.NewBucket[permission.PermissionMap](st, "Permission")
	queryBucket := db.NewBucket[query.Response](st, "Query")

	// FILES
	certFiles := files.NewStorage("cert", ".crt")
	configFiles := files.NewStorage("config", ".json")
	queryLogFiles := files.NewStorage("query", ".jsonl")

	// REPOS
	clusterRepo := cluster.NewRepository(clusterBucket)
	certRepo := cert.NewRepository(certBucket, certFiles)
	tagRepo := tag.NewRepository(tagBucket)
	secretRepo := secret.NewRepository(secretBucket)
	vaultRepo := vault.NewRepository(vaultBucket)
	permissionRepo := permission.NewRepository(permissionBucket)
	queryRepo := query.NewRepository(queryBucket, queryLogFiles)

	// CLIENTS
	httpClient := http.NewClient()
	sshClient := ssh.NewClient()

	// ADAPTERS
	patroniAdapter := patroni.NewAdapter(httpClient)
	postgresAdapter := postgres.NewAdapter()
	onpremAdapter := onprem.NewAdapter(sshClient)

	// REGISTRY (we cannot use Factory pattern in clients package because of cycle dependencies)
	keeperPlugins := keeper2.NewPluginRegistry()
	keeperPlugins.Register(keeper2.PATRONI, patroniAdapter)
	dbPlugins := database2.NewPluginRegistry()
	dbPlugins.Register(database2.POSTGRES, postgresAdapter)
	platformPlugins := platform2.NewPluginRegistry()
	platformPlugins.Register(platform2.Onprem, onpremAdapter)

	// AUTH PROVIDER
	basicProvider := basic.NewProvider()
	ldapProvider := ldap.NewProvider()
	oidcProvider := oidc.NewProvider()

	// SERVICES
	encryptionService := encryption.NewService()
	secretService := secret.NewService(secretRepo, encryptionService)
	vaultService := vault.NewService(vaultRepo, sshClient, secretService, encryptionService)
	permissionService := permission.NewService(permissionRepo)
	certService := cert.NewService(certRepo)
	nodeService := node.NewService(platformPlugins, keeperPlugins, vaultService, certService)
	tagService := tag.NewService(tagRepo)
	toolsService := tools.NewService(vaultService)
	queryService := query.NewService(queryRepo, dbPlugins, vaultService, certService, secretService, appEnv.Version.Label)
	clusterService := cluster.NewService(clusterRepo, nodeService, tagService, queryService, toolsService, vaultService)
	authService := auth.NewService(secretService, basicProvider, ldapProvider, oidcProvider, permissionService)
	configService := config.NewService(configFiles, encryptionService, secretService, authService, permissionService, basicProvider, ldapProvider, oidcProvider)
	backupService := backup.NewService(clusterService, queryService, permissionService)
	managementService := management.NewService(
		appEnv,
		authService,
		vaultService,
		clusterService,
		certService,
		tagService,
		toolsService,
		queryService,
		nodeService,
		secretService,
		configService,
		permissionService,
		backupService,
	)

	return &Context{
		env:              appEnv,
		authRouter:       auth.NewRouter(authService, appEnv.Config.UrlPath, appEnv.Config.TlsEnabled),
		clusterRouter:    cluster.NewRouter(clusterService),
		toolsRouter:      tools.NewRouter(toolsService),
		certRouter:       cert.NewRouter(certService),
		secretRouter:     secret.NewRouter(secretService),
		vaultRouter:      vault.NewRouter(vaultService),
		permissionRouter: permission.NewRouter(permissionService),
		tagRouter:        tag.NewRouter(tagService),
		nodeRouter:       node.NewRouter(nodeService),
		queryRouter:      query.NewRouter(queryService, configService),
		managementRouter: management.NewRouter(managementService),
		configRouter:     config.NewRouter(configService),
	}
}
