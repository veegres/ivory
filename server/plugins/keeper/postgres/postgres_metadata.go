package postgres

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
		env.ManageNodeKeeperConfigUpdate: false,
		env.ManageNodeKeeperSwitchover:   false,
		env.ManageNodeKeeperReinitialize: false,
		env.ManageNodeKeeperRestart:      false,
		env.ManageNodeKeeperReload:       true,
		env.ManageNodeKeeperFailover:     true,
		env.ManageNodeKeeperActivation:   false,
	}
}

func (a *Adapter) DeploymentSpec() keeper.DeploymentSpec {
	return keeper.DeploymentSpec{
		DefaultImage:  "postgres:18",
		DefaultValues: map[string]string{"dcs": "empty"},
		Ports:         []string{"{{dbPort}}"},
		Volumes: []keeper.VolumeSpec{
			{HostPath: "/data/postgres", ContainerPath: "/var/lib/postgresql/data"},
		},
		Env: []keeper.EnvVar{
			{Name: "PGPORT", Value: `"{{dbPort}}"`},
			{Name: "POSTGRES_USER", Value: `"{{dbUser}}"`},
			{Name: "POSTGRES_PASSWORD", Value: `"{{dbPass}}"`},
		},
	}
}
