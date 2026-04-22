package node

import (
	"ivory/src/clients/keeper"
	sshclient "ivory/src/clients/ssh"
	"ivory/src/features/cert"
	"ivory/src/features/password"
	"ivory/src/features/vm"
)

type Service struct {
	keeperClient    keeper.Client
	sshClient       sshclient.Client
	passwordService *password.Service
	certService     *cert.Service
	vmService       *vm.Service
}

func NewService(
	keeperClient keeper.Client,
	sshClient sshclient.Client,
	passwordService *password.Service,
	certService *cert.Service,
	vmService *vm.Service,
) *Service {
	return &Service{
		keeperClient:    keeperClient,
		sshClient:       sshClient,
		passwordService: passwordService,
		certService:     certService,
		vmService:       vmService,
	}
}

func (s *Service) mapRequest(instance Request) (keeper.Request, error) {
	request := keeper.Request{Keeper: instance.Keeper, Body: instance.Body}
	if instance.Certs != nil {
		err := s.certService.EnrichTLSConfig(&request.TlsConfig, instance.Certs)
		if err != nil {
			return request, err
		}
	}
	if instance.CredentialId != nil {
		pass, err := s.passwordService.GetDecrypted(*instance.CredentialId)
		if err != nil {
			return request, err
		}
		request.Credentials = &keeper.Credentials{Username: pass.Username, Password: pass.Password}
	}
	return request, nil
}
