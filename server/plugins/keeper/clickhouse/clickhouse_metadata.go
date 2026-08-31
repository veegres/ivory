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
//
// The script is the image's arguments behind --entrypoint sh, not a tail after
// the image, because the image's own entrypoint ends in exec "$@": a script
// passed the ordinary way is the last thing it runs, so the config file would
// land after the entrypoint has already started a server on the untouched
// config. --entrypoint is a docker option, so the command text still supplies
// options only, and it is the only way to get the file written first.
//
// CLICKHOUSE_DB is deliberately not set. It is the sole trigger for an
// initialization pass whose client is hardcoded to --host 127.0.0.1 with no
// --port, so on a VM running more than one node it talks to another node's
// server and kills this one; the default database exists in any case.
// CLICKHOUSE_USER and CLICKHOUSE_PASSWORD are unaffected - the entrypoint
// writes users.d/default-user.xml earlier and outside that block.
//
// interserver_http_host is stated rather than left to default to the machine's
// hostname: a replica registers its fetch endpoint in zookeeper under that
// name, and where it does not resolve the DDL and the inserts still succeed
// while the other replicas silently never catch up. The shard's replicas carry
// credentials for the same class of reason - without them a distributed query
// fails authentication as soon as a password is set. They are read from env the
// command itself sets, never interpolated: this heredoc is parsed by a shell
// the command starts, where a password holding a `$` or a backtick would be
// expanded or executed.

const deployMultiHost = `docker run -d
  --name {{name}}
  --hostname {{host}}
  --restart unless-stopped
  -p {{dbPort}}:{{dbPort}}
  -v /data/clickhouse:/var/lib/clickhouse
  -e CLICKHOUSE_USER="{{dbUser}}"
  -e CLICKHOUSE_PASSWORD="{{dbPass}}"
  --entrypoint sh
  clickhouse/clickhouse-server:24
  -c '
cat > /etc/clickhouse-server/config.d/ivory-cluster.xml <<IVORYEOF
<clickhouse>
  <interserver_http_host>{{host}}</interserver_http_host>
  <remote_servers>
    <ivory_cluster>
      <shard>
        <replica><host>10.0.0.1</host><port>9000</port><user>$CLICKHOUSE_USER</user><password>$CLICKHOUSE_PASSWORD</password></replica>
        <replica><host>10.0.0.2</host><port>9000</port><user>$CLICKHOUSE_USER</user><password>$CLICKHOUSE_PASSWORD</password></replica>
        <replica><host>10.0.0.3</host><port>9000</port><user>$CLICKHOUSE_USER</user><password>$CLICKHOUSE_PASSWORD</password></replica>
      </shard>
    </ivory_cluster>
  </remote_servers>
  <zookeeper>
    <node><host>10.0.0.1</host><port>2181</port></node>
    <node><host>10.0.0.2</host><port>2181</port></node>
    <node><host>10.0.0.3</host><port>2181</port></node>
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
// later config.d entry overrides config.xml - so each node states its five
// listening ports there and no custom image is needed. http, interserver, mysql
// and postgresql are literal per node because only this template knows they
// have to differ; the native port follows {{dbPort}}, which is the one Ivory
// connects on. mysql and postgresql are there despite nothing using them: the
// image binds 9004 and 9005 whether or not they are wanted, so three nodes
// sharing one port namespace collide on them like any other.
//
// The zookeeper ensemble is the one Ivory's own zookeeper single-host template
// deploys - the same VM, three client ports - so the two shipped templates work
// side by side rather than naming hosts nothing deploys.

const deploySingleHostNode1 = `docker run -d
  --name {{name}}
  --network host
  -e CLICKHOUSE_USER="{{dbUser}}"
  -e CLICKHOUSE_PASSWORD="{{dbPass}}"
  --entrypoint sh
  clickhouse/clickhouse-server:24
  -c '
cat > /etc/clickhouse-server/config.d/ivory-cluster.xml <<IVORYEOF
<clickhouse>
  <tcp_port>{{dbPort}}</tcp_port>
  <http_port>8123</http_port>
  <interserver_http_port>9009</interserver_http_port>
  <mysql_port>9004</mysql_port>
  <postgresql_port>9005</postgresql_port>
  <interserver_http_host>{{host}}</interserver_http_host>
  <remote_servers>
    <ivory_cluster>
      <shard>
        <replica><host>{{host}}</host><port>9000</port><user>$CLICKHOUSE_USER</user><password>$CLICKHOUSE_PASSWORD</password></replica>
        <replica><host>{{host}}</host><port>9001</port><user>$CLICKHOUSE_USER</user><password>$CLICKHOUSE_PASSWORD</password></replica>
        <replica><host>{{host}}</host><port>9002</port><user>$CLICKHOUSE_USER</user><password>$CLICKHOUSE_PASSWORD</password></replica>
      </shard>
    </ivory_cluster>
  </remote_servers>
  <zookeeper>
    <node><host>{{host}}</host><port>2181</port></node>
    <node><host>{{host}}</host><port>2182</port></node>
    <node><host>{{host}}</host><port>2183</port></node>
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
  --entrypoint sh
  clickhouse/clickhouse-server:24
  -c '
cat > /etc/clickhouse-server/config.d/ivory-cluster.xml <<IVORYEOF
<clickhouse>
  <tcp_port>{{dbPort}}</tcp_port>
  <http_port>8124</http_port>
  <interserver_http_port>9010</interserver_http_port>
  <mysql_port>9014</mysql_port>
  <postgresql_port>9015</postgresql_port>
  <interserver_http_host>{{host}}</interserver_http_host>
  <remote_servers>
    <ivory_cluster>
      <shard>
        <replica><host>{{host}}</host><port>9000</port><user>$CLICKHOUSE_USER</user><password>$CLICKHOUSE_PASSWORD</password></replica>
        <replica><host>{{host}}</host><port>9001</port><user>$CLICKHOUSE_USER</user><password>$CLICKHOUSE_PASSWORD</password></replica>
        <replica><host>{{host}}</host><port>9002</port><user>$CLICKHOUSE_USER</user><password>$CLICKHOUSE_PASSWORD</password></replica>
      </shard>
    </ivory_cluster>
  </remote_servers>
  <zookeeper>
    <node><host>{{host}}</host><port>2181</port></node>
    <node><host>{{host}}</host><port>2182</port></node>
    <node><host>{{host}}</host><port>2183</port></node>
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
  --entrypoint sh
  clickhouse/clickhouse-server:24
  -c '
cat > /etc/clickhouse-server/config.d/ivory-cluster.xml <<IVORYEOF
<clickhouse>
  <tcp_port>{{dbPort}}</tcp_port>
  <http_port>8125</http_port>
  <interserver_http_port>9011</interserver_http_port>
  <mysql_port>9024</mysql_port>
  <postgresql_port>9025</postgresql_port>
  <interserver_http_host>{{host}}</interserver_http_host>
  <remote_servers>
    <ivory_cluster>
      <shard>
        <replica><host>{{host}}</host><port>9000</port><user>$CLICKHOUSE_USER</user><password>$CLICKHOUSE_PASSWORD</password></replica>
        <replica><host>{{host}}</host><port>9001</port><user>$CLICKHOUSE_USER</user><password>$CLICKHOUSE_PASSWORD</password></replica>
        <replica><host>{{host}}</host><port>9002</port><user>$CLICKHOUSE_USER</user><password>$CLICKHOUSE_PASSWORD</password></replica>
      </shard>
    </ivory_cluster>
  </remote_servers>
  <zookeeper>
    <node><host>{{host}}</host><port>2181</port></node>
    <node><host>{{host}}</host><port>2182</port></node>
    <node><host>{{host}}</host><port>2183</port></node>
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
			Description: "Three clickhouse replicas in one shard, one per VM, coordinated through a ZooKeeper ensemble on the same three VMs. Replace 10.0.0.1-3 in the replica and keeper lists with the VM addresses.",
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
			Description: "Three clickhouse replicas on one VM, each on its own native, http, interserver, mysql and postgresql port, coordinated through the ZooKeeper ensemble Ivory's own single-host template deploys on the same VM. Deploy that one first, then fill in the host and deploy this.",
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
