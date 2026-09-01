package etcd

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
		config.ViewNodeKeeperConfig:         false,
		config.ManageNodeKeeperConfigUpdate: false,
		config.ManageNodeKeeperSwitchover:   true,
		config.ManageNodeKeeperReinitialize: false,
		config.ManageNodeKeeperRestart:      false,
		config.ManageNodeKeeperReload:       false,
		config.ManageNodeKeeperFailover:     false,
		config.ManageNodeKeeperActivation:   false,
	}
}

func (p *Plugin) HasLeader() bool { return true }

// The peer port, the member list and every other value only the operator knows
// are written literally: they are plain text to read and edit, not variables.
// Only the node's own identity and endpoints are interpolated.
//
// The multi-host member list is addresses, not container names: each member is
// a plain docker run on its own VM with no shared network between them, so
// nothing resolves the others by name. The 10.0.0.x are example text the
// operator replaces with their own VMs.

const deployMultiHost = `docker run -d
  --name {{name}}
  --hostname {{host}}
  --restart unless-stopped
  -p {{dbPort}}:{{dbPort}}
  -p 2380:2380
  -v /data/etcd:/data/etcd
  -e ETCD_NAME="{{name}}"
  -e ETCD_DATA_DIR="/data/etcd"
  -e ETCD_INITIAL_CLUSTER="etcd1=http://10.0.0.1:2380,etcd2=http://10.0.0.2:2380,etcd3=http://10.0.0.3:2380"
  -e ETCD_INITIAL_CLUSTER_STATE="new"
  -e ETCD_INITIAL_CLUSTER_TOKEN="{{cluster}}"
  -e ETCD_LISTEN_CLIENT_URLS="http://0.0.0.0:{{dbPort}}"
  -e ETCD_ADVERTISE_CLIENT_URLS="http://{{host}}:{{dbPort}}"
  -e ETCD_LISTEN_PEER_URLS="http://0.0.0.0:2380"
  -e ETCD_INITIAL_ADVERTISE_PEER_URLS="http://{{host}}:2380"
  quay.io/coreos/etcd:v3.6.5`

// The single-host commands drop --hostname: docker rejects it outright
// alongside --network host, which every member here needs to reach its peers
// on the host's own interface. The member list addresses them by {{host}}
// rather than by container name for the same reason - host networking joins no
// docker network, so there is no embedded dns to resolve a name against.
//
// The ports are 2479-2484 rather than the conventional 2379/2380, and neither
// family may sit on an etcd default. etcd replaces an advertise url that is
// byte-identical to its own built-in default with the detected default host, so
// with {{host}}=localhost the first member - and only it - has both of its urls
// rewritten to the machine's address. On the peer url that is fatal: it no
// longer matches its own entry in the member list and the node dies with an
// error that never mentions localhost. On the client url it is quieter and
// worse: the member comes up and registers a client address the cluster was
// never configured with, so Ivory's overview reports that node unreachable and
// invents a second row for the address it actually advertised. Any non-default
// port sidesteps both. Multi-host keeps 2379/2380 - it sets --hostname {{host}},
// so its urls are never the default to begin with.

const deploySingleHostNode1 = `docker run -d
  --name {{name}}
  --network host
  -e ETCD_NAME="{{name}}"
  -e ETCD_DATA_DIR="/data/etcd"
  -e ETCD_INITIAL_CLUSTER="etcd1=http://{{host}}:2480,etcd2=http://{{host}}:2482,etcd3=http://{{host}}:2484"
  -e ETCD_INITIAL_CLUSTER_STATE="new"
  -e ETCD_INITIAL_CLUSTER_TOKEN="{{cluster}}"
  -e ETCD_LISTEN_CLIENT_URLS="http://0.0.0.0:{{dbPort}}"
  -e ETCD_ADVERTISE_CLIENT_URLS="http://{{host}}:{{dbPort}}"
  -e ETCD_LISTEN_PEER_URLS="http://0.0.0.0:2480"
  -e ETCD_INITIAL_ADVERTISE_PEER_URLS="http://{{host}}:2480"
  quay.io/coreos/etcd:v3.6.5`

const deploySingleHostNode2 = `docker run -d
  --name {{name}}
  --network host
  -e ETCD_NAME="{{name}}"
  -e ETCD_DATA_DIR="/data/etcd"
  -e ETCD_INITIAL_CLUSTER="etcd1=http://{{host}}:2480,etcd2=http://{{host}}:2482,etcd3=http://{{host}}:2484"
  -e ETCD_INITIAL_CLUSTER_STATE="new"
  -e ETCD_INITIAL_CLUSTER_TOKEN="{{cluster}}"
  -e ETCD_LISTEN_CLIENT_URLS="http://0.0.0.0:{{dbPort}}"
  -e ETCD_ADVERTISE_CLIENT_URLS="http://{{host}}:{{dbPort}}"
  -e ETCD_LISTEN_PEER_URLS="http://0.0.0.0:2482"
  -e ETCD_INITIAL_ADVERTISE_PEER_URLS="http://{{host}}:2482"
  quay.io/coreos/etcd:v3.6.5`

const deploySingleHostNode3 = `docker run -d
  --name {{name}}
  --network host
  -e ETCD_NAME="{{name}}"
  -e ETCD_DATA_DIR="/data/etcd"
  -e ETCD_INITIAL_CLUSTER="etcd1=http://{{host}}:2480,etcd2=http://{{host}}:2482,etcd3=http://{{host}}:2484"
  -e ETCD_INITIAL_CLUSTER_STATE="new"
  -e ETCD_INITIAL_CLUSTER_TOKEN="{{cluster}}"
  -e ETCD_LISTEN_CLIENT_URLS="http://0.0.0.0:{{dbPort}}"
  -e ETCD_ADVERTISE_CLIENT_URLS="http://{{host}}:{{dbPort}}"
  -e ETCD_LISTEN_PEER_URLS="http://0.0.0.0:2484"
  -e ETCD_INITIAL_ADVERTISE_PEER_URLS="http://{{host}}:2484"
  quay.io/coreos/etcd:v3.6.5`

// Both templates ship unauthenticated, and name no keeper or database user. A
// cluster that enables auth in a post-script is a cluster the shipped patroni
// template cannot use as its DCS: patroni sends no etcd credentials, and etcd
// answers an unauthenticated client with a bare gRPC InvalidArgument that never
// mentions authentication. Adding the accounts afterwards is one etcdctl user
// add / auth enable away, and is the operator's call rather than a default.
func (p *Plugin) DefaultTemplates() []keeper.DeploymentTemplate {
	return []keeper.DeploymentTemplate{
		{
			Platform:    platform.Docker,
			Name:        "Etcd (Multi Host)",
			Description: "Three-member static etcd cluster, one member per VM, without authentication. Replace 10.0.0.1-3 with the VM addresses, and name the nodes etcd1..3 or edit the member list to match.",
			Commands: []keeper.DeploymentCommand{
				{
					Command:  deployMultiHost,
					Defaults: keeper.DeploymentCommandDefaults{Name: "etcd1", SshPort: 22, KeeperPort: 2379, DbPort: 2379},
				},
				{
					Command:  deployMultiHost,
					Defaults: keeper.DeploymentCommandDefaults{Name: "etcd2", SshPort: 22, KeeperPort: 2379, DbPort: 2379},
				},
				{
					Command:  deployMultiHost,
					Defaults: keeper.DeploymentCommandDefaults{Name: "etcd3", SshPort: 22, KeeperPort: 2379, DbPort: 2379},
				},
			},
		},
		{
			Platform:    platform.Docker,
			Name:        "Etcd (Single Host)",
			Description: "Three-member etcd cluster on one VM, each member on its own client and peer port (2479-2484, off etcd's own defaults so localhost works as the host), without authentication. Fill in the host and deploy.",
			Commands: []keeper.DeploymentCommand{
				{
					Command:  deploySingleHostNode1,
					Defaults: keeper.DeploymentCommandDefaults{Name: "etcd1", Host: "localhost", SshPort: 22, KeeperPort: 2479, DbPort: 2479},
				},
				{
					Command:  deploySingleHostNode2,
					Defaults: keeper.DeploymentCommandDefaults{Name: "etcd2", Host: "localhost", SshPort: 22, KeeperPort: 2481, DbPort: 2481},
				},
				{
					Command:  deploySingleHostNode3,
					Defaults: keeper.DeploymentCommandDefaults{Name: "etcd3", Host: "localhost", SshPort: 22, KeeperPort: 2483, DbPort: 2483},
				},
			},
		},
	}
}
