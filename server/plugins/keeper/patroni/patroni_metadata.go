package patroni

import (
	"ivory/core/config"
	"ivory/plugins/keeper"
)

// NOTE: validate that is matches interface in compile-time
var _ keeper.Metadata = (*Adapter)(nil)

func (a *Adapter) SupportedFeatures() map[env.Feature]bool {
	return map[env.Feature]bool{
		env.ViewNodeKeeperOverview:       true,
		env.ViewNodeKeeperConfig:         true,
		env.ManageNodeKeeperConfigUpdate: true,
		env.ManageNodeKeeperSwitchover:   true,
		env.ManageNodeKeeperReinitialize: true,
		env.ManageNodeKeeperRestart:      true,
		env.ManageNodeKeeperReload:       true,
		env.ManageNodeKeeperFailover:     true,
		env.ManageNodeKeeperActivation:   true,
	}
}

func (a *Adapter) DeploymentSpec() keeper.DeploymentSpec {
	return keeper.DeploymentSpec{
		DefaultImage:  "ghcr.io/zalando/spilo-18:4.1-p2",
		DefaultValues: map[string]string{"username": "postgres"},
		Ports:         []string{"{{keeperPort}}", "{{dbPort}}"},
		Volumes: []keeper.VolumeSpec{
			{HostPath: "/data/postgres", ContainerPath: "/home/postgres/pgdata"},
		},
		Env: []keeper.EnvVar{
			{Name: "SCOPE", Value: `"{{cluster}}"`},
			{Name: "PATRONI_NAME", Value: `"{{host}}"`},
			{Name: "ETCD3_HOSTS", Value: `"{{dcs}}"`},
			{Name: "PGPORT", Value: `{{dbPort}}`},
			{Name: "APIPORT", Value: `{{keeperPort}}`},
			{Name: "PGPASSWORD_SUPERUSER", Value: `"{{dbPass}}"`},
			{Name: "RESTAPI_CONNECT_ADDRESS", Value: `"{{host}}:{{keeperPort}}"`},
			{Name: "SPILO_CONFIGURATION", Value: `'{"postgresql":{"connect_address":"{{host}}:{{dbPort}}"},"bootstrap":{"dcs":{"primary_start_timeout":999}}}'`},
		},
	}
}
