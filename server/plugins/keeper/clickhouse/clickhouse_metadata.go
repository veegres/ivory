package clickhouse

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
		config.ManageNodeKeeperFailover:     false,
		config.ManageNodeKeeperActivation:   false,
	}
}

// DeploymentSpec deploys a bare, standalone clickhouse-server per node, the
// same posture as native postgres (see that plugin's Adapter doc): real
// multi-node replication needs a config.xml <remote_servers>/<zookeeper>
// section that cannot be expressed through env vars or CLI flags, so Ivory
// does not attempt to wire clustering automatically here - List/Config/
// Reload work against each node's own connection regardless. The official
// image reads CLICKHOUSE_USER/CLICKHOUSE_PASSWORD from Env on every node
// including the primary, so - unlike redis - no EntryScript is needed.
// Unlike VarDbPort's other users, the official image has no env var to
// change the native tcp_port away from its config.xml default (9000): it is
// still declared here because Defaults always requires it, but a deploy
// that overrides it would need a custom image with its own config.xml.
func (a *Adapter) DeploymentSpec() keeper.DeploymentSpec {
	return keeper.DeploymentSpec{
		DefaultImage: "clickhouse/clickhouse-server:24",
		// NOTE: the empty username means credentials are consumed but the
		// username is the user's choice
		Defaults: map[keeper.Var]string{
			keeper.VarDbPort: "9000",
			keeper.VarDbUser: "",
		},
		Ports: []string{string(keeper.VarDbPort)},
		Volumes: []keeper.VolumeSpec{
			{HostPath: "/data/clickhouse", ContainerPath: "/var/lib/clickhouse"},
		},
		Env: []keeper.EnvVar{
			{Name: "CLICKHOUSE_USER", Value: `"` + string(keeper.VarDbUser) + `"`},
			{Name: "CLICKHOUSE_PASSWORD", Value: `"` + string(keeper.VarDbPass) + `"`},
			{Name: "CLICKHOUSE_DB", Value: `"default"`},
		},
	}
}
