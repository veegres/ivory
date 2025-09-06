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
	env              *env.AppInfo
	authRouter       *auth.AuthRouter
	clusterRouter    *cluster.ClusterRouter
	bloatRouter      *bloat.BloatRouter
	certRouter       *cert.CertRouter
	secretRouter     *secret.SecretRouter
	passwordRouter   *password.PasswordRouter
	tagRouter        *tag.TagRouter
	instanceRouter   *instance.Router
	queryRouter      *query.QueryRouter
	managementRouter *management.Router
	configRouter     *config.Router
}

func NewContext() *Context {
	info := env.NewAppInfo()

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
	clusterRepo := cluster.NewClusterRepository(clusterBucket)
	bloatRepo := bloat.NewBloatRepository(compactTableBucket, compactTableFiles)
	certRepo := cert.NewCertRepository(certBucket, certFiles)
	tagRepo := tag.NewTagRepository(tagBucket)
	secretRepo := secret.NewSecretRepository(secretBucket)
	passwordRepo := password.NewPasswordRepository(passwordBucket)
	queryLogRepo := query.NewLogRepository(queryLogFiles)
	queryRepo := query.NewRepository(queryBucket, queryLogRepo)

	// CLIENTS
	sidecarGateway := sidecar.NewSidecarGateway()
	patroniClient := patroni.NewPatroniClient(sidecarGateway)
	postgresClient := postgres.NewPostgresClient(info.Version.Label)

	// SERVICES
	encryptionService := encryption.NewEncryptionService()
	secretService := secret.NewSecretService(secretRepo, encryptionService)
	configService := config.NewService(configFiles, encryptionService, secretService)
	passwordService := password.NewPasswordService(passwordRepo, secretService, encryptionService)
	certService := cert.NewCertService(certRepo)
	instanceService := instance.NewService(patroniClient, passwordService, certService)
	tagService := tag.NewTagService(tagRepo)
	clusterService := cluster.NewClusterService(clusterRepo, instanceService, tagService)
	queryService := query.NewService(queryRepo, secretService)
	queryLogService := query.NewLogService(queryLogRepo)
	queryRunService := query.NewRunService(queryRepo, postgresClient, queryLogService, passwordService, certService)
	bloatService := bloat.NewBloatService(bloatRepo, passwordService)
	authService := auth.NewAuthService(secretService, encryptionService)
	managementService := management.NewService(
		info,
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
		env:              info,
		authRouter:       auth.NewAuthRouter(authService, configService),
		clusterRouter:    cluster.NewClusterRouter(clusterService),
		bloatRouter:      bloat.NewBloatRouter(bloatService),
		certRouter:       cert.NewCertRouter(certService),
		secretRouter:     secret.NewSecretRouter(secretService),
		passwordRouter:   password.NewPasswordRouter(passwordService),
		tagRouter:        tag.NewTagRouter(tagService),
		instanceRouter:   instance.NewInstanceRouter(instanceService),
		queryRouter:      query.NewQueryRouter(queryService, queryRunService, queryLogService, configService),
		managementRouter: management.NewRouter(managementService),
		configRouter:     config.NewRouter(configService),
	}
}
