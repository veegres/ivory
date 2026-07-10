package node

import (
	"ivory/plugins/keeper"
	"reflect"
	"testing"
)

func TestMapKeeperDeploymentFields(t *testing.T) {
	tests := []struct {
		name     string
		spec     keeper.DeploymentSpec
		expected DeployFieldsResponse
	}{
		{
			name: "separate keeper port with a manual field",
			spec: keeper.DeploymentSpec{
				Defaults: map[keeper.Var]string{
					keeper.VarKeeperPort: "8008",
					keeper.VarDbPort:     "5432",
					keeper.VarDbUser:     "postgres",
				},
				Fields: []keeper.FieldSpec{{Name: "{{dcs}}", Label: "DCS", Example: "etcd1:2379", Type: keeper.FieldText}},
			},
			expected: DeployFieldsResponse{
				Defaults: map[string]string{
					keeper.VarKeeperPort: "8008",
					keeper.VarDbPort:     "5432",
					keeper.VarDbUser:     "postgres",
				},
				Fields: []DeployFieldResponse{{Name: "{{dcs}}", Label: "DCS", Example: "etcd1:2379", Type: "text"}},
			},
		},
		{
			name: "no keeper port and no extra fields",
			spec: keeper.DeploymentSpec{
				Defaults: map[keeper.Var]string{keeper.VarDbPort: "5432", keeper.VarDbUser: ""},
			},
			expected: DeployFieldsResponse{
				Defaults: map[string]string{keeper.VarDbPort: "5432", keeper.VarDbUser: ""},
				Fields:   []DeployFieldResponse{},
			},
		},
		{
			name: "port field with a derived value",
			spec: keeper.DeploymentSpec{
				Defaults: map[keeper.Var]string{keeper.VarDbPort: "2379"},
				Fields: []keeper.FieldSpec{
					{Name: "{{peerPort}}", Label: "Peer Port", Type: keeper.FieldPort, Default: "2380"},
					{Name: "{{initialCluster}}", Label: "Initial Cluster", Type: keeper.FieldText, Template: "{{host}}=http://{{host}}:{{peerPort}}", Separator: ","},
				},
			},
			expected: DeployFieldsResponse{
				Defaults: map[string]string{keeper.VarDbPort: "2379"},
				Fields: []DeployFieldResponse{
					{Name: "{{peerPort}}", Label: "Peer Port", Type: "port", Default: "2380"},
					{Name: "{{initialCluster}}", Label: "Initial Cluster", Type: "text", Derived: true},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapKeeperDeploymentFields(tt.spec)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("mapKeeperDeploymentFields() = %+v, want %+v", got, tt.expected)
			}
		})
	}
}
