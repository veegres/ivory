package basic

import (
	"errors"
	"testing"
)

var errFakeCredentials = errors.New("credentials are not correct")

type fakeVerifier struct {
	users  map[string]string
	called int
}

func (v *fakeVerifier) VerifyPassword(username string, password string) error {
	v.called++
	stored, ok := v.users[username]
	if !ok || stored != password {
		return errFakeCredentials
	}
	return nil
}

func createTestProvider() (*Provider, *fakeVerifier) {
	verifier := &fakeVerifier{users: map[string]string{"admin": "password123", "Admin": "Password123"}}
	return NewProvider(verifier), verifier
}

func TestProvider_Configured(t *testing.T) {
	t.Run("should return false when not configured", func(t *testing.T) {
		provider, _ := createTestProvider()

		if provider.Configured() {
			t.Error("Expected provider to not be configured")
		}
	})

	t.Run("should return true after configuration", func(t *testing.T) {
		provider, _ := createTestProvider()

		if err := provider.SetConfig(Config{}); err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if !provider.Configured() {
			t.Error("Expected provider to be configured")
		}
	})
}

func TestProvider_SetConfig(t *testing.T) {
	t.Run("should accept the config, since it carries no credentials any more", func(t *testing.T) {
		provider, _ := createTestProvider()

		err := provider.SetConfig(Config{})

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if !provider.Configured() {
			t.Error("Expected provider to be configured")
		}
	})
}

func TestProvider_DeleteConfig(t *testing.T) {
	t.Run("should delete config", func(t *testing.T) {
		provider, _ := createTestProvider()

		if err := provider.SetConfig(Config{}); err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		provider.DeleteConfig()

		if provider.Configured() {
			t.Error("Expected provider to not be configured after deletion")
		}
	})

	t.Run("should handle deleting when not configured", func(t *testing.T) {
		provider, _ := createTestProvider()

		// Should not panic
		provider.DeleteConfig()

		if provider.Configured() {
			t.Error("Expected provider to not be configured")
		}
	})
}

func TestProvider_Verify(t *testing.T) {
	t.Run("should verify correct credentials", func(t *testing.T) {
		provider, verifier := createTestProvider()
		if err := provider.SetConfig(Config{}); err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

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
		if verifier.called != 1 {
			t.Errorf("Expected the verifier to be asked once, got %d", verifier.called)
		}
	})

	t.Run("should reject unknown username", func(t *testing.T) {
		provider, _ := createTestProvider()
		if err := provider.SetConfig(Config{}); err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		_, err := provider.Verify(Login{
			Username: "wronguser",
			Password: "password123",
		})

		if !errors.Is(err, errFakeCredentials) {
			t.Errorf("Expected the verifier error, got: %v", err)
		}
	})

	t.Run("should reject wrong password", func(t *testing.T) {
		provider, _ := createTestProvider()
		if err := provider.SetConfig(Config{}); err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		_, err := provider.Verify(Login{
			Username: "admin",
			Password: "wrongpassword",
		})

		if !errors.Is(err, errFakeCredentials) {
			t.Errorf("Expected the verifier error, got: %v", err)
		}
	})

	t.Run("should reject an empty username without asking the verifier", func(t *testing.T) {
		provider, verifier := createTestProvider()
		if err := provider.SetConfig(Config{}); err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		_, err := provider.Verify(Login{
			Username: "",
			Password: "password123",
		})

		if !errors.Is(err, ErrUsernameNotSpecified) {
			t.Fatalf("Expected ErrUsernameNotSpecified, got: %v", err)
		}
		if verifier.called != 0 {
			t.Errorf("Expected the verifier not to be asked, got %d calls", verifier.called)
		}
	})

	t.Run("should reject an empty password without asking the verifier", func(t *testing.T) {
		provider, verifier := createTestProvider()
		if err := provider.SetConfig(Config{}); err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		_, err := provider.Verify(Login{
			Username: "admin",
			Password: "",
		})

		if !errors.Is(err, ErrPasswordNotSpecified) {
			t.Fatalf("Expected ErrPasswordNotSpecified, got: %v", err)
		}
		if verifier.called != 0 {
			t.Errorf("Expected the verifier not to be asked, got %d calls", verifier.called)
		}
	})

	t.Run("should be case-sensitive for username and password", func(t *testing.T) {
		provider, _ := createTestProvider()
		if err := provider.SetConfig(Config{}); err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if _, err := provider.Verify(Login{Username: "admin", Password: "Password123"}); err == nil {
			t.Fatal("Expected error for case-sensitive password mismatch, got nil")
		}

		username, err := provider.Verify(Login{Username: "Admin", Password: "Password123"})
		if err != nil {
			t.Fatalf("Expected no error for correct case, got: %v", err)
		}
		if username != "Admin" {
			t.Errorf("Expected username 'Admin', got '%s'", username)
		}
	})

	t.Run("should return error when not configured", func(t *testing.T) {
		provider, verifier := createTestProvider()

		_, err := provider.Verify(Login{
			Username: "admin",
			Password: "password123",
		})

		if !errors.Is(err, ErrConfigNotConfigured) {
			t.Fatalf("Expected ErrConfigNotConfigured, got: %v", err)
		}
		if verifier.called != 0 {
			t.Errorf("Expected the verifier not to be asked, got %d calls", verifier.called)
		}
	})
}

func TestProvider_Connect(t *testing.T) {
	t.Run("should return error for obsolete method", func(t *testing.T) {
		provider, _ := createTestProvider()

		err := provider.Connect(Config{})

		if !errors.Is(err, ErrConnectionObsolete) {
			t.Fatalf("Expected ErrConnectionObsolete, got: %v", err)
		}
	})
}
