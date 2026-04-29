package node

import (
	"crypto/ed25519"
	"errors"
	"ivory/src/clients/ssh"
	"ivory/src/features"
	"ivory/src/features/cert"
	"ivory/src/features/vault"
	"ivory/src/plugins/keeper"
	"ivory/src/plugins/os"
)

var ErrSshKeyNotSpecified = errors.New("ssh key is not specified")

type Service struct {
	osRegistry     *os.PluginRegistry
	keeperRegistry *keeper.PluginRegistry
	vaultService   *vault.Service
	certService    *cert.Service

	dbFeatures map[features.Feature]bool
}

func NewService(
	osRegistry *os.PluginRegistry,
	keeperRegistry *keeper.PluginRegistry,
	vaultService *vault.Service,
	certService *cert.Service,
) *Service {
	return &Service{
		osRegistry:     osRegistry,
		keeperRegistry: keeperRegistry,
		vaultService:   vaultService,
		certService:    certService,

		dbFeatures: make(map[features.Feature]bool),
	}
}

func (s *Service) getOSAdapter(connection Connection) (os.Adapter, *ssh.Connection, error) {
	sshConn, err := s.getSshConnection(connection)
	if err != nil {
		return nil, nil, err
	}
	adapter, err := s.osRegistry.Get(os.Linux)
	return adapter, sshConn, err
}

func (s *Service) getKeeperAdapter(r KeeperRequest) (keeper.Adapter, keeper.Request, error) {
	client, errClient := s.keeperRegistry.Get(r.Type)
	if errClient != nil {
		return nil, keeper.Request{}, errClient
	}
	request, err := s.getKeeperRequest(r)
	return client, request, err
}

func (s *Service) SupportedFeatures(t keeper.Type) []features.Feature {
	c, e := s.keeperRegistry.Get(t)
	if e != nil {
		return []features.Feature{}
	}
	return c.SupportedFeatures()
}

func (s *Service) getKeeperRequest(request KeeperRequest) (keeper.Request, error) {
	keeperRequest := keeper.Request{Host: request.Host, Port: request.Port, Body: request.Body}
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
		Host:       connection.Host,
		Port:       sshPort,
		Username:   cred.Username,
		PrivateKey: ed25519.PrivateKey(cred.Secret),
	}, nil
}
