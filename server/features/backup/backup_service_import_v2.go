package backup

import (
	"encoding/json"
	"errors"
	"fmt"
	"ivory/features/query"
)

// importV2 restores data from a V2 backup file. Every collection is the V2
// shape, mapped onto the current models by this version's own mappers - V1's
// are left to the files that were written in V1.
func (s *Service) importV2(data []byte) error {
	var bkp BackupV2
	if err := json.Unmarshal(data, &bkp); err != nil {
		return err
	}

	var err error
	// NOTE: the users come first because they are who everything else is about,
	// and they arrive without a password - a restored user who signs in with one
	// waits for a registration from whoever restored the file
	for i, bu := range bkp.Users {
		if _, errMut := s.userService.CreateOutright(bu.Username, "", bu.toAuthTypes(), bu.Superuser); errMut != nil {
			err = errors.Join(err, fmt.Errorf("%s[%d]: %w", "user", i, errMut))
		}
	}
	for i, bc := range bkp.Clusters {
		if _, errMut := s.clusterService.Update(bc.toCluster()); errMut != nil {
			err = errors.Join(err, fmt.Errorf("%s[%d]: %w", "cluster", i, errMut))
		}
	}
	for i, bq := range bkp.Queries {
		queryModel, errMap := bq.toQuery()
		if errMap != nil {
			err = errors.Join(err, fmt.Errorf("%s[%d]: %w", "query", i, errMap))
			continue
		}
		if _, _, errMut := s.queryService.Create(query.Manual, queryModel); errMut != nil {
			err = errors.Join(err, fmt.Errorf("%s[%d]: %w", "query", i, errMut))
		}
	}
	for i, bp := range bkp.Permissions {
		permModel := bp.toUserPermissions()
		if errMut := s.permissionService.UpdateUserPermissions(permModel.Username, permModel.Permissions); errMut != nil {
			err = errors.Join(err, fmt.Errorf("%s[%d]: %w", "permission", i, errMut))
		}
	}
	for i, bd := range bkp.Deployments {
		if _, errMut := s.deploymentService.Create(bd.toTemplateRequest()); errMut != nil {
			err = errors.Join(err, fmt.Errorf("%s[%d]: %w", "deployment", i, errMut))
		}
	}

	return err
}
