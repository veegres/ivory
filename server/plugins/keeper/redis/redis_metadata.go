package redis

import (
	"ivory/core/config"
	"ivory/plugins/keeper"
	"ivory/plugins/platform"
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

func (a *Adapter) Requirements() keeper.Requirements {
	return keeper.Requirements{
		KeeperCredentials: true,
		KeeperUser:        "default",
		DbCredentials:     true,
		DbUser:            "default",
	}
}

// The bitnami image is used rather than the official one because that takes
// port and password as redis-server flags only, so the leader - which runs no
// replica bootstrap - could not be configured through environment variables at
// all. The leader/replica difference is just a different command, and the
// leader's host is written literally: only the operator knows which node it is.

const deployMultiHostLeader = `docker run -d
  --name {{name}}
  --hostname {{host}}
  --restart unless-stopped
  -p {{dbPort}}:{{dbPort}}
  -v /data/redis:/bitnami/redis/data
  -e REDIS_PORT_NUMBER="{{dbPort}}"
  -e REDIS_PASSWORD="{{dbPass}}"
  -e ALLOW_EMPTY_PASSWORD="no"
  bitnami/redis:7.4`

const deployMultiHostReplica = `docker run -d
  --name {{name}}
  --hostname {{host}}
  --restart unless-stopped
  -p {{dbPort}}:{{dbPort}}
  -v /data/redis:/bitnami/redis/data
  -e REDIS_PORT_NUMBER="{{dbPort}}"
  -e REDIS_PASSWORD="{{dbPass}}"
  -e ALLOW_EMPTY_PASSWORD="no"
  -e REDIS_REPLICATION_MODE="slave"
  -e REDIS_MASTER_HOST="redis1"
  -e REDIS_MASTER_PORT_NUMBER="6379"
  -e REDIS_MASTER_PASSWORD="{{dbPass}}"
  bitnami/redis:7.4`

// The single-host commands drop --hostname, which docker rejects alongside
// --network host, and the replicas follow the leader at {{host}}:6379 - host
// networking resolves no container names, and 6379 is the first node's port
// rather than the replica's own.

const deploySingleHostLeader = `docker run -d
  --name {{name}}
  --network host
  -e REDIS_PORT_NUMBER="{{dbPort}}"
  -e REDIS_PASSWORD="{{dbPass}}"
  -e ALLOW_EMPTY_PASSWORD="no"
  bitnami/redis:7.4`

const deploySingleHostReplica = `docker run -d
  --name {{name}}
  --network host
  -e REDIS_PORT_NUMBER="{{dbPort}}"
  -e REDIS_PASSWORD="{{dbPass}}"
  -e ALLOW_EMPTY_PASSWORD="no"
  -e REDIS_REPLICATION_MODE="slave"
  -e REDIS_MASTER_HOST="{{host}}"
  -e REDIS_MASTER_PORT_NUMBER="6379"
  -e REDIS_MASTER_PASSWORD="{{dbPass}}"
  bitnami/redis:7.4`

func (a *Adapter) DefaultTemplates() []keeper.DeploymentTemplate {
	return []keeper.DeploymentTemplate{
		{
			Platform:    platform.Docker,
			Name:        "Redis (Multi Host)",
			Description: "One redis leader and two replicas, one per VM. Name the leader redis1 or edit REDIS_MASTER_HOST to match.",
			Commands: []keeper.DeploymentCommand{
				{
					Command:  deployMultiHostLeader,
					Defaults: keeper.DeploymentDefaults{Name: "redis1", KeeperPort: 6379, DbPort: 6379},
				},
				{
					Command:  deployMultiHostReplica,
					Defaults: keeper.DeploymentDefaults{Name: "redis2", KeeperPort: 6379, DbPort: 6379},
				},
				{
					Command:  deployMultiHostReplica,
					Defaults: keeper.DeploymentDefaults{Name: "redis3", KeeperPort: 6379, DbPort: 6379},
				},
			},
		},
		{
			Platform:    platform.Docker,
			Name:        "Redis (Single Host)",
			Description: "One redis leader and two replicas on one VM, each on its own port. The replicas follow the leader on 6379. Fill in the host and deploy.",
			Commands: []keeper.DeploymentCommand{
				{
					Command:  deploySingleHostLeader,
					Defaults: keeper.DeploymentDefaults{Name: "redis1", KeeperPort: 6379, DbPort: 6379},
				},
				{
					Command:  deploySingleHostReplica,
					Defaults: keeper.DeploymentDefaults{Name: "redis2", KeeperPort: 6380, DbPort: 6380},
				},
				{
					Command:  deploySingleHostReplica,
					Defaults: keeper.DeploymentDefaults{Name: "redis3", KeeperPort: 6381, DbPort: 6381},
				},
			},
		},
	}
}
