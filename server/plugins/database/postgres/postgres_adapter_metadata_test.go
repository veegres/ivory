package postgres

import (
	"ivory/core/config"
	"ivory/plugins/database"
	"testing"
)

func TestSupportedFeaturesAllSupported(t *testing.T) {
	features := NewAdapter().SupportedFeatures()

	all := []env.Feature{
		env.ViewQueryDbInfo, env.ViewQueryDbChart, env.ManageQueryDbTemplate,
		env.ManageQueryDbConsole, env.ManageQueryDbCancel, env.ManageQueryDbTerminate,
	}
	for _, feature := range all {
		if !features[feature] {
			t.Errorf("feature %v must be supported for postgres", feature)
		}
	}
}

func TestSystemChartsCoverAllChartTypes(t *testing.T) {
	charts := NewAdapter().SystemCharts()

	expected := []database.SystemChartType{
		database.Databases, database.Connections, database.DatabaseSize, database.DatabaseUptime,
		database.Schemas, database.TablesSize, database.IndexesSize, database.TotalSize,
	}
	for _, chartType := range expected {
		if charts[chartType] == "" {
			t.Errorf("expected a query for chart type %v", chartType)
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
		if len(request.Params) > 0 && len(request.Varieties) == 0 {
			t.Errorf("system request %q takes params but declares no varieties", request.Name)
		}
	}
}
