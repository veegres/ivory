package postgres

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
		config.ManageNodeKeeperReload:       true,
		config.ManageNodeKeeperFailover:     true,
		config.ManageNodeKeeperActivation:   false,
	}
}

func (p *Plugin) HasLeader() bool { return true }

// A replica rebases from the leader before postgres ever starts: streaming
// replication ships changes to an existing copy, it cannot build the initial
// database, and a fresh initdb has its own system identifier the leader would
// refuse to stream to. PG_VERSION only exists once the data directory is
// initialized, so its absence marks the container's first boot. -R writes
// standby.signal and primary_conninfo; application_name is threaded through a
// conninfo -d (pg_basebackup has no flag for it) so the replica appears under
// its own name, which is what lets Ivory tell sync from async standbys apart
// later. The -h/-p/-U flags are appended after that conninfo and win over it,
// so the connection itself is still stated as flags. The leader's host is
// written literally: only the operator knows which node it is.
//
// The volume is mounted at /var/lib/postgresql, the directory the image itself
// declares, rather than at .../data: postgres:18 moved its default PGDATA to
// /var/lib/postgresql/18/docker, so a volume on the pre-18 path caught nothing
// the database actually writes - pg_basebackup -D "$PGDATA" rebuilt a replica
// into an unmounted directory that the next container recreation threw away.
//
// The script reads the credentials from the env the command itself sets rather
// than interpolating {{dbUser}}/{{dbPass}} a second time: this text is parsed
// again by the container's own shell, where a password holding a `$` or a
// backtick would be expanded or executed. They stay out of the conninfo for a
// second reason - libpq applies its own quoting rules there, so a space in
// either would terminate the value early.
//
// The leader drops an init hook rather than starting straight away, because
// the image's generated pg_hba.conf grants replication to loopback only: its
// catch-all last line is "host all all all <method>", and the database keyword
// all deliberately excludes the replication pseudo-database, so a replica
// connecting from anywhere else matches nothing. The hook is a plain file in
// /docker-entrypoint-initdb.d, which the entrypoint sources as the postgres
// user after it writes its own hba block and before the real server starts.
// The heredoc delimiter is backslash-quoted so $PGDATA is written literally
// and expands inside the hook, where it is set, rather than here, where it is
// not.
//
// set -e is what makes a failed rebase visible: without it the replica script
// falls through to the entrypoint, which finds an empty data directory and
// initdb's a brand-new standalone primary - a cluster of three primaries that
// reports itself as a successful deploy. The chown is for the other half of
// the same step: pg_basebackup runs as root, so creating $PGDATA also creates
// its parent root-owned and mode 700, and the entrypoint's gosu phase can then
// not even traverse it.

const deployMultiHostLeader = `docker run -d
  --name {{name}}
  --hostname {{host}}
  --restart unless-stopped
  -p {{dbPort}}:{{dbPort}}
  -v /data/postgres:/var/lib/postgresql
  -e PGPORT="{{dbPort}}"
  -e POSTGRES_USER="{{dbUser}}"
  -e POSTGRES_PASSWORD="{{dbPass}}"
  postgres:18
  sh -c '
cat > /docker-entrypoint-initdb.d/00-ivory-replication.sh <<\IVORYEOF
echo "host replication all all scram-sha-256" >> "$PGDATA/pg_hba.conf"
IVORYEOF
exec docker-entrypoint.sh postgres
'`

const deployMultiHostReplica = `docker run -d
  --name {{name}}
  --hostname {{host}}
  --restart unless-stopped
  -p {{dbPort}}:{{dbPort}}
  -v /data/postgres:/var/lib/postgresql
  -e PGPORT="{{dbPort}}"
  -e POSTGRES_USER="{{dbUser}}"
  -e POSTGRES_PASSWORD="{{dbPass}}"
  postgres:18
  sh -c '
set -e
if [ ! -s "$PGDATA/PG_VERSION" ]; then
  export PGPASSWORD="$POSTGRES_PASSWORD"
  until pg_isready -h 10.0.0.1 -p 5432 -U "$POSTGRES_USER"; do sleep 1; done
  pg_basebackup -d "application_name={{name}}" -h 10.0.0.1 -p 5432 -U "$POSTGRES_USER" -D "$PGDATA" -Fp -R -X stream -c fast
  chown -R postgres:postgres /var/lib/postgresql
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
  postgres:18
  sh -c '
cat > /docker-entrypoint-initdb.d/00-ivory-replication.sh <<\IVORYEOF
echo "host replication all all scram-sha-256" >> "$PGDATA/pg_hba.conf"
IVORYEOF
exec docker-entrypoint.sh postgres
'`

const deploySingleHostReplica = `docker run -d
  --name {{name}}
  --network host
  -e PGPORT="{{dbPort}}"
  -e POSTGRES_USER="{{dbUser}}"
  -e POSTGRES_PASSWORD="{{dbPass}}"
  postgres:18
  sh -c '
set -e
if [ ! -s "$PGDATA/PG_VERSION" ]; then
  export PGPASSWORD="$POSTGRES_PASSWORD"
  until pg_isready -h {{host}} -p 5432 -U "$POSTGRES_USER"; do sleep 1; done
  pg_basebackup -d "application_name={{name}}" -h {{host}} -p 5432 -U "$POSTGRES_USER" -D "$PGDATA" -Fp -R -X stream -c fast
  chown -R postgres:postgres /var/lib/postgresql
fi
exec docker-entrypoint.sh postgres
'`

func (p *Plugin) DefaultTemplates() []keeper.DeploymentTemplate {
	return []keeper.DeploymentTemplate{
		{
			Platform:    platform.Docker,
			Name:        "Postgres (Multi Host)",
			Description: "One postgres leader and two streaming replicas, one per VM. Replace 10.0.0.1 in the replica connection with the leader's VM address.",
			Defaults:    keeper.DeploymentTemplateDefaults{KeeperUser: "postgres", DbUser: "postgres"},
			Commands: []keeper.DeploymentCommand{
				{
					Command:  deployMultiHostLeader,
					Defaults: keeper.DeploymentCommandDefaults{Name: "postgres1", SshPort: 22, KeeperPort: 5432, DbPort: 5432},
				},
				{
					Command:  deployMultiHostReplica,
					Defaults: keeper.DeploymentCommandDefaults{Name: "postgres2", SshPort: 22, KeeperPort: 5432, DbPort: 5432},
				},
				{
					Command:  deployMultiHostReplica,
					Defaults: keeper.DeploymentCommandDefaults{Name: "postgres3", SshPort: 22, KeeperPort: 5432, DbPort: 5432},
				},
			},
		},
		{
			Platform:    platform.Docker,
			Name:        "Postgres (Single Host)",
			Description: "One postgres leader and two streaming replicas on one VM, each on its own port. The replicas rebase from the leader on 5432. Fill in the host and deploy.",
			Defaults:    keeper.DeploymentTemplateDefaults{KeeperUser: "postgres", DbUser: "postgres"},
			Commands: []keeper.DeploymentCommand{
				{
					Command:  deploySingleHostLeader,
					Defaults: keeper.DeploymentCommandDefaults{Name: "postgres1", SshPort: 22, KeeperPort: 5432, DbPort: 5432},
				},
				{
					Command:  deploySingleHostReplica,
					Defaults: keeper.DeploymentCommandDefaults{Name: "postgres2", SshPort: 22, KeeperPort: 5433, DbPort: 5433},
				},
				{
					Command:  deploySingleHostReplica,
					Defaults: keeper.DeploymentCommandDefaults{Name: "postgres3", SshPort: 22, KeeperPort: 5434, DbPort: 5434},
				},
			},
		},
	}
}
