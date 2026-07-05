package etcd

import (
	"ivory/core/config"
	"ivory/plugins/keeper"
)

// NOTE: validate that is matches interface in compile-time
var _ keeper.Metadata = (*Adapter)(nil)

func (a *Adapter) SupportedFeatures() []env.Feature {
	return []env.Feature{
		env.ViewNodeKeeperOverview,
		env.ManageNodeKeeperSwitchover,
	}
}

// DeploymentSpec bootstraps a static etcd cluster: the client port maps to
// the node database port and the peer port stays on 2380. The deploy "DCS"
// field carries the initial cluster string (name=http://host:2380,...).
func (a *Adapter) DeploymentSpec() keeper.DeploymentSpec {
	return keeper.DeploymentSpec{
		DefaultImage: "quay.io/coreos/etcd:v3.6.5",
		Ports:        []string{"{{dbPort}}", "2380"},
		Volumes: []keeper.VolumeSpec{
			{HostPath: "/data/etcd", ContainerPath: "/data/etcd"},
		},
		Env: []keeper.EnvVar{
			{Name: "ETCD_NAME", Value: `"{{host}}"`},
			{Name: "ETCD_DATA_DIR", Value: `"/data/etcd"`},
			{Name: "ETCD_INITIAL_CLUSTER", Value: `"{{dcs}}"`},
			{Name: "ETCD_INITIAL_CLUSTER_STATE", Value: `"new"`},
			{Name: "ETCD_INITIAL_CLUSTER_TOKEN", Value: `"{{cluster}}"`},
			{Name: "ETCD_LISTEN_CLIENT_URLS", Value: `"http://0.0.0.0:{{dbPort}}"`},
			{Name: "ETCD_ADVERTISE_CLIENT_URLS", Value: `"http://{{host}}:{{dbPort}}"`},
			{Name: "ETCD_LISTEN_PEER_URLS", Value: `"http://0.0.0.0:2380"`},
			{Name: "ETCD_INITIAL_ADVERTISE_PEER_URLS", Value: `"http://{{host}}:2380"`},
		},
	}
}
