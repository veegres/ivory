package redis

import (
	"ivory/core/config"
	"testing"
)

func TestSupportedFeaturesExclusions(t *testing.T) {
	features := NewAdapter().SupportedFeatures()

	supported := []config.Feature{config.ManageQueryDbConsole, config.ManageQueryDbTemplate, config.ManageQueryDbTerminate}
	for _, feature := range supported {
		if !features[feature] {
			t.Errorf("feature %v must be supported for redis", feature)
		}
	}

	excluded := []config.Feature{config.ViewQueryDbInfo, config.ViewQueryDbChart, config.ManageQueryDbCancel}
	for _, feature := range excluded {
		if features[feature] {
			t.Errorf("feature %v must not be supported for redis", feature)
		}
	}
}

func TestSystemChartsEmpty(t *testing.T) {
	if len(NewAdapter().SystemCharts()) != 0 {
		t.Fatal("expected no system charts for redis")
	}
}

func TestSystemRequestsCatalog(t *testing.T) {
	requests := NewAdapter().SystemRequests()
	if len(requests) == 0 {
		t.Fatal("expected at least one system request")
	}
	// NOTE: a name is what seeding matches an already-stored system query
	// against (see query.Service.initializeSystemQueries), so two requests
	// sharing one would leave the second permanently unseeded.
	seen := map[string]bool{}
	for _, request := range requests {
		if request.Name == "" {
			t.Errorf("system request has empty name: %+v", request)
		}
		if seen[request.Name] {
			t.Errorf("system request name %q is used twice", request.Name)
		}
		seen[request.Name] = true
		if request.Query == "" {
			t.Errorf("system request %q has empty query", request.Name)
		}
	}
}
