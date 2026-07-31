package clickhouse

import (
	"ivory/core/config"
	"testing"
)

func TestSupportedFeaturesExclusions(t *testing.T) {
	features := NewAdapter().SupportedFeatures()

	supported := []config.Feature{config.ViewQueryDbInfo, config.ViewQueryDbChart, config.ManageQueryDbConsole, config.ManageQueryDbTemplate}
	for _, feature := range supported {
		if !features[feature] {
			t.Errorf("feature %v must be supported for clickhouse", feature)
		}
	}

	excluded := []config.Feature{config.ManageQueryDbCancel, config.ManageQueryDbTerminate}
	for _, feature := range excluded {
		if features[feature] {
			t.Errorf("feature %v must not be supported for clickhouse", feature)
		}
	}
}

func TestSystemChartsCatalog(t *testing.T) {
	charts := NewAdapter().SystemCharts()
	if len(charts) == 0 {
		t.Fatal("expected at least one system chart")
	}
	for chartType, query := range charts {
		if query == "" {
			t.Errorf("chart %v has empty query", chartType)
		}
	}
}

func TestSystemRequestsCatalog(t *testing.T) {
	requests := NewAdapter().SystemRequests()
	if len(requests) == 0 {
		t.Fatal("expected at least one system request")
	}
	for _, request := range requests {
		if request.Name == "" {
			t.Errorf("system request has empty name: %+v", request)
		}
		if request.Query == "" {
			t.Errorf("system request %q has empty query", request.Name)
		}
	}
}
