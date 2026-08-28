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
		DbPort:      5432,
		Credentials: true,
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
  until pg_isready -d "host=postgres-1 port=5432 user={{dbUser}} password={{dbPass}}"; do sleep 1; done
  pg_basebackup -d "host=postgres-1 port=5432 user={{dbUser}} password={{dbPass}} application_name={{name}}" -D "$PGDATA" -Fp -R -X stream -c fast
fi
exec docker-entrypoint.sh postgres
'`

const deploySingleHostLeader = `docker run -d
  --name {{name}}
  --hostname {{host}}
  --network host
  -e PGPORT="{{dbPort}}"
  -e POSTGRES_USER="{{dbUser}}"
  -e POSTGRES_PASSWORD="{{dbPass}}"
  postgres:18`

const deploySingleHostReplica = `docker run -d
  --name {{name}}
  --hostname {{host}}
  --network host
  -e PGPORT="{{dbPort}}"
  -e POSTGRES_USER="{{dbUser}}"
  -e POSTGRES_PASSWORD="{{dbPass}}"
  postgres:18
  sh -c '
if [ ! -s "$PGDATA/PG_VERSION" ]; then
  until pg_isready -d "host=postgres-1 port=5432 user={{dbUser}} password={{dbPass}}"; do sleep 1; done
  pg_basebackup -d "host=postgres-1 port=5432 user={{dbUser}} password={{dbPass}} application_name={{name}}" -D "$PGDATA" -Fp -R -X stream -c fast
fi
exec docker-entrypoint.sh postgres
'`

func (a *Adapter) DefaultTemplates() []keeper.DeploymentTemplate {
	return []keeper.DeploymentTemplate{
		{
			Platform:    platform.Docker,
			Name:        "Postgres (Multi Host)",
			Description: "One postgres leader and two streaming replicas, one per VM. Name the leader postgres-1 or edit the replica connection to match.",
			Commands: []keeper.DeploymentCommand{
				{Command: deployMultiHostLeader},
				{Command: deployMultiHostReplica},
				{Command: deployMultiHostReplica},
			},
		},
		{
			Platform:    platform.Docker,
			Name:        "Postgres (Single Host)",
			Description: "One postgres leader and two streaming replicas on one VM. Give each node its own database port in the deploy form, and point the replicas at the leader's.",
			Commands: []keeper.DeploymentCommand{
				{Command: deploySingleHostLeader},
				{Command: deploySingleHostReplica},
				{Command: deploySingleHostReplica},
			},
		},
	}
}
