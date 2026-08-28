package backup

import (
	"encoding/json"
	"ivory/features/deployment"
	"ivory/plugins/keeper"
	"ivory/plugins/platform"
	"testing"
)

func testBackupDeployment(name string) backupDeploymentV2 {
	return backupDeploymentV2{
		Name:        name,
		Description: "three etcd members",
		Keeper:      string(keeper.NATIVE_ETCD),
		Platform:    string(platform.Docker),
		Commands: []backupDeploymentCommandV2{
			{Command: "docker run -d --name {{name}} etcd"},
			{Command: "docker run -d --name {{name}} etcd", PostScript: "etcdctl auth enable"},
		},
	}
}

func TestImportV2(t *testing.T) {
	s := createTestBackupService(t)

	data, errMarshal := json.Marshal(BackupV2{
		Clusters:    []backupClusterV1{{Name: "cluster1", Sidecars: []backupSidecarV1{{Host: "h1", Port: 8008}}}},
		Deployments: []backupDeploymentV2{testBackupDeployment("mine")},
	})
	if errMarshal != nil {
		t.Fatalf("failed to marshal backup: %v", errMarshal)
	}

	if err := s.importV2(data); err != nil {
		t.Fatalf("importV2() error = %v", err)
	}

	clusters, _ := s.clusterService.List()
	if len(clusters) != 1 {
		t.Fatalf("expected the cluster to be restored too, got %v", clusters)
	}

	templates, err := s.deploymentService.List(deployment.ListRequest{})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	restored := 0
	for _, tpl := range templates {
		if tpl.Creation == deployment.Manual {
			restored++
			if tpl.Name != "mine" || len(tpl.Commands) != 2 {
				t.Errorf("got %+v, want the exported template", tpl)
			}
		}
	}
	if restored != 1 {
		t.Fatalf("expected exactly one restored template, got %d", restored)
	}
}

// NOTE: a template restores as the user's own, never as a shipped one - the
// shipped ones are computed and cannot be written at all
func TestImportV2AggregatesTemplateErrorsButContinues(t *testing.T) {
	s := createTestBackupService(t)

	data, errMarshal := json.Marshal(BackupV2{
		Deployments: []backupDeploymentV2{
			{Name: "", Commands: []backupDeploymentCommandV2{{Command: "docker run"}}},
			testBackupDeployment("good"),
		},
	})
	if errMarshal != nil {
		t.Fatalf("failed to marshal backup: %v", errMarshal)
	}

	if err := s.importV2(data); err == nil {
		t.Fatal("expected the nameless template to be reported")
	}

	templates, _ := s.deploymentService.List(deployment.ListRequest{Keeper: nil, Platform: nil})
	found := false
	for _, tpl := range templates {
		if tpl.Name == "good" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected the valid template to be restored despite the broken one")
	}
}

func TestImportV2RejectsMalformedJSON(t *testing.T) {
	s := createTestBackupService(t)
	if err := s.importV2([]byte("not-json")); err == nil {
		t.Fatal("expected an error for malformed backup data")
	}
}
