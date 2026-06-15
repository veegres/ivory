package node

import (
	"testing"
)

func TestService_normalizeDatabaseOptions(t *testing.T) {
	s := &Service{}
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single line",
			input:    "--name test",
			expected: "--name test",
		},
		{
			name:     "multiple lines",
			input:    "--name test\n--restart always",
			expected: "--name test --restart always",
		},
		{
			name:     "multiple lines with spaces and tabs",
			input:    "--name test \n\t --restart always  \r\n -p 80:80",
			expected: "--name test --restart always -p 80:80",
		},
		{
			name:     "quoted strings with spaces",
			input:    "-e SCOPE=\"my cluster\"\n-e NAME=test",
			expected: "-e SCOPE=\"my cluster\" -e NAME=test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.normalizeDatabaseOptions(tt.input)
			if got != tt.expected {
				t.Errorf("normalizeDatabaseOptions() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestService_getInterpolatedStringDeployKeys(t *testing.T) {
	s := &Service{}

	values := ImageOptions{
		Cluster:    "main",
		Dcs:        "etcd1:2379",
		Host:       "db1",
		DbUser:     "postgres",
		DbPass:     "secret",
		KeeperPort: "8008",
		DbPort:     "5432",
	}

	got, err := s.getInterpolatedString(
		"{{cluster}} {{dcs}} {{host}} {{keeperPort}} {{dbPort}} {{dbUser}} {{dbPass}}",
		values,
	)
	if err != nil {
		t.Fatalf("getInterpolatedString failed: %v", err)
	}
	want := "main etcd1:2379 db1 8008 5432 postgres secret"
	if got != want {
		t.Errorf("getInterpolatedString() = %q, want %q", got, want)
	}
}
