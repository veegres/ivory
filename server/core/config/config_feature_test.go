package config

import "testing"

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
