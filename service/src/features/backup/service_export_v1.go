package backup

import (
	"ivory/src/features/cluster"
	"ivory/src/features/permission"
	"ivory/src/features/query"
)

// exportV1 builds the fixed V1 backup representation from current internal
// models.
//
// This function defines how the current system is projected into the frozen V1
// wire schema. Once a new backup format is introduced, exportV1 must remain
// unchanged so Ivory can continue producing or reasoning about legacy V1 data if
// needed, while new exports should move to exportV2 or later.
func (s *Service) exportV1() (*BackupV1, error) {
	clusters, errCluster := s.clusterService.List()
	if errCluster != nil {
		return nil, errCluster
	}
	queries, errQuery := s.queryService.GetList(nil)
	if errQuery != nil {
		return nil, errQuery
	}
	permissions, errPermission := s.permissionService.GetAllUserPermissions()
	if errPermission != nil {
		return nil, errPermission
	}

	backupClusters := make([]backupClusterV1, 0)
	for _, c := range clusters {
		backupClusters = append(backupClusters, clusterToBackupV1(c))
	}

	backupQueries := make([]backupQueryV1, 0)
	for _, q := range queries {
		bq, err := queryToBackupV1(q)
		if err != nil {
			return nil, err
		}
		if bq != nil {
			backupQueries = append(backupQueries, *bq)
		}
	}

	backupPermissions := make([]backupPermissionsV1, 0)
	for _, p := range permissions {
		bp, err := userPermissionsToBackupV1(p)
		if err != nil {
			return nil, err
		}
		backupPermissions = append(backupPermissions, *bp)
	}

	return &BackupV1{
		Clusters:    backupClusters,
		Queries:     backupQueries,
		Permissions: backupPermissions,
	}, nil
}

func clusterToBackupV1(c cluster.Response) backupClusterV1 {
	sidecars := make([]backupSidecarV1, len(c.Nodes))
	for i, n := range c.Nodes {
		sidecars[i] = backupSidecarV1{
			Host: n.Host,
			Port: *n.KeeperPort,
		}
	}
	return backupClusterV1{
		Name:     c.Name,
		Tags:     c.Tags,
		Sidecars: sidecars,
	}
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
	return &backupPermissionsV1{
		Username:    up.Username,
		Permissions: perms,
	}, nil
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
