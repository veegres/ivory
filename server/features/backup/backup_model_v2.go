package backup

import "ivory/features/deployment"

// BackupV2 is the current backup format: everything V1 carried, plus the
// deployment templates a user has written.
//
// SACRED RULE: like BackupV1, this structure and its subtypes MUST NOT change
// once released. A further schema change means a BackupV3, not an edit here.
//
// It reuses the V1 cluster, query and permission shapes deliberately: those are
// frozen by the same rule, so borrowing them cannot drift, and duplicating them
// would only make two identical schemas to keep in step.
type BackupV2 struct {
	Clusters    []backupClusterV1     `json:"clusters"`
	Queries     []backupQueryV1       `json:"queries"`
	Permissions []backupPermissionsV1 `json:"permissions"`
	Deployments []backupDeploymentV2  `json:"deployments"`
}

// backupDeploymentV2 is one deployment template. It carries no id or creation
// type: an import creates a new template, and only a user's own are ever
// exported - the shipped ones are computed from the plugins on every request
// and would be recreated as editable copies of themselves.
type backupDeploymentV2 struct {
	Name        string                      `json:"name"`
	Description string                      `json:"description"`
	Keeper      string                      `json:"keeper"`
	Platform    string                      `json:"platform"`
	Commands    []backupDeploymentCommandV2 `json:"commands"`
}

type backupDeploymentCommandV2 struct {
	Command    string `json:"command"`
	PostScript string `json:"postScript"`
}

// Export mapper: domain → backup V2 schema

func deploymentToBackupV2(t deployment.Template) backupDeploymentV2 {
	commands := make([]backupDeploymentCommandV2, len(t.Commands))
	for i, c := range t.Commands {
		commands[i] = backupDeploymentCommandV2{Command: c.Command, PostScript: c.PostScript}
	}
	return backupDeploymentV2{
		Name:        t.Name,
		Description: t.Description,
		Keeper:      string(t.Keeper),
		Platform:    string(t.Platform),
		Commands:    commands,
	}
}

// Import mapper: backup V2 schema → domain

func (b backupDeploymentV2) toTemplateRequest() deployment.TemplateRequest {
	commands := make([]deployment.TemplateCommand, len(b.Commands))
	for i, c := range b.Commands {
		commands[i] = deployment.TemplateCommand{Command: c.Command, PostScript: c.PostScript}
	}
	return deployment.TemplateRequest{
		Name:        b.Name,
		Description: b.Description,
		Keeper:      deployment.KeeperPlugin(b.Keeper),
		Platform:    deployment.PlatformPlugin(b.Platform),
		Commands:    commands,
	}
}
