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

// Requirements declares only the database pair: the shipped spilo deployment
// leaves patroni's rest api unauthenticated. Spilo names its superuser
// postgres, so the password is the only free choice.
func (a *Adapter) Requirements() keeper.Requirements {
	return keeper.Requirements{
		DbCredentials: true,
		DbUser:        "postgres",
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
  -e ETCD3_HOSTS="etcd1:2379,etcd2:2379,etcd3:2379"
  -e PGPORT={{dbPort}}
  -e APIPORT={{keeperPort}}
  -e PGPASSWORD_SUPERUSER="{{dbPass}}"
  -e RESTAPI_CONNECT_ADDRESS="{{host}}:{{keeperPort}}"
  -e SPILO_CONFIGURATION='{"postgresql":{"connect_address":"{{host}}:{{dbPort}}"},"bootstrap":{"dcs":{"primary_start_timeout":999}}}'
  ghcr.io/zalando/spilo-18:4.1-p2`

// The single-host command drops --hostname: docker rejects it outright
// alongside --network host, which every node here needs to reach the DCS and
// its peers on the host's own interface. Each node is told its own PGPORT and
// APIPORT, which is what keeps three spilos out of each other's way on one
// port namespace.

const deploySingleHost = `docker run -d
  --name {{name}}
  --network host
  -e SCOPE="{{cluster}}"
  -e PATRONI_NAME="{{name}}"
  -e ETCD3_HOSTS="etcd1:2379,etcd2:2379,etcd3:2379"
  -e PGPORT={{dbPort}}
  -e APIPORT={{keeperPort}}
  -e PGPASSWORD_SUPERUSER="{{dbPass}}"
  -e RESTAPI_CONNECT_ADDRESS="{{host}}:{{keeperPort}}"
  -e SPILO_CONFIGURATION='{"postgresql":{"connect_address":"{{host}}:{{dbPort}}"},"bootstrap":{"dcs":{"primary_start_timeout":999}}}'
  ghcr.io/zalando/spilo-18:4.1-p2`

func (a *Adapter) DefaultTemplates() []keeper.DeploymentTemplate {
	return []keeper.DeploymentTemplate{
		{
			Platform:    platform.Docker,
			Name:        "Patroni (Multi Host)",
			Description: "Three spilo nodes, one per VM, coordinating through an external DCS. Point ETCD3_HOSTS at the DCS you run.",
			Commands: []keeper.DeploymentCommand{
				{
					Command:  deployMultiHost,
					Defaults: keeper.DeploymentDefaults{Name: "patroni1", KeeperPort: 8008, DbPort: 5432},
				},
				{
					Command:  deployMultiHost,
					Defaults: keeper.DeploymentDefaults{Name: "patroni2", KeeperPort: 8008, DbPort: 5432},
				},
				{
					Command:  deployMultiHost,
					Defaults: keeper.DeploymentDefaults{Name: "patroni3", KeeperPort: 8008, DbPort: 5432},
				},
			},
		},
		{
			Platform:    platform.Docker,
			Name:        "Patroni (Single Host)",
			Description: "Three spilo nodes on one VM, each on its own keeper and database port, coordinating through an external DCS. Point ETCD3_HOSTS at the DCS you run, fill in the host and deploy.",
			Commands: []keeper.DeploymentCommand{
				{
					Command:  deploySingleHost,
					Defaults: keeper.DeploymentDefaults{Name: "patroni1", KeeperPort: 8008, DbPort: 5432},
				},
				{
					Command:  deploySingleHost,
					Defaults: keeper.DeploymentDefaults{Name: "patroni2", KeeperPort: 8009, DbPort: 5433},
				},
				{
					Command:  deploySingleHost,
					Defaults: keeper.DeploymentDefaults{Name: "patroni3", KeeperPort: 8010, DbPort: 5434},
				},
			},
		},
	}
}
