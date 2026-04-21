package node

import (
	"ivory/src/clients/sidecar"
	sshclient "ivory/src/clients/ssh"
	"ivory/src/features/cert"
	"ivory/src/features/password"
	"ivory/src/features/vm"
)

type Service struct {
	sidecarClient   sidecar.Client
	sshClient       sshclient.Client
	passwordService *password.Service
	certService     *cert.Service
	vmService       *vm.Service
}

func NewService(
	sidecarClient sidecar.Client,
	sshClient sshclient.Client,
	passwordService *password.Service,
	certService *cert.Service,
	vmService *vm.Service,
) *Service {
	return &Service{
		sidecarClient:   sidecarClient,
		sshClient:       sshClient,
		passwordService: passwordService,
		certService:     certService,
		vmService:       vmService,
	}
}

func (s *Service) mapRequest(instance Request) (sidecar.Request, error) {
	request := sidecar.Request{Sidecar: instance.Sidecar, Body: instance.Body}
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
		request.Credentials = &sidecar.Credentials{Username: pass.Username, Password: pass.Password}
	}
	return request, nil
}
