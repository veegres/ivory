package node

import (
	"ivory/src/clients/keeper"
	"ivory/src/clients/ssh"
	"ivory/src/features/cert"
	"ivory/src/features/password"
	"ivory/src/features/vm"

	"github.com/google/uuid"
)

type Service struct {
	keeperClient    keeper.Client
	sshClient       ssh.Client
	passwordService *password.Service
	certService     *cert.Service
	vmService       *vm.Service
}

func NewService(
	keeperClient keeper.Client,
	sshClient ssh.Client,
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

func (s *Service) getKeeperRequest(request Request) (keeper.Request, error) {
	keeperRequest := keeper.Request{Keeper: request.Keeper, Body: request.Body}
	if request.Certs != nil {
		err := s.certService.EnrichTLSConfig(&keeperRequest.TlsConfig, request.Certs)
		if err != nil {
			return keeperRequest, err
		}
	}
	if request.CredentialId != nil {
		pass, err := s.passwordService.GetDecrypted(*request.CredentialId)
		if err != nil {
			return keeperRequest, err
		}
		keeperRequest.Credentials = &keeper.Credentials{Username: pass.Username, Password: pass.Password}
	}
	return keeperRequest, nil
}

func (s *Service) getSshConnection(id uuid.UUID) (*ssh.Connection, error) {
	connection, err := s.vmService.GetDecrypted(id)
	if err != nil {
		return nil, err
	}
	return &ssh.Connection{
		Host:     connection.Host,
		Port:     connection.Port,
		Username: connection.Username,
		Key:      connection.SshKey,
	}, nil
}
