package clickhouse

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
		config.ManageNodeKeeperFailover:     false,
		config.ManageNodeKeeperActivation:   false,
	}
}

// Every replica accepts writes and they coordinate through ClickHouse
// Keeper/ZooKeeper, so there is no leader to elect.
func (p *Plugin) HasLeader() bool { return false }

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
        <replica><host>clickhouse1</host><port>9000</port></replica>
        <replica><host>clickhouse2</host><port>9000</port></replica>
        <replica><host>clickhouse3</host><port>9000</port></replica>
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

// The single-host commands drop --hostname, which docker rejects alongside
// --network host. The native port is fixed in the image's config.xml and has no
// env var, but the startup script is already writing a config.d file, and a
// later config.d entry overrides config.xml - so each node states its three
// listening ports there and no custom image is needed. http and interserver are
// literal per node because only this template knows they have to differ; the
// native port follows {{dbPort}}, which is the one Ivory connects on.

const deploySingleHostNode1 = `docker run -d
  --name {{name}}
  --network host
  -e CLICKHOUSE_USER="{{dbUser}}"
  -e CLICKHOUSE_PASSWORD="{{dbPass}}"
  -e CLICKHOUSE_DB="default"
  clickhouse/clickhouse-server:24
  sh -c '
cat > /etc/clickhouse-server/config.d/ivory-cluster.xml <<IVORYEOF
<clickhouse>
  <tcp_port>{{dbPort}}</tcp_port>
  <http_port>8123</http_port>
  <interserver_http_port>9009</interserver_http_port>
  <remote_servers>
    <ivory_cluster>
      <shard>
        <replica><host>{{host}}</host><port>9000</port></replica>
        <replica><host>{{host}}</host><port>9001</port></replica>
        <replica><host>{{host}}</host><port>9002</port></replica>
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

const deploySingleHostNode2 = `docker run -d
  --name {{name}}
  --network host
  -e CLICKHOUSE_USER="{{dbUser}}"
  -e CLICKHOUSE_PASSWORD="{{dbPass}}"
  -e CLICKHOUSE_DB="default"
  clickhouse/clickhouse-server:24
  sh -c '
cat > /etc/clickhouse-server/config.d/ivory-cluster.xml <<IVORYEOF
<clickhouse>
  <tcp_port>{{dbPort}}</tcp_port>
  <http_port>8124</http_port>
  <interserver_http_port>9010</interserver_http_port>
  <remote_servers>
    <ivory_cluster>
      <shard>
        <replica><host>{{host}}</host><port>9000</port></replica>
        <replica><host>{{host}}</host><port>9001</port></replica>
        <replica><host>{{host}}</host><port>9002</port></replica>
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

const deploySingleHostNode3 = `docker run -d
  --name {{name}}
  --network host
  -e CLICKHOUSE_USER="{{dbUser}}"
  -e CLICKHOUSE_PASSWORD="{{dbPass}}"
  -e CLICKHOUSE_DB="default"
  clickhouse/clickhouse-server:24
  sh -c '
cat > /etc/clickhouse-server/config.d/ivory-cluster.xml <<IVORYEOF
<clickhouse>
  <tcp_port>{{dbPort}}</tcp_port>
  <http_port>8125</http_port>
  <interserver_http_port>9011</interserver_http_port>
  <remote_servers>
    <ivory_cluster>
      <shard>
        <replica><host>{{host}}</host><port>9000</port></replica>
        <replica><host>{{host}}</host><port>9001</port></replica>
        <replica><host>{{host}}</host><port>9002</port></replica>
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

func (p *Plugin) DefaultTemplates() []keeper.DeploymentTemplate {
	return []keeper.DeploymentTemplate{
		{
			Platform:    platform.Docker,
			Name:        "ClickHouse (Multi Host)",
			Description: "Three clickhouse replicas in one shard, coordinated through an external ZooKeeper or ClickHouse Keeper ensemble. Edit the replica and keeper lists to match your hosts.",
			Defaults:    keeper.DeploymentTemplateDefaults{KeeperUser: "default", DbUser: "default"},
			Commands: []keeper.DeploymentCommand{
				{
					Command:  deployMultiHost,
					Defaults: keeper.DeploymentCommandDefaults{Name: "clickhouse1", SshPort: 22, KeeperPort: 9000, DbPort: 9000},
				},
				{
					Command:  deployMultiHost,
					Defaults: keeper.DeploymentCommandDefaults{Name: "clickhouse2", SshPort: 22, KeeperPort: 9000, DbPort: 9000},
				},
				{
					Command:  deployMultiHost,
					Defaults: keeper.DeploymentCommandDefaults{Name: "clickhouse3", SshPort: 22, KeeperPort: 9000, DbPort: 9000},
				},
			},
		},
		{
			Platform:    platform.Docker,
			Name:        "ClickHouse (Single Host)",
			Description: "Three clickhouse replicas on one VM, each on its own native, http and interserver port, coordinated through an external ZooKeeper or ClickHouse Keeper ensemble. Edit the keeper list to match your hosts, fill in the host and deploy.",
			Defaults:    keeper.DeploymentTemplateDefaults{KeeperUser: "default", DbUser: "default"},
			Commands: []keeper.DeploymentCommand{
				{
					Command:  deploySingleHostNode1,
					Defaults: keeper.DeploymentCommandDefaults{Name: "clickhouse1", SshPort: 22, KeeperPort: 9000, DbPort: 9000},
				},
				{
					Command:  deploySingleHostNode2,
					Defaults: keeper.DeploymentCommandDefaults{Name: "clickhouse2", SshPort: 22, KeeperPort: 9001, DbPort: 9001},
				},
				{
					Command:  deploySingleHostNode3,
					Defaults: keeper.DeploymentCommandDefaults{Name: "clickhouse3", SshPort: 22, KeeperPort: 9002, DbPort: 9002},
				},
			},
		},
	}
}
