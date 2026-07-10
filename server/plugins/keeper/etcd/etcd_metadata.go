package etcd

import (
	"ivory/core/config"
	"ivory/plugins/keeper"
)

// NOTE: validate that is matches interface in compile-time
var _ keeper.Metadata = (*Adapter)(nil)

func (a *Adapter) SupportedFeatures() map[env.Feature]bool {
	return map[env.Feature]bool{
		env.ViewNodeKeeperOverview:       true,
		env.ViewNodeKeeperConfig:         false,
		env.ManageNodeKeeperConfigUpdate: false,
		env.ManageNodeKeeperSwitchover:   true,
		env.ManageNodeKeeperReinitialize: false,
		env.ManageNodeKeeperRestart:      false,
		env.ManageNodeKeeperReload:       false,
		env.ManageNodeKeeperFailover:     false,
		env.ManageNodeKeeperActivation:   false,
	}
}

// DeploymentSpec bootstraps a static etcd cluster: the client port maps to
// the node database port, the {{peerPort}} field (default 2380, unique per
// node in single-host mode) is the peer listener, and the {{initialCluster}}
// member list (name=http://host:peerPort,...) is derived from the node list.
// Etcd has no bootstrap-time credentials: authentication is enabled after the
// cluster is up via the PostDeploy etcdctl commands, using the database
// credentials as the root user.
func (a *Adapter) DeploymentSpec() keeper.DeploymentSpec {
	return keeper.DeploymentSpec{
		DefaultImage: "quay.io/coreos/etcd:v3.6.5",
		// NOTE: etcd requires a user named root before auth can be enabled
		Defaults: map[keeper.Var]string{
			keeper.VarDbPort: "2379",
			keeper.VarDbUser: "root",
		},
		Fields: []keeper.FieldSpec{
			{Name: keeper.VarPeerPort, Label: "Peer Port", Type: keeper.FieldPort, Default: "2380"},
			{Name: keeper.VarInitialCluster, Label: "Initial Cluster", Type: keeper.FieldText, Template: keeper.VarHost + "=http://" + keeper.VarHost + ":" + keeper.VarPeerPort, Separator: ","},
		},
		PostDeploy: []string{
			"etcdctl --endpoints=http://localhost:" + keeper.VarDbPort + " user add '" + keeper.VarDbUser + ":" + keeper.VarDbPass + "'",
			"etcdctl --endpoints=http://localhost:" + keeper.VarDbPort + " user grant-role " + keeper.VarDbUser + " root",
			"etcdctl --endpoints=http://localhost:" + keeper.VarDbPort + " auth enable",
		},
		Ports: []string{keeper.VarDbPort, keeper.VarPeerPort},
		Volumes: []keeper.VolumeSpec{
			{HostPath: "/data/etcd", ContainerPath: "/data/etcd"},
		},
		Env: []keeper.EnvVar{
			{Name: "ETCD_NAME", Value: `"` + keeper.VarHost + `"`},
			{Name: "ETCD_DATA_DIR", Value: `"/data/etcd"`},
			{Name: "ETCD_INITIAL_CLUSTER", Value: `"` + keeper.VarInitialCluster + `"`},
			{Name: "ETCD_INITIAL_CLUSTER_STATE", Value: `"new"`},
			{Name: "ETCD_INITIAL_CLUSTER_TOKEN", Value: `"` + keeper.VarCluster + `"`},
			{Name: "ETCD_LISTEN_CLIENT_URLS", Value: `"http://0.0.0.0:` + keeper.VarDbPort + `"`},
			{Name: "ETCD_ADVERTISE_CLIENT_URLS", Value: `"http://` + keeper.VarHost + `:` + keeper.VarDbPort + `"`},
			{Name: "ETCD_LISTEN_PEER_URLS", Value: `"http://0.0.0.0:` + keeper.VarPeerPort + `"`},
			{Name: "ETCD_INITIAL_ADVERTISE_PEER_URLS", Value: `"http://` + keeper.VarHost + `:` + keeper.VarPeerPort + `"`},
		},
	}
}
