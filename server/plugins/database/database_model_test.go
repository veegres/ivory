package database

import "testing"

func TestPluginTypeValid(t *testing.T) {
	for _, plugin := range PluginTypes {
		if !plugin.Valid() {
			t.Errorf("%s should be valid", plugin)
		}
	}
	for _, plugin := range []PluginType{"", "oracle"} {
		if plugin.Valid() {
			t.Errorf("%q should not be valid", plugin)
		}
	}
}

func TestSystemRequestCategoryValid(t *testing.T) {
	for _, category := range []SystemRequestCategory{BLOAT, ACTIVITY, REPLICATION, STATISTIC, OTHER} {
		if !category.Valid() {
			t.Errorf("%d should be valid", category)
		}
	}
	for _, category := range []SystemRequestCategory{-1, OTHER + 1} {
		if category.Valid() {
			t.Errorf("%d should not be valid", category)
		}
	}
}

func TestSystemRequestVarietyValid(t *testing.T) {
	for _, variety := range []SystemRequestVariety{DatabaseSensitive, MasterOnly, ReplicaRecommended} {
		if !variety.Valid() {
			t.Errorf("%d should be valid", variety)
		}
	}
	for _, variety := range []SystemRequestVariety{-1, ReplicaRecommended + 1} {
		if variety.Valid() {
			t.Errorf("%d should not be valid", variety)
		}
	}
}

func TestSystemChartTypeValid(t *testing.T) {
	for _, chart := range SystemChartTypes {
		if !chart.Valid() {
			t.Errorf("%s should be valid", chart)
		}
	}
	if SystemChartType("Uptime").Valid() {
		t.Error("an unknown chart should not be valid")
	}
}
