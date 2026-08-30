// Package plugins wires up every keeper/database/platform plugin adapter
// into the registries the rest of Ivory looks them up from by Plugin type.
package plugins

import (
	"ivory/clients/console/ssh"
	"ivory/clients/http"
	"ivory/core/utils"
	"ivory/plugins/database"
	chdb "ivory/plugins/database/clickhouse"
	etcddb "ivory/plugins/database/etcd"
	mongodb "ivory/plugins/database/mongo"
	"ivory/plugins/database/postgres"
	redisdb "ivory/plugins/database/redis"
	zkdb "ivory/plugins/database/zookeeper"
	"ivory/plugins/keeper"
	chkeeper "ivory/plugins/keeper/clickhouse"
	etcdkeeper "ivory/plugins/keeper/etcd"
	mongokeeper "ivory/plugins/keeper/mongo"
	"ivory/plugins/keeper/patroni"
	pgkeeper "ivory/plugins/keeper/postgres"
	rediskeeper "ivory/plugins/keeper/redis"
	zkkeeper "ivory/plugins/keeper/zookeeper"
	"ivory/plugins/platform"
	"ivory/plugins/platform/docker"
)

// Context holds every plugin registry, keyed by each plugin family's own
// PluginType, that the rest of Ivory resolves plugins through. A plugin is
// registered once: keeper.Plugin and platform.Plugin carry both the adapter
// and the metadata half, and consumers depend on the narrower view they need.
type Context struct {
	KeeperRegistry   *utils.Registry[keeper.PluginType, keeper.Plugin]
	DatabaseRegistry *utils.Registry[database.PluginType, database.Adapter]
	PlatformRegistry *utils.Registry[platform.PluginType, platform.Plugin]
}

func NewContext(httpClient *http.Client, sshClient *ssh.Client) *Context {
	// ADAPTERS
	patroniPlugin := patroni.NewPlugin(httpClient)
	pgKeeperPlugin := pgkeeper.NewPlugin()
	etcdKeeperPlugin := etcdkeeper.NewPlugin()
	redisKeeperPlugin := rediskeeper.NewPlugin()
	clickhouseKeeperPlugin := chkeeper.NewPlugin()
	zookeeperKeeperPlugin := zkkeeper.NewPlugin()
	mongoKeeperPlugin := mongokeeper.NewPlugin()
	postgresAdapter := postgres.NewAdapter()
	etcdDbAdapter := etcddb.NewAdapter()
	redisDbAdapter := redisdb.NewAdapter()
	clickhouseDbAdapter := chdb.NewAdapter()
	zookeeperDbAdapter := zkdb.NewAdapter()
	mongoDbAdapter := mongodb.NewAdapter()
	dockerPlugin := docker.NewPlugin(sshClient)

	// REGISTRY
	keeperRegistry := utils.NewRegistry[keeper.PluginType, keeper.Plugin]()
	keeperRegistry.Register(keeper.PATRONI_POSTGRES, patroniPlugin)
	keeperRegistry.Register(keeper.NATIVE_POSTGRES, pgKeeperPlugin)
	keeperRegistry.Register(keeper.NATIVE_ETCD, etcdKeeperPlugin)
	keeperRegistry.Register(keeper.NATIVE_REDIS, redisKeeperPlugin)
	keeperRegistry.Register(keeper.NATIVE_CLICKHOUSE, clickhouseKeeperPlugin)
	keeperRegistry.Register(keeper.NATIVE_ZOOKEEPER, zookeeperKeeperPlugin)
	keeperRegistry.Register(keeper.NATIVE_MONGO, mongoKeeperPlugin)
	databaseRegistry := utils.NewRegistry[database.PluginType, database.Adapter]()
	databaseRegistry.Register(database.POSTGRES, postgresAdapter)
	databaseRegistry.Register(database.ETCD, etcdDbAdapter)
	databaseRegistry.Register(database.REDIS, redisDbAdapter)
	databaseRegistry.Register(database.CLICKHOUSE, clickhouseDbAdapter)
	databaseRegistry.Register(database.ZOOKEEPER, zookeeperDbAdapter)
	databaseRegistry.Register(database.MONGO, mongoDbAdapter)
	platformRegistry := utils.NewRegistry[platform.PluginType, platform.Plugin]()
	platformRegistry.Register(platform.Docker, dockerPlugin)

	return &Context{
		KeeperRegistry:   keeperRegistry,
		DatabaseRegistry: databaseRegistry,
		PlatformRegistry: platformRegistry,
	}
}
