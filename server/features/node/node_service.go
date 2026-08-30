package node

import (
	"crypto/tls"
	"ivory/core/config"
	"ivory/core/service/cert"
	"ivory/core/service/job"
	"ivory/core/service/vault"
	"ivory/plugins/keeper"
	"ivory/plugins/platform"
)

// platformRegistry is the narrow view node needs of the platform plugin
// registry: resolving one plugin by key, never registering.
type platformRegistry interface {
	Get(platform.PluginType) (platform.Plugin, error)
}

// keeperRegistry is the same narrow view of the keeper plugin registry.
type keeperRegistry interface {
	Get(keeper.PluginType) (keeper.Plugin, error)
}

type Service struct {
	platformRegistry platformRegistry
	keeperRegistry   keeperRegistry
	vaultService     *vault.Service
	certService      *cert.Service
	jobManager       *job.Service

	dbFeatures map[config.Feature]bool
}

func NewService(
	platformRegistry platformRegistry,
	keeperRegistry keeperRegistry,
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

		dbFeatures: make(map[config.Feature]bool),
	}
}

func (s *Service) SupportedFeatures(t keeper.PluginType) map[config.Feature]bool {
	c, e := s.keeperRegistry.Get(t)
	if e != nil {
		return map[config.Feature]bool{}
	}
	return c.SupportedFeatures()
}

// KeeperHasLeader reports whether the keeper elects a single primary at all.
// An unknown plugin is treated as electing one: that is what every keeper but
// clickhouse does, and it keeps a warning rather than silently dropping it.
func (s *Service) KeeperHasLeader(t keeper.PluginType) bool {
	c, e := s.keeperRegistry.Get(t)
	if e != nil {
		return true
	}
	return c.HasLeader()
}

// PlatformSupportedFeatures reports what the platform itself can be asked
// for, which is a separate question from what the keeper supports: a platform
// that only ever addresses a scheduler has no node of its own to show.
func (s *Service) PlatformSupportedFeatures(p PlatformPlugin) map[config.Feature]bool {
	c, e := s.platformRegistry.Get(p)
	if e != nil {
		return map[config.Feature]bool{}
	}
	return c.SupportedFeatures()
}

func (s *Service) getPlatformAdapter(c PlatformVaultConnection) (platform.Adapter, platform.Connection, error) {
	adapter, err := s.platformRegistry.Get(c.PlatformOrDefault())
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
