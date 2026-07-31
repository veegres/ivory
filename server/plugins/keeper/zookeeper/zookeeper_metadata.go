package zookeeper

import (
	"ivory/core/config"
	"ivory/plugins/keeper"
)

// varLeaderElectionPort is zookeeper's own third port (its ZAB leader
// election listener, distinct from the client port and the peer/quorum
// port keeper.VarPeerPort already covers). It stays local to this package
// rather than joining the shared keeper vocabulary, since zookeeper is the
// only plugin needing it today - see AGENTS.md's rule of only promoting a
// Var once a second real plugin needs the identical mechanism.
const varLeaderElectionPort keeper.Var = "{{leaderElectionPort}}"

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
		config.ManageNodeKeeperFailover:     false,
		config.ManageNodeKeeperActivation:   false,
	}
}

// DeploymentSpec bootstraps a static zookeeper ensemble entirely through Env
// - no EntryScript is needed, unlike postgres/redis, because every node's
// startup config (including the leader's) is expressible as plain
// environment variables:
//   - {{index}} (see keeper.VarIndex) becomes ZOO_MY_ID, since zookeeper
//     requires a genuinely unique small integer per member - a hostname
//     alone (as etcd/clickhouse use) isn't a valid member id here.
//   - {{clusterHosts}} (see keeper.VarClusterHosts) builds ZOO_SERVERS, one
//     "server.<index>=<host>:<peerPort>:<leaderElectionPort>;<dbPort>" entry
//     per node - the official image has no separate client-port env var, so
//     the port must be embedded in each entry via the ";<port>" suffix.
//   - ZOO_4LW_COMMANDS_WHITELIST must explicitly list "mntr"/"conf": recent
//     zookeeper versions whitelist only "srvr" by default, and
//     zookeeper_adapter.go's List/Config depend on the other two.
//
// This deploys a self-forming ensemble the same way native etcd does -
// unlike patroni or clickhouse, which point at an externally-run coordinator
// via {{dcs}} instead. A deployed zookeeper ensemble is exactly the kind of
// coordinator clickhouse's own {{dcs}} field can then point at.
func (a *Adapter) DeploymentSpec() keeper.DeploymentSpec {
	return keeper.DeploymentSpec{
		DefaultImage: "zookeeper:3.9",
		Defaults: map[keeper.Var]string{
			keeper.VarDbPort: "2181",
		},
		Fields: []keeper.FieldSpec{
			{Name: keeper.VarPeerPort, Label: "Peer Port (quorum)", Type: keeper.FieldPort, Default: "2888"},
			{Name: varLeaderElectionPort, Label: "Leader Election Port", Type: keeper.FieldPort, Default: "3888"},
			{
				Name: keeper.VarClusterHosts, Label: "Ensemble Servers (auto)", Type: keeper.FieldText,
				Template:  "server." + string(keeper.VarIndex) + "=" + string(keeper.VarHost) + ":" + string(keeper.VarPeerPort) + ":" + string(varLeaderElectionPort) + ";" + string(keeper.VarDbPort),
				Separator: " ",
			},
		},
		Ports: []string{string(keeper.VarDbPort), string(keeper.VarPeerPort), string(varLeaderElectionPort)},
		Volumes: []keeper.VolumeSpec{
			{HostPath: "/data/zookeeper/data", ContainerPath: "/data"},
			{HostPath: "/data/zookeeper/datalog", ContainerPath: "/datalog"},
		},
		Env: []keeper.EnvVar{
			{Name: "ZOO_MY_ID", Value: `"` + string(keeper.VarIndex) + `"`},
			{Name: "ZOO_SERVERS", Value: `"` + string(keeper.VarClusterHosts) + `"`},
			{Name: "ZOO_4LW_COMMANDS_WHITELIST", Value: `"mntr,conf,ruok,srvr"`},
		},
	}
}
