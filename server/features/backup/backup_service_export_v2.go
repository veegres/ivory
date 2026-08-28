package backup

import "ivory/features/deployment"

// exportV2 builds the V2 backup representation from current internal models.
// Like exportV1 it must stop changing once released; a new format means an
// exportV3 beside it rather than an edit here.
func (s *Service) exportV2() (*BackupV2, error) {
	v1, err := s.exportV1()
	if err != nil {
		return nil, err
	}
	templates, errTemplate := s.deploymentService.List(deployment.ListRequest{})
	if errTemplate != nil {
		return nil, errTemplate
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
		Clusters:    v1.Clusters,
		Queries:     v1.Queries,
		Permissions: v1.Permissions,
		Deployments: backupDeployments,
	}, nil
}
