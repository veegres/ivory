package app

import (
	"ivory/src/clients/auth/basic"
	"ivory/src/clients/auth/ldap"
	"ivory/src/clients/auth/oidc"
	"ivory/src/clients/database"
	"ivory/src/clients/database/postgres"
	"ivory/src/clients/http"
	"ivory/src/clients/keeper"
	"ivory/src/clients/keeper/patroni"
	"ivory/src/clients/ssh"
	"ivory/src/features/auth"
	"ivory/src/features/backup"
	"ivory/src/features/bloat"
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
	"ivory/src/features/vault"
	"ivory/src/storage/db"
	"ivory/src/storage/env"
	"ivory/src/storage/files"
)

type Context struct {
	env              *env.AppEnv
	authRouter       *auth.Router
	clusterRouter    *cluster.Router
	bloatRouter      *bloat.Router
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
	clusterBucket := db.NewBucket[cluster.Cluster](st, "Cluster")
	compactTableBucket := db.NewBucket[bloat.Bloat](st, "CompactTable")
	certBucket := db.NewBucket[cert.Cert](st, "Cert")
	tagBucket := db.NewBucket[[]string](st, "Tag")
	secretBucket := db.NewBucket[string](st, "Secret")
	vaultBucket := db.NewBucket[vault.Vault](st, "Vault")
	permissionBucket := db.NewBucket[permission.PermissionMap](st, "Permission")
	queryBucket := db.NewBucket[query.Query](st, "Query")

	// FILES
	compactTableFiles := files.NewStorage("pgcompacttable", ".log")
	certFiles := files.NewStorage("cert", ".crt")
	configFiles := files.NewStorage("config", ".json")
	queryLogFiles := files.NewStorage("query", ".jsonl")

	// REPOS
	clusterRepo := cluster.NewRepository(clusterBucket)
	bloatRepo := bloat.NewRepository(compactTableBucket, compactTableFiles)
	certRepo := cert.NewRepository(certBucket, certFiles)
	tagRepo := tag.NewRepository(tagBucket)
	secretRepo := secret.NewRepository(secretBucket)
	vaultRepo := vault.NewRepository(vaultBucket)
	permissionRepo := permission.NewRepository(permissionBucket)
	queryRepo := query.NewRepository(queryBucket, queryLogFiles)

	// CLIENTS
	httpClient := http.NewClient()
	patroniClient := patroni.NewClient(httpClient)
	postgresClient := postgres.NewClient(appEnv.Version.Label)
	sshClient := ssh.NewClient()

	// REGISTRY (we cannot use Factory pattern in clients package because of cycle dependencies)
	keeperRegistry := keeper.NewRegistry()
	keeperRegistry.Register(keeper.PATRONI, patroniClient)
	dbRegistry := database.NewRegistry()
	dbRegistry.Register(database.POSTGRES, postgresClient)

	// AUTH PROVIDER
	basicProvider := basic.NewProvider()
	ldapProvider := ldap.NewProvider()
	oidcProvider := oidc.NewProvider()

	// SERVICES
	encryptionService := encryption.NewService()
	secretService := secret.NewService(secretRepo, encryptionService)
	vaultService := vault.NewService(vaultRepo, secretService, encryptionService)
	permissionService := permission.NewService(permissionRepo)
	certService := cert.NewService(certRepo)
	nodeService := node.NewService(keeperRegistry, sshClient, vaultService, certService)
	tagService := tag.NewService(tagRepo)
	bloatService := bloat.NewService(bloatRepo, vaultService)
	queryService := query.NewService(queryRepo, dbRegistry, vaultService, certService, secretService)
	clusterService := cluster.NewService(clusterRepo, nodeService, tagService, queryService)
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
		bloatService,
		queryService,
		secretService,
		configService,
		permissionService,
		backupService,
	)

	return &Context{
		env:              appEnv,
		authRouter:       auth.NewRouter(authService, appEnv.Config.UrlPath, appEnv.Config.TlsEnabled),
		clusterRouter:    cluster.NewRouter(clusterService),
		bloatRouter:      bloat.NewRouter(bloatService),
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
