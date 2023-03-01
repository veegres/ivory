package app

import (
	"ivory/src/config"
	. "ivory/src/model"
	"ivory/src/persistence"
	"ivory/src/router"
	"ivory/src/service"
)

type Context struct {
	env            *config.Env
	infoRouter     *router.InfoRouter
	clusterRouter  *router.ClusterRouter
	bloatRouter    *router.BloatRouter
	certRouter     *router.CertRouter
	secretRouter   *router.SecretRouter
	passwordRouter *router.PasswordRouter
	tagRouter      *router.TagRouter
	instanceRouter *router.InstanceRouter
	queryRouter    *router.QueryRouter
	eraseRouter    *router.EraseRouter
}

func NewContext() *Context {
	env := config.NewEnv()

	db := config.NewBoltDB("ivory.db")
	clusterBucket := config.NewBoltBucket[ClusterModel](db, "Cluster")
	compactTableBucket := config.NewBoltBucket[CompactTableModel](db, "CompactTable")
	certBucket := config.NewBoltBucket[CertModel](db, "Cert")
	tagBucket := config.NewBoltBucket[[]string](db, "Tag")
	secretBucket := config.NewBoltBucket[string](db, "Secret")
	passwordBucket := config.NewBoltBucket[Credential](db, "Password")
	queryBucket := config.NewBoltBucket[Query](db, "Query")

	compactTableFiles := config.NewFilePersistence("pgcompacttable", ".log")
	certFiles := config.NewFilePersistence("cert", ".crt")

	clusterRepo := persistence.NewClusterRepository(clusterBucket)
	compactTableRepo := persistence.NewCompactTableRepository(compactTableBucket, compactTableFiles)
	certRepo := persistence.NewCertRepository(certBucket, certFiles)
	tagRepo := persistence.NewTagRepository(tagBucket)
	secretRepo := persistence.NewSecretRepository(secretBucket)
	passwordRepo := persistence.NewPasswordRepository(passwordBucket)
	queryRepo := persistence.NewQueryRepository(queryBucket)

	encryption := service.NewEncryption()
	proxy := service.NewProxy(clusterRepo, passwordRepo, certRepo, certFiles)

	secretService := service.NewSecretService(secretRepo, encryption)
	passwordService := service.NewPasswordService(passwordRepo, secretService, encryption)
	queryService := service.NewQueryService(queryRepo, clusterRepo, passwordService, secretService)
	bloatService := service.NewBloatService(compactTableRepo, passwordRepo, compactTableFiles, secretService, encryption)
	patroniService := service.NewPatroniService(proxy)
	eraseService := service.NewEraseService(passwordRepo, clusterRepo, certRepo, tagRepo, compactTableRepo, queryService, secretService)

	// TODO refactor: router shouldn't use repos, change to service usage
	return &Context{
		env:            env,
		infoRouter:     router.NewInfoRouter(env, secretService),
		clusterRouter:  router.NewClusterRouter(clusterRepo, tagRepo),
		bloatRouter:    router.NewBloatRouter(bloatService, compactTableRepo),
		certRouter:     router.NewCertRouter(certRepo),
		secretRouter:   router.NewSecretRouter(secretService, passwordService),
		passwordRouter: router.NewPasswordRouter(passwordService),
		tagRouter:      router.NewTagRouter(tagRepo),
		instanceRouter: router.NewInstanceRouter(patroniService),
		queryRouter:    router.NewQueryRouter(queryService),
		eraseRouter:    router.NewEraseRouter(eraseService),
	}
}
