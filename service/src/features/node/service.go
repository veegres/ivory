package node

import (
	"errors"
	"ivory/src/clients/keeper"
	"ivory/src/clients/ssh"
	"ivory/src/features"
	"ivory/src/features/cert"
	"ivory/src/features/vault"
)

var ErrSshKeyNotSpecified = errors.New("ssh key is not specified")

type Service struct {
	keeperRegistry *keeper.Registry
	sshClient      ssh.Client
	vaultService   *vault.Service
	certService    *cert.Service

	dbFeatures map[features.Feature]bool
}

func NewService(
	keeperRegistry *keeper.Registry,
	sshClient ssh.Client,
	vaultService *vault.Service,
	certService *cert.Service,
) *Service {
	return &Service{
		keeperRegistry: keeperRegistry,
		sshClient:      sshClient,
		vaultService:   vaultService,
		certService:    certService,

		dbFeatures: make(map[features.Feature]bool),
	}
}

func (s *Service) SupportedFeatures(t keeper.Type) []features.Feature {
	c, e := s.keeperRegistry.Get(t)
	if e != nil {
		return []features.Feature{}
	}
	return c.SupportedFeatures()
}

func (s *Service) getKeeperRequest(request KeeperRequest) (keeper.Request, error) {
	keeperRequest := keeper.Request{Host: request.Connection.Host, Port: request.Connection.KeeperPort, Body: request.Body}
	if request.Certs != nil {
		err := s.certService.EnrichTLSConfig(&keeperRequest.TlsConfig, request.Certs)
		if err != nil {
			return keeperRequest, err
		}
	}
	if request.VaultId != nil {
		cred, err := s.vaultService.GetDecrypted(*request.VaultId)
		if err != nil {
			return keeperRequest, err
		}
		keeperRequest.Credentials = &keeper.Credentials{Username: cred.Username, Password: cred.Secret}
	}
	return keeperRequest, nil
}

func (s *Service) getSshConnection(connection Connection) (*ssh.Connection, error) {
	if connection.SshKeyId == nil {
		return nil, ErrSshKeyNotSpecified
	}
	cred, err := s.vaultService.GetDecrypted(*connection.SshKeyId)
	if err != nil {
		return nil, err
	}

	sshPort := connection.SshPort
	if sshPort == 0 {
		sshPort = 22
	}

	return &ssh.Connection{
		Host:     connection.Host,
		Port:     sshPort,
		Username: cred.Username,
		Key:      cred.Secret,
	}, nil
}
