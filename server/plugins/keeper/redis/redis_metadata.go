package redis

import (
	"ivory/core/config"
	"ivory/plugins/keeper"
)

// NOTE: validate that is matches interface in compile-time
var _ keeper.Metadata = (*Adapter)(nil)

func (a *Adapter) SupportedFeatures() map[config.Feature]bool {
	return map[config.Feature]bool{
		config.ViewNodeKeeperOverview:       true,
		config.ViewNodeKeeperConfig:         true,
		config.ManageNodeKeeperConfigUpdate: false,
		config.ManageNodeKeeperSwitchover:   false,
		config.ManageNodeKeeperReinitialize: false,
		config.ManageNodeKeeperRestart:      false,
		config.ManageNodeKeeperReload:       false,
		config.ManageNodeKeeperFailover:     true,
		config.ManageNodeKeeperActivation:   false,
	}
}

// entryScript makes a non-primary node start as a replica of
// {{primaryHost}} instead of coming up as its own standalone master.
//
// The official redis image only takes port/password as redis-server CLI
// flags, not env vars, so it cannot be configured through DeploymentSpec.Env
// alone the way postgres/etcd are - the primary node (which never gets an
// EntryScript, see keeper.DeploymentSpec.EntryScript) would be left with no
// password and the default port. bitnami/redis is used instead because it
// reads REDIS_PASSWORD/REDIS_PORT_NUMBER from Env on every node including
// the primary, and only the replica-specific settings
// (REDIS_REPLICATION_MODE/REDIS_MASTER_*) need to be layered on top here,
// exactly the same division of labor as postgres' own EntryScript (Env
// configures every node identically, EntryScript adds what only a replica
// needs) - just with env vars exported before exec instead of a basebackup.
const entryScript = `sh -c '
export REDIS_REPLICATION_MODE="slave"
export REDIS_MASTER_HOST="` + string(keeper.VarPrimaryHost) + `"
export REDIS_MASTER_PORT_NUMBER="` + string(keeper.VarDbPort) + `"
export REDIS_MASTER_PASSWORD="` + string(keeper.VarDbPass) + `"
exec /opt/bitnami/scripts/redis/entrypoint.sh /opt/bitnami/scripts/redis/run.sh
'`

func (a *Adapter) DeploymentSpec() keeper.DeploymentSpec {
	return keeper.DeploymentSpec{
		DefaultImage: "bitnami/redis:7.4",
		Defaults: map[keeper.Var]string{
			keeper.VarDbPort: "6379",
			keeper.VarDbUser: "default",
		},
		Ports: []string{string(keeper.VarDbPort)},
		Volumes: []keeper.VolumeSpec{
			{HostPath: "/data/redis", ContainerPath: "/bitnami/redis/data"},
		},
		Env: []keeper.EnvVar{
			{Name: "REDIS_PORT_NUMBER", Value: `"` + string(keeper.VarDbPort) + `"`},
			{Name: "REDIS_PASSWORD", Value: `"` + string(keeper.VarDbPass) + `"`},
			{Name: "ALLOW_EMPTY_PASSWORD", Value: `"no"`},
		},
		EntryScript:             entryScript,
		EntryScriptReplicasOnly: true,
	}
}
