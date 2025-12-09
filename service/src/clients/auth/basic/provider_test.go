package basic

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
			Username: "admin",
			Password: "password123",
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
			Username: "admin",
			Password: "password123",
		}

		err := provider.SetConfig(config)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if !provider.Configured() {
			t.Error("Expected provider to be configured")
		}
	})

	t.Run("should return error for empty username", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			Username: "",
			Password: "password123",
		}

		err := provider.SetConfig(config)

		if err == nil {
			t.Fatal("Expected error for empty username, got nil")
		}
		if err.Error() != "username is not specified" {
			t.Errorf("Expected 'username is not specified', got: %v", err)
		}
	})

	t.Run("should return error for empty password", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			Username: "admin",
			Password: "",
		}

		err := provider.SetConfig(config)

		if err == nil {
			t.Fatal("Expected error for empty password, got nil")
		}
		if err.Error() != "password is not specified" {
			t.Errorf("Expected 'password is not specified', got: %v", err)
		}
	})

	t.Run("should return error for both empty", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			Username: "",
			Password: "",
		}

		err := provider.SetConfig(config)

		if err == nil {
			t.Fatal("Expected error for empty credentials, got nil")
		}
		// Will fail on username check first
		if err.Error() != "username is not specified" {
			t.Errorf("Expected 'username is not specified', got: %v", err)
		}
	})

	t.Run("should overwrite previous config", func(t *testing.T) {
		provider := NewProvider()
		config1 := Config{
			Username: "admin",
			Password: "password123",
		}
		config2 := Config{
			Username: "user",
			Password: "newpassword",
		}

		provider.SetConfig(config1)
		provider.SetConfig(config2)

		// Verify new config by trying to login
		_, err := provider.Verify(Login{Username: "user", Password: "newpassword"})
		if err != nil {
			t.Errorf("Expected new config to be active, got error: %v", err)
		}

		// Old credentials should fail
		_, err = provider.Verify(Login{Username: "admin", Password: "password123"})
		if err == nil {
			t.Error("Expected old credentials to fail")
		}
	})
}

func TestProvider_DeleteConfig(t *testing.T) {
	t.Run("should delete config", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			Username: "admin",
			Password: "password123",
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

func TestProvider_Verify(t *testing.T) {
	t.Run("should verify correct credentials", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			Username: "admin",
			Password: "password123",
		}
		provider.SetConfig(config)

		username, err := provider.Verify(Login{
			Username: "admin",
			Password: "password123",
		})

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if username != "admin" {
			t.Errorf("Expected username 'admin', got '%s'", username)
		}
	})

	t.Run("should reject wrong username", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			Username: "admin",
			Password: "password123",
		}
		provider.SetConfig(config)

		_, err := provider.Verify(Login{
			Username: "wronguser",
			Password: "password123",
		})

		if err == nil {
			t.Fatal("Expected error for wrong username, got nil")
		}
		if err.Error() != "credentials are not correct" {
			t.Errorf("Expected 'credentials are not correct', got: %v", err)
		}
	})

	t.Run("should reject wrong password", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			Username: "admin",
			Password: "password123",
		}
		provider.SetConfig(config)

		_, err := provider.Verify(Login{
			Username: "admin",
			Password: "wrongpassword",
		})

		if err == nil {
			t.Fatal("Expected error for wrong password, got nil")
		}
		if err.Error() != "credentials are not correct" {
			t.Errorf("Expected 'credentials are not correct', got: %v", err)
		}
	})

	t.Run("should reject both wrong", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			Username: "admin",
			Password: "password123",
		}
		provider.SetConfig(config)

		_, err := provider.Verify(Login{
			Username: "wronguser",
			Password: "wrongpassword",
		})

		if err == nil {
			t.Fatal("Expected error for wrong credentials, got nil")
		}
		if err.Error() != "credentials are not correct" {
			t.Errorf("Expected 'credentials are not correct', got: %v", err)
		}
	})

	t.Run("should handle empty username in login", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			Username: "admin",
			Password: "password123",
		}
		provider.SetConfig(config)

		_, err := provider.Verify(Login{
			Username: "",
			Password: "password123",
		})

		if err == nil {
			t.Fatal("Expected error for empty username, got nil")
		}
	})

	t.Run("should handle empty password in login", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			Username: "admin",
			Password: "password123",
		}
		provider.SetConfig(config)

		_, err := provider.Verify(Login{
			Username: "admin",
			Password: "",
		})

		if err == nil {
			t.Fatal("Expected error for empty password, got nil")
		}
	})

	t.Run("should be case-sensitive for username", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			Username: "Admin",
			Password: "password123",
		}
		provider.SetConfig(config)

		// Wrong case
		_, err := provider.Verify(Login{
			Username: "admin",
			Password: "password123",
		})

		if err == nil {
			t.Fatal("Expected error for case-sensitive username mismatch, got nil")
		}

		// Correct case
		username, err := provider.Verify(Login{
			Username: "Admin",
			Password: "password123",
		})

		if err != nil {
			t.Fatalf("Expected no error for correct case, got: %v", err)
		}
		if username != "Admin" {
			t.Errorf("Expected username 'Admin', got '%s'", username)
		}
	})

	t.Run("should be case-sensitive for password", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			Username: "admin",
			Password: "Password123",
		}
		provider.SetConfig(config)

		// Wrong case
		_, err := provider.Verify(Login{
			Username: "admin",
			Password: "password123",
		})

		if err == nil {
			t.Fatal("Expected error for case-sensitive password mismatch, got nil")
		}

		// Correct case
		_, err = provider.Verify(Login{
			Username: "admin",
			Password: "Password123",
		})

		if err != nil {
			t.Fatalf("Expected no error for correct case, got: %v", err)
		}
	})

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

func TestProvider_Connect(t *testing.T) {
	t.Run("should return error for obsolete method", func(t *testing.T) {
		provider := NewProvider()
		config := Config{
			Username: "admin",
			Password: "password123",
		}

		err := provider.Connect(config)

		if err == nil {
			t.Fatal("Expected error for obsolete method, got nil")
		}
		if err.Error() != "connection is obsolete" {
			t.Errorf("Expected 'connection is obsolete', got: %v", err)
		}
	})

	t.Run("should return error even with empty config", func(t *testing.T) {
		provider := NewProvider()
		config := Config{}

		err := provider.Connect(config)

		if err == nil {
			t.Fatal("Expected error for obsolete method, got nil")
		}
	})
}
