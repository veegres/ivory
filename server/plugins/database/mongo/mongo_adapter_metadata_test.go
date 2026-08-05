package mongo

import (
	"ivory/core/config"
	"testing"
)

func TestSupportedFeaturesExclusions(t *testing.T) {
	features := NewAdapter().SupportedFeatures()

	supported := []config.Feature{config.ManageQueryDbTemplate, config.ManageQueryDbConsole, config.ManageQueryDbTerminate}
	for _, feature := range supported {
		if !features[feature] {
			t.Errorf("feature %v must be supported for mongo", feature)
		}
	}

	excluded := []config.Feature{config.ViewQueryDbInfo, config.ViewQueryDbChart, config.ManageQueryDbCancel}
	for _, feature := range excluded {
		if features[feature] {
			t.Errorf("feature %v must not be supported for mongo", feature)
		}
	}
}

func TestSystemRequestsAreParseable(t *testing.T) {
	for _, request := range NewAdapter().SystemRequests() {
		if _, err := parseCommand(request.Query); err != nil {
			t.Errorf("system request %q has an unparseable query %q: %v", request.Name, request.Query, err)
		}
	}
}

func TestSystemChartsEmpty(t *testing.T) {
	if charts := NewAdapter().SystemCharts(); len(charts) != 0 {
		t.Errorf("expected no system charts, got %+v", charts)
	}
}
