package core

import (
	"ivory/clients/console/ssh"
	"ivory/core/service/cert"
	"ivory/core/service/encryption"
	"ivory/core/service/job"
	"ivory/core/service/secret"
	"ivory/core/service/vault"
	"ivory/core/store"
)

type Service struct {
	Vault *vault.Service
	Job   *job.Service
	Cert  *cert.Service
}

type Router struct {
	Vault  *vault.Router
	Cert   *cert.Router
	Secret *secret.Router
}

type Context struct {
	Service *Service
	Router  *Router
}

func NewContext(sshClient *ssh.Client) *Context {
	st := store.NewDbStorage("core.db")

	secretBucket := store.NewDbBucket[string](st, "Secret")
	vaultBucket := store.NewDbBucket[vault.Vault](st, "Vault")
	certBucket := store.NewDbBucket[cert.Cert](st, "Cert")

	certFiles := store.NewFileStorage("cert", ".crt")
	jobFiles := store.NewFileStorage("job", ".log")

	secretRepo := secret.NewRepository(secretBucket)
	vaultRepo := vault.NewRepository(vaultBucket)
	certRepo := cert.NewRepository(certBucket, certFiles)

	encryptionService := encryption.NewService()
	secretService := secret.NewService(secretRepo, encryptionService)
	vaultService := vault.NewService(vaultRepo, sshClient, secretService, encryptionService)
	jobService := job.NewService(jobFiles)
	certService := cert.NewService(certRepo)

	return &Context{
		Service: &Service{
			Vault: vaultService,
			Job:   jobService,
			Cert:  certService,
		},
		Router: &Router{
			Vault:  vault.NewRouter(vaultService),
			Cert:   cert.NewRouter(certService),
			Secret: secret.NewRouter(secretService),
		},
	}
}
