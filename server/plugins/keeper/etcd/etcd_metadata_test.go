package etcd

import (
	"ivory/core/config"
	"ivory/plugins/keeper"
	"strings"
	"testing"
)

func TestSupportedFeaturesExclusions(t *testing.T) {
	features := NewAdapter().SupportedFeatures()

	supported := []config.Feature{config.ViewNodeKeeperOverview, config.ManageNodeKeeperSwitchover}
	for _, feature := range supported {
		if !features[feature] {
			t.Errorf("feature %v must be supported for etcd", feature)
		}
	}

	excluded := []config.Feature{
		config.ViewNodeKeeperConfig, config.ManageNodeKeeperConfigUpdate, config.ManageNodeKeeperReinitialize,
		config.ManageNodeKeeperRestart, config.ManageNodeKeeperReload, config.ManageNodeKeeperFailover, config.ManageNodeKeeperActivation,
	}
	for _, feature := range excluded {
		if features[feature] {
			t.Errorf("feature %v must not be supported for etcd", feature)
		}
	}
}

func TestRequirements(t *testing.T) {
	req := NewAdapter().Requirements()

	if req.DbPort != 2379 {
		t.Errorf("expected client port 2379, got %d", req.DbPort)
	}
	if req.KeeperPort != 2379 {
		t.Errorf("expected the client port 2379 declared as the keeper endpoint too, got %d", req.KeeperPort)
	}
	if !req.KeeperCredentials || req.KeeperUser != "root" {
		t.Errorf("expected keeper credentials locked to root, got %v/%q", req.KeeperCredentials, req.KeeperUser)
	}
	if !req.DbCredentials {
		t.Error("expected etcd to consume database credentials")
	}
	if req.DbUser != "root" {
		t.Errorf("expected the etcd-required username root (auth enable needs it), got %q", req.DbUser)
	}
}

// TestDefaultTemplates covers etcd's bootstrap: it has no bootstrap-time
// credentials, so the root user can only be created once the whole cluster is
// running - which is why the auth script sits on the last command, not the
// first.
func TestDefaultTemplates(t *testing.T) {
	templates := NewAdapter().DefaultTemplates()

	if len(templates) != 2 {
		t.Fatalf("expected a multi-host and a single-host template, got %d", len(templates))
	}

	for _, template := range templates {
		t.Run(template.Name, func(t *testing.T) {
			if len(template.Commands) != 3 {
				t.Fatalf("expected a three-member cluster, got %d commands", len(template.Commands))
			}

			last := template.Commands[len(template.Commands)-1]
			if last.PostScript == "" {
				t.Fatal("expected the last command to enable authentication")
			}
			for _, fragment := range []string{"user add", "{{dbUser}}:{{dbPass}}", "user grant-role {{dbUser}} root", "auth enable"} {
				if !strings.Contains(last.PostScript, fragment) {
					t.Errorf("expected the post script to contain %q, got %q", fragment, last.PostScript)
				}
			}
			for i, command := range template.Commands[:len(template.Commands)-1] {
				if command.PostScript != "" {
					t.Errorf("command %d must not enable auth before every member is up", i)
				}
			}

			// NOTE: the member list is literal text now, so a shipped
			// template has to carry a complete one to be deployable as-is
			for i, command := range template.Commands {
				if !strings.Contains(command.Command, "ETCD_INITIAL_CLUSTER=\"etcd-1=") {
					t.Errorf("command %d has no initial cluster list to edit", i)
				}
				if len(keeper.UnknownPlaceholders(command.Command)) > 0 {
					t.Errorf("command %d references a variable outside the vocabulary", i)
				}
			}
		})
	}
}

// TestDefaultTemplatesSingleHostPeerPortsDiffer covers the collision the
// deleted singleHost flag used to compute away: three members on one VM cannot
// share a peer listener.
func TestDefaultTemplatesSingleHostPeerPortsDiffer(t *testing.T) {
	for _, template := range NewAdapter().DefaultTemplates() {
		if !strings.Contains(template.Name, "Single Host") {
			continue
		}
		seen := make(map[string]bool)
		for i, command := range template.Commands {
			for _, port := range []string{"2380", "2382", "2384"} {
				if !strings.Contains(command.Command, "http://0.0.0.0:"+port) {
					continue
				}
				if seen[port] {
					t.Errorf("command %d reuses peer port %s on the same VM", i, port)
				}
				seen[port] = true
			}
		}
		if len(seen) != len(template.Commands) {
			t.Errorf("expected one distinct peer port per member, got %d", len(seen))
		}
	}
}
