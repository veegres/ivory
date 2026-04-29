package backup

import (
	"encoding/json"
	"errors"
	"fmt"
	"ivory/src/features"
	"ivory/src/features/cluster"
	"ivory/src/features/node"
	"ivory/src/features/permission"
	"ivory/src/features/query"
	"ivory/src/plugins/database"
	"ivory/src/plugins/keeper"
)

// importV1 restores data from a V1 backup file.
// It maps the fixed BackupV1 schema to the current internal models.
// If internal structures evolve, this method must provide sensible defaults
// for new fields that didn't exist in V1 (e.g., setting Cluster.Type to POSTGRES_PATRONI).
func (s *Service) importV1(data []byte) error {
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

func backupToQueryV1(bq backupQueryV1) (query.Request, error) {
	varieties := make([]query.VarietyType, 0, len(bq.Varieties))
	for _, v := range bq.Varieties {
		variety, err := syncQueryVarietyV1(v)
		if err == nil {
			varieties = append(varieties, variety)
		}
	}

	queryType, err := syncQueryTypeV1(bq.Type)
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

func syncQueryTypeV1(bqt backupQueryTypeV1) (query.Type, error) {
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

func syncQueryVarietyV1(bqv backupQueryVarietyV1) (query.VarietyType, error) {
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
