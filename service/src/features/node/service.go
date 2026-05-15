package node

import (
	"crypto/ed25519"
	"crypto/tls"
	"errors"
	"ivory/src/clients/ssh"
	"ivory/src/features"
	"ivory/src/features/cert"
	"ivory/src/features/vault"
	"ivory/src/plugins/cloud"
	"ivory/src/plugins/keeper"
)

var ErrSshKeyNotSpecified = errors.New("ssh key is not specified")

type Service struct {
	cloudRegistry  *cloud.PluginRegistry
	keeperRegistry *keeper.PluginRegistry
	vaultService   *vault.Service
	certService    *cert.Service

	dbFeatures map[features.Feature]bool
}

func NewService(
	cloudRegistry *cloud.PluginRegistry,
	keeperRegistry *keeper.PluginRegistry,
	vaultService *vault.Service,
	certService *cert.Service,
) *Service {
	return &Service{
		cloudRegistry:  cloudRegistry,
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

func (s *Service) getCloudAdapter(c CloudVaultConnection) (cloud.Adapter, *ssh.Connection, error) {
	adapter, err := s.cloudRegistry.Get(cloud.Onprem)
	if err != nil {
		return nil, nil, err
	}
	cloudConn, err := s.getCloudVaultConnection(c)
	if err != nil {
		return nil, nil, err
	}
	return adapter, cloudConn, err
}

func (s *Service) getKeeperAdapter(c KeeperOptions) (keeper.Adapter, *tls.Config, *keeper.Credentials, error) {
	adapter, err := s.keeperRegistry.Get(c.Plugin)
	if err != nil {
		return nil, nil, nil, err
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
	return adapter, tlsConfig, cred, nil
}

func (s *Service) getCloudVaultConnection(c CloudVaultConnection) (*ssh.Connection, error) {
	if c.VaultId == nil {
		return nil, ErrSshKeyNotSpecified
	}
	cred, err := s.vaultService.GetDecrypted(*c.VaultId)
	if err != nil {
		return nil, err
	}
	prvKey := ed25519.PrivateKey(cred.Secret)
	return &ssh.Connection{
		Host:       c.Host,
		Port:       c.Port,
		Username:   cred.Username,
		PrivateKey: &prvKey,
	}, nil
}

func (s *Service) getCloudCredConnection(c CloudCredConnection) *ssh.Connection {
	return &ssh.Connection{
		Host:     c.Host,
		Port:     c.Port,
		Username: c.Username,
		Password: c.Password,
	}
}
