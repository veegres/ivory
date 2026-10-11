package config

import (
	"slices"
	"strings"
	"testing"
)

// TestV1PermissionsStillRecognised splits every permission key v1.4.2 wrote to
// disk into the ones that are still features and the ones that are not. A key
// in the first group must keep resolving, or an upgrade silently drops a grant
// that never changed name. A key in the second group is meant to be dropped:
// the feature was renamed, nothing forwards the old name any more, and whoever
// upgrades grants it again under the name it goes by now.
func TestV1PermissionsStillRecognised(t *testing.T) {
	unchanged := []Feature{
		"view.cluster.list", "view.cluster.item", "view.cluster.overview",
		"manage.cluster.create", "manage.cluster.update", "manage.cluster.delete",
		"view.tag.list",
		"view.query.log.list", "manage.query.log.delete",
		"view.cert.list", "manage.cert.create", "manage.cert.delete",
		"view.permission.list", "manage.permission.update", "manage.permission.delete",
		"view.management.secret", "manage.management.secret", "manage.management.erase",
		"manage.management.free", "manage.management.backup",
	}
	renamedSince := []Feature{
		"view.node.platform", "manage.node.platform",
		"view.instance.overview", "view.instance.config", "manage.instance.config.update",
		"manage.instance.switchover", "manage.instance.reinitialize", "manage.instance.restart",
		"manage.instance.reload", "manage.instance.failover", "manage.instance.activation",
		"view.query.list", "manage.query.create", "manage.query.update", "manage.query.delete",
		"view.query.execute.info", "view.query.execute.chart", "manage.query.execute.template",
		"manage.query.execute.console", "manage.query.execute.cancel", "manage.query.execute.terminate",
		"view.password.list", "manage.password.create", "manage.password.update", "manage.password.delete",
		"view.bloat.list", "view.bloat.item", "view.bloat.logs", "manage.bloat.job",
	}

	for _, feature := range unchanged {
		t.Run("kept: "+string(feature), func(t *testing.T) {
			if !slices.Contains(All, feature) {
				t.Errorf("v1 permission %q is not a feature any more, so an upgrade would drop a grant that never changed name", feature)
			}
		})
	}

	for _, feature := range renamedSince {
		t.Run("dropped: "+string(feature), func(t *testing.T) {
			if slices.Contains(All, feature) {
				t.Errorf("v1 permission %q is in All, so it was never renamed and belongs in the kept group", feature)
			}
		})
	}
}

func TestFeaturesWithheldWithoutAuth(t *testing.T) {
	t.Run("nothing is withheld when authentication is on", func(t *testing.T) {
		if withheld := FeaturesWithheldWithoutAuth(true); len(withheld) != 0 {
			t.Errorf("expected nothing withheld with auth enabled, got %v", withheld)
		}
	})

	t.Run("only user and permission features are withheld without authentication", func(t *testing.T) {
		withheld := FeaturesWithheldWithoutAuth(false)
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
