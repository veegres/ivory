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
	"ivory/src/storage/bolt"
	"ivory/src/storage/env"
	"ivory/src/storage/files"
)

type Context struct {
	env              *env.Env
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
	v := env.NewEnv()

	// DB
	db := bolt.NewBoltDB("ivory.db")
	clusterBucket := bolt.NewBoltBucket[cluster.Cluster](db, "Cluster")
	compactTableBucket := bolt.NewBoltBucket[bloat.Bloat](db, "CompactTable")
	certBucket := bolt.NewBoltBucket[cert.Cert](db, "Cert")
	tagBucket := bolt.NewBoltBucket[[]string](db, "Tag")
	secretBucket := bolt.NewBoltBucket[string](db, "Secret")
	passwordBucket := bolt.NewBoltBucket[password.Password](db, "Password")
	queryBucket := bolt.NewBoltBucket[query.Query](db, "Query")

	// FILES
	compactTableFiles := files.NewFileGateway("pgcompacttable", ".log")
	certFiles := files.NewFileGateway("cert", ".crt")
	configFiles := files.NewFileGateway("config", ".json")
	queryLogFiles := files.NewFileGateway("query", ".jsonl")

	// REPOS
	clusterRepo := cluster.NewClusterRepository(clusterBucket)
	bloatRepo := bloat.NewBloatRepository(compactTableBucket, compactTableFiles)
	certRepo := cert.NewCertRepository(certBucket, certFiles)
	tagRepo := tag.NewTagRepository(tagBucket)
	secretRepo := secret.NewSecretRepository(secretBucket)
	passwordRepo := password.NewPasswordRepository(passwordBucket)
	queryRepo := query.NewQueryRepository(queryBucket, queryLogFiles)

	// CLIENTS
	sidecarGateway := sidecar.NewSidecarGateway()
	patroniClient := patroni.NewPatroniClient(sidecarGateway)
	postgresClient := postgres.NewPostgresClient(v.Version.Label)

	// SERVICES
	encryptionService := encryption.NewEncryptionService()
	secretService := secret.NewSecretService(secretRepo, encryptionService)
	configService := config.NewService(configFiles, encryptionService, secretService)
	passwordService := password.NewPasswordService(passwordRepo, secretService, encryptionService)
	certService := cert.NewCertService(certRepo)
	instanceService := instance.NewService(patroniClient, passwordService, certService)
	tagService := tag.NewTagService(tagRepo)
	clusterService := cluster.NewClusterService(clusterRepo, instanceService, tagService)
	queryService := query.NewQueryService(queryRepo, postgresClient, secretService, passwordService, certService)
	bloatService := bloat.NewBloatService(bloatRepo, passwordService)
	authService := auth.NewAuthService(secretService, encryptionService)
	managementService := management.NewService(
		v,
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
		env:              v,
		authRouter:       auth.NewAuthRouter(authService, configService),
		clusterRouter:    cluster.NewClusterRouter(clusterService),
		bloatRouter:      bloat.NewBloatRouter(bloatService),
		certRouter:       cert.NewCertRouter(certService),
		secretRouter:     secret.NewSecretRouter(secretService),
		passwordRouter:   password.NewPasswordRouter(passwordService),
		tagRouter:        tag.NewTagRouter(tagService),
		instanceRouter:   instance.NewInstanceRouter(instanceService),
		queryRouter:      query.NewQueryRouter(queryService, configService),
		managementRouter: management.NewRouter(managementService),
		configRouter:     config.NewRouter(configService),
	}
}
