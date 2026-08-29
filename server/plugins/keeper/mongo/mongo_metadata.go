package mongo

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
		config.ManageNodeKeeperSwitchover:   false,
		config.ManageNodeKeeperReinitialize: false,
		config.ManageNodeKeeperRestart:      false,
		config.ManageNodeKeeperReload:       false,
		config.ManageNodeKeeperFailover:     true,
		config.ManageNodeKeeperActivation:   false,
	}
}

// Requirements deliberately declares no credentials: enabling client auth on a
// replica set also requires internal authentication between members (a shared
// keyfile mounted into every container), which the deploy model has no
// mechanism for yet. Adapter.connect still accepts credentials when an operator
// configures auth on the deployed cluster themselves.
func (a *Adapter) Requirements() keeper.Requirements {
	return keeper.Requirements{}
}

// Every member starts as a plain mongod that already knows its replica set
// name; who becomes primary is decided later, by deployInitiate. The official
// image has no env var for --replSet, so this replaces the image's default
// command while still going through docker-entrypoint.sh so its normal startup
// housekeeping still runs.

const deployMultiHost = `docker run -d
  --name {{name}}
  --hostname {{host}}
  --restart unless-stopped
  -p {{dbPort}}:{{dbPort}}
  -v /data/mongo:/data/db
  mongo:8
  sh -c '
exec docker-entrypoint.sh mongod --replSet "{{cluster}}" --port {{dbPort}} --bind_ip_all
'`

// The single-host command drops --hostname, which docker rejects alongside
// --network host, and each member is given its own --port so three mongods can
// share the one port namespace.

const deploySingleHost = `docker run -d
  --name {{name}}
  --network host
  mongo:8
  sh -c '
exec docker-entrypoint.sh mongod --replSet "{{cluster}}" --port {{dbPort}} --bind_ip_all
'`

// deployInitiate turns the standalone mongod processes into an actual replica
// set. It sits on the last command because it needs every member already
// running, and the member list is written literally since only the operator
// knows the hosts. The double quotes inside --eval are backslash-escaped
// because the whole script is itself a single-quoted sh -c argument, where a
// literal single quote would terminate it early.
const deployInitiate = `sh -c '
until mongosh --quiet --port {{dbPort}} --eval "1" >/dev/null 2>&1; do sleep 1; done
mongosh --quiet --port {{dbPort}} --eval "rs.initiate({_id: \"{{cluster}}\", members: [{_id: 0, host: \"mongo1:27017\"}, {_id: 1, host: \"mongo2:27017\"}, {_id: 2, host: \"mongo3:27017\"}]})"
'`

// deploySingleHostInitiate is deployInitiate's twin for one VM: the members
// differ only in port, and they are addressed by {{host}} because host
// networking resolves no container names.
const deploySingleHostInitiate = `sh -c '
until mongosh --quiet --port {{dbPort}} --eval "1" >/dev/null 2>&1; do sleep 1; done
mongosh --quiet --port {{dbPort}} --eval "rs.initiate({_id: \"{{cluster}}\", members: [{_id: 0, host: \"{{host}}:27017\"}, {_id: 1, host: \"{{host}}:27018\"}, {_id: 2, host: \"{{host}}:27019\"}]})"
'`

func (a *Adapter) DefaultTemplates() []keeper.DeploymentTemplate {
	return []keeper.DeploymentTemplate{
		{
			Platform:    platform.Docker,
			Name:        "Mongo (Multi Host)",
			Description: "Three-member replica set, one per VM. The last command initiates the set once every member is running - name the nodes mongo1..3 or edit the member list to match.",
			Commands: []keeper.DeploymentCommand{
				{
					Command:  deployMultiHost,
					Defaults: keeper.DeploymentDefaults{Name: "mongo1", KeeperPort: 27017, DbPort: 27017},
				},
				{
					Command:  deployMultiHost,
					Defaults: keeper.DeploymentDefaults{Name: "mongo2", KeeperPort: 27017, DbPort: 27017},
				},
				{
					Command:     deployMultiHost,
					PostScripts: []string{deployInitiate},
					Defaults:    keeper.DeploymentDefaults{Name: "mongo3", KeeperPort: 27017, DbPort: 27017},
				},
			},
		},
		{
			Platform:    platform.Docker,
			Name:        "Mongo (Single Host)",
			Description: "Three-member replica set on one VM, each member on its own port. The last command initiates the set once every member is running. Fill in the host and deploy.",
			Commands: []keeper.DeploymentCommand{
				{
					Command:  deploySingleHost,
					Defaults: keeper.DeploymentDefaults{Name: "mongo1", KeeperPort: 27017, DbPort: 27017},
				},
				{
					Command:  deploySingleHost,
					Defaults: keeper.DeploymentDefaults{Name: "mongo2", KeeperPort: 27018, DbPort: 27018},
				},
				{
					Command:     deploySingleHost,
					PostScripts: []string{deploySingleHostInitiate},
					Defaults:    keeper.DeploymentDefaults{Name: "mongo3", KeeperPort: 27019, DbPort: 27019},
				},
			},
		},
	}
}
