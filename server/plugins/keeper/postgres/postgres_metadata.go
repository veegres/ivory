package postgres

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
		config.ManageNodeKeeperReload:       true,
		config.ManageNodeKeeperFailover:     true,
		config.ManageNodeKeeperActivation:   false,
	}
}

func (a *Adapter) Requirements() keeper.Requirements {
	return keeper.Requirements{
		DbPort:            5432,
		KeeperPort:        5432,
		KeeperCredentials: true,
		DbCredentials:     true,
	}
}

// A replica rebases from the leader before postgres ever starts: streaming
// replication ships changes to an existing copy, it cannot build the initial
// database, and a fresh initdb has its own system identifier the leader would
// refuse to stream to. PG_VERSION only exists once the data directory is
// initialized, so its absence marks the container's first boot. -R writes
// standby.signal and primary_conninfo; application_name is threaded through
// the connection string (pg_basebackup has no flag for it) so the replica
// appears under its own name, which is what lets Ivory tell sync from async
// standbys apart later. The leader's host is written literally: only the
// operator knows which node it is.

const deployMultiHostLeader = `docker run -d
  --name {{name}}
  --hostname {{host}}
  --restart unless-stopped
  -p {{dbPort}}:{{dbPort}}
  -v /data/postgres:/var/lib/postgresql/data
  -e PGPORT="{{dbPort}}"
  -e POSTGRES_USER="{{dbUser}}"
  -e POSTGRES_PASSWORD="{{dbPass}}"
  postgres:18`

const deployMultiHostReplica = `docker run -d
  --name {{name}}
  --hostname {{host}}
  --restart unless-stopped
  -p {{dbPort}}:{{dbPort}}
  -v /data/postgres:/var/lib/postgresql/data
  -e PGPORT="{{dbPort}}"
  -e POSTGRES_USER="{{dbUser}}"
  -e POSTGRES_PASSWORD="{{dbPass}}"
  postgres:18
  sh -c '
if [ ! -s "$PGDATA/PG_VERSION" ]; then
  until pg_isready -d "host=postgres1 port=5432 user={{dbUser}} password={{dbPass}}"; do sleep 1; done
  pg_basebackup -d "host=postgres1 port=5432 user={{dbUser}} password={{dbPass}} application_name={{name}}" -D "$PGDATA" -Fp -R -X stream -c fast
fi
exec docker-entrypoint.sh postgres
'`

// The single-host commands drop --hostname: docker rejects it outright
// alongside --network host, which each node needs to answer on its own port of
// the host's one port namespace. The replicas rebase from {{host}} rather than
// from a container name - host networking joins no docker network, so there is
// no embedded dns to resolve one - and the leader's port is literal 5432
// because that is the first node's port, not the replica's own.

const deploySingleHostLeader = `docker run -d
  --name {{name}}
  --network host
  -e PGPORT="{{dbPort}}"
  -e POSTGRES_USER="{{dbUser}}"
  -e POSTGRES_PASSWORD="{{dbPass}}"
  postgres:18`

const deploySingleHostReplica = `docker run -d
  --name {{name}}
  --network host
  -e PGPORT="{{dbPort}}"
  -e POSTGRES_USER="{{dbUser}}"
  -e POSTGRES_PASSWORD="{{dbPass}}"
  postgres:18
  sh -c '
if [ ! -s "$PGDATA/PG_VERSION" ]; then
  until pg_isready -d "host={{host}} port=5432 user={{dbUser}} password={{dbPass}}"; do sleep 1; done
  pg_basebackup -d "host={{host}} port=5432 user={{dbUser}} password={{dbPass}} application_name={{name}}" -D "$PGDATA" -Fp -R -X stream -c fast
fi
exec docker-entrypoint.sh postgres
'`

func (a *Adapter) DefaultTemplates() []keeper.DeploymentTemplate {
	return []keeper.DeploymentTemplate{
		{
			Platform:    platform.Docker,
			Name:        "Postgres (Multi Host)",
			Description: "One postgres leader and two streaming replicas, one per VM. Name the leader postgres1 or edit the replica connection to match.",
			Commands: []keeper.DeploymentCommand{
				{
					Command:  deployMultiHostLeader,
					Defaults: keeper.DeploymentDefaults{Name: "postgres1", KeeperPort: 5432, DbPort: 5432},
				},
				{
					Command:  deployMultiHostReplica,
					Defaults: keeper.DeploymentDefaults{Name: "postgres2", KeeperPort: 5432, DbPort: 5432},
				},
				{
					Command:  deployMultiHostReplica,
					Defaults: keeper.DeploymentDefaults{Name: "postgres3", KeeperPort: 5432, DbPort: 5432},
				},
			},
		},
		{
			Platform:    platform.Docker,
			Name:        "Postgres (Single Host)",
			Description: "One postgres leader and two streaming replicas on one VM, each on its own port. The replicas rebase from the leader on 5432. Fill in the host and deploy.",
			Commands: []keeper.DeploymentCommand{
				{
					Command:  deploySingleHostLeader,
					Defaults: keeper.DeploymentDefaults{Name: "postgres1", KeeperPort: 5432, DbPort: 5432},
				},
				{
					Command:  deploySingleHostReplica,
					Defaults: keeper.DeploymentDefaults{Name: "postgres2", KeeperPort: 5433, DbPort: 5433},
				},
				{
					Command:  deploySingleHostReplica,
					Defaults: keeper.DeploymentDefaults{Name: "postgres3", KeeperPort: 5434, DbPort: 5434},
				},
			},
		},
	}
}
