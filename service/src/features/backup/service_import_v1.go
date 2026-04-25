package backup

import (
	"encoding/json"
	"errors"
	"fmt"
	"ivory/src/clients/database"
	"ivory/src/clients/keeper"
	"ivory/src/features"
	"ivory/src/features/cluster"
	"ivory/src/features/node"
	"ivory/src/features/permission"
	"ivory/src/features/query"
)

// ImportV1 restores data from a V1 backup file.
// It maps the fixed BackupV1 schema to the current internal models.
// If internal structures evolve, this method must provide sensible defaults
// for new fields that didn't exist in V1 (e.g., setting Cluster.Type to POSTGRES_PATRONI).
func (s *Service) ImportV1(data []byte) error {
	var bkp BackupV1
	if err := json.Unmarshal(data, &bkp); err != nil {
		return err
	}

	var err error
	// Save clusters
	for i, bc := range bkp.Clusters {
		clusterModel := backupToClusterV1(bc)
		_, errMut := s.clusterService.Update(clusterModel)
		if errMut != nil {
			err = errors.Join(err, fmt.Errorf("%s[%d]: %w", "cluster", i, errMut))
		}
	}
	// Save queries
	for i, bq := range bkp.Queries {
		queryModel, errMap := backupToQueryV1(bq)
		if errMap != nil {
			continue
		}
		_, _, errMut := s.queryService.Create(query.Manual, queryModel)
		if errMut != nil {
			err = errors.Join(err, fmt.Errorf("%s[%d]: %w", "query", i, errMut))
		}
	}
	// Save permissions
	for i, bp := range bkp.Permissions {
		permModel := backupToUserPermissionsV1(bp)
		errMut := s.permissionService.UpdateUserPermissions(permModel.Username, permModel.Permissions)
		if errMut != nil {
			err = errors.Join(err, fmt.Errorf("%s[%d]: %w", "permission", i, errMut))
		}
	}

	return err
}

func backupToClusterV1(bc backupClusterV1) cluster.Cluster {
	nodes := make([]node.Connection, len(bc.Sidecars))
	for i, k := range bc.Sidecars {
		nodes[i] = node.Connection{
			Host:       k.Host,
			KeeperPort: k.Port,
		}
	}
	return cluster.Cluster{
		Name: bc.Name,
		ClusterOptions: cluster.ClusterOptions{
			DbType:     database.POSTGRES,
			KeeperType: keeper.PATRONI,
			Tags:       bc.Tags,
		},
		Nodes: nodes,
	}
}

func backupToQueryV1(bq backupQueryV1) (database.Query, error) {
	varieties := make([]database.QueryVariety, 0, len(bq.Varieties))
	for _, v := range bq.Varieties {
		variety, err := syncQueryVarietyV1(v)
		if err == nil {
			varieties = append(varieties, variety)
		}
	}

	queryType, err := syncQueryTypeV1(bq.Type)
	if err != nil {
		return database.Query{}, err
	}

	return database.Query{
		Name:        bq.Name,
		Type:        &queryType,
		Description: bq.Description,
		Query:       bq.Default,
		Varieties:   varieties,
		Params:      bq.Params,
	}, nil
}

func syncQueryTypeV1(bqt backupQueryTypeV1) (database.QueryType, error) {
	switch bqt {
	case BLOAT_V1:
		return database.BLOAT, nil
	case ACTIVITY_V1:
		return database.ACTIVITY, nil
	case REPLICATION_V1:
		return database.REPLICATION, nil
	case STATISTIC_V1:
		return database.STATISTIC, nil
	case OTHER_V1:
		return database.OTHER, nil
	default:
		return 0, ErrInvalidQueryType
	}
}

func syncQueryVarietyV1(bqv backupQueryVarietyV1) (database.QueryVariety, error) {
	switch bqv {
	case DatabaseSensitiveV1:
		return database.DatabaseSensitive, nil
	case MasterOnlyV1:
		return database.MasterOnly, nil
	case ReplicaRecommendedV1:
		return database.ReplicaRecommended, nil
	default:
		return 0, ErrInvalidQueryVariety
	}
}

func backupToUserPermissionsV1(bp backupPermissionsV1) permission.UserPermissions {
	perms := make(permission.PermissionMap)
	for k, v := range bp.Permissions {
		perm, err := syncPermissionV1(k)
		if err != nil {
			continue
		}
		status, err := syncPermissionStatusV1(v)
		if err != nil {
			continue
		}
		perms[perm] = status
	}
	return permission.UserPermissions{
		Username:    bp.Username,
		Permissions: perms,
	}
}

func syncPermissionV1(p string) (features.Feature, error) {
	for _, validFeature := range features.All {
		if string(validFeature) == p {
			return validFeature, nil
		}
	}
	return "", ErrInvalidFeature
}

func syncPermissionStatusV1(bpt backupPermissionTypeV1) (permission.Status, error) {
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
