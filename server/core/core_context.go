package core

import (
	"ivory/clients/console/ssh"
	"ivory/clients/storage"
	"ivory/core/service/cert"
	"ivory/core/service/encryption"
	"ivory/core/service/job"
	"ivory/core/service/secret"
	"ivory/core/service/vault"
)

type Service struct {
	Vault *vault.Service
	Job   *job.Service
	Cert  *cert.Service
}

type HttpRouter struct {
	Vault  *vault.HttpRouter
	Cert   *cert.HttpRouter
	Secret *secret.HttpRouter
}

type Context struct {
	Service    *Service
	HttpRouter *HttpRouter
}

// NewContext TODO probably we should move User, Permission and Auth to this folder it looks like core more - think
func NewContext(sshClient *ssh.Client) *Context {
	st := storage.NewDbStorage("core.db")

	secretBucket := storage.NewDbBucket[string](st, "Secret")
	vaultBucket := storage.NewDbBucket[vault.Vault](st, "Vault")
	certBucket := storage.NewDbBucket[cert.Cert](st, "Cert")

	certFiles := storage.NewFileStorage("cert", ".crt")
	jobFiles := storage.NewFileStorage("job", ".log")

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
		HttpRouter: &HttpRouter{
			Vault:  vault.NewHttpRouter(vaultService),
			Cert:   cert.NewHttpRouter(certService),
			Secret: secret.NewHttpRouter(secretService),
		},
	}
}
