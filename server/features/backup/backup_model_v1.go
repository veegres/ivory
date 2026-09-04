package backup

import (
	"fmt"
	"ivory/core/config"
	"ivory/features/cluster"
	"ivory/features/permission"
	"ivory/features/query"
	"ivory/plugins/database"
	"ivory/plugins/keeper"
	"strings"
)

// BackupV1 represents the legacy backup format used in Ivory v1.
//
// SACRED RULE: This structure and all its subtypes MUST NOT be changed.
// If the internal system models change, the mapping logic in service_import_v1.go
// must be updated to translate these fixed structures into the new internal models.
// If a new backup format is required, create a new BackupV2 instead.
type BackupV1 struct {
	Clusters    []backupClusterV1     `json:"clusters"`
	Queries     []backupQueryV1       `json:"queries"`
	Permissions []backupPermissionsV1 `json:"permissions"`
}

type backupClusterV1 struct {
	Name     string            `json:"name"`
	Tags     []string          `json:"tags"`
	Sidecars []backupSidecarV1 `json:"sidecars"`
}

type backupSidecarV1 struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type backupQueryV1 struct {
	Name        string                 `json:"name"`
	Type        backupQueryTypeV1      `json:"type"`
	Varieties   []backupQueryVarietyV1 `json:"varieties"`
	Params      []string               `json:"params"`
	Description *string                `json:"description"`
	Default     string                 `json:"default"`
	Custom      string                 `json:"custom"`
}

type backupQueryTypeV1 int8

const (
	BLOAT_V1 backupQueryTypeV1 = iota
	ACTIVITY_V1
	REPLICATION_V1
	STATISTIC_V1
	OTHER_V1
)

type backupQueryVarietyV1 int8

const (
	DatabaseSensitiveV1 backupQueryVarietyV1 = iota
	MasterOnlyV1
	ReplicaRecommendedV1
)

type backupPermissionsV1 struct {
	Username    string                            `json:"username"`
	Permissions map[string]backupPermissionTypeV1 `json:"permissions"`
}

type backupPermissionTypeV1 int

const (
	NOT_PERMITTED_V1 backupPermissionTypeV1 = iota
	PENDING_V1
	GRANTED_V1
)

// Export mappers: domain → backup V1 schema

func clusterToBackupV1(c cluster.Response) backupClusterV1 {
	sidecars := make([]backupSidecarV1, len(c.Nodes))
	for i, n := range c.Nodes {
		port := 0
		if n.KeeperPort != nil {
			port = *n.KeeperPort
		}
		sidecars[i] = backupSidecarV1{Host: n.Host, Port: port}
	}
	return backupClusterV1{Name: c.Name, Tags: c.Tags, Sidecars: sidecars}
}

func queryToBackupV1(q query.Response) (*backupQueryV1, error) {
	if q.Creation == query.System {
		return nil, nil
	}
	varieties := make([]backupQueryVarietyV1, len(q.Varieties))
	for i, v := range q.Varieties {
		variety, err := queryVarietyToBackupV1(v)
		if err != nil {
			return nil, err
		}
		varieties[i] = variety
	}
	queryType, err := queryTypeToBackupV1(q.Type)
	if err != nil {
		return nil, err
	}
	return &backupQueryV1{
		Name:        q.Name,
		Type:        queryType,
		Varieties:   varieties,
		Params:      q.Params,
		Description: q.Description,
		Default:     q.Default,
		Custom:      q.Custom,
	}, nil
}

func queryTypeToBackupV1(qt query.Type) (backupQueryTypeV1, error) {
	switch qt {
	case query.BLOAT:
		return BLOAT_V1, nil
	case query.ACTIVITY:
		return ACTIVITY_V1, nil
	case query.REPLICATION:
		return REPLICATION_V1, nil
	case query.STATISTIC:
		return STATISTIC_V1, nil
	case query.OTHER:
		return OTHER_V1, nil
	default:
		return 0, ErrInvalidQueryType
	}
}

func queryVarietyToBackupV1(qv query.VarietyType) (backupQueryVarietyV1, error) {
	switch qv {
	case query.DatabaseSensitive:
		return DatabaseSensitiveV1, nil
	case query.MasterOnly:
		return MasterOnlyV1, nil
	case query.ReplicaRecommended:
		return ReplicaRecommendedV1, nil
	default:
		return 0, ErrInvalidQueryVariety
	}
}

func userPermissionsToBackupV1(up permission.UserPermissions) (*backupPermissionsV1, error) {
	perms := make(map[string]backupPermissionTypeV1)
	for k, v := range up.Permissions {
		status, err := permissionStatusToBackupV1(v)
		if err != nil {
			return nil, err
		}
		perms[string(k)] = status
	}
	return &backupPermissionsV1{Username: up.Username, Permissions: perms}, nil
}

func permissionStatusToBackupV1(ps permission.Status) (backupPermissionTypeV1, error) {
	switch ps {
	case permission.NOT_PERMITTED:
		return NOT_PERMITTED_V1, nil
	case permission.PENDING:
		return PENDING_V1, nil
	case permission.GRANTED:
		return GRANTED_V1, nil
	default:
		return 0, ErrInvalidStatus
	}
}

// Import mappers: backup V1 schema → domain

// toCluster names every node after its host: V1 knew no node names, and a
// host is what used to identify a deployment. Repeats are suffixed because a
// V1 cluster could have several nodes on one VM, and names must be unique
// within a cluster - mapping them all to the same host made a file holding
// one such cluster restore no clusters at all.
func (bc backupClusterV1) toCluster() cluster.Request {
	nodes := make([]cluster.NodeConfig, len(bc.Sidecars))
	taken := make(map[string]bool, len(bc.Sidecars))
	for i, k := range bc.Sidecars {
		nodes[i] = cluster.NodeConfig{Name: uniqueNodeNameV1(k.Host, taken), Host: k.Host, KeeperPort: &k.Port}
	}
	return cluster.Request{
		Name: bc.Name,
		Options: cluster.Options{
			Plugins: cluster.Plugins{Keeper: keeper.PATRONI_POSTGRES, Database: database.POSTGRES},
			Tags:    bc.Tags,
		},
		Nodes: nodes,
	}
}

// uniqueNodeNameV1 returns host, or the first free host-2, host-3, ... when it
// is already spoken for, and records the result as taken.
func uniqueNodeNameV1(host string, taken map[string]bool) string {
	name := host
	for i := 2; taken[name]; i++ {
		name = fmt.Sprintf("%s-%d", host, i)
	}
	taken[name] = true
	return name
}

func (bq backupQueryV1) toQuery() (query.Request, error) {
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
		Description: bq.Description,
		Query:       bq.Default,
		Varieties:   varieties,
		Params:      bq.Params,
	}, nil
}

func (bqt backupQueryTypeV1) toQueryType() (query.Type, error) {
	switch bqt {
	case BLOAT_V1:
		return query.BLOAT, nil
	case ACTIVITY_V1:
		return query.ACTIVITY, nil
	case REPLICATION_V1:
		return query.REPLICATION, nil
	case STATISTIC_V1:
		return query.STATISTIC, nil
	case OTHER_V1:
		return query.OTHER, nil
	default:
		return 0, ErrInvalidQueryType
	}
}

func (bqv backupQueryVarietyV1) toQueryVariety() (query.VarietyType, error) {
	switch bqv {
	case DatabaseSensitiveV1:
		return query.DatabaseSensitive, nil
	case MasterOnlyV1:
		return query.MasterOnly, nil
	case ReplicaRecommendedV1:
		return query.ReplicaRecommended, nil
	default:
		return 0, ErrInvalidQueryVariety
	}
}

func (bp backupPermissionsV1) toUserPermissions() permission.UserPermissions {
	perms := make(permission.PermissionMap)
	for k, v := range bp.Permissions {
		perm, err := syncPermissionV1(k)
		if err != nil {
			continue
		}
		status, err := v.toPermissionStatus()
		if err != nil {
			continue
		}
		perms[perm] = status
	}
	return permission.UserPermissions{Username: syncUsernameV1(bp.Username), Permissions: perms}
}

// syncUsernameV1 drops the authority prefix a permission key used to be stored
// under. A record is keyed by the username alone now - one person is one
// identity, whichever way they sign in - and every V1 file names them
// "basic:alice". A name carrying no known prefix is left alone.
func syncUsernameV1(username string) string {
	for _, prefix := range []string{"basic:", "ldap:", "oidc:", "superuser:"} {
		if after, found := strings.CutPrefix(username, prefix); found {
			return after
		}
	}
	return username
}

// syncPermissionV1 looks up a stored permission key against the current set of
// valid features, resolving one saved under a name that has since been renamed;
// its input is a plain string (the backup's map key), not a named local type,
// so unlike its siblings it cannot become a method.
func syncPermissionV1(p string) (config.Feature, error) {
	stored := config.Feature(p).Current()
	for _, validFeature := range config.All {
		if validFeature == stored {
			return validFeature, nil
		}
	}
	return "", ErrInvalidFeature
}

func (bpt backupPermissionTypeV1) toPermissionStatus() (permission.Status, error) {
	switch bpt {
	case NOT_PERMITTED_V1:
		return permission.NOT_PERMITTED, nil
	case PENDING_V1:
		return permission.PENDING, nil
	case GRANTED_V1:
		return permission.GRANTED, nil
	default:
		return 0, ErrInvalidStatus
	}
}
