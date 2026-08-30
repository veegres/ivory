package redis

import (
	"ivory/core/config"
	"ivory/plugins/keeper"
	"ivory/plugins/platform"
)

// NOTE: validate that is matches interface in compile-time
var _ keeper.Metadata = (*Plugin)(nil)

func (p *Plugin) SupportedFeatures() map[config.Feature]bool {
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

func (p *Plugin) HasLeader() bool { return true }

// The official image takes its port and password as redis-server flags rather
// than environment variables, so the command states them after the image - the
// container runs redis-server directly, with no shell in between to reparse
// them. bitnami/redis was used for this once and was retired from Docker Hub;
// a frozen bitnamilegacy tag would only postpone the same failure.
//
// The leader/replica difference is just a different command. The leader's
// address is written literally, as example text the operator replaces: only
// they know which VM it is, and a container name resolves on neither a plain
// docker run across VMs nor host networking.

const deployMultiHostLeader = `docker run -d
  --name {{name}}
  --hostname {{host}}
  --restart unless-stopped
  -p {{dbPort}}:{{dbPort}}
  -v /data/redis:/data
  redis:7.4
  redis-server --port {{dbPort}} --requirepass "{{dbPass}}" --appendonly yes`

const deployMultiHostReplica = `docker run -d
  --name {{name}}
  --hostname {{host}}
  --restart unless-stopped
  -p {{dbPort}}:{{dbPort}}
  -v /data/redis:/data
  redis:7.4
  redis-server --port {{dbPort}} --requirepass "{{dbPass}}" --appendonly yes --replicaof 10.0.0.1 6379 --masterauth "{{dbPass}}"`

// The single-host commands drop --hostname, which docker rejects alongside
// --network host, and publish no port: each node answers on its own port of
// the host's one namespace. The replicas follow the leader at {{host}}:6379 -
// host networking resolves no container names, and 6379 is the first node's
// port rather than the replica's own.

const deploySingleHostLeader = `docker run -d
  --name {{name}}
  --network host
  redis:7.4
  redis-server --port {{dbPort}} --requirepass "{{dbPass}}"`

const deploySingleHostReplica = `docker run -d
  --name {{name}}
  --network host
  redis:7.4
  redis-server --port {{dbPort}} --requirepass "{{dbPass}}" --replicaof {{host}} 6379 --masterauth "{{dbPass}}"`

func (p *Plugin) DefaultTemplates() []keeper.DeploymentTemplate {
	return []keeper.DeploymentTemplate{
		{
			Platform:    platform.Docker,
			Name:        "Redis (Multi Host)",
			Description: "One redis leader and two replicas, one per VM. Replace 10.0.0.1 in the replica commands with the leader's address.",
			Defaults:    keeper.DeploymentTemplateDefaults{KeeperUser: "default", DbUser: "default"},
			Commands: []keeper.DeploymentCommand{
				{
					Command:  deployMultiHostLeader,
					Defaults: keeper.DeploymentCommandDefaults{Name: "redis1", KeeperPort: 6379, DbPort: 6379},
				},
				{
					Command:  deployMultiHostReplica,
					Defaults: keeper.DeploymentCommandDefaults{Name: "redis2", KeeperPort: 6379, DbPort: 6379},
				},
				{
					Command:  deployMultiHostReplica,
					Defaults: keeper.DeploymentCommandDefaults{Name: "redis3", KeeperPort: 6379, DbPort: 6379},
				},
			},
		},
		{
			Platform:    platform.Docker,
			Name:        "Redis (Single Host)",
			Description: "One redis leader and two replicas on one VM, each on its own port. The replicas follow the leader on 6379. Fill in the host and deploy.",
			Defaults:    keeper.DeploymentTemplateDefaults{KeeperUser: "default", DbUser: "default"},
			Commands: []keeper.DeploymentCommand{
				{
					Command:  deploySingleHostLeader,
					Defaults: keeper.DeploymentCommandDefaults{Name: "redis1", KeeperPort: 6379, DbPort: 6379},
				},
				{
					Command:  deploySingleHostReplica,
					Defaults: keeper.DeploymentCommandDefaults{Name: "redis2", KeeperPort: 6380, DbPort: 6380},
				},
				{
					Command:  deploySingleHostReplica,
					Defaults: keeper.DeploymentCommandDefaults{Name: "redis3", KeeperPort: 6381, DbPort: 6381},
				},
			},
		},
	}
}
