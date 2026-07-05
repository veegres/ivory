package oidc

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
			IssuerURL:    "https://accounts.google.com",
			ClientID:     "client-id-123",
			ClientSecret: "client-secret-456",
			RedirectURL:  "https://example.com/callback",
		}

		provider.SetConfig(config)

		if !provider.Configured() {
			t.Error("Expected provider to be configured")
		}
	})
}

func TestProvider_SetConfig(t *testing.T) {
	t.Run("should accept valid config", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			IssuerURL:    "https://accounts.google.com",
			ClientID:     "client-id-123",
			ClientSecret: "client-secret-456",
			RedirectURL:  "https://example.com/callback",
		}

		err := provider.SetConfig(config)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if !provider.Configured() {
			t.Error("Expected provider to be configured")
		}
	})

	t.Run("should return error for empty IssuerURL", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			IssuerURL:    "",
			ClientID:     "client-id-123",
			ClientSecret: "client-secret-456",
			RedirectURL:  "https://example.com/callback",
		}

		err := provider.SetConfig(config)

		if err == nil {
			t.Fatal("Expected error for empty IssuerURL, got nil")
		}
		if err.Error() != "IssuerURL is not specified" {
			t.Errorf("Expected 'IssuerURL is not specified', got: %v", err)
		}
	})

	t.Run("should return error for empty ClientID", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			IssuerURL:    "https://accounts.google.com",
			ClientID:     "",
			ClientSecret: "client-secret-456",
			RedirectURL:  "https://example.com/callback",
		}

		err := provider.SetConfig(config)

		if err == nil {
			t.Fatal("Expected error for empty ClientID, got nil")
		}
		if err.Error() != "ClientID is not specified" {
			t.Errorf("Expected 'ClientID is not specified', got: %v", err)
		}
	})

	t.Run("should return error for empty ClientSecret", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			IssuerURL:    "https://accounts.google.com",
			ClientID:     "client-id-123",
			ClientSecret: "",
			RedirectURL:  "https://example.com/callback",
		}

		err := provider.SetConfig(config)

		if err == nil {
			t.Fatal("Expected error for empty ClientSecret, got nil")
		}
		if err.Error() != "ClientSecret is not specified" {
			t.Errorf("Expected 'ClientSecret is not specified', got: %v", err)
		}
	})

	t.Run("should return error for empty RedirectURL", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			IssuerURL:    "https://accounts.google.com",
			ClientID:     "client-id-123",
			ClientSecret: "client-secret-456",
			RedirectURL:  "",
		}

		err := provider.SetConfig(config)

		if err == nil {
			t.Fatal("Expected error for empty RedirectURL, got nil")
		}
		if err.Error() != "RedirectURL is not specified" {
			t.Errorf("Expected 'RedirectURL is not specified', got: %v", err)
		}
	})

	t.Run("should validate fields in order", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			IssuerURL:    "",
			ClientID:     "",
			ClientSecret: "",
			RedirectURL:  "",
		}

		err := provider.SetConfig(config)

		if err == nil {
			t.Fatal("Expected error for empty config, got nil")
		}
		// Should fail on first validation (IssuerURL)
		if err.Error() != "IssuerURL is not specified" {
			t.Errorf("Expected 'IssuerURL is not specified', got: %v", err)
		}
	})

	t.Run("should overwrite previous config", func(t *testing.T) {
		provider := NewProvider()
		config1 := Config{
			IssuerURL:    "https://accounts.google.com",
			ClientID:     "client-id-1",
			ClientSecret: "client-secret-1",
			RedirectURL:  "https://example.com/callback1",
		}
		config2 := Config{
			IssuerURL:    "https://auth.example.com",
			ClientID:     "client-id-2",
			ClientSecret: "client-secret-2",
			RedirectURL:  "https://example.com/callback2",
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
			IssuerURL:    "https://accounts.google.com",
			ClientID:     "client-id-123",
			ClientSecret: "client-secret-456",
			RedirectURL:  "https://example.com/callback",
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

	t.Run("should delete all internal state", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			IssuerURL:    "https://accounts.google.com",
			ClientID:     "client-id-123",
			ClientSecret: "client-secret-456",
			RedirectURL:  "https://example.com/callback",
		}

		provider.SetConfig(config)
		provider.DeleteConfig()

		// Verify config is nil
		if provider.config != nil {
			t.Error("Expected config to be nil after deletion")
		}
		// Verify oauthConfig is nil
		if provider.oauthConfig != nil {
			t.Error("Expected oauthConfig to be nil after deletion")
		}
		// Verify oauthVerifier is nil
		if provider.oauthVerifier != nil {
			t.Error("Expected oauthVerifier to be nil after deletion")
		}
	})
}

func TestProvider_initialize(t *testing.T) {
	t.Run("should return error when not configured", func(t *testing.T) {
		provider := NewProvider()

		err := provider.initialize()

		if err == nil {
			t.Fatal("Expected error when not configured, got nil")
		}
		if err.Error() != "config is not configured" {
			t.Errorf("Expected 'config is not configured', got: %v", err)
		}
	})
}
