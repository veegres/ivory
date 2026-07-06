// Package plugins wires up every keeper/database/platform plugin adapter
// into the registries the rest of Ivory looks them up from by Plugin type.
package plugins

import (
	"ivory/clients/console/ssh"
	"ivory/clients/http"
	"ivory/core/utils"
	"ivory/plugins/database"
	etcddb "ivory/plugins/database/etcd"
	"ivory/plugins/database/postgres"
	"ivory/plugins/keeper"
	etcdkeeper "ivory/plugins/keeper/etcd"
	"ivory/plugins/keeper/patroni"
	pgkeeper "ivory/plugins/keeper/postgres"
	"ivory/plugins/platform"
	"ivory/plugins/platform/linux"
)

// Context holds every plugin registry, keyed by each plugin family's own
// Plugin type, that the rest of Ivory resolves adapters through.
type Context struct {
	KeeperRegistry         *utils.Registry[keeper.Plugin, keeper.Adapter]
	KeeperMetadataRegistry *utils.Registry[keeper.Plugin, keeper.Metadata]
	DatabaseRegistry       *utils.Registry[database.Plugin, database.Adapter]
	PlatformRegistry       *utils.Registry[platform.Plugin, platform.Adapter]
}

func NewContext(httpClient *http.Client, sshClient *ssh.Client) *Context {
	// ADAPTERS
	patroniAdapter := patroni.NewAdapter(httpClient)
	pgKeeperAdapter := pgkeeper.NewAdapter()
	etcdKeeperAdapter := etcdkeeper.NewAdapter()
	postgresAdapter := postgres.NewAdapter()
	etcdDbAdapter := etcddb.NewAdapter()
	linuxAdapter := linux.NewAdapter(sshClient)

	// REGISTRY
	keeperRegistry := utils.NewRegistry[keeper.Plugin, keeper.Adapter]()
	keeperRegistry.Register(keeper.PATRONI_POSTGRES, patroniAdapter)
	keeperRegistry.Register(keeper.NATIVE_POSTGRES, pgKeeperAdapter)
	keeperRegistry.Register(keeper.NATIVE_ETCD, etcdKeeperAdapter)
	keeperMetadataRegistry := utils.NewRegistry[keeper.Plugin, keeper.Metadata]()
	keeperMetadataRegistry.Register(keeper.PATRONI_POSTGRES, patroniAdapter)
	keeperMetadataRegistry.Register(keeper.NATIVE_POSTGRES, pgKeeperAdapter)
	keeperMetadataRegistry.Register(keeper.NATIVE_ETCD, etcdKeeperAdapter)
	databaseRegistry := utils.NewRegistry[database.Plugin, database.Adapter]()
	databaseRegistry.Register(database.POSTGRES, postgresAdapter)
	databaseRegistry.Register(database.ETCD, etcdDbAdapter)
	platformRegistry := utils.NewRegistry[platform.Plugin, platform.Adapter]()
	platformRegistry.Register(platform.Linux, linuxAdapter)

	return &Context{
		KeeperRegistry:         keeperRegistry,
		KeeperMetadataRegistry: keeperMetadataRegistry,
		DatabaseRegistry:       databaseRegistry,
		PlatformRegistry:       platformRegistry,
	}
}
