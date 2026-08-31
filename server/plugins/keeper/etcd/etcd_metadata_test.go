package etcd

import (
	"ivory/core/config"
	"ivory/plugins/keeper"
	"strings"
	"testing"
)

func TestSupportedFeaturesExclusions(t *testing.T) {
	features := NewPlugin().SupportedFeatures()

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
// the deployment. Both templates ship unauthenticated and so name no accounts -
// a cluster that enabled auth in a post-script could not serve as the DCS for
// the shipped patroni template, which sends no etcd credentials.
func TestDefaultTemplateDefaults(t *testing.T) {
	for _, template := range NewPlugin().DefaultTemplates() {
		t.Run(template.Name, func(t *testing.T) {
			if template.Defaults.KeeperUser != "" {
				t.Errorf("expected no keeper user, got %q", template.Defaults.KeeperUser)
			}
			if template.Defaults.DbUser != "" {
				t.Errorf("expected no database user, got %q", template.Defaults.DbUser)
			}
		})
	}
}

// TestDefaultTemplates covers etcd's bootstrap: it has no bootstrap-time
// credentials, so nothing about a shipped deployment may depend on an account
// existing. Neither template enables auth, and neither one names a credential
// variable it would then have to be given.
func TestDefaultTemplates(t *testing.T) {
	templates := NewPlugin().DefaultTemplates()

	if len(templates) != 2 {
		t.Fatalf("expected a multi-host and a single-host template, got %d", len(templates))
	}

	for _, template := range templates {
		t.Run(template.Name, func(t *testing.T) {
			if len(template.Commands) != 3 {
				t.Fatalf("expected a three-member cluster, got %d commands", len(template.Commands))
			}

			for i, command := range template.Commands {
				if len(command.PostScripts) > 0 {
					t.Errorf("command %d initializes the cluster, which leaves patroni unable to use it", i)
				}
				for _, v := range []keeper.Var{keeper.VarKeeperUser, keeper.VarKeeperPass, keeper.VarDbUser, keeper.VarDbPass} {
					if strings.Contains(command.Command, string(v)) {
						t.Errorf("command %d references %s, which an unauthenticated deployment is never given", i, v)
					}
				}
				// NOTE: the member list is literal text now, so a shipped
				// template has to carry a complete one to be deployable as-is
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
// share a peer listener. Neither the peer nor the client port may be etcd's own
// default either: etcd replaces an advertise url identical to its built-in
// default with the detected default host, so under {{host}}=localhost the
// member using it dies on the peer url and, on the client url, comes up
// advertising an address the cluster was never configured with.
func TestDefaultTemplatesSingleHostPeerPortsDiffer(t *testing.T) {
	for _, template := range NewPlugin().DefaultTemplates() {
		if !strings.Contains(template.Name, "Single Host") {
			continue
		}
		seen := make(map[string]bool)
		for i, command := range template.Commands {
			if strings.Contains(command.Command, ":2380") {
				t.Errorf("command %d peers on etcd's default port, which localhost rewrites", i)
			}
			if command.Defaults.DbPort == 2379 {
				t.Errorf("command %d serves clients on etcd's default port, which localhost rewrites", i)
			}
			for _, port := range []string{"2480", "2482", "2484"} {
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
