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

	compactTableFiles := config.NewFilePersistence("pgcompacttable", ".log")
	certFiles := config.NewFilePersistence("cert", ".crt")

	clusterRepo := persistence.NewClusterRepository(clusterBucket)
	compactTableRepo := persistence.NewCompactTableRepository(compactTableBucket, compactTableFiles)
	certRepo := persistence.NewCertRepository(certBucket, certFiles)
	tagRepo := persistence.NewTagRepository(tagBucket)
	secretRepo := persistence.NewSecretRepository(secretBucket)
	passwordRepo := persistence.NewPasswordRepository(passwordBucket)

	encryption := service.NewEncryption()
	proxy := service.NewProxy(clusterRepo, passwordRepo, certRepo, certFiles)
	secretService := service.NewSecretService(passwordRepo, secretRepo, clusterRepo, certRepo, tagRepo, compactTableRepo, encryption)
	bloatService := service.NewBloatService(compactTableRepo, passwordRepo, compactTableFiles, secretService, encryption)
	patroniService := service.NewPatroniService(proxy)

	return &Context{
		env:            env,
		infoRouter:     router.NewInfoRouter(env, secretService),
		clusterRouter:  router.NewClusterRouter(clusterRepo, tagRepo),
		bloatRouter:    router.NewBloatRouter(bloatService, compactTableRepo),
		certRouter:     router.NewCertRouter(certRepo),
		secretRouter:   router.NewSecretRouter(secretService),
		passwordRouter: router.NewPasswordRouter(passwordRepo, secretService, encryption),
		tagRouter:      router.NewTagRouter(tagRepo),
		instanceRouter: router.NewInstanceRouter(patroniService),
	}
}
