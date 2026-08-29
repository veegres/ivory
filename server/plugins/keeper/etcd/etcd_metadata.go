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

// Requirements reports etcd's client port as the database endpoint; etcd has no
// separate management API, so the keeper endpoint is the database itself. Auth
// can only be enabled through a user named root.
func (a *Adapter) Requirements() keeper.Requirements {
	return keeper.Requirements{
		DbPort:            2379,
		KeeperPort:        2379,
		KeeperCredentials: true,
		KeeperUser:        "root",
		DbCredentials:     true,
		DbUser:            "root",
	}
}

// The peer port, the member list and every other value only the operator knows
// are written literally: they are plain text to read and edit, not variables.
// Only the node's own identity and endpoints are interpolated.

const deployMultiHost = `docker run -d
  --name {{name}}
  --hostname {{host}}
  --restart unless-stopped
  -p {{dbPort}}:{{dbPort}}
  -p 2380:2380
  -v /data/etcd:/data/etcd
  -e ETCD_NAME="{{name}}"
  -e ETCD_DATA_DIR="/data/etcd"
  -e ETCD_INITIAL_CLUSTER="etcd-1=http://etcd-1:2380,etcd-2=http://etcd-2:2380,etcd-3=http://etcd-3:2380"
  -e ETCD_INITIAL_CLUSTER_STATE="new"
  -e ETCD_INITIAL_CLUSTER_TOKEN="{{cluster}}"
  -e ETCD_LISTEN_CLIENT_URLS="http://0.0.0.0:{{dbPort}}"
  -e ETCD_ADVERTISE_CLIENT_URLS="http://{{host}}:{{dbPort}}"
  -e ETCD_LISTEN_PEER_URLS="http://0.0.0.0:2380"
  -e ETCD_INITIAL_ADVERTISE_PEER_URLS="http://{{host}}:2380"
  quay.io/coreos/etcd:v3.6.5`

const deploySingleHostNode1 = `docker run -d
  --name {{name}}
  --hostname {{host}}
  --network host
  -e ETCD_NAME="{{name}}"
  -e ETCD_DATA_DIR="/data/etcd"
  -e ETCD_INITIAL_CLUSTER="etcd-1=http://etcd-1:2380,etcd-2=http://etcd-2:2382,etcd-3=http://etcd-3:2384"
  -e ETCD_INITIAL_CLUSTER_STATE="new"
  -e ETCD_INITIAL_CLUSTER_TOKEN="{{cluster}}"
  -e ETCD_LISTEN_CLIENT_URLS="http://0.0.0.0:{{dbPort}}"
  -e ETCD_ADVERTISE_CLIENT_URLS="http://{{host}}:{{dbPort}}"
  -e ETCD_LISTEN_PEER_URLS="http://0.0.0.0:2380"
  -e ETCD_INITIAL_ADVERTISE_PEER_URLS="http://{{host}}:2380"
  quay.io/coreos/etcd:v3.6.5`

const deploySingleHostNode2 = `docker run -d
  --name {{name}}
  --hostname {{host}}
  --network host
  -e ETCD_NAME="{{name}}"
  -e ETCD_DATA_DIR="/data/etcd"
  -e ETCD_INITIAL_CLUSTER="etcd-1=http://etcd-1:2380,etcd-2=http://etcd-2:2382,etcd-3=http://etcd-3:2384"
  -e ETCD_INITIAL_CLUSTER_STATE="new"
  -e ETCD_INITIAL_CLUSTER_TOKEN="{{cluster}}"
  -e ETCD_LISTEN_CLIENT_URLS="http://0.0.0.0:{{dbPort}}"
  -e ETCD_ADVERTISE_CLIENT_URLS="http://{{host}}:{{dbPort}}"
  -e ETCD_LISTEN_PEER_URLS="http://0.0.0.0:2382"
  -e ETCD_INITIAL_ADVERTISE_PEER_URLS="http://{{host}}:2382"
  quay.io/coreos/etcd:v3.6.5`

const deploySingleHostNode3 = `docker run -d
  --name {{name}}
  --hostname {{host}}
  --network host
  -e ETCD_NAME="{{name}}"
  -e ETCD_DATA_DIR="/data/etcd"
  -e ETCD_INITIAL_CLUSTER="etcd-1=http://etcd-1:2380,etcd-2=http://etcd-2:2382,etcd-3=http://etcd-3:2384"
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
// rather than the first. The steps are chained with "&&" so a failure partway
// through stops the rest instead of silently continuing.
const deployAuth = `sh -c '
etcdctl --endpoints=http://localhost:{{dbPort}} user add "{{dbUser}}:{{dbPass}}" &&
etcdctl --endpoints=http://localhost:{{dbPort}} user grant-role {{dbUser}} root &&
etcdctl --endpoints=http://localhost:{{dbPort}} auth enable
'`

func (a *Adapter) DefaultTemplates() []keeper.DeploymentTemplate {
	return []keeper.DeploymentTemplate{
		{
			Platform:    platform.Docker,
			Name:        "Etcd (Multi Host)",
			Description: "Three-member static etcd cluster, one member per VM. Name the nodes etcd-1..3 or edit the member list to match.",
			Commands: []keeper.DeploymentCommand{
				{Command: deployMultiHost},
				{Command: deployMultiHost},
				{Command: deployMultiHost, PostScript: deployAuth},
			},
		},
		{
			Platform:    platform.Docker,
			Name:        "Etcd (Single Host)",
			Description: "Three-member etcd cluster on one VM. Each member peers on its own port; give each its own client port in the deploy form.",
			Commands: []keeper.DeploymentCommand{
				{Command: deploySingleHostNode1},
				{Command: deploySingleHostNode2},
				{Command: deploySingleHostNode3, PostScript: deployAuth},
			},
		},
	}
}
