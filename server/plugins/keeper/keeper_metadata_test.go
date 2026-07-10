package keeper

import (
	"reflect"
	"testing"
)

func TestDeploymentSpecUnknownVariables(t *testing.T) {
	tests := []struct {
		name     string
		spec     DeploymentSpec
		expected []string
	}{
		{
			name: "built-in variables are known",
			spec: DeploymentSpec{
				Ports: []string{"{{keeperPort}}", "{{dbPort}}"},
				Env: []EnvVar{
					{Name: "SCOPE", Value: `"{{cluster}}"`},
					{Name: "NAME", Value: `"{{host}}"`},
					{Name: "PASSWORD", Value: `"{{dbPass}}"`},
					{Name: "USER", Value: `"{{dbUser}}"`},
				},
			},
			expected: []string{},
		},
		{
			name: "declared field names are known",
			spec: DeploymentSpec{
				Fields: []FieldSpec{
					{Name: "{{peerPort}}", Type: FieldPort, Default: "2380"},
					{Name: "{{initialCluster}}", Type: FieldText, Template: "{{host}}=http://{{host}}:{{peerPort}}"},
				},
				Env: []EnvVar{
					{Name: "MEMBERS", Value: `"{{initialCluster}}"`},
					{Name: "PEER", Value: `"{{peerPort}}"`},
				},
			},
			expected: []string{},
		},
		{
			name: "misspelled variables are reported once",
			spec: DeploymentSpec{
				Ports:      []string{"{{dbPortt}}"},
				PostDeploy: []string{"login {{dbUesr}}:{{dbPass}}", "again {{dbUesr}}"},
				Env: []EnvVar{
					{Name: "PASSWORD", Value: `"{{dbPassword}}"`},
				},
			},
			expected: []string{"{{dbPortt}}", "{{dbUesr}}", "{{dbPassword}}"},
		},
		{
			name: "undeclared field template variables are reported",
			spec: DeploymentSpec{
				Fields: []FieldSpec{
					{Name: "{{members}}", Type: FieldText, Template: "{{host}}:{{gossipPort}}"},
				},
			},
			expected: []string{"{{gossipPort}}"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.spec.UnknownVariables()
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("UnknownVariables() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestInterpolate(t *testing.T) {
	values := map[string]string{
		"{{cluster}}":    "main",
		"{{dcs}}":        "etcd1:2379",
		"{{host}}":       "db1",
		"{{dbUser}}":     "postgres",
		"{{dbPass}}":     "secret",
		"{{keeperPort}}": "8008",
		"{{dbPort}}":     "5432",
	}

	t.Run("deploy keys", func(t *testing.T) {
		got := Interpolate(
			"{{cluster}} {{dcs}} {{host}} {{keeperPort}} {{dbPort}} {{dbUser}} {{dbPass}}",
			values,
		)
		want := "main etcd1:2379 db1 8008 5432 postgres secret"
		if got != want {
			t.Errorf("Interpolate() = %q, want %q", got, want)
		}
	})

	t.Run("aux ports", func(t *testing.T) {
		got := Interpolate(
			"{{host}}:{{peerPort}} {{dbPort}}",
			map[string]string{"{{host}}": "db1", "{{peerPort}}": "2380", "{{dbPort}}": "5432"},
		)
		want := "db1:2380 5432"
		if got != want {
			t.Errorf("Interpolate() = %q, want %q", got, want)
		}
	})

	t.Run("missing and empty values keep placeholders unresolved", func(t *testing.T) {
		got := Interpolate(
			"{{cluster}} {{dcs}} {{peerPort}}",
			map[string]string{"{{cluster}}": "main", "{{dcs}}": ""},
		)
		want := "main {{dcs}} {{peerPort}}"
		if got != want {
			t.Errorf("Interpolate() = %q, want %q", got, want)
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
			input:    "-p {{peerPort}}:{{peerPort}} -e ETCD_INITIAL_CLUSTER=\"{{initialCluster}}\"",
			expected: []string{"{{peerPort}}", "{{initialCluster}}"},
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
