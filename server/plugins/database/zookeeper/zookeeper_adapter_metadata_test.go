package zookeeper

import (
	"ivory/core/config"
	"testing"
)

func TestSupportedFeaturesExclusions(t *testing.T) {
	features := NewAdapter().SupportedFeatures()

	supported := []config.Feature{config.ManageQueryDbConsole, config.ManageQueryDbTemplate}
	for _, feature := range supported {
		if !features[feature] {
			t.Errorf("feature %v must be supported for zookeeper", feature)
		}
	}

	excluded := []config.Feature{config.ViewQueryDbInfo, config.ViewQueryDbChart, config.ManageQueryDbCancel, config.ManageQueryDbTerminate}
	for _, feature := range excluded {
		if features[feature] {
			t.Errorf("feature %v must not be supported for zookeeper", feature)
		}
	}
}

func TestSystemChartsEmpty(t *testing.T) {
	if len(NewAdapter().SystemCharts()) != 0 {
		t.Fatal("expected no system charts for zookeeper")
	}
}

// TestSystemRequestsAreParseable holds the shipped queries to the same verb set
// the console itself accepts, so a request naming a znode command this adapter
// never implemented cannot reach a user as a ready-made query.
func TestSystemRequestsAreParseable(t *testing.T) {
	for _, request := range NewAdapter().SystemRequests() {
		if _, err := parseCommand(request.Query); err != nil {
			t.Errorf("system request %q has an unparseable query %q: %v", request.Name, request.Query, err)
		}
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
