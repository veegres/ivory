package backup

import (
	"encoding/json"
	"errors"
	"fmt"
)

// importV2 restores data from a V2 backup file. The clusters, queries and
// permissions are the V1 shapes, so it reuses V1's restore for them and only
// adds the deployment templates V1 had no concept of.
func (s *Service) importV2(data []byte) error {
	var bkp BackupV2
	if err := json.Unmarshal(data, &bkp); err != nil {
		return err
	}

	err := s.restoreV1(BackupV1{
		Clusters:    bkp.Clusters,
		Queries:     bkp.Queries,
		Permissions: bkp.Permissions,
	})

	for i, bd := range bkp.Deployments {
		if _, errMut := s.deploymentService.Create(bd.toTemplateRequest()); errMut != nil {
			err = errors.Join(err, fmt.Errorf("%s[%d]: %w", "deployment", i, errMut))
		}
	}

	return err
}
