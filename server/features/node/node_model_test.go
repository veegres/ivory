package node

import (
	"ivory/core/utils"
	"ivory/plugins/platform"
	"strings"
	"testing"

	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

func TestPlatformVaultConnectionPlatformOrDefault(t *testing.T) {
	tests := []struct {
		name       string
		connection PlatformVaultConnection
		expected   PlatformPlugin
	}{
		{
			// NOTE: clusters stored before platforms were selectable have no
			// platform at all, and must keep resolving to the only adapter
			// that existed then
			name:       "an empty platform falls back to docker",
			connection: PlatformVaultConnection{Host: "db1"},
			expected:   platform.Docker,
		},
		{
			name:       "an explicit platform is kept",
			connection: PlatformVaultConnection{Host: "db1", Platform: platform.Docker},
			expected:   platform.Docker,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.connection.PlatformOrDefault(); got != tt.expected {
				t.Errorf("PlatformOrDefault() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// TestRequestBinding covers the shape rules node's request types state in
// their binding tags - including the ones the ?request= query routes used to
// skip, since they decoded without validating.
func TestRequestBinding(t *testing.T) {
	utils.RegisterValidators()
	connection := func() PlatformVaultConnection {
		return PlatformVaultConnection{Host: "db1", Port: 22, VaultId: uuid.New()}
	}
	vaults := Vaults{SshKeyId: uuid.New()}

	tests := []struct {
		name    string
		request any
		want    string
	}{
		{name: "keeper request passes", request: &KeeperOneRequest{KeeperConnection: KeeperConnection{Host: "db1", Port: 8008}, KeeperOptions: KeeperOptions{Plugin: "patroni_postgres"}}},
		{name: "keeper request names a known plugin", request: &KeeperOneRequest{KeeperConnection: KeeperConnection{Host: "db1", Port: 8008}, KeeperOptions: KeeperOptions{Plugin: "pgpool"}}, want: "'KeeperOneRequest.KeeperOptions.Plugin' Error:Field validation for 'Plugin' failed on the 'valid' tag"},
		{name: "keeper request needs a host", request: &KeeperOneRequest{KeeperConnection: KeeperConnection{Port: 8008}, KeeperOptions: KeeperOptions{Plugin: "patroni_postgres"}}, want: "'KeeperOneRequest.KeeperConnection.Host' Error:Field validation for 'Host' failed on the 'required' tag"},
		{name: "keeper port is a port", request: &KeeperOneRequest{KeeperConnection: KeeperConnection{Host: "db1", Port: 99999}, KeeperOptions: KeeperOptions{Plugin: "patroni_postgres"}}, want: "'KeeperOneRequest.KeeperConnection.Port' Error:Field validation for 'Port' failed on the 'max' tag"},
		{name: "connection may leave the platform out", request: &PlatformActionRequest{Connection: connection(), Name: "pg1"}},
		{name: "connection platform must be known", request: &PlatformActionRequest{Connection: PlatformVaultConnection{Host: "db1", Port: 22, VaultId: uuid.New(), Platform: "k8s"}, Name: "pg1"}, want: "'PlatformActionRequest.Connection.Platform' Error:Field validation for 'Platform' failed on the 'valid' tag"},
		{name: "connection needs a vault", request: &PlatformActionRequest{Connection: PlatformVaultConnection{Host: "db1", Port: 22}, Name: "pg1"}, want: "'PlatformActionRequest.Connection.VaultId' Error:Field validation for 'VaultId' failed on the 'required' tag"},
		{name: "action names a container", request: &PlatformActionRequest{Connection: connection(), Name: " "}, want: "'PlatformActionRequest.Name' Error:Field validation for 'Name' failed on the 'notblank' tag"},
		{name: "logs need a path", request: &PlatformLogsRequest{Connection: connection()}, want: "'PlatformLogsRequest.Path' Error:Field validation for 'Path' failed on the 'required' tag"},
		{name: "logs tail is not negative", request: &PlatformLogsRequest{Connection: connection(), Path: "/var/log", Tail: -1}, want: "'PlatformLogsRequest.Tail' Error:Field validation for 'Tail' failed on the 'min' tag"},
		{name: "up needs a command", request: &PlatformUpRequest{Connection: connection(), Vaults: vaults, Command: " "}, want: "'PlatformUpRequest.Command' Error:Field validation for 'Command' failed on the 'notblank' tag"},
		{name: "up needs an ssh vault", request: &PlatformUpRequest{Connection: connection(), Command: "-d img"}, want: "'PlatformUpRequest.Vaults.SshKeyId' Error:Field validation for 'SshKeyId' failed on the 'required' tag"},
		{name: "copy id needs a public key", request: &PlatformCopyIdRequest{PlatformCredConnection: PlatformCredConnection{Host: "db1", Port: 22, Username: "root", Password: "x", Platform: platform.Docker}}, want: "'PlatformCopyIdRequest.PublicKey' Error:Field validation for 'PublicKey' failed on the 'required' tag"},
		{
			name:    "keeper deploy states both ports",
			request: &KeeperDeployRequest{Plugin: "native_etcd", Name: "etcd1", Command: "-d etcd", KeeperPort: 2379, Connection: connection(), Vaults: vaults},
			want:    "'KeeperDeployRequest.DbPort' Error:Field validation for 'DbPort' failed on the 'required' tag",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := binding.Validator.ValidateStruct(test.request)
			if test.want == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("expected an error containing %q, got %v", test.want, err)
			}
		})
	}
}
