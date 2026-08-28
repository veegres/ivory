package patroni

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
		config.ViewNodeKeeperConfig:         true,
		config.ManageNodeKeeperConfigUpdate: true,
		config.ManageNodeKeeperSwitchover:   true,
		config.ManageNodeKeeperReinitialize: true,
		config.ManageNodeKeeperRestart:      true,
		config.ManageNodeKeeperReload:       true,
		config.ManageNodeKeeperFailover:     true,
		config.ManageNodeKeeperActivation:   true,
	}
}

// Requirements reports spilo's endpoints: patroni's REST API on its own port
// and postgres beneath it. Spilo names its superuser postgres, so the password
// is the only free choice.
func (a *Adapter) Requirements() keeper.Requirements {
	keeperPort := 8008
	return keeper.Requirements{
		DbPort:      5432,
		KeeperPort:  &keeperPort,
		Credentials: true,
		DbUser:      "postgres",
	}
}

// The DCS address is written literally: patroni coordinates through one the
// user already runs, so only they know where it is.

const deployMultiHost = `docker run -d
  --name {{name}}
  --hostname {{host}}
  --restart unless-stopped
  -p {{keeperPort}}:{{keeperPort}}
  -p {{dbPort}}:{{dbPort}}
  -v /data/postgres:/home/postgres/pgdata
  -e SCOPE="{{cluster}}"
  -e PATRONI_NAME="{{name}}"
  -e ETCD3_HOSTS="etcd-1:2379,etcd-2:2379,etcd-3:2379"
  -e PGPORT={{dbPort}}
  -e APIPORT={{keeperPort}}
  -e PGPASSWORD_SUPERUSER="{{dbPass}}"
  -e RESTAPI_CONNECT_ADDRESS="{{host}}:{{keeperPort}}"
  -e SPILO_CONFIGURATION='{"postgresql":{"connect_address":"{{host}}:{{dbPort}}"},"bootstrap":{"dcs":{"primary_start_timeout":999}}}'
  ghcr.io/zalando/spilo-18:4.1-p2`

const deploySingleHost = `docker run -d
  --name {{name}}
  --hostname {{host}}
  --network host
  -e SCOPE="{{cluster}}"
  -e PATRONI_NAME="{{name}}"
  -e ETCD3_HOSTS="etcd-1:2379,etcd-2:2379,etcd-3:2379"
  -e PGPORT={{dbPort}}
  -e APIPORT={{keeperPort}}
  -e PGPASSWORD_SUPERUSER="{{dbPass}}"
  -e RESTAPI_CONNECT_ADDRESS="{{host}}:{{keeperPort}}"
  -e SPILO_CONFIGURATION='{"postgresql":{"connect_address":"{{host}}:{{dbPort}}"},"bootstrap":{"dcs":{"primary_start_timeout":999}}}'
  ghcr.io/zalando/spilo-18:4.1-p2`

func (a *Adapter) DefaultTemplates() []keeper.DeploymentTemplate {
	return []keeper.DeploymentTemplate{
		{
			Platform:    platform.Linux,
			Name:        "Patroni (Multi Host)",
			Description: "Three spilo nodes, one per VM, coordinating through an external DCS. Point ETCD3_HOSTS at the DCS you run.",
			Commands: []keeper.DeploymentCommand{
				{Command: deployMultiHost},
				{Command: deployMultiHost},
				{Command: deployMultiHost},
			},
		},
		{
			Platform:    platform.Linux,
			Name:        "Patroni (Single Host)",
			Description: "Three spilo nodes on one VM. Give each node its own keeper and database port in the deploy form.",
			Commands: []keeper.DeploymentCommand{
				{Command: deploySingleHost},
				{Command: deploySingleHost},
				{Command: deploySingleHost},
			},
		},
	}
}
