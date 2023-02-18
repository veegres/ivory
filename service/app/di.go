package app

import (
	"ivory/config"
	"ivory/model"
	"ivory/persistence"
	"ivory/router"
	"ivory/service"
)

type DI struct {
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

func NewDI() *DI {
	env := config.NewEnv()

	db := config.NewBoltDB("ivory.db")
	clusterBucket := config.NewBoltBucket[model.ClusterModel](db, "Cluster")
	compactTableBucket := config.NewBoltBucket[model.CompactTableModel](db, "CompactTable")
	certBucket := config.NewBoltBucket[model.CertModel](db, "Cert")
	tagBucket := config.NewBoltBucket[[]string](db, "Tag")
	secretBucket := config.NewBoltBucket[string](db, "Secret")
	passwordBucket := config.NewBoltBucket[model.Credential](db, "Password")

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

	return &DI{
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
