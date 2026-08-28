package backup

import "ivory/features/deployment"

// exportV2 builds the V2 backup representation from current internal models.
// Like exportV1 it must stop changing once released; a new format means an
// exportV3 beside it rather than an edit here.
//
// It builds every collection from the current models directly: V2 borrows no
// shape from V1, so there is nothing to project through - and V1 could not
// carry a plugin, a node name or a node's own ports in the first place.
func (s *Service) exportV2() (*BackupV2, error) {
	clusters, errCluster := s.clusterService.List()
	if errCluster != nil {
		return nil, errCluster
	}
	queries, errQuery := s.queryService.GetList(nil, nil)
	if errQuery != nil {
		return nil, errQuery
	}
	permissions, errPermission := s.permissionService.GetAllUserPermissions()
	if errPermission != nil {
		return nil, errPermission
	}
	templates, errTemplate := s.deploymentService.List(deployment.ListRequest{})
	if errTemplate != nil {
		return nil, errTemplate
	}

	backupClusters := make([]backupClusterV2, 0)
	for _, c := range clusters {
		backupClusters = append(backupClusters, clusterToBackupV2(c))
	}

	backupQueries := make([]backupQueryV2, 0)
	for _, q := range queries {
		bq, err := queryToBackupV2(q)
		if err != nil {
			return nil, err
		}
		if bq != nil {
			backupQueries = append(backupQueries, *bq)
		}
	}

	backupPermissions := make([]backupPermissionsV2, 0)
	for _, p := range permissions {
		bp, err := userPermissionsToBackupV2(p)
		if err != nil {
			return nil, err
		}
		backupPermissions = append(backupPermissions, *bp)
	}

	backupDeployments := make([]backupDeploymentV2, 0)
	for _, t := range templates {
		// NOTE: the shipped templates are computed from the keeper plugins on
		// every request, so exporting one would import it back as an editable
		// copy sitting next to the original
		if t.Creation != deployment.Manual {
			continue
		}
		backupDeployments = append(backupDeployments, deploymentToBackupV2(t))
	}

	return &BackupV2{
		Clusters:    backupClusters,
		Queries:     backupQueries,
		Permissions: backupPermissions,
		Deployments: backupDeployments,
	}, nil
}
