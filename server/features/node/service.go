package node

import (
	"crypto/tls"
	"ivory/core/config"
	"ivory/core/service/cert"
	"ivory/core/service/job"
	"ivory/core/service/vault"
	"ivory/core/utils"
	"ivory/plugins/keeper"
	"ivory/plugins/platform"
)

type Service struct {
	platformRegistry *utils.Registry[platform.Plugin, platform.Adapter]
	keeperRegistry   *utils.Registry[keeper.Plugin, keeper.Adapter]
	vaultService     *vault.Service
	certService      *cert.Service
	jobManager       *job.Service

	dbFeatures map[env.Feature]bool
}

func NewService(
	platformRegistry *utils.Registry[platform.Plugin, platform.Adapter],
	keeperRegistry *utils.Registry[keeper.Plugin, keeper.Adapter],
	vaultService *vault.Service,
	certService *cert.Service,
	jobManager *job.Service,
) *Service {
	return &Service{
		platformRegistry: platformRegistry,
		keeperRegistry:   keeperRegistry,
		vaultService:     vaultService,
		certService:      certService,
		jobManager:       jobManager,

		dbFeatures: make(map[env.Feature]bool),
	}
}

func (s *Service) SupportedFeatures(t keeper.Plugin) []env.Feature {
	c, e := s.keeperRegistry.Get(t)
	if e != nil {
		return []env.Feature{}
	}
	return c.SupportedFeatures()
}

func (s *Service) getPlatformAdapter(c PlatformVaultConnection) (platform.Adapter, platform.Connection, error) {
	adapter, err := s.platformRegistry.Get(platform.Linux)
	if err != nil {
		return nil, platform.Connection{}, err
	}
	platformConn, err := s.getPlatformVaultConnection(c)
	if err != nil {
		return nil, platform.Connection{}, err
	}
	return adapter, platformConn, err
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

func (s *Service) getPlatformVaultConnection(c PlatformVaultConnection) (platform.Connection, error) {
	cred, err := s.vaultService.GetDecrypted(c.VaultId)
	if err != nil {
		return platform.Connection{}, err
	}
	return platform.Connection{
		Host:       c.Host,
		Port:       c.Port,
		Username:   cred.Username,
		PrivateKey: []byte(cred.Secret),
	}, nil
}

func (s *Service) getPlatformCredConnection(c PlatformCredConnection) platform.Connection {
	return platform.Connection{
		Host:     c.Host,
		Port:     c.Port,
		Username: c.Username,
		Password: c.Password,
	}
}
