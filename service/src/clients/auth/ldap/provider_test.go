package ldap

import (
	"testing"
)

func TestProvider_Configured(t *testing.T) {
	t.Run("should return false when not configured", func(t *testing.T) {
		provider := NewProvider()

		if provider.Configured() {
			t.Error("Expected provider to not be configured")
		}
	})

	t.Run("should return true after configuration", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			Url:      "ldap://localhost:389",
			BindDN:   "cn=admin,dc=example,dc=com",
			BindPass: "password",
			BaseDN:   "dc=example,dc=com",
			Filter:   "(uid=%s)",
		}

		provider.SetConfig(config)

		if !provider.Configured() {
			t.Error("Expected provider to be configured")
		}
	})
}

func TestProvider_SetConfig(t *testing.T) {
	t.Run("should accept valid config without TLS", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			Url:      "ldap://localhost:389",
			BindDN:   "cn=admin,dc=example,dc=com",
			BindPass: "password",
			BaseDN:   "dc=example,dc=com",
			Filter:   "(uid=%s)",
		}

		err := provider.SetConfig(config)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if !provider.Configured() {
			t.Error("Expected provider to be configured")
		}
	})

	t.Run("should accept valid config with TLS", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			Url:      "ldaps://localhost:636",
			BindDN:   "cn=admin,dc=example,dc=com",
			BindPass: "password",
			BaseDN:   "dc=example,dc=com",
			Filter:   "(uid=%s)",
			Tls: &TlsConfig{
				CaCert: "-----BEGIN CERTIFICATE-----\nMIIC...\n-----END CERTIFICATE-----",
			},
		}

		err := provider.SetConfig(config)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
	})

	t.Run("should return error for empty URL", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			Url:      "",
			BindDN:   "cn=admin,dc=example,dc=com",
			BindPass: "password",
			BaseDN:   "dc=example,dc=com",
			Filter:   "(uid=%s)",
		}

		err := provider.SetConfig(config)

		if err == nil {
			t.Fatal("Expected error for empty URL, got nil")
		}
		if err.Error() != "url is not specified" {
			t.Errorf("Expected 'url is not specified', got: %v", err)
		}
	})

	t.Run("should return error for empty BindDN", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			Url:      "ldap://localhost:389",
			BindDN:   "",
			BindPass: "password",
			BaseDN:   "dc=example,dc=com",
			Filter:   "(uid=%s)",
		}

		err := provider.SetConfig(config)

		if err == nil {
			t.Fatal("Expected error for empty BindDN, got nil")
		}
		if err.Error() != "BindDN is not specified" {
			t.Errorf("Expected 'BindDN is not specified', got: %v", err)
		}
	})

	t.Run("should return error for empty BindPass", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			Url:      "ldap://localhost:389",
			BindDN:   "cn=admin,dc=example,dc=com",
			BindPass: "",
			BaseDN:   "dc=example,dc=com",
			Filter:   "(uid=%s)",
		}

		err := provider.SetConfig(config)

		if err == nil {
			t.Fatal("Expected error for empty BindPass, got nil")
		}
		if err.Error() != "BindPass is not specified" {
			t.Errorf("Expected 'BindPass is not specified', got: %v", err)
		}
	})

	t.Run("should return error for empty BaseDN", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			Url:      "ldap://localhost:389",
			BindDN:   "cn=admin,dc=example,dc=com",
			BindPass: "password",
			BaseDN:   "",
			Filter:   "(uid=%s)",
		}

		err := provider.SetConfig(config)

		if err == nil {
			t.Fatal("Expected error for empty BaseDN, got nil")
		}
		if err.Error() != "BaseDN is not specified" {
			t.Errorf("Expected 'BaseDN is not specified', got: %v", err)
		}
	})

	t.Run("should return error for empty Filter", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			Url:      "ldap://localhost:389",
			BindDN:   "cn=admin,dc=example,dc=com",
			BindPass: "password",
			BaseDN:   "dc=example,dc=com",
			Filter:   "",
		}

		err := provider.SetConfig(config)

		if err == nil {
			t.Fatal("Expected error for empty Filter, got nil")
		}
		if err.Error() != "filter is not specified" {
			t.Errorf("Expected 'filter is not specified', got: %v", err)
		}
	})

	t.Run("should return error for TLS config with non-ldaps scheme", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			Url:      "ldap://localhost:389",
			BindDN:   "cn=admin,dc=example,dc=com",
			BindPass: "password",
			BaseDN:   "dc=example,dc=com",
			Filter:   "(uid=%s)",
			Tls: &TlsConfig{
				CaCert: "-----BEGIN CERTIFICATE-----\nMIIC...\n-----END CERTIFICATE-----",
			},
		}

		err := provider.SetConfig(config)

		if err == nil {
			t.Fatal("Expected error for non-ldaps scheme with TLS, got nil")
		}
		if err.Error() != "scheme is not ldaps" {
			t.Errorf("Expected 'scheme is not ldaps', got: %v", err)
		}
	})

	t.Run("should accept ldaps:// with TLS config", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			Url:      "ldaps://localhost:636",
			BindDN:   "cn=admin,dc=example,dc=com",
			BindPass: "password",
			BaseDN:   "dc=example,dc=com",
			Filter:   "(uid=%s)",
			Tls: &TlsConfig{
				CaCert: "-----BEGIN CERTIFICATE-----\nMIIC...\n-----END CERTIFICATE-----",
			},
		}

		err := provider.SetConfig(config)

		if err != nil {
			t.Fatalf("Expected no error for ldaps:// with TLS, got: %v", err)
		}
	})

	t.Run("should accept ldaps:// without TLS config", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			Url:      "ldaps://localhost:636",
			BindDN:   "cn=admin,dc=example,dc=com",
			BindPass: "password",
			BaseDN:   "dc=example,dc=com",
			Filter:   "(uid=%s)",
			Tls:      nil,
		}

		err := provider.SetConfig(config)

		if err != nil {
			t.Fatalf("Expected no error for ldaps:// without TLS config, got: %v", err)
		}
	})

	t.Run("should accept TLS config with empty CaCert", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			Url:      "ldaps://localhost:636",
			BindDN:   "cn=admin,dc=example,dc=com",
			BindPass: "password",
			BaseDN:   "dc=example,dc=com",
			Filter:   "(uid=%s)",
			Tls: &TlsConfig{
				CaCert: "",
			},
		}

		err := provider.SetConfig(config)

		if err != nil {
			t.Fatalf("Expected no error for empty CaCert, got: %v", err)
		}
	})

	t.Run("should overwrite previous config", func(t *testing.T) {
		provider := NewProvider()
		config1 := Config{
			Url:      "ldap://localhost:389",
			BindDN:   "cn=admin,dc=example,dc=com",
			BindPass: "password1",
			BaseDN:   "dc=example,dc=com",
			Filter:   "(uid=%s)",
		}
		config2 := Config{
			Url:      "ldap://server2:389",
			BindDN:   "cn=user,dc=test,dc=com",
			BindPass: "password2",
			BaseDN:   "dc=test,dc=com",
			Filter:   "(cn=%s)",
		}

		provider.SetConfig(config1)
		err := provider.SetConfig(config2)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if !provider.Configured() {
			t.Error("Expected provider to be configured with new config")
		}
	})
}

func TestProvider_DeleteConfig(t *testing.T) {
	t.Run("should delete config", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			Url:      "ldap://localhost:389",
			BindDN:   "cn=admin,dc=example,dc=com",
			BindPass: "password",
			BaseDN:   "dc=example,dc=com",
			Filter:   "(uid=%s)",
		}

		provider.SetConfig(config)
		provider.DeleteConfig()

		if provider.Configured() {
			t.Error("Expected provider to not be configured after deletion")
		}
	})

	t.Run("should handle deleting when not configured", func(t *testing.T) {
		provider := NewProvider()

		// Should not panic
		provider.DeleteConfig()

		if provider.Configured() {
			t.Error("Expected provider to not be configured")
		}
	})
}

func TestProvider_Verify_NotConfigured(t *testing.T) {
	t.Run("should return error when not configured", func(t *testing.T) {
		provider := NewProvider()

		_, err := provider.Verify(Login{
			Username: "testuser",
			Password: "password",
		})

		if err == nil {
			t.Fatal("Expected error when not configured, got nil")
		}
		if err.Error() != "config is not configured" {
			t.Errorf("Expected 'config is not configured', got: %v", err)
		}
	})
}
