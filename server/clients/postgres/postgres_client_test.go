package postgres

import (
	"crypto/tls"
	"errors"
	"ivory/clients"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectedUrl string
		expectedErr error
	}{
		{
			name:        "default database",
			config:      Config{Host: "db1", Port: 5432, Username: "postgres", Password: "secret"},
			expectedUrl: "postgres://db1:5432/postgres",
		},
		{
			name:        "custom database",
			config:      Config{Host: "db1", Port: 5433, Database: "app", Username: "postgres", Password: "secret"},
			expectedUrl: "postgres://db1:5433/app",
		},
		{
			name:        "tls appends verify-ca",
			config:      Config{Host: "db1", Port: 5432, Username: "postgres", Password: "secret", TLS: &tls.Config{}},
			expectedUrl: "postgres://db1:5432/postgres?sslmode=verify-ca",
		},
		{
			name:        "missing host",
			config:      Config{Port: 5432},
			expectedErr: ErrHostOrPortNotSpecified,
		},
		{
			name:        "dash host",
			config:      Config{Host: "-", Port: 5432},
			expectedErr: ErrHostOrPortNotSpecified,
		},
		{
			name:        "missing port",
			config:      Config{Host: "db1"},
			expectedErr: ErrHostOrPortNotSpecified,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conConfig, url, err := Parse(tt.config)
			if tt.expectedErr != nil {
				if !errors.Is(err, tt.expectedErr) {
					t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if url != tt.expectedUrl {
				t.Errorf("expected url %q, got %q", tt.expectedUrl, url)
			}
			if conConfig.User != tt.config.Username {
				t.Errorf("expected user %q, got %q", tt.config.Username, conConfig.User)
			}
			if conConfig.Password != tt.config.Password {
				t.Errorf("expected password to be set")
			}
			if conConfig.RuntimeParams["application_name"] != tt.config.AppName {
				t.Errorf("expected application_name %q, got %q", tt.config.AppName, conConfig.RuntimeParams["application_name"])
			}
			if conConfig.ConnectTimeout != clients.IntegrationTimeout {
				t.Errorf("expected default connect timeout %v, got %v", clients.IntegrationTimeout, conConfig.ConnectTimeout)
			}
		})
	}
}
