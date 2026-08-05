package mongo

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
		config.ManageNodeKeeperConfigUpdate: true,
		config.ManageNodeKeeperSwitchover:   false,
		config.ManageNodeKeeperReinitialize: false,
		config.ManageNodeKeeperRestart:      false,
		config.ManageNodeKeeperReload:       false,
		config.ManageNodeKeeperFailover:     true,
		config.ManageNodeKeeperActivation:   false,
	}
}

// entryScript starts every node - including the first, since a replica set
// has no primary/replica asymmetry at container-start time (all of them come
// up as plain standalone members of an as-yet-uninitiated set; who becomes
// primary is decided later, by postScript's rs.initiate()) - as a mongod
// listening with --replSet already set. The official image has no env var
// for --replSet, so this replaces the image's default command the same way
// clickhouse's own EntryScript does, still going through
// docker-entrypoint.sh so the image's normal startup housekeeping still runs.
const entryScript = `sh -c '
exec docker-entrypoint.sh mongod --replSet "` + string(keeper.VarCluster) + `" --port ` + string(keeper.VarDbPort) + ` --bind_ip_all
'`

// postScript runs once, only on the deploy's first node (see
// keeper.DeploymentSpec.PostScript), and turns the set of standalone mongod
// processes EntryScript started into an actual replica set: it waits for the
// local mongod to accept connections, then calls rs.initiate() with every
// node from {{clusterHosts}} (built by the Fields entry below) as members.
// Double quotes inside the mongosh --eval argument are backslash-escaped
// rather than using literal single quotes, since the whole script is itself
// wrapped in a single-quoted `+"`sh -c '...'`"+` argument (see EntryScript and
// e.g. the postgres/redis plugins for the same convention) where a literal
// single quote would terminate the script early.
const postScript = `sh -c '
until mongosh --quiet --port ` + string(keeper.VarDbPort) + ` --eval "1" >/dev/null 2>&1; do sleep 1; done
mongosh --quiet --port ` + string(keeper.VarDbPort) + ` --eval "rs.initiate({_id: \"` + string(keeper.VarCluster) + `\", members: [` + string(keeper.VarClusterHosts) + `]})"
'`

// DeploymentSpec deploys a self-forming replica set the same way native etcd
// does, rather than pointing at an externally-run coordinator - there is no
// {{dcs}} field here because mongo's replica set protocol is the coordinator.
// It deliberately declares no credentials (no VarDbUser/VarDbPass in
// Defaults): enabling client auth on a replica set also requires internal
// authentication between members (a shared keyfile mounted into every
// container), which Ivory's deploy model has no mechanism for yet - the same
// posture as native zookeeper, which also ships with no auth. Operators who
// need auth can still configure it themselves on the deployed cluster and
// supply keeper credentials manually; Adapter.connect already accepts them
// when present.
func (a *Adapter) DeploymentSpec() keeper.DeploymentSpec {
	return keeper.DeploymentSpec{
		DefaultImage: "mongo:8",
		Defaults: map[keeper.Var]string{
			keeper.VarDbPort: "27017",
		},
		Fields: []keeper.FieldSpec{
			{
				Name: keeper.VarClusterHosts, Label: "Replica Set Members (auto)", Type: keeper.FieldText,
				Template:  `{_id: ` + string(keeper.VarIndex) + `, host: \"` + string(keeper.VarHost) + `:` + string(keeper.VarDbPort) + `\"}`,
				Separator: ", ",
			},
		},
		Ports: []string{string(keeper.VarDbPort)},
		Volumes: []keeper.VolumeSpec{
			{HostPath: "/data/mongo", ContainerPath: "/data/db"},
		},
		EntryScript: entryScript,
		PostScript:  postScript,
	}
}
