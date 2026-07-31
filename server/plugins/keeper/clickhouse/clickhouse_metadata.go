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

// entryScript generates a config.d/ivory-cluster.xml at container startup so
// every node - including the primary, since clickhouse has no primary/replica
// asymmetry and EntryScript applies uniformly by default (see
// keeper.DeploymentSpec.EntryScriptReplicasOnly) - joins the same
// <remote_servers> shard and points at the same external coordinator, and
// then hands off to the image's normal startup:
//   - {{clusterHosts}} is already-resolved <replica> XML for every node in
//     this deploy (built by KeeperDeployPlan from the Fields entry below), so
//     it is embedded directly, no runtime work needed.
//   - {{dcs}} is a single comma-separated host:port list the user provides
//     (same free-text shape as patroni's own {{dcs}} field - see that
//     plugin's Adapter doc), so unlike {{clusterHosts}} it still needs
//     splitting into <node> entries at container startup, since Ivory only
//     ever sees it as one opaque string.
//   - <macros> gives every node's replicated tables a stable per-replica
//     identity (its own {{host}}) for ON CLUSTER / {replica} substitutions.
const entryScript = `sh -c '
dcs="` + string(keeper.VarDcs) + `"
zk=""
oldifs="$IFS"
IFS=","
for hp in $dcs; do
  h=${hp%%:*}
  p=${hp##*:}
  zk="$zk<node><host>$h</host><port>$p</port></node>"
done
IFS="$oldifs"
cat > /etc/clickhouse-server/config.d/ivory-cluster.xml <<IVORYEOF
<clickhouse>
  <remote_servers>
    <ivory_cluster>
      <shard>
        ` + string(keeper.VarClusterHosts) + `
      </shard>
    </ivory_cluster>
  </remote_servers>
  <zookeeper>
    $zk
  </zookeeper>
  <macros>
    <cluster>` + string(keeper.VarCluster) + `</cluster>
    <replica>` + string(keeper.VarHost) + `</replica>
  </macros>
</clickhouse>
IVORYEOF
exec /entrypoint.sh
'`

// DeploymentSpec clusters every node into one <remote_servers> shard
// coordinated through an externally-run ZooKeeper or ClickHouse Keeper
// ensemble - the same posture as patroni pointing at an external DCS via
// {{dcs}} (see patroni's Adapter doc), rather than Ivory deploying and
// managing that coordinator itself. EntryScriptReplicasOnly is deliberately
// left false (the default): clickhouse has no primary/replica asymmetry, so
// every node - including the first - needs EntryScript to generate the same
// cluster config file, unlike postgres/etcd/redis which skip it on node 0.
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
		Fields: []keeper.FieldSpec{
			{Name: keeper.VarDcs, Label: "ZooKeeper / ClickHouse Keeper", Example: "keeper1:9181, keeper2:9181, keeper3:9181", Type: keeper.FieldText},
			{
				Name: keeper.VarClusterHosts, Label: "Cluster Replicas (auto)", Type: keeper.FieldText,
				Template: "<replica><host>" + string(keeper.VarHost) + "</host><port>" + string(keeper.VarDbPort) + "</port></replica>",
			},
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
		EntryScript: entryScript,
	}
}
