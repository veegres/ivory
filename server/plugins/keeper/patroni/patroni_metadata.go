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

// The DCS address is written literally: patroni coordinates through one the
// user already runs, so only they know where it is. The 10.0.0.x are example
// text the operator replaces - a container name resolves on neither a plain
// docker run across VMs nor host networking.
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
  -e ETCD3_HOSTS="10.0.0.1:2379,10.0.0.2:2379,10.0.0.3:2379"
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
// declared command, which the script hands back over to. The DCS ports are the
// ones etcd's own single-host template listens on.

const deploySingleHost = `docker run -d
  --name {{name}}
  --network host
  -e SCOPE="{{cluster}}"
  -e ETCD3_HOSTS="{{host}}:2479,{{host}}:2481,{{host}}:2483"
  -e PGPORT={{dbPort}}
  -e APIPORT={{keeperPort}}
  -e PGPASSWORD_SUPERUSER="{{dbPass}}"
  -e RESTAPI_CONNECT_ADDRESS="{{host}}:{{keeperPort}}"
  -e SPILO_CONFIGURATION='{"name":"{{name}}","postgresql":{"connect_address":"{{host}}:{{dbPort}}"},"bootstrap":{"dcs":{"primary_start_timeout":999}}}'
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
			Description: "Three spilo nodes, one per VM, coordinating through an external DCS. Point ETCD3_HOSTS at the DCS you run - the 10.0.0.x are examples.",
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
			Description: "Three spilo nodes on one VM, each on its own keeper and database port, coordinating through an external DCS - the ports are etcd's own single-host template, which ships unauthenticated. Fill in the host and deploy.",
			Defaults:    keeper.DeploymentTemplateDefaults{DbUser: "postgres"},
			Commands: []keeper.DeploymentCommand{
				{
					Command:  deploySingleHost,
					Defaults: keeper.DeploymentCommandDefaults{Name: "patroni1", SshPort: 22, KeeperPort: 8008, DbPort: 5432},
				},
				{
					Command:  deploySingleHost,
					Defaults: keeper.DeploymentCommandDefaults{Name: "patroni2", SshPort: 22, KeeperPort: 8009, DbPort: 5433},
				},
				{
					Command:  deploySingleHost,
					Defaults: keeper.DeploymentCommandDefaults{Name: "patroni3", SshPort: 22, KeeperPort: 8010, DbPort: 5434},
				},
			},
		},
	}
}
