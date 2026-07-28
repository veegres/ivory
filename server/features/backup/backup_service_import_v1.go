package backup

import (
	"encoding/json"
	"errors"
	"fmt"
	"ivory/features/query"
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
		clusterModel := bc.toCluster()
		_, errMut := s.clusterService.Update(clusterModel)
		if errMut != nil {
			err = errors.Join(err, fmt.Errorf("%s[%d]: %w", "cluster", i, errMut))
		}
	}
	// Save queries
	for i, bq := range bkp.Queries {
		queryModel, errMap := bq.toQuery()
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
		permModel := bp.toUserPermissions()
		errMut := s.permissionService.UpdateUserPermissions(permModel.Username, permModel.Permissions)
		if errMut != nil {
			err = errors.Join(err, fmt.Errorf("%s[%d]: %w", "permission", i, errMut))
		}
	}

	return err
}
