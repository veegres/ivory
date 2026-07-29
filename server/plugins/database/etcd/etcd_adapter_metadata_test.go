package etcd

import (
	"ivory/core/config"
	"testing"
)

func TestSupportedFeaturesExclusions(t *testing.T) {
	features := NewAdapter().SupportedFeatures()

	supported := []config.Feature{config.ManageQueryDbConsole, config.ManageQueryDbTemplate}
	for _, feature := range supported {
		if !features[feature] {
			t.Errorf("feature %v must be supported for etcd", feature)
		}
	}

	excluded := []config.Feature{config.ViewQueryDbInfo, config.ViewQueryDbChart, config.ManageQueryDbCancel, config.ManageQueryDbTerminate}
	for _, feature := range excluded {
		if features[feature] {
			t.Errorf("feature %v must not be supported for etcd", feature)
		}
	}
}

func TestSystemChartsEmpty(t *testing.T) {
	if len(NewAdapter().SystemCharts()) != 0 {
		t.Fatal("expected no system charts for etcd")
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
