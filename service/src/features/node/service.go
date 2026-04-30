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

func (s *Service) SupportedFeatures(t keeper.Plugin) []features.Feature {
	c, e := s.keeperRegistry.Get(t)
	if e != nil {
		return []features.Feature{}
	}
	return c.SupportedFeatures()
}

func (s *Service) getOSAdapter(c SshConnection) (os.Adapter, *ssh.Connection, error) {
	sshConn, err := s.getSshConnection(c)
	if err != nil {
		return nil, nil, err
	}
	adapter, err := s.osRegistry.Get(os.Linux)
	return adapter, sshConn, err
}

func (s *Service) getKeeperAdapter(c KeeperConnection) (keeper.Adapter, keeper.Request, error) {
	client, errClient := s.keeperRegistry.Get(c.Plugin)
	if errClient != nil {
		return nil, keeper.Request{}, errClient
	}
	request, err := s.getKeeperRequest(c)
	return client, request, err
}

func (s *Service) getKeeperRequest(c KeeperConnection) (keeper.Request, error) {
	keeperRequest := keeper.Request{Host: c.Host, Port: c.Port, Body: c.Body}
	if c.Certs != nil {
		err := s.certService.EnrichTLSConfig(&keeperRequest.TlsConfig, c.Certs)
		if err != nil {
			return keeperRequest, err
		}
	}
	if c.VaultId != nil {
		cred, err := s.vaultService.GetDecrypted(*c.VaultId)
		if err != nil {
			return keeperRequest, err
		}
		keeperRequest.Credentials = &keeper.Credentials{Username: cred.Username, Password: cred.Secret}
	}
	return keeperRequest, nil
}

func (s *Service) getSshConnection(c SshConnection) (*ssh.Connection, error) {
	if c.VaultId == nil {
		return nil, ErrSshKeyNotSpecified
	}
	cred, err := s.vaultService.GetDecrypted(*c.VaultId)
	if err != nil {
		return nil, err
	}

	return &ssh.Connection{
		Host:       c.Host,
		Port:       c.Port,
		Username:   cred.Username,
		PrivateKey: ed25519.PrivateKey(cred.Secret),
	}, nil
}
