package backup

import (
	"ivory/features/deployment"
	"ivory/plugins/keeper"
	"ivory/plugins/platform"
	"testing"
)

func testTemplateRequest(name string) deployment.TemplateRequest {
	return deployment.TemplateRequest{
		Name:        name,
		Description: "three etcd members",
		Keeper:      keeper.NATIVE_ETCD,
		Platform:    platform.Linux,
		Commands: []deployment.TemplateCommand{
			{Command: "docker run -d --name {{name}} etcd"},
			{Command: "docker run -d --name {{name}} etcd", PostScript: "etcdctl auth enable"},
		},
	}
}

func TestExportV2(t *testing.T) {
	s := createTestBackupService(t)
	if _, err := s.deploymentService.Create(testTemplateRequest("mine")); err != nil {
		t.Fatalf("failed to seed template: %v", err)
	}

	backupModel, err := s.exportV2()
	if err != nil {
		t.Fatalf("exportV2() error = %v", err)
	}
	if len(backupModel.Deployments) != 1 {
		t.Fatalf("expected exactly the one manual template, got %+v", backupModel.Deployments)
	}

	got := backupModel.Deployments[0]
	if got.Name != "mine" || got.Description != "three etcd members" {
		t.Errorf("got %+v, want the stored name and description", got)
	}
	if got.Keeper != string(keeper.NATIVE_ETCD) || got.Platform != string(platform.Linux) {
		t.Errorf("got keeper %q platform %q, want the stored pair", got.Keeper, got.Platform)
	}
	if len(got.Commands) != 2 || got.Commands[1].PostScript != "etcdctl auth enable" {
		t.Errorf("got commands %+v, want both, post script included", got.Commands)
	}
}

// NOTE: the shipped templates are computed from the keeper plugins on every
// request, so exporting one would import it back as an editable copy sitting
// next to the original
func TestExportV2SkipsSystemTemplates(t *testing.T) {
	s := createTestBackupService(t)
	list, err := s.deploymentService.List(deployment.ListRequest{})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(list) == 0 {
		t.Fatal("expected the etcd defaults to be available")
	}

	backupModel, errExport := s.exportV2()
	if errExport != nil {
		t.Fatalf("exportV2() error = %v", errExport)
	}
	if len(backupModel.Deployments) != 0 {
		t.Fatalf("expected no template exported, got %+v", backupModel.Deployments)
	}
}
