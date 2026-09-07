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

// executableVerbs mirrors the switch in executeCollectionVerb, plus the one
// verb db itself takes. Parsing alone proves nothing about a shipped query:
// parseCommand accepts any verb and only execution rejects an unknown one, so
// a query naming a verb this adapter cannot run would otherwise ship looking
// fine and fail on every use.
var executableVerbs = map[string]bool{
	"find": true, "findOne": true, "insertOne": true, "insertMany": true,
	"updateOne": true, "updateMany": true, "deleteOne": true, "deleteMany": true,
	"countDocuments": true, "distinct": true, "aggregate": true,
}

func TestSystemRequestsAreParseable(t *testing.T) {
	// NOTE: a name is what seeding matches an already-stored system query
	// against (see query.Service.initializeSystemQueries), so two requests
	// sharing one would leave the second permanently unseeded.
	seen := map[string]bool{}
	for _, request := range NewAdapter().SystemRequests() {
		parsed, err := parseCommand(request.Query)
		if err != nil {
			t.Errorf("system request %q has an unparseable query %q: %v", request.Name, request.Query, err)
		}
		if err == nil {
			if parsed.Collection == "db" {
				if parsed.Verb != "runCommand" {
					t.Errorf("system request %q names db.%s, but db only takes runCommand", request.Name, parsed.Verb)
				}
			} else if !executableVerbs[parsed.Verb] {
				t.Errorf("system request %q names verb %q, which executeCollectionVerb cannot run", request.Name, parsed.Verb)
			}
		}
		if seen[request.Name] {
			t.Errorf("system request name %q is used twice", request.Name)
		}
		seen[request.Name] = true
	}
}

func TestSystemChartsEmpty(t *testing.T) {
	if charts := NewAdapter().SystemCharts(); len(charts) != 0 {
		t.Errorf("expected no system charts, got %+v", charts)
	}
}
