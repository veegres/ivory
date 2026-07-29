package etcd

import (
	"ivory/core/config"
	"ivory/plugins/keeper"
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
			{Name: keeper.VarInitialCluster, Label: "Initial Cluster", Type: keeper.FieldText, Template: string(keeper.VarHost) + "=http://" + string(keeper.VarHost) + ":" + string(keeper.VarPeerPort), Separator: ","},
		},
		PostDeploy: []string{
			"etcdctl --endpoints=http://localhost:" + string(keeper.VarDbPort) + " user add '" + string(keeper.VarDbUser) + ":" + string(keeper.VarDbPass) + "'",
			"etcdctl --endpoints=http://localhost:" + string(keeper.VarDbPort) + " user grant-role " + string(keeper.VarDbUser) + " root",
			"etcdctl --endpoints=http://localhost:" + string(keeper.VarDbPort) + " auth enable",
		},
		Ports: []string{string(keeper.VarDbPort), string(keeper.VarPeerPort)},
		Volumes: []keeper.VolumeSpec{
			{HostPath: "/data/etcd", ContainerPath: "/data/etcd"},
		},
		Env: []keeper.EnvVar{
			{Name: "ETCD_NAME", Value: `"` + string(keeper.VarHost) + `"`},
			{Name: "ETCD_DATA_DIR", Value: `"/data/etcd"`},
			{Name: "ETCD_INITIAL_CLUSTER", Value: `"` + string(keeper.VarInitialCluster) + `"`},
			{Name: "ETCD_INITIAL_CLUSTER_STATE", Value: `"new"`},
			{Name: "ETCD_INITIAL_CLUSTER_TOKEN", Value: `"` + string(keeper.VarCluster) + `"`},
			{Name: "ETCD_LISTEN_CLIENT_URLS", Value: `"http://0.0.0.0:` + string(keeper.VarDbPort) + `"`},
			{Name: "ETCD_ADVERTISE_CLIENT_URLS", Value: `"http://` + string(keeper.VarHost) + `:` + string(keeper.VarDbPort) + `"`},
			{Name: "ETCD_LISTEN_PEER_URLS", Value: `"http://0.0.0.0:` + string(keeper.VarPeerPort) + `"`},
			{Name: "ETCD_INITIAL_ADVERTISE_PEER_URLS", Value: `"http://` + string(keeper.VarHost) + `:` + string(keeper.VarPeerPort) + `"`},
		},
	}
}
