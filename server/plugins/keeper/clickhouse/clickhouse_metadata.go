package clickhouse

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
		config.ManageNodeKeeperFailover:     false,
		config.ManageNodeKeeperActivation:   false,
	}
}

// Requirements reports the native tcp protocol port. The official image has no
// env var to move it away from its config.xml default, so a deploy overriding
// it needs a custom image. The username is the user's own choice.
func (a *Adapter) Requirements() keeper.Requirements {
	return keeper.Requirements{
		DbPort:      9000,
		Credentials: true,
	}
}

// Every node runs the same command, including the first: clickhouse has no
// leader/replica asymmetry. The startup script writes a cluster config file and
// then hands off to the image's own entrypoint. Both the shard's replica list
// and the coordinator address are written literally - only the operator knows
// which hosts those are - and <macros> gives replicated tables a stable
// per-replica identity for ON CLUSTER / {replica}.

const deployMultiHost = `docker run -d
  --name {{name}}
  --hostname {{host}}
  --restart unless-stopped
  -p {{dbPort}}:{{dbPort}}
  -v /data/clickhouse:/var/lib/clickhouse
  -e CLICKHOUSE_USER="{{dbUser}}"
  -e CLICKHOUSE_PASSWORD="{{dbPass}}"
  -e CLICKHOUSE_DB="default"
  clickhouse/clickhouse-server:24
  sh -c '
cat > /etc/clickhouse-server/config.d/ivory-cluster.xml <<IVORYEOF
<clickhouse>
  <remote_servers>
    <ivory_cluster>
      <shard>
        <replica><host>clickhouse-1</host><port>9000</port></replica>
        <replica><host>clickhouse-2</host><port>9000</port></replica>
        <replica><host>clickhouse-3</host><port>9000</port></replica>
      </shard>
    </ivory_cluster>
  </remote_servers>
  <zookeeper>
    <node><host>keeper-1</host><port>9181</port></node>
    <node><host>keeper-2</host><port>9181</port></node>
    <node><host>keeper-3</host><port>9181</port></node>
  </zookeeper>
  <macros>
    <cluster>{{cluster}}</cluster>
    <replica>{{name}}</replica>
  </macros>
</clickhouse>
IVORYEOF
exec /entrypoint.sh
'`

const deploySingleHost = `docker run -d
  --name {{name}}
  --hostname {{host}}
  --network host
  -e CLICKHOUSE_USER="{{dbUser}}"
  -e CLICKHOUSE_PASSWORD="{{dbPass}}"
  -e CLICKHOUSE_DB="default"
  clickhouse/clickhouse-server:24
  sh -c '
cat > /etc/clickhouse-server/config.d/ivory-cluster.xml <<IVORYEOF
<clickhouse>
  <remote_servers>
    <ivory_cluster>
      <shard>
        <replica><host>clickhouse-1</host><port>9000</port></replica>
        <replica><host>clickhouse-2</host><port>9000</port></replica>
        <replica><host>clickhouse-3</host><port>9000</port></replica>
      </shard>
    </ivory_cluster>
  </remote_servers>
  <zookeeper>
    <node><host>keeper-1</host><port>9181</port></node>
    <node><host>keeper-2</host><port>9181</port></node>
    <node><host>keeper-3</host><port>9181</port></node>
  </zookeeper>
  <macros>
    <cluster>{{cluster}}</cluster>
    <replica>{{name}}</replica>
  </macros>
</clickhouse>
IVORYEOF
exec /entrypoint.sh
'`

func (a *Adapter) DefaultTemplates() []keeper.DeploymentTemplate {
	return []keeper.DeploymentTemplate{
		{
			Platform:    platform.Linux,
			Name:        "ClickHouse (Multi Host)",
			Description: "Three clickhouse replicas in one shard, coordinated through an external ZooKeeper or ClickHouse Keeper ensemble. Edit the replica and keeper lists to match your hosts.",
			Commands: []keeper.DeploymentCommand{
				{Command: deployMultiHost},
				{Command: deployMultiHost},
				{Command: deployMultiHost},
			},
		},
		{
			Platform:    platform.Linux,
			Name:        "ClickHouse (Single Host)",
			Description: "Three clickhouse replicas on one VM. The native port is fixed by the image's config.xml, so this needs a custom image to avoid collisions.",
			Commands: []keeper.DeploymentCommand{
				{Command: deploySingleHost},
				{Command: deploySingleHost},
				{Command: deploySingleHost},
			},
		},
	}
}
