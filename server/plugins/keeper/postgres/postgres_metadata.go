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
		DefaultImage: "postgres:18",
		// NOTE: the empty username means credentials are consumed but the
		// username is the user's choice
		Defaults: map[keeper.Var]string{
			keeper.VarDbPort: "5432",
			keeper.VarDbUser: "",
		},
		Ports: []string{keeper.VarDbPort},
		Volumes: []keeper.VolumeSpec{
			{HostPath: "/data/postgres", ContainerPath: "/var/lib/postgresql/data"},
		},
		Env: []keeper.EnvVar{
			{Name: "PGPORT", Value: `"` + keeper.VarDbPort + `"`},
			{Name: "POSTGRES_USER", Value: `"` + keeper.VarDbUser + `"`},
			{Name: "POSTGRES_PASSWORD", Value: `"` + keeper.VarDbPass + `"`},
		},
	}
}
