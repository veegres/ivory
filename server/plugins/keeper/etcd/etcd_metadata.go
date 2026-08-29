package etcd

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

// Requirements reports both credential pairs: etcd has no separate management
// API, so the keeper endpoint is the database itself and is asked twice. Auth
// can only be enabled through a user named root.
func (a *Adapter) Requirements() keeper.Requirements {
	return keeper.Requirements{
		KeeperCredentials: true,
		KeeperUser:        "root",
		DbCredentials:     true,
		DbUser:            "root",
	}
}

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

const deploySingleHostNode1 = `docker run -d
  --name {{name}}
  --network host
  -e ETCD_NAME="{{name}}"
  -e ETCD_DATA_DIR="/data/etcd"
  -e ETCD_INITIAL_CLUSTER="etcd1=http://{{host}}:2380,etcd2=http://{{host}}:2382,etcd3=http://{{host}}:2384"
  -e ETCD_INITIAL_CLUSTER_STATE="new"
  -e ETCD_INITIAL_CLUSTER_TOKEN="{{cluster}}"
  -e ETCD_LISTEN_CLIENT_URLS="http://0.0.0.0:{{dbPort}}"
  -e ETCD_ADVERTISE_CLIENT_URLS="http://{{host}}:{{dbPort}}"
  -e ETCD_LISTEN_PEER_URLS="http://0.0.0.0:2380"
  -e ETCD_INITIAL_ADVERTISE_PEER_URLS="http://{{host}}:2380"
  quay.io/coreos/etcd:v3.6.5`

const deploySingleHostNode2 = `docker run -d
  --name {{name}}
  --network host
  -e ETCD_NAME="{{name}}"
  -e ETCD_DATA_DIR="/data/etcd"
  -e ETCD_INITIAL_CLUSTER="etcd1=http://{{host}}:2380,etcd2=http://{{host}}:2382,etcd3=http://{{host}}:2384"
  -e ETCD_INITIAL_CLUSTER_STATE="new"
  -e ETCD_INITIAL_CLUSTER_TOKEN="{{cluster}}"
  -e ETCD_LISTEN_CLIENT_URLS="http://0.0.0.0:{{dbPort}}"
  -e ETCD_ADVERTISE_CLIENT_URLS="http://{{host}}:{{dbPort}}"
  -e ETCD_LISTEN_PEER_URLS="http://0.0.0.0:2382"
  -e ETCD_INITIAL_ADVERTISE_PEER_URLS="http://{{host}}:2382"
  quay.io/coreos/etcd:v3.6.5`

const deploySingleHostNode3 = `docker run -d
  --name {{name}}
  --network host
  -e ETCD_NAME="{{name}}"
  -e ETCD_DATA_DIR="/data/etcd"
  -e ETCD_INITIAL_CLUSTER="etcd1=http://{{host}}:2380,etcd2=http://{{host}}:2382,etcd3=http://{{host}}:2384"
  -e ETCD_INITIAL_CLUSTER_STATE="new"
  -e ETCD_INITIAL_CLUSTER_TOKEN="{{cluster}}"
  -e ETCD_LISTEN_CLIENT_URLS="http://0.0.0.0:{{dbPort}}"
  -e ETCD_ADVERTISE_CLIENT_URLS="http://{{host}}:{{dbPort}}"
  -e ETCD_LISTEN_PEER_URLS="http://0.0.0.0:2384"
  -e ETCD_INITIAL_ADVERTISE_PEER_URLS="http://{{host}}:2384"
  quay.io/coreos/etcd:v3.6.5`

// deployAuth enables authentication once the whole cluster is running: etcd
// has no bootstrap-time credentials, so the root user can only be created
// against an already-formed cluster. That is why it sits on the last command
// rather than the first.
//
// It is three commands rather than one script chained with "&&" because the
// etcd image ships no shell at all - only etcd, etcdctl and etcdutl - so there
// is nothing to interpret a chain. Each runs as its own exec, which is also
// what lets {{dbUser}}/{{dbPass}} be used directly: every argument is filled
// in after the command is split, so nothing parses the value.
const (
	deployAuthAddUser   = `etcdctl --endpoints=http://localhost:{{dbPort}} user add {{dbUser}}:{{dbPass}}`
	deployAuthGrantRole = `etcdctl --endpoints=http://localhost:{{dbPort}} user grant-role {{dbUser}} root`
	deployAuthEnable    = `etcdctl --endpoints=http://localhost:{{dbPort}} auth enable`
)

var deployAuth = []string{deployAuthAddUser, deployAuthGrantRole, deployAuthEnable}

func (a *Adapter) DefaultTemplates() []keeper.DeploymentTemplate {
	return []keeper.DeploymentTemplate{
		{
			Platform:    platform.Docker,
			Name:        "Etcd (Multi Host)",
			Description: "Three-member static etcd cluster, one member per VM. Replace 10.0.0.1-3 with the VM addresses, and name the nodes etcd1..3 or edit the member list to match.",
			Commands: []keeper.DeploymentCommand{
				{
					Command:  deployMultiHost,
					Defaults: keeper.DeploymentDefaults{Name: "etcd1", KeeperPort: 2379, DbPort: 2379},
				},
				{
					Command:  deployMultiHost,
					Defaults: keeper.DeploymentDefaults{Name: "etcd2", KeeperPort: 2379, DbPort: 2379},
				},
				{
					Command:     deployMultiHost,
					PostScripts: deployAuth,
					Defaults:    keeper.DeploymentDefaults{Name: "etcd3", KeeperPort: 2379, DbPort: 2379},
				},
			},
		},
		{
			Platform:    platform.Docker,
			Name:        "Etcd (Single Host)",
			Description: "Three-member etcd cluster on one VM, each member on its own client and peer port. Fill in the host and deploy.",
			Commands: []keeper.DeploymentCommand{
				{
					Command:  deploySingleHostNode1,
					Defaults: keeper.DeploymentDefaults{Name: "etcd1", KeeperPort: 2379, DbPort: 2379},
				},
				{
					Command:  deploySingleHostNode2,
					Defaults: keeper.DeploymentDefaults{Name: "etcd2", KeeperPort: 2381, DbPort: 2381},
				},
				{
					Command:     deploySingleHostNode3,
					PostScripts: deployAuth,
					Defaults:    keeper.DeploymentDefaults{Name: "etcd3", KeeperPort: 2383, DbPort: 2383},
				},
			},
		},
	}
}
