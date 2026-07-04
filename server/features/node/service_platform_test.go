package node

import (
	"ivory/clients/console/ssh"
	"ivory/core/utils"
	"ivory/plugins/keeper"
	"ivory/plugins/keeper/patroni"
	pgkeeper "ivory/plugins/keeper/postgres"
	"ivory/plugins/platform"
	"ivory/plugins/platform/linux"
	"reflect"
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

// TestService_PlatformContainerDeployOptions pins the rendered deploy options
// byte-for-byte to the templates that previously lived hardcoded in the web
// app (DatabaseImageOptions), so moving them to the backend changes nothing
// for existing users.
func TestService_PlatformContainerDeployOptions(t *testing.T) {
	platformRegistry := utils.NewRegistry[platform.Plugin, platform.Adapter]()
	platformRegistry.Register(platform.Linux, linux.NewAdapter(ssh.NewClient()))
	keeperMetadataRegistry := utils.NewRegistry[keeper.Plugin, keeper.Metadata]()
	keeperMetadataRegistry.Register(keeper.PATRONI, patroni.NewAdapter(nil))
	keeperMetadataRegistry.Register(keeper.POSTGRES, pgkeeper.NewAdapter())

	s := &Service{
		platformRegistry:       platformRegistry,
		keeperMetadataRegistry: keeperMetadataRegistry,
	}

	tests := []struct {
		name     string
		plugin   KeeperPlugin
		expected PlatformDeployOptionsResponse
	}{
		{
			name:   "patroni",
			plugin: keeper.PATRONI,
			expected: PlatformDeployOptionsResponse{
				Uri:           "ghcr.io/zalando/spilo-18:4.1-p2",
				DefaultValues: map[string]string{"username": "postgres"},
				Options: "--name {{host}}\n" +
					"--hostname {{host}}\n" +
					"--restart unless-stopped\n" +
					"-p {{keeperPort}}:{{keeperPort}}\n" +
					"-p {{dbPort}}:{{dbPort}}\n" +
					"-v /data/postgres:/home/postgres/pgdata\n" +
					"-e SCOPE=\"{{cluster}}\"\n" +
					"-e PATRONI_NAME=\"{{host}}\"\n" +
					"-e ETCD3_HOSTS=\"{{dcs}}\"\n" +
					"-e PGPORT={{dbPort}}\n" +
					"-e APIPORT={{keeperPort}}\n" +
					"-e PGPASSWORD_SUPERUSER=\"{{dbPass}}\"\n" +
					"-e RESTAPI_CONNECT_ADDRESS=\"{{host}}:{{keeperPort}}\"\n" +
					`-e SPILO_CONFIGURATION='{"postgresql":{"connect_address":"{{host}}:{{dbPort}}"},"bootstrap":{"dcs":{"primary_start_timeout":999}}}'`,
				OptionsSingleHost: "--name {{host}}\n" +
					"--hostname {{host}}\n" +
					"--network host\n" +
					"-e SCOPE=\"{{cluster}}\"\n" +
					"-e PATRONI_NAME=\"{{host}}\"\n" +
					"-e ETCD3_HOSTS=\"{{dcs}}\"\n" +
					"-e PGPORT={{dbPort}}\n" +
					"-e APIPORT={{keeperPort}}\n" +
					"-e PGPASSWORD_SUPERUSER=\"{{dbPass}}\"\n" +
					"-e RESTAPI_CONNECT_ADDRESS=\"{{host}}:{{keeperPort}}\"\n" +
					`-e SPILO_CONFIGURATION='{"postgresql":{"connect_address":"{{host}}:{{dbPort}}"},"bootstrap":{"dcs":{"primary_start_timeout":999}}}'`,
			},
		},
		{
			name:   "postgres",
			plugin: keeper.POSTGRES,
			expected: PlatformDeployOptionsResponse{
				Uri:           "postgres:18",
				DefaultValues: map[string]string{"dcs": "empty"},
				Options: "--name {{host}}\n" +
					"--hostname {{host}}\n" +
					"--restart unless-stopped\n" +
					"-p {{dbPort}}:{{dbPort}}\n" +
					"-v /data/postgres:/var/lib/postgresql/data\n" +
					"-e PGPORT=\"{{dbPort}}\"\n" +
					"-e POSTGRES_USER=\"{{dbUser}}\"\n" +
					"-e POSTGRES_PASSWORD=\"{{dbPass}}\"",
				OptionsSingleHost: "--name {{host}}\n" +
					"--hostname {{host}}\n" +
					"--network host\n" +
					"-e PGPORT=\"{{dbPort}}\"\n" +
					"-e POSTGRES_USER=\"{{dbUser}}\"\n" +
					"-e POSTGRES_PASSWORD=\"{{dbPass}}\"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.PlatformContainerDeployOptions(PlatformDeployOptionsRequest{Plugin: tt.plugin})
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if got.Uri != tt.expected.Uri {
				t.Errorf("Uri = %q, want %q", got.Uri, tt.expected.Uri)
			}
			if !reflect.DeepEqual(got.DefaultValues, tt.expected.DefaultValues) {
				t.Errorf("DefaultValues = %v, want %v", got.DefaultValues, tt.expected.DefaultValues)
			}
			if got.Options != tt.expected.Options {
				t.Errorf("Options mismatch\ngot:\n%s\nwant:\n%s", got.Options, tt.expected.Options)
			}
			if got.OptionsSingleHost != tt.expected.OptionsSingleHost {
				t.Errorf("OptionsSingleHost mismatch\ngot:\n%s\nwant:\n%s", got.OptionsSingleHost, tt.expected.OptionsSingleHost)
			}
		})
	}
}

func TestService_PlatformContainerDeployOptionsUnknownPlugin(t *testing.T) {
	s := &Service{
		platformRegistry:       utils.NewRegistry[platform.Plugin, platform.Adapter](),
		keeperMetadataRegistry: utils.NewRegistry[keeper.Plugin, keeper.Metadata](),
	}

	_, err := s.PlatformContainerDeployOptions(PlatformDeployOptionsRequest{Plugin: "unknown"})
	if err == nil {
		t.Fatal("expected error for unknown plugin, got nil")
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
