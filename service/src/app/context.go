package app

import (
	"ivory/src/clients/auth/basic"
	"ivory/src/clients/auth/ldap"
	"ivory/src/clients/auth/oidc"
	"ivory/src/clients/http"
	"ivory/src/clients/ssh"
	"ivory/src/features/auth"
	"ivory/src/features/backup"
	"ivory/src/features/cert"
	"ivory/src/features/cluster"
	"ivory/src/features/config"
	"ivory/src/features/encryption"
	"ivory/src/features/management"
	"ivory/src/features/node"
	"ivory/src/features/permission"
	"ivory/src/features/query"
	"ivory/src/features/secret"
	"ivory/src/features/tag"
	"ivory/src/features/tools"
	"ivory/src/features/vault"
	"ivory/src/plugins/database"
	"ivory/src/plugins/database/postgres"
	"ivory/src/plugins/keeper"
	"ivory/src/plugins/keeper/patroni"
	"ivory/src/plugins/os"
	"ivory/src/plugins/os/linux"
	"ivory/src/storage/db"
	"ivory/src/storage/env"
	"ivory/src/storage/files"
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
	linuxAdapter := linux.NewAdapter(sshClient)

	// REGISTRY (we cannot use Factory pattern in clients package because of cycle dependencies)
	keeperPlugins := keeper.NewPluginRegistry()
	keeperPlugins.Register(keeper.PATRONI, patroniAdapter)
	dbPlugins := database.NewPluginRegistry()
	dbPlugins.Register(database.POSTGRES, postgresAdapter)
	osPlugins := os.NewPluginRegistry()
	osPlugins.Register(os.Linux, linuxAdapter)

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
	nodeService := node.NewService(osPlugins, keeperPlugins, vaultService, certService)
	tagService := tag.NewService(tagRepo)
	toolsService := tools.NewService(vaultService)
	queryService := query.NewService(queryRepo, dbPlugins, vaultService, certService, secretService, appEnv.Version.Label)
	clusterService := cluster.NewService(clusterRepo, nodeService, tagService, queryService, toolsService)
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
