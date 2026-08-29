package keeper

import (
	"reflect"
	"slices"
	"strings"
	"testing"
)

// TestValuesCoverEveryVar is the exhaustiveness check the type system cannot
// make on its own: adding a Var without giving it a home in Values (or
// Defaults) would leave it permanently unresolvable, and adding a field
// without a Var would leave it unreachable from any command.
func TestValuesCoverEveryVar(t *testing.T) {
	lookup := Values{}.lookup()

	t.Run("every var is looked up", func(t *testing.T) {
		for _, v := range Vars {
			if _, ok := lookup[v]; !ok {
				t.Errorf("var %s has no field in Values", v)
			}
		}
	})

	t.Run("every looked up key is a declared var", func(t *testing.T) {
		for v := range lookup {
			if !slices.Contains(Vars, v) {
				t.Errorf("Values maps %s, which is not in Vars", v)
			}
		}
	})

	t.Run("every struct field is reachable", func(t *testing.T) {
		if fields := reflect.TypeOf(Values{}).NumField(); fields != len(Vars) {
			t.Errorf("Values has %d fields but there are %d vars", fields, len(Vars))
		}
	})
}

func TestInterpolate(t *testing.T) {
	values := Values{
		Cluster:    "main",
		Name:       "db-1",
		Host:       "10.0.0.1",
		SshPort:    "22",
		KeeperPort: "8008",
		DbPort:     "5432",
		KeeperUser: "root",
		KeeperPass: "keeper-secret",
		DbUser:     "postgres",
		DbPass:     "secret",
	}

	t.Run("every variable resolves", func(t *testing.T) {
		template := make([]string, 0, len(Vars))
		for _, v := range Vars {
			template = append(template, string(v))
		}
		got := Interpolate(strings.Join(template, " "), values)
		if unresolved := UnresolvedPlaceholders(got); len(unresolved) > 0 {
			t.Errorf("Interpolate() left %v unresolved", unresolved)
		}
	})

	t.Run("node values", func(t *testing.T) {
		got := Interpolate("{{cluster}} {{name}} {{host}} {{keeperPort}} {{dbPort}} {{keeperUser}} {{keeperPass}} {{dbUser}} {{dbPass}}", values)
		want := "main db-1 10.0.0.1 8008 5432 root keeper-secret postgres secret"
		if got != want {
			t.Errorf("Interpolate() = %q, want %q", got, want)
		}
	})

	t.Run("missing and empty values keep placeholders unresolved", func(t *testing.T) {
		got := Interpolate("{{cluster}} {{name}} {{dbPort}}", Values{Cluster: "main"})
		want := "main {{name}} {{dbPort}}"
		if got != want {
			t.Errorf("Interpolate() = %q, want %q", got, want)
		}
	})

	// NOTE: values belong to one deployment, so a second node resolves its own
	// and never sees the first's
	t.Run("one command's values never reach another", func(t *testing.T) {
		other := Values{Cluster: "main", Name: "db-2"}
		if got := Interpolate("{{name}}", other); got != "db-2" {
			t.Errorf("Interpolate() = %q, want the command's own value db-2", got)
		}
		if got := Interpolate("{{host}}", other); got != "{{host}}" {
			t.Errorf("Interpolate() = %q, want an unresolved placeholder: no shared scope exists", got)
		}
	})
}

func TestUnresolvedPlaceholders(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "no placeholders",
			input:    "--name db1 -p 5432:5432",
			expected: nil,
		},
		{
			name:     "single placeholder",
			input:    "--name db1 -e ETCD3_HOSTS=\"{{dcs}}\"",
			expected: []string{"{{dcs}}"},
		},
		{
			name:     "several placeholders deduplicated",
			input:    "-p {{dbPort}}:{{dbPort}} -e ETCD_NAME=\"{{name}}\"",
			expected: []string{"{{dbPort}}", "{{name}}"},
		},
		{
			name:     "docker template syntax is not matched",
			input:    "--format '{{json .}}' --name {{.Name}}",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UnresolvedPlaceholders(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("UnresolvedPlaceholders() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestUnknownPlaceholders(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "every declared variable is known",
			input:    "{{cluster}} {{name}} {{host}} {{sshPort}} {{keeperPort}} {{dbPort}} {{dbUser}} {{dbPass}}",
			expected: nil,
		},
		{
			name:     "misspellings are reported once, not treated as new variables",
			input:    "login {{dbUesr}}:{{dbPass}}\nagain {{dbUesr}} -p {{dbPortt}}",
			expected: []string{"{{dbUesr}}", "{{dbPortt}}"},
		},
		{
			// NOTE: a peer port, a member list, a coordinator address and the
			// leader's host are written literally now - they are values only
			// the operator knows, not variables
			name:     "variables retired from the vocabulary are unknown",
			input:    "-p {{peerPort}} -e HOSTS={{clusterHosts}} -e DCS={{dcs}} -h {{leaderHost}}",
			expected: []string{"{{peerPort}}", "{{clusterHosts}}", "{{dcs}}", "{{leaderHost}}"},
		},
		{
			name:     "docker template syntax is not matched",
			input:    "--format '{{json .}}'",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UnknownPlaceholders(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("UnknownPlaceholders() = %v, want %v", got, tt.expected)
			}
		})
	}
}
