package postgres

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
		config.ManageNodeKeeperReload:       true,
		config.ManageNodeKeeperFailover:     true,
		config.ManageNodeKeeperActivation:   false,
	}
}

// entryScript makes a non-primary node rebase itself as a streaming replica
// of {{primaryHost}} before postgres ever starts, instead of coming up as
// its own fresh standalone instance:
//   - PG_VERSION only exists once the data directory has actually been
//     initialized, so its absence means this is the container's very first
//     boot; on every later restart the check is false and startup goes
//     straight to the real entrypoint.
//   - a streaming replica can only ever be created from a real base backup
//     of the primary - WAL streaming ships changes to an existing copy, it
//     cannot build the initial database from nothing, and a fresh initdb'd
//     directory has its own unrelated system identifier that postgres will
//     refuse to stream to. So the first thing a fresh node does is wait for
//     the primary to accept connections, then pg_basebackup from it.
//   - -R makes pg_basebackup write standby.signal and primary_conninfo
//     itself; application_name is threaded through the same connection
//     string (rather than a bare --application-name flag, which
//     pg_basebackup has no such option for) so the replica shows up under
//     its own {{host}} in primary_conninfo - see the postgres Adapter's
//     mapSyncStandby doc for why that convention is what lets Ivory tell
//     sync from async standbys apart later.
const entryScript = `sh -c '
if [ ! -s "$PGDATA/PG_VERSION" ]; then
  until pg_isready -d "host=` + string(keeper.VarPrimaryHost) + ` port=` + string(keeper.VarDbPort) + ` user=` + string(keeper.VarDbUser) + ` password=` + string(keeper.VarDbPass) + `"; do sleep 1; done
  pg_basebackup -d "host=` + string(keeper.VarPrimaryHost) + ` port=` + string(keeper.VarDbPort) + ` user=` + string(keeper.VarDbUser) + ` password=` + string(keeper.VarDbPass) + ` application_name=` + string(keeper.VarHost) + `" -D "$PGDATA" -Fp -R -X stream -c fast
fi
exec docker-entrypoint.sh postgres
'`

func (a *Adapter) DeploymentSpec() keeper.DeploymentSpec {
	return keeper.DeploymentSpec{
		DefaultImage: "postgres:18",
		// NOTE: the empty username means credentials are consumed but the
		// username is the user's choice
		Defaults: map[keeper.Var]string{
			keeper.VarDbPort: "5432",
			keeper.VarDbUser: "",
		},
		Ports: []string{string(keeper.VarDbPort)},
		Volumes: []keeper.VolumeSpec{
			{HostPath: "/data/postgres", ContainerPath: "/var/lib/postgresql/data"},
		},
		Env: []keeper.EnvVar{
			{Name: "PGPORT", Value: `"` + string(keeper.VarDbPort) + `"`},
			{Name: "POSTGRES_USER", Value: `"` + string(keeper.VarDbUser) + `"`},
			{Name: "POSTGRES_PASSWORD", Value: `"` + string(keeper.VarDbPass) + `"`},
		},
		EntryScript: entryScript,
	}
}
