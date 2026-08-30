package zookeeper

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
		config.ManageNodeKeeperReload:       false,
		config.ManageNodeKeeperFailover:     false,
		config.ManageNodeKeeperActivation:   false,
	}
}

func (a *Adapter) HasLeader() bool { return true }

// ZOO_MY_ID is written literally per command: zookeeper needs a genuinely
// unique small integer per member, which no host-derived value can provide -
// which is why all three commands are spelled out rather than shared. The
// server list has to be identical on every member, or the ensemble splits; its
// trailing ";2181" is the client port, since the image has no separate env var
// for it. ZOO_4LW_COMMANDS_WHITELIST must list mntr and conf explicitly:
// recent versions whitelist only srvr, and the adapter's List/Config need both.

const deployMultiHostNode1 = `docker run -d
  --name {{name}}
  --hostname {{host}}
  --restart unless-stopped
  -p {{dbPort}}:{{dbPort}}
  -p 2888:2888
  -p 3888:3888
  -v /data/zookeeper/data:/data
  -v /data/zookeeper/datalog:/datalog
  -e ZOO_MY_ID="1"
  -e ZOO_SERVERS="server.1=zookeeper1:2888:3888;2181 server.2=zookeeper2:2888:3888;2181 server.3=zookeeper3:2888:3888;2181"
  -e ZOO_4LW_COMMANDS_WHITELIST="mntr,conf,ruok,srvr"
  zookeeper:3.9`

const deployMultiHostNode2 = `docker run -d
  --name {{name}}
  --hostname {{host}}
  --restart unless-stopped
  -p {{dbPort}}:{{dbPort}}
  -p 2888:2888
  -p 3888:3888
  -v /data/zookeeper/data:/data
  -v /data/zookeeper/datalog:/datalog
  -e ZOO_MY_ID="2"
  -e ZOO_SERVERS="server.1=zookeeper1:2888:3888;2181 server.2=zookeeper2:2888:3888;2181 server.3=zookeeper3:2888:3888;2181"
  -e ZOO_4LW_COMMANDS_WHITELIST="mntr,conf,ruok,srvr"
  zookeeper:3.9`

const deployMultiHostNode3 = `docker run -d
  --name {{name}}
  --hostname {{host}}
  --restart unless-stopped
  -p {{dbPort}}:{{dbPort}}
  -p 2888:2888
  -p 3888:3888
  -v /data/zookeeper/data:/data
  -v /data/zookeeper/datalog:/datalog
  -e ZOO_MY_ID="3"
  -e ZOO_SERVERS="server.1=zookeeper1:2888:3888;2181 server.2=zookeeper2:2888:3888;2181 server.3=zookeeper3:2888:3888;2181"
  -e ZOO_4LW_COMMANDS_WHITELIST="mntr,conf,ruok,srvr"
  zookeeper:3.9`

// The single-host commands drop --hostname, which docker rejects alongside
// --network host, and address the ensemble by {{host}} - host networking joins
// no docker network, so a container name resolves to nothing. Three members
// sharing one port namespace need three of everything the image binds:
// ZOO_CFG_EXTRA carries clientPort because the entrypoint hardcodes 2181 and
// the config is read as java properties, where the later key wins; the admin
// server is switched off rather than moved, since it binds 8080 on every
// member and the adapter reads the ensemble over 4lw on the client port
// anyway. The client port is left off the server specs for the same
// last-writer reason - zookeeper rejects a config that sets it in both places.

const deploySingleHostNode1 = `docker run -d
  --name {{name}}
  --network host
  -e ZOO_MY_ID="1"
  -e ZOO_SERVERS="server.1={{host}}:2888:3888 server.2={{host}}:2890:3890 server.3={{host}}:2892:3892"
  -e ZOO_CFG_EXTRA="clientPort={{dbPort}}"
  -e ZOO_ADMINSERVER_ENABLED="false"
  -e ZOO_4LW_COMMANDS_WHITELIST="mntr,conf,ruok,srvr"
  zookeeper:3.9`

const deploySingleHostNode2 = `docker run -d
  --name {{name}}
  --network host
  -e ZOO_MY_ID="2"
  -e ZOO_SERVERS="server.1={{host}}:2888:3888 server.2={{host}}:2890:3890 server.3={{host}}:2892:3892"
  -e ZOO_CFG_EXTRA="clientPort={{dbPort}}"
  -e ZOO_ADMINSERVER_ENABLED="false"
  -e ZOO_4LW_COMMANDS_WHITELIST="mntr,conf,ruok,srvr"
  zookeeper:3.9`

const deploySingleHostNode3 = `docker run -d
  --name {{name}}
  --network host
  -e ZOO_MY_ID="3"
  -e ZOO_SERVERS="server.1={{host}}:2888:3888 server.2={{host}}:2890:3890 server.3={{host}}:2892:3892"
  -e ZOO_CFG_EXTRA="clientPort={{dbPort}}"
  -e ZOO_ADMINSERVER_ENABLED="false"
  -e ZOO_4LW_COMMANDS_WHITELIST="mntr,conf,ruok,srvr"
  zookeeper:3.9`

func (a *Adapter) DefaultTemplates() []keeper.DeploymentTemplate {
	return []keeper.DeploymentTemplate{
		{
			Platform:    platform.Docker,
			Name:        "ZooKeeper (Multi Host)",
			Description: "Three-node zookeeper ensemble, one per VM. Name the nodes zookeeper1..3 or edit the server list to match.",
			Commands: []keeper.DeploymentCommand{
				{
					Command:  deployMultiHostNode1,
					Defaults: keeper.DeploymentCommandDefaults{Name: "zookeeper1", KeeperPort: 2181, DbPort: 2181},
				},
				{
					Command:  deployMultiHostNode2,
					Defaults: keeper.DeploymentCommandDefaults{Name: "zookeeper2", KeeperPort: 2181, DbPort: 2181},
				},
				{
					Command:  deployMultiHostNode3,
					Defaults: keeper.DeploymentCommandDefaults{Name: "zookeeper3", KeeperPort: 2181, DbPort: 2181},
				},
			},
		},
		{
			Platform:    platform.Docker,
			Name:        "ZooKeeper (Single Host)",
			Description: "Three-node zookeeper ensemble on one VM, each node on its own client, quorum and election ports. Fill in the host and deploy.",
			Commands: []keeper.DeploymentCommand{
				{
					Command:  deploySingleHostNode1,
					Defaults: keeper.DeploymentCommandDefaults{Name: "zookeeper1", KeeperPort: 2181, DbPort: 2181},
				},
				{
					Command:  deploySingleHostNode2,
					Defaults: keeper.DeploymentCommandDefaults{Name: "zookeeper2", KeeperPort: 2182, DbPort: 2182},
				},
				{
					Command:  deploySingleHostNode3,
					Defaults: keeper.DeploymentCommandDefaults{Name: "zookeeper3", KeeperPort: 2183, DbPort: 2183},
				},
			},
		},
	}
}
