package patroni

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
		config.ManageNodeKeeperConfigUpdate: true,
		config.ManageNodeKeeperSwitchover:   true,
		config.ManageNodeKeeperReinitialize: true,
		config.ManageNodeKeeperRestart:      true,
		config.ManageNodeKeeperReload:       true,
		config.ManageNodeKeeperFailover:     true,
		config.ManageNodeKeeperActivation:   true,
	}
}

func (p *Plugin) HasLeader() bool { return true }

// The DCS address comes in through {{dcs}}: patroni coordinates through a store
// the user already runs, so only they know where it is, and it is one address
// for the whole cluster rather than something a node knows about itself. The
// multi-host template states no default for it, for the same reason its
// commands state no host - the store runs on machines a template cannot know
// ahead of time, and example text there is a value the operator has to notice
// is wrong. A container name is never one either: it resolves on neither a
// plain docker run across VMs nor host networking.
//
// The node's name goes into SPILO_CONFIGURATION rather than PATRONI_NAME:
// configure_spilo.py never reads that variable, so all three nodes used to
// register under whatever socket.gethostname() returned and collide. It is a
// plain top-level patroni key, and spilo merges this JSON over its generated
// config.

const deployMultiHost = `docker run -d
  --name {{name}}
  --hostname {{host}}
  --restart unless-stopped
  -p {{keeperPort}}:{{keeperPort}}
  -p {{dbPort}}:{{dbPort}}
  -v /data/postgres:/home/postgres/pgdata
  -e SCOPE="{{cluster}}"
  -e ETCD3_HOSTS="{{dcs}}"
  -e PGPORT={{dbPort}}
  -e APIPORT={{keeperPort}}
  -e PGPASSWORD_SUPERUSER="{{dbPass}}"
  -e RESTAPI_CONNECT_ADDRESS="{{host}}:{{keeperPort}}"
  -e SPILO_CONFIGURATION='{"name":"{{name}}","postgresql":{"connect_address":"{{host}}:{{dbPort}}"},"bootstrap":{"dcs":{"primary_start_timeout":999}}}'
  ghcr.io/zalando/spilo-18:4.1-p2`

// The single-host command drops --hostname: docker rejects it outright
// alongside --network host, which every node here needs to reach the DCS and
// its peers on the host's own interface. Each node is told its own PGPORT and
// APIPORT, which is what keeps three spilos out of each other's way on one
// port namespace.
//
// bg_mon is switched off rather than moved: unlike PGPORT/APIPORT, its listen
// port is not configurable through an env var, only by dropping it from
// shared_preload_libraries, and it binds the same fixed port on every node
// under host networking. The override restates spilo's own default list
// (postgres-appliance/scripts/configure_spilo.py) minus bg_mon, so
// pg_stat_statements, pgextwlist, pg_auth_mon and set_user keep loading.
//
// The startup line is what keeps spilo alive without --hostname. Host
// networking leaves the container answering to the VM's own hostname, and
// configure_spilo.py resolves it - getaddrinfo(gethostname()) on its first
// line, before it reads any configuration - so an unresolvable one kills the
// container outright, which is the usual case: most distributions never put
// the machine's own name in /etc/hosts. The container's /etc/hosts is a
// private copy even under host networking, so mapping the name there costs
// nothing and touches nothing on the VM; the address it resolves to is never
// used, only that it resolves. This replaces an --add-host carrying example
// text the operator had to replace with the VM's real hostname, which nothing
// but that operator could know. /bin/sh /launch.sh init is the image's own
// declared command, which the script hands back over to. The {{dcs}} default
// names the ports etcd's own single-host template listens on, through localhost
// rather than the node's host: --network host makes the container's loopback
// the VM's own, which is where that etcd runs.

const deploySingleHost = `docker run -d
  --name {{name}}
  --network host
  -e SCOPE="{{cluster}}"
  -e ETCD3_HOSTS="{{dcs}}"
  -e PGPORT={{dbPort}}
  -e APIPORT={{keeperPort}}
  -e PGPASSWORD_SUPERUSER="{{dbPass}}"
  -e RESTAPI_CONNECT_ADDRESS="{{host}}:{{keeperPort}}"
  -e SPILO_CONFIGURATION='{"name":"{{name}}","postgresql":{"connect_address":"{{host}}:{{dbPort}}"},"bootstrap":{"dcs":{"primary_start_timeout":999,"postgresql":{"parameters":{"shared_preload_libraries":"pg_stat_statements,pgextwlist,pg_auth_mon,set_user"}}}}}'
  ghcr.io/zalando/spilo-18:4.1-p2
  sh -c '
echo "127.0.0.1 $(hostname)" >> /etc/hosts
exec /bin/sh /launch.sh init
'`

func (p *Plugin) DefaultTemplates() []keeper.DeploymentTemplate {
	return []keeper.DeploymentTemplate{
		{
			Platform:    platform.Docker,
			Name:        "Patroni (Multi Host)",
			Description: "Three spilo nodes, one per VM, coordinating through an external DCS. Give the DCS address of the store you run - one host:port per member, comma separated.",
			Defaults:    keeper.DeploymentTemplateDefaults{DbUser: "postgres"},
			Commands: []keeper.DeploymentCommand{
				{
					Command:  deployMultiHost,
					Defaults: keeper.DeploymentCommandDefaults{Name: "patroni1", SshPort: 22, KeeperPort: 8008, DbPort: 5432},
				},
				{
					Command:  deployMultiHost,
					Defaults: keeper.DeploymentCommandDefaults{Name: "patroni2", SshPort: 22, KeeperPort: 8008, DbPort: 5432},
				},
				{
					Command:  deployMultiHost,
					Defaults: keeper.DeploymentCommandDefaults{Name: "patroni3", SshPort: 22, KeeperPort: 8008, DbPort: 5432},
				},
			},
		},
		{
			Platform:    platform.Docker,
			Name:        "Patroni (Single Host)",
			Description: "Three spilo nodes on one VM, each on its own keeper and database port, coordinating through an external DCS - the DCS address is etcd's own single-host template, which ships unauthenticated. Fill in the host and deploy.",
			Defaults: keeper.DeploymentTemplateDefaults{
				DbUser: "postgres",
				Dcs:    "localhost:2479,localhost:2481,localhost:2483",
			},
			Commands: []keeper.DeploymentCommand{
				{
					Command:  deploySingleHost,
					Defaults: keeper.DeploymentCommandDefaults{Name: "patroni1", Host: "localhost", SshPort: 22, KeeperPort: 8008, DbPort: 5442},
				},
				{
					Command:  deploySingleHost,
					Defaults: keeper.DeploymentCommandDefaults{Name: "patroni2", Host: "localhost", SshPort: 22, KeeperPort: 8009, DbPort: 5443},
				},
				{
					Command:  deploySingleHost,
					Defaults: keeper.DeploymentCommandDefaults{Name: "patroni3", Host: "localhost", SshPort: 22, KeeperPort: 8010, DbPort: 5444},
				},
			},
		},
	}
}
