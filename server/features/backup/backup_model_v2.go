package backup

import (
	"ivory/core/config"
	"ivory/features/cluster"
	"ivory/features/deployment"
	"ivory/features/permission"
	"ivory/features/query"
	"ivory/plugins/database"
	"ivory/plugins/keeper"
)

// BackupV2 is the current backup format. Its shapes mirror today's root models
// - a cluster with its plugins and nodes, a query with the engine it was
// written for, a permission map, a deployment template - carrying only the
// fields a restore actually needs.
//
// SACRED RULE: like BackupV1, this structure and its subtypes MUST NOT change
// once released. A further schema change means a BackupV3, not an edit here.
//
// It borrows nothing from V1. The V1 shapes describe what Ivory looked like
// when it knew one engine: no plugin anywhere, a node that is a host and a
// keeper port, enums stored as the ordinals of a list that has since grown.
// They live on for reading old files, where importV1 maps them onto the
// current models, and a new format that reused them would inherit their
// blind spots instead.
type BackupV2 struct {
	Clusters    []backupClusterV2     `json:"clusters"`
	Queries     []backupQueryV2       `json:"queries"`
	Permissions []backupPermissionsV2 `json:"permissions"`
	Deployments []backupDeploymentV2  `json:"deployments"`
}

// backupClusterV2 is cluster.Response minus what points outside the file: the
// vault and certificate ids reference secrets and files no backup carries, so
// restoring them would leave a cluster pointing at nothing. Tls stays because
// it is the cluster's own configuration, and dropping it would quietly restore
// a TLS cluster as a plaintext one.
type backupClusterV2 struct {
	Name     string         `json:"name"`
	Keeper   string         `json:"keeper"`
	Database string         `json:"database"`
	Tls      backupTlsV2    `json:"tls"`
	Tags     []string       `json:"tags"`
	Nodes    []backupNodeV2 `json:"nodes"`
}

type backupTlsV2 struct {
	Keeper   bool `json:"keeper"`
	Database bool `json:"database"`
}

// backupNodeV2 is cluster.NodeConfig whole: a node has a name of its own now,
// and the ports it is reached on are what makes a restored cluster usable.
type backupNodeV2 struct {
	Name       string `json:"name"`
	Host       string `json:"host"`
	SshPort    *int   `json:"sshPort"`
	KeeperPort *int   `json:"keeperPort"`
	DbPort     *int   `json:"dbPort"`
}

// backupQueryV2 is query.Response minus its id, creation type and timestamp,
// which a restore assigns itself. Query is the text as it stands now
// (Response.Custom): a manual query's original is only ever the text it was
// created with, so keeping both would back up an edit's history rather than
// the query.
type backupQueryV2 struct {
	Name        string                 `json:"name"`
	Type        backupQueryTypeV2      `json:"type"`
	Plugin      string                 `json:"plugin"`
	Description *string                `json:"description"`
	Query       string                 `json:"query"`
	Varieties   []backupQueryVarietyV2 `json:"varieties"`
	Params      []string               `json:"params"`
}

// NOTE: spelled out rather than stored as the ordinals the domain enums use -
// a frozen file must keep meaning what it said, and an ordinal only means
// anything next to the list it was written against
type backupQueryTypeV2 string

const (
	queryTypeBloatV2       backupQueryTypeV2 = "bloat"
	queryTypeActivityV2    backupQueryTypeV2 = "activity"
	queryTypeReplicationV2 backupQueryTypeV2 = "replication"
	queryTypeStatisticV2   backupQueryTypeV2 = "statistic"
	queryTypeOtherV2       backupQueryTypeV2 = "other"
)

type backupQueryVarietyV2 string

const (
	queryVarietyDatabaseSensitiveV2  backupQueryVarietyV2 = "database_sensitive"
	queryVarietyMasterOnlyV2         backupQueryVarietyV2 = "master_only"
	queryVarietyReplicaRecommendedV2 backupQueryVarietyV2 = "replica_recommended"
)

type backupPermissionsV2 struct {
	Username    string                              `json:"username"`
	Permissions map[string]backupPermissionStatusV2 `json:"permissions"`
}

type backupPermissionStatusV2 string

const (
	permissionNotPermittedV2 backupPermissionStatusV2 = "not_permitted"
	permissionPendingV2      backupPermissionStatusV2 = "pending"
	permissionGrantedV2      backupPermissionStatusV2 = "granted"
)

// backupDeploymentV2 is one deployment template. It carries no id or creation
// type: an import creates a new template, and only a user's own are ever
// exported - the shipped ones are computed from the plugins on every request
// and would be recreated as editable copies of themselves.
type backupDeploymentV2 struct {
	Name        string                      `json:"name"`
	Description string                      `json:"description"`
	Keeper      string                      `json:"keeper"`
	Platform    string                      `json:"platform"`
	Defaults    backupDeploymentDefaultsV2  `json:"defaults"`
	Commands    []backupDeploymentCommandV2 `json:"commands"`
}

type backupDeploymentCommandV2 struct {
	Command     string                            `json:"command"`
	PostScripts []string                          `json:"postScripts"`
	Defaults    backupDeploymentCommandDefaultsV2 `json:"defaults"`
}

// backupDeploymentCommandDefaultsV2 is what the command fills its node card in
// with. It is part of the template rather than of a node: a file restored
// without it would put every node of a single-host template back on one port.
type backupDeploymentCommandDefaultsV2 struct {
	Name       string `json:"name"`
	SshPort    int    `json:"sshPort"`
	KeeperPort int    `json:"keeperPort"`
	DbPort     int    `json:"dbPort"`
}

// backupDeploymentDefaultsV2 is what the template fills the deploy screen's
// credential fields in with. It holds usernames only - a template never carries
// a password, so neither does a file full of templates.
type backupDeploymentDefaultsV2 struct {
	KeeperUser string `json:"keeperUser"`
	DbUser     string `json:"dbUser"`
}

// Export mappers: domain → backup V2 schema

func clusterToBackupV2(c cluster.Response) backupClusterV2 {
	nodes := make([]backupNodeV2, len(c.Nodes))
	for i, n := range c.Nodes {
		nodes[i] = backupNodeV2{
			Name:       n.Name,
			Host:       n.Host,
			SshPort:    n.SshPort,
			KeeperPort: n.KeeperPort,
			DbPort:     n.DbPort,
		}
	}
	return backupClusterV2{
		Name:     c.Name,
		Keeper:   string(c.Plugins.Keeper),
		Database: string(c.Plugins.Database),
		Tls:      backupTlsV2{Keeper: c.Tls.Keeper, Database: c.Tls.Database},
		Tags:     c.Tags,
		Nodes:    nodes,
	}
}

func queryToBackupV2(q query.Response) (*backupQueryV2, error) {
	if q.Creation == query.System {
		return nil, nil
	}
	varieties := make([]backupQueryVarietyV2, len(q.Varieties))
	for i, v := range q.Varieties {
		variety, err := queryVarietyToBackupV2(v)
		if err != nil {
			return nil, err
		}
		varieties[i] = variety
	}
	queryType, err := queryTypeToBackupV2(q.Type)
	if err != nil {
		return nil, err
	}
	return &backupQueryV2{
		Name:        q.Name,
		Type:        queryType,
		Plugin:      string(q.Plugin),
		Description: q.Description,
		Query:       q.Custom,
		Varieties:   varieties,
		Params:      q.Params,
	}, nil
}

func queryTypeToBackupV2(qt query.Type) (backupQueryTypeV2, error) {
	switch qt {
	case query.BLOAT:
		return queryTypeBloatV2, nil
	case query.ACTIVITY:
		return queryTypeActivityV2, nil
	case query.REPLICATION:
		return queryTypeReplicationV2, nil
	case query.STATISTIC:
		return queryTypeStatisticV2, nil
	case query.OTHER:
		return queryTypeOtherV2, nil
	default:
		return "", ErrInvalidQueryType
	}
}

func queryVarietyToBackupV2(qv query.VarietyType) (backupQueryVarietyV2, error) {
	switch qv {
	case query.DatabaseSensitive:
		return queryVarietyDatabaseSensitiveV2, nil
	case query.MasterOnly:
		return queryVarietyMasterOnlyV2, nil
	case query.ReplicaRecommended:
		return queryVarietyReplicaRecommendedV2, nil
	default:
		return "", ErrInvalidQueryVariety
	}
}

func userPermissionsToBackupV2(up permission.UserPermissions) (*backupPermissionsV2, error) {
	perms := make(map[string]backupPermissionStatusV2)
	for k, v := range up.Permissions {
		status, err := permissionStatusToBackupV2(v)
		if err != nil {
			return nil, err
		}
		perms[string(k)] = status
	}
	return &backupPermissionsV2{Username: up.Username, Permissions: perms}, nil
}

func permissionStatusToBackupV2(ps permission.Status) (backupPermissionStatusV2, error) {
	switch ps {
	case permission.NOT_PERMITTED:
		return permissionNotPermittedV2, nil
	case permission.PENDING:
		return permissionPendingV2, nil
	case permission.GRANTED:
		return permissionGrantedV2, nil
	default:
		return "", ErrInvalidStatus
	}
}

func deploymentToBackupV2(t deployment.Template) backupDeploymentV2 {
	commands := make([]backupDeploymentCommandV2, len(t.Commands))
	for i, c := range t.Commands {
		commands[i] = backupDeploymentCommandV2{
			Command:     c.Command,
			PostScripts: c.PostScripts,
			Defaults: backupDeploymentCommandDefaultsV2{
				Name:       c.Defaults.Name,
				SshPort:    c.Defaults.SshPort,
				KeeperPort: c.Defaults.KeeperPort,
				DbPort:     c.Defaults.DbPort,
			},
		}
	}
	return backupDeploymentV2{
		Name:        t.Name,
		Description: t.Description,
		Keeper:      string(t.Keeper),
		Platform:    string(t.Platform),
		Defaults: backupDeploymentDefaultsV2{
			KeeperUser: t.Defaults.KeeperUser,
			DbUser:     t.Defaults.DbUser,
		},
		Commands: commands,
	}
}

// Import mappers: backup V2 schema → domain

func (bc backupClusterV2) toCluster() cluster.Request {
	nodes := make([]cluster.NodeConfig, len(bc.Nodes))
	for i, n := range bc.Nodes {
		nodes[i] = cluster.NodeConfig{
			Name:       n.Name,
			Host:       n.Host,
			SshPort:    n.SshPort,
			KeeperPort: n.KeeperPort,
			DbPort:     n.DbPort,
		}
	}
	return cluster.Request{
		Name: bc.Name,
		Options: cluster.Options{
			Plugins: cluster.Plugins{Keeper: keeper.PluginType(bc.Keeper), Database: database.PluginType(bc.Database)},
			Tls:     cluster.Tls{Keeper: bc.Tls.Keeper, Database: bc.Tls.Database},
			Tags:    bc.Tags,
		},
		Nodes: nodes,
	}
}

func (bq backupQueryV2) toQuery() (query.Request, error) {
	varieties := make([]query.VarietyType, 0, len(bq.Varieties))
	for _, v := range bq.Varieties {
		variety, err := v.toQueryVariety()
		if err == nil {
			varieties = append(varieties, variety)
		}
	}
	queryType, err := bq.Type.toQueryType()
	if err != nil {
		return query.Request{}, err
	}
	return query.Request{
		Name:        bq.Name,
		Type:        &queryType,
		Plugin:      database.PluginType(bq.Plugin),
		Description: bq.Description,
		Query:       bq.Query,
		Varieties:   varieties,
		Params:      bq.Params,
	}, nil
}

func (bqt backupQueryTypeV2) toQueryType() (query.Type, error) {
	switch bqt {
	case queryTypeBloatV2:
		return query.BLOAT, nil
	case queryTypeActivityV2:
		return query.ACTIVITY, nil
	case queryTypeReplicationV2:
		return query.REPLICATION, nil
	case queryTypeStatisticV2:
		return query.STATISTIC, nil
	case queryTypeOtherV2:
		return query.OTHER, nil
	default:
		return 0, ErrInvalidQueryType
	}
}

func (bqv backupQueryVarietyV2) toQueryVariety() (query.VarietyType, error) {
	switch bqv {
	case queryVarietyDatabaseSensitiveV2:
		return query.DatabaseSensitive, nil
	case queryVarietyMasterOnlyV2:
		return query.MasterOnly, nil
	case queryVarietyReplicaRecommendedV2:
		return query.ReplicaRecommended, nil
	default:
		return 0, ErrInvalidQueryVariety
	}
}

func (bp backupPermissionsV2) toUserPermissions() permission.UserPermissions {
	perms := make(permission.PermissionMap)
	for k, v := range bp.Permissions {
		perm, err := syncPermissionV2(k)
		if err != nil {
			continue
		}
		status, err := v.toPermissionStatus()
		if err != nil {
			continue
		}
		perms[perm] = status
	}
	return permission.UserPermissions{Username: bp.Username, Permissions: perms}
}

// syncPermissionV2 looks up a stored permission key against the current set of
// valid features, resolving one saved under a name that has since been renamed;
// its input is a plain string (the backup's map key), not a named local type,
// so unlike its siblings it cannot become a method.
func syncPermissionV2(p string) (config.Feature, error) {
	stored := config.Feature(p).Current()
	for _, validFeature := range config.All {
		if validFeature == stored {
			return validFeature, nil
		}
	}
	return "", ErrInvalidFeature
}

func (bps backupPermissionStatusV2) toPermissionStatus() (permission.Status, error) {
	switch bps {
	case permissionNotPermittedV2:
		return permission.NOT_PERMITTED, nil
	case permissionPendingV2:
		return permission.PENDING, nil
	case permissionGrantedV2:
		return permission.GRANTED, nil
	default:
		return 0, ErrInvalidStatus
	}
}

func (b backupDeploymentV2) toTemplateRequest() deployment.TemplateRequest {
	commands := make([]deployment.TemplateCommand, len(b.Commands))
	for i, c := range b.Commands {
		commands[i] = deployment.TemplateCommand{
			Command:     c.Command,
			PostScripts: c.PostScripts,
			Defaults: deployment.CommandDefaults{
				Name:       c.Defaults.Name,
				SshPort:    c.Defaults.SshPort,
				KeeperPort: c.Defaults.KeeperPort,
				DbPort:     c.Defaults.DbPort,
			},
		}
	}
	return deployment.TemplateRequest{
		Name:        b.Name,
		Description: b.Description,
		Keeper:      deployment.KeeperPlugin(b.Keeper),
		Platform:    deployment.PlatformPlugin(b.Platform),
		Defaults: deployment.TemplateDefaults{
			KeeperUser: b.Defaults.KeeperUser,
			DbUser:     b.Defaults.DbUser,
		},
		Commands: commands,
	}
}
