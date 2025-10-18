package app

import (
	"ivory/src/clients/database/postgres"
	"ivory/src/clients/sidecar"
	"ivory/src/clients/sidecar/patroni"
	"ivory/src/features/auth"
	"ivory/src/features/bloat"
	"ivory/src/features/cert"
	"ivory/src/features/cluster"
	"ivory/src/features/config"
	"ivory/src/features/encryption"
	"ivory/src/features/instance"
	"ivory/src/features/management"
	"ivory/src/features/password"
	"ivory/src/features/query"
	"ivory/src/features/secret"
	"ivory/src/features/tag"
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
	passwordRouter   *password.Router
	tagRouter        *tag.Router
	instanceRouter   *instance.Router
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
	passwordBucket := db.NewBucket[password.Password](st, "Password")
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
	passwordRepo := password.NewRepository(passwordBucket)
	queryLogRepo := query.NewLogRepository(queryLogFiles)
	queryRepo := query.NewRepository(queryBucket, queryLogRepo)

	// CLIENTS
	sidecarGateway := sidecar.NewGateway()
	patroniClient := patroni.NewClient(sidecarGateway)
	postgresClient := postgres.NewClient(appEnv.Version.Label)

	// SERVICES
	encryptionService := encryption.NewService()
	secretService := secret.NewService(secretRepo, encryptionService)
	configService := config.NewService(configFiles, encryptionService, secretService)
	passwordService := password.NewService(passwordRepo, secretService, encryptionService)
	certService := cert.NewService(certRepo)
	instanceService := instance.NewService(patroniClient, passwordService, certService)
	tagService := tag.NewService(tagRepo)
	clusterService := cluster.NewService(clusterRepo, instanceService, tagService)
	queryService := query.NewService(queryRepo, secretService)
	queryLogService := query.NewLogService(queryLogRepo)
	queryRunService := query.NewRunService(queryRepo, postgresClient, queryLogService, passwordService, certService)
	bloatService := bloat.NewService(bloatRepo, passwordService)
	authService := auth.NewService(secretService, configService)
	managementService := management.NewService(
		appEnv,
		authService,
		passwordService,
		clusterService,
		certService,
		tagService,
		bloatService,
		queryService,
		secretService,
		configService,
	)

	return &Context{
		env:              appEnv,
		authRouter:       auth.NewRouter(authService, appEnv.Config.UrlPath, appEnv.Config.TlsEnabled),
		clusterRouter:    cluster.NewRouter(clusterService),
		bloatRouter:      bloat.NewRouter(bloatService),
		certRouter:       cert.NewRouter(certService),
		secretRouter:     secret.NewRouter(secretService),
		passwordRouter:   password.NewRouter(passwordService),
		tagRouter:        tag.NewRouter(tagService),
		instanceRouter:   instance.NewRouter(instanceService),
		queryRouter:      query.NewRouter(queryService, queryRunService, queryLogService, configService),
		managementRouter: management.NewRouter(managementService),
		configRouter:     config.NewRouter(configService),
	}
}
