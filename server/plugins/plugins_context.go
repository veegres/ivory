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
// Plugin type, that the rest of Ivory resolves adapters through.
type Context struct {
	KeeperRegistry           *utils.Registry[keeper.Plugin, keeper.Adapter]
	KeeperMetadataRegistry   *utils.Registry[keeper.Plugin, keeper.Metadata]
	DatabaseRegistry         *utils.Registry[database.Plugin, database.Adapter]
	PlatformRegistry         *utils.Registry[platform.Plugin, platform.Adapter]
	PlatformMetadataRegistry *utils.Registry[platform.Plugin, platform.Metadata]
}

func NewContext(httpClient *http.Client, sshClient *ssh.Client) *Context {
	// ADAPTERS
	patroniAdapter := patroni.NewAdapter(httpClient)
	pgKeeperAdapter := pgkeeper.NewAdapter()
	etcdKeeperAdapter := etcdkeeper.NewAdapter()
	redisKeeperAdapter := rediskeeper.NewAdapter()
	clickhouseKeeperAdapter := chkeeper.NewAdapter()
	zookeeperKeeperAdapter := zkkeeper.NewAdapter()
	mongoKeeperAdapter := mongokeeper.NewAdapter()
	postgresAdapter := postgres.NewAdapter()
	etcdDbAdapter := etcddb.NewAdapter()
	redisDbAdapter := redisdb.NewAdapter()
	clickhouseDbAdapter := chdb.NewAdapter()
	zookeeperDbAdapter := zkdb.NewAdapter()
	mongoDbAdapter := mongodb.NewAdapter()
	dockerAdapter := docker.NewAdapter(sshClient)

	// REGISTRY
	keeperRegistry := utils.NewRegistry[keeper.Plugin, keeper.Adapter]()
	keeperRegistry.Register(keeper.PATRONI_POSTGRES, patroniAdapter)
	keeperRegistry.Register(keeper.NATIVE_POSTGRES, pgKeeperAdapter)
	keeperRegistry.Register(keeper.NATIVE_ETCD, etcdKeeperAdapter)
	keeperRegistry.Register(keeper.NATIVE_REDIS, redisKeeperAdapter)
	keeperRegistry.Register(keeper.NATIVE_CLICKHOUSE, clickhouseKeeperAdapter)
	keeperRegistry.Register(keeper.NATIVE_ZOOKEEPER, zookeeperKeeperAdapter)
	keeperRegistry.Register(keeper.NATIVE_MONGO, mongoKeeperAdapter)
	keeperMetadataRegistry := utils.NewRegistry[keeper.Plugin, keeper.Metadata]()
	keeperMetadataRegistry.Register(keeper.PATRONI_POSTGRES, patroniAdapter)
	keeperMetadataRegistry.Register(keeper.NATIVE_POSTGRES, pgKeeperAdapter)
	keeperMetadataRegistry.Register(keeper.NATIVE_ETCD, etcdKeeperAdapter)
	keeperMetadataRegistry.Register(keeper.NATIVE_REDIS, redisKeeperAdapter)
	keeperMetadataRegistry.Register(keeper.NATIVE_CLICKHOUSE, clickhouseKeeperAdapter)
	keeperMetadataRegistry.Register(keeper.NATIVE_ZOOKEEPER, zookeeperKeeperAdapter)
	keeperMetadataRegistry.Register(keeper.NATIVE_MONGO, mongoKeeperAdapter)
	databaseRegistry := utils.NewRegistry[database.Plugin, database.Adapter]()
	databaseRegistry.Register(database.POSTGRES, postgresAdapter)
	databaseRegistry.Register(database.ETCD, etcdDbAdapter)
	databaseRegistry.Register(database.REDIS, redisDbAdapter)
	databaseRegistry.Register(database.CLICKHOUSE, clickhouseDbAdapter)
	databaseRegistry.Register(database.ZOOKEEPER, zookeeperDbAdapter)
	databaseRegistry.Register(database.MONGO, mongoDbAdapter)
	platformRegistry := utils.NewRegistry[platform.Plugin, platform.Adapter]()
	platformRegistry.Register(platform.Docker, dockerAdapter)
	platformMetadataRegistry := utils.NewRegistry[platform.Plugin, platform.Metadata]()
	platformMetadataRegistry.Register(platform.Docker, dockerAdapter)

	return &Context{
		KeeperRegistry:           keeperRegistry,
		KeeperMetadataRegistry:   keeperMetadataRegistry,
		DatabaseRegistry:         databaseRegistry,
		PlatformRegistry:         platformRegistry,
		PlatformMetadataRegistry: platformMetadataRegistry,
	}
}
