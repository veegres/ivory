package plugins

import (
	"ivory/clients/console/ssh"
	"ivory/clients/http"
	"ivory/core/utils"
	"ivory/plugins/database"
	"ivory/plugins/database/postgres"
	"ivory/plugins/keeper"
	"ivory/plugins/keeper/patroni"
	pgkeeper "ivory/plugins/keeper/postgres"
	"ivory/plugins/platform"
	"ivory/plugins/platform/linux"
)

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
	postgresAdapter := postgres.NewAdapter()
	linuxAdapter := linux.NewAdapter(sshClient)

	// REGISTRY
	keeperRegistry := utils.NewRegistry[keeper.Plugin, keeper.Adapter]()
	keeperRegistry.Register(keeper.PATRONI, patroniAdapter)
	keeperMetadataRegistry := utils.NewRegistry[keeper.Plugin, keeper.Metadata]()
	keeperMetadataRegistry.Register(keeper.PATRONI, patroniAdapter)
	keeperMetadataRegistry.Register(keeper.POSTGRES, pgKeeperAdapter)
	databaseRegistry := utils.NewRegistry[database.Plugin, database.Adapter]()
	databaseRegistry.Register(database.POSTGRES, postgresAdapter)
	platformRegistry := utils.NewRegistry[platform.Plugin, platform.Adapter]()
	platformRegistry.Register(platform.Linux, linuxAdapter)

	return &Context{
		KeeperRegistry:         keeperRegistry,
		KeeperMetadataRegistry: keeperMetadataRegistry,
		DatabaseRegistry:       databaseRegistry,
		PlatformRegistry:       platformRegistry,
	}
}
