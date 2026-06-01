package plugins

import (
	"ivory/clients/console/ssh"
	"ivory/clients/http"
	"ivory/core/utils"
	"ivory/plugins/database"
	"ivory/plugins/database/postgres"
	"ivory/plugins/keeper"
	"ivory/plugins/keeper/patroni"
	"ivory/plugins/platform"
	"ivory/plugins/platform/onprem"
)

type Context struct {
	KeeperRegistry   *utils.Registry[keeper.Plugin, keeper.Adapter]
	DatabaseRegistry *utils.Registry[database.Plugin, database.Adapter]
	PlatformRegistry *utils.Registry[platform.Plugin, platform.Adapter]
}

func NewContext(httpClient *http.Client, sshClient *ssh.Client) *Context {
	// ADAPTERS
	patroniAdapter := patroni.NewAdapter(httpClient)
	postgresAdapter := postgres.NewAdapter()
	onpremAdapter := onprem.NewAdapter(sshClient)

	// REGISTRY
	keeperRegistry := utils.NewRegistry[keeper.Plugin, keeper.Adapter]()
	keeperRegistry.Register(keeper.PATRONI, patroniAdapter)
	databaseRegistry := utils.NewRegistry[database.Plugin, database.Adapter]()
	databaseRegistry.Register(database.POSTGRES, postgresAdapter)
	platformRegistry := utils.NewRegistry[platform.Plugin, platform.Adapter]()
	platformRegistry.Register(platform.Onprem, onpremAdapter)

	return &Context{
		KeeperRegistry:   keeperRegistry,
		DatabaseRegistry: databaseRegistry,
		PlatformRegistry: platformRegistry,
	}
}
