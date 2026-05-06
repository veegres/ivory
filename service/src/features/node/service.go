package node

import (
	"crypto/ed25519"
	"crypto/tls"
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
	adapter, err := s.osRegistry.Get(os.Linux)
	sshConn, err := s.getSshConnection(c)
	if err != nil {
		return nil, nil, err
	}
	return adapter, sshConn, err
}

func (s *Service) getKeeperAdapter(c KeeperOptions) (keeper.Adapter, *tls.Config, *keeper.Credentials, error) {
	client, errClient := s.keeperRegistry.Get(c.Plugin)
	if errClient != nil {
		return nil, nil, nil, errClient
	}
	var tlsConfig *tls.Config
	if c.Certs != nil {
		err := s.certService.EnrichTLSConfig(&tlsConfig, c.Certs)
		if err != nil {
			return nil, nil, nil, err
		}
	}
	var cred *keeper.Credentials
	if c.VaultId != nil {
		d, err := s.vaultService.GetDecrypted(*c.VaultId)
		if err != nil {
			return nil, nil, nil, err
		}
		cred = &keeper.Credentials{Username: d.Username, Password: d.Secret}
	}
	return client, tlsConfig, cred, nil
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
