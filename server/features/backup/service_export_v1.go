package backup

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
