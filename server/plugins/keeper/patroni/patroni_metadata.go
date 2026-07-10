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
		DefaultImage: "ghcr.io/zalando/spilo-18:4.1-p2",
		// NOTE: spilo names its superuser postgres, the password is the only
		// free choice
		Defaults: map[keeper.Var]string{
			keeper.VarKeeperPort: "8008",
			keeper.VarDbPort:     "5432",
			keeper.VarDbUser:     "postgres",
		},
		// NOTE: patroni needs the address of an external DCS it coordinates
		// through; only the user knows where it runs
		Fields: []keeper.FieldSpec{
			{Name: keeper.VarDcs, Label: "DCS (etcd, zookeper, etc)", Example: "etcd1:2379, etcd2:2379, etcd3:2379", Type: keeper.FieldText},
		},
		Ports: []string{keeper.VarKeeperPort, keeper.VarDbPort},
		Volumes: []keeper.VolumeSpec{
			{HostPath: "/data/postgres", ContainerPath: "/home/postgres/pgdata"},
		},
		Env: []keeper.EnvVar{
			{Name: "SCOPE", Value: `"` + keeper.VarCluster + `"`},
			{Name: "PATRONI_NAME", Value: `"` + keeper.VarHost + `"`},
			{Name: "ETCD3_HOSTS", Value: `"` + keeper.VarDcs + `"`},
			{Name: "PGPORT", Value: keeper.VarDbPort},
			{Name: "APIPORT", Value: keeper.VarKeeperPort},
			{Name: "PGPASSWORD_SUPERUSER", Value: `"` + keeper.VarDbPass + `"`},
			{Name: "RESTAPI_CONNECT_ADDRESS", Value: `"` + keeper.VarHost + `:` + keeper.VarKeeperPort + `"`},
			{Name: "SPILO_CONFIGURATION", Value: `'{"postgresql":{"connect_address":"` + keeper.VarHost + `:` + keeper.VarDbPort + `"},"bootstrap":{"dcs":{"primary_start_timeout":999}}}'`},
		},
	}
}
