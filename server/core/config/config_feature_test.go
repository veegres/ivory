package config

import (
	"slices"
	"strings"
	"testing"
)

func TestFeatureCurrent(t *testing.T) {
	tests := []struct {
		name     string
		stored   Feature
		expected Feature
	}{
		{name: "a renamed feature resolves to its current name", stored: "view.node.platform", expected: ViewNodeSystem},
		{name: "its manage counterpart resolves too", stored: "manage.node.platform", expected: ManageNodeSystem},
		{name: "a current feature is left untouched", stored: ViewClusterList, expected: ViewClusterList},
		{name: "an unknown feature is left untouched", stored: "view.node.nonsense", expected: "view.node.nonsense"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.stored.Current(); got != tt.expected {
				t.Errorf("Current() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// TestRenamedFeaturesResolveToKnownFeatures keeps the rename table honest: a
// target dropped from All would silently turn every stored grant for it into an
// invalid feature, which is exactly what the table exists to prevent.
func TestRenamedFeaturesResolveToKnownFeatures(t *testing.T) {
	for stored, current := range renamedFeatures {
		found := false
		for _, feature := range All {
			if feature == current {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("feature %q renames to %q, which is not in All", stored, current)
		}
		if _, ok := renamedFeatures[current]; ok {
			t.Errorf("feature %q renames to %q, which is itself renamed", stored, current)
		}
	}
}

// TestV1PermissionsResolve pins every permission key v1.4.2 wrote to disk -
// into the same Permission bucket an upgrade reads, and into every v1 backup
// file. A key that no longer resolves is not a missing entry in a table: it is
// a granted admin losing that permission on the first startup after the
// upgrade, since normalizeDatabase rewrites what it cannot recognise.
func TestV1PermissionsResolve(t *testing.T) {
	stored := []Feature{
		"view.cluster.list", "view.cluster.item", "view.cluster.overview",
		"manage.cluster.create", "manage.cluster.update", "manage.cluster.delete",
		"view.tag.list",
		"view.instance.overview", "view.instance.config", "manage.instance.config.update",
		"manage.instance.switchover", "manage.instance.reinitialize", "manage.instance.restart",
		"manage.instance.reload", "manage.instance.failover", "manage.instance.activation",
		"view.query.list", "manage.query.create", "manage.query.update", "manage.query.delete",
		"view.query.execute.info", "view.query.execute.chart", "manage.query.execute.template",
		"manage.query.execute.console", "manage.query.execute.cancel", "manage.query.execute.terminate",
		"view.query.log.list", "manage.query.log.delete",
		"view.password.list", "manage.password.create", "manage.password.update", "manage.password.delete",
		"view.cert.list", "manage.cert.create", "manage.cert.delete",
		"view.permission.list", "manage.permission.update", "manage.permission.delete",
		"view.bloat.list", "view.bloat.item", "view.bloat.logs", "manage.bloat.job",
		"view.management.secret", "manage.management.secret", "manage.management.erase",
		"manage.management.free", "manage.management.backup",
	}

	for _, feature := range stored {
		t.Run(string(feature), func(t *testing.T) {
			current := feature.Current()
			for _, known := range All {
				if known == current {
					return
				}
			}
			t.Errorf("v1 permission %q resolves to %q, which is not a feature any more", feature, current)
		})
	}
}

func TestWithheld(t *testing.T) {
	t.Run("nothing is withheld when authentication is on", func(t *testing.T) {
		if withheld := Withheld(true); len(withheld) != 0 {
			t.Errorf("expected nothing withheld with auth enabled, got %v", withheld)
		}
	})

	t.Run("only user and permission features are withheld without authentication", func(t *testing.T) {
		withheld := Withheld(false)
		if len(withheld) == 0 {
			t.Fatal("expected the session-only features to be withheld with auth disabled")
		}
		for _, feature := range withheld {
			if !slices.Contains(All, feature) {
				t.Errorf("withheld feature %q is not a feature at all", feature)
			}
			name := string(feature)
			if !strings.Contains(name, ".user.") && !strings.Contains(name, ".permission.") {
				t.Errorf("unexpected withheld feature %q", feature)
			}
		}
	})
}
