package keeper

import (
	"crypto/tls"
	"testing"
)

func TestMap(t *testing.T) {
	t.Run("maps host, port, path and body", func(t *testing.T) {
		request := Request{Host: "db1", Port: 8008, Body: map[string]string{"key": "value"}}

		result := Map(request, "/cluster")

		if result.Host != "db1" || result.Port != 8008 || result.Path != "/cluster" {
			t.Errorf("expected host/port/path to carry over, got %+v", result)
		}
		if result.Credentials != nil {
			t.Errorf("expected nil credentials, got %+v", result.Credentials)
		}
	})

	t.Run("maps credentials when present", func(t *testing.T) {
		request := Request{
			Host:        "db1",
			Port:        8008,
			Credentials: &Credentials{Username: "user", Password: "pass"},
		}

		result := Map(request, "/config")

		if result.Credentials == nil {
			t.Fatal("expected credentials to be mapped")
		}
		if result.Credentials.Username != "user" || result.Credentials.Password != "pass" {
			t.Errorf("expected username/password to carry over, got %+v", result.Credentials)
		}
	})

	t.Run("maps TLS config", func(t *testing.T) {
		tlsConfig := &tls.Config{}
		request := Request{Host: "db1", Port: 8008, TlsConfig: tlsConfig}

		result := Map(request, "/config")

		if result.TLSConfig != tlsConfig {
			t.Errorf("expected TLS config to carry over unchanged")
		}
	})
}
