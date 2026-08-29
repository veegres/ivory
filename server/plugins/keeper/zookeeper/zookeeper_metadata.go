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

// Requirements reports the client port only. Zookeeper ships with no auth, so
// the deployment consumes no credentials - operators who need it configure it
// on the deployed ensemble themselves.
func (a *Adapter) Requirements() keeper.Requirements {
	return keeper.Requirements{DbPort: 2181, KeeperPort: 2181}
}

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
  -e ZOO_SERVERS="server.1=zookeeper-1:2888:3888;2181 server.2=zookeeper-2:2888:3888;2181 server.3=zookeeper-3:2888:3888;2181"
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
  -e ZOO_SERVERS="server.1=zookeeper-1:2888:3888;2181 server.2=zookeeper-2:2888:3888;2181 server.3=zookeeper-3:2888:3888;2181"
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
  -e ZOO_SERVERS="server.1=zookeeper-1:2888:3888;2181 server.2=zookeeper-2:2888:3888;2181 server.3=zookeeper-3:2888:3888;2181"
  -e ZOO_4LW_COMMANDS_WHITELIST="mntr,conf,ruok,srvr"
  zookeeper:3.9`

const deploySingleHostNode1 = `docker run -d
  --name {{name}}
  --hostname {{host}}
  --network host
  -e ZOO_MY_ID="1"
  -e ZOO_SERVERS="server.1=zookeeper-1:2888:3888;2181 server.2=zookeeper-2:2890:3890;2181 server.3=zookeeper-3:2892:3892;2181"
  -e ZOO_4LW_COMMANDS_WHITELIST="mntr,conf,ruok,srvr"
  zookeeper:3.9`

const deploySingleHostNode2 = `docker run -d
  --name {{name}}
  --hostname {{host}}
  --network host
  -e ZOO_MY_ID="2"
  -e ZOO_SERVERS="server.1=zookeeper-1:2888:3888;2181 server.2=zookeeper-2:2890:3890;2181 server.3=zookeeper-3:2892:3892;2181"
  -e ZOO_4LW_COMMANDS_WHITELIST="mntr,conf,ruok,srvr"
  zookeeper:3.9`

const deploySingleHostNode3 = `docker run -d
  --name {{name}}
  --hostname {{host}}
  --network host
  -e ZOO_MY_ID="3"
  -e ZOO_SERVERS="server.1=zookeeper-1:2888:3888;2181 server.2=zookeeper-2:2890:3890;2181 server.3=zookeeper-3:2892:3892;2181"
  -e ZOO_4LW_COMMANDS_WHITELIST="mntr,conf,ruok,srvr"
  zookeeper:3.9`

func (a *Adapter) DefaultTemplates() []keeper.DeploymentTemplate {
	return []keeper.DeploymentTemplate{
		{
			Platform:    platform.Docker,
			Name:        "ZooKeeper (Multi Host)",
			Description: "Three-node zookeeper ensemble, one per VM. Name the nodes zookeeper-1..3 or edit the server list to match.",
			Commands: []keeper.DeploymentCommand{
				{Command: deployMultiHostNode1},
				{Command: deployMultiHostNode2},
				{Command: deployMultiHostNode3},
			},
		},
		{
			Platform:    platform.Docker,
			Name:        "ZooKeeper (Single Host)",
			Description: "Three-node zookeeper ensemble on one VM. Each node uses its own quorum and election ports; give each its own client port in the deploy form.",
			Commands: []keeper.DeploymentCommand{
				{Command: deploySingleHostNode1},
				{Command: deploySingleHostNode2},
				{Command: deploySingleHostNode3},
			},
		},
	}
}
