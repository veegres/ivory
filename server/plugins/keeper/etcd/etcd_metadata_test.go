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

// TestDefaultTemplateDefaults covers what replaced keeper.Requirements: the
// deploy screen's credential fields are filled in by the template that creates
// the deployment, because etcd is its own keeper and can only enable auth
// through a user named root.
func TestDefaultTemplateDefaults(t *testing.T) {
	for _, template := range NewAdapter().DefaultTemplates() {
		t.Run(template.Name, func(t *testing.T) {
			if template.Defaults.KeeperUser != "root" {
				t.Errorf("expected keeper user %q, got %q", "root", template.Defaults.KeeperUser)
			}
			if template.Defaults.DbUser != "root" {
				t.Errorf("expected database user %q, got %q", "root", template.Defaults.DbUser)
			}
		})
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
			// NOTE: three separate steps, not one script chained with "&&":
			// the etcd image ships no shell at all, so there is nothing to
			// interpret a chain
			if len(last.PostScripts) != 3 {
				t.Fatalf("expected three steps to enable authentication, got %d", len(last.PostScripts))
			}
			for i, fragment := range []string{"user add", "user grant-role", "auth enable"} {
				if !strings.Contains(last.PostScripts[i], fragment) {
					t.Errorf("step %d: expected %q, got %q", i, fragment, last.PostScripts[i])
				}
			}
			for i, script := range last.PostScripts {
				if strings.Contains(script, "sh -c") {
					t.Errorf("step %d assumes a shell the etcd image does not ship", i)
				}
				if !strings.Contains(script, string(keeper.VarDbPort)) {
					t.Errorf("step %d does not address the node's own client port", i)
				}
			}
			for i, command := range template.Commands[:len(template.Commands)-1] {
				if len(command.PostScripts) > 0 {
					t.Errorf("command %d must not enable auth before every member is up", i)
				}
			}

			// NOTE: the member list is literal text now, so a shipped
			// template has to carry a complete one to be deployable as-is
			for i, command := range template.Commands {
				if !strings.Contains(command.Command, "ETCD_INITIAL_CLUSTER=\"etcd1=") {
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
