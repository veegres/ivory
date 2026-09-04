package auth

import (
	"errors"
	"ivory/clients/auth/basic"
	"ivory/clients/auth/ldap"
	"ivory/clients/auth/oidc"
	"ivory/clients/storage"
	"ivory/core/service/encryption"
	"ivory/core/service/secret"
	"ivory/core/service/token"
	"ivory/features/permission"
	"path/filepath"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/golang-jwt/jwt/v5"
)

type fakeUserVerifier struct {
	users map[string]string
}

func (v *fakeUserVerifier) VerifyPassword(username string, password string) error {
	stored, ok := v.users[username]
	if !ok || stored != password {
		return errors.New("credentials are not correct")
	}
	return nil
}

func createTestBasicProvider(t *testing.T) *basic.Provider {
	t.Helper()
	provider := basic.NewProvider(&fakeUserVerifier{users: map[string]string{"admin": "password123"}})
	if err := provider.SetConfig(basic.Config{}); err != nil {
		t.Fatalf("failed to configure basic provider: %v", err)
	}
	return provider
}

func createTestAuthService(t *testing.T) *Service {
	t.Helper()

	db, errOpen := bolt.Open(filepath.Join(t.TempDir(), "test.db"), 0600, nil)
	if errOpen != nil {
		t.Fatalf("failed to open test database: %v", errOpen)
	}
	t.Cleanup(func() {
		db.Close()
	})

	secretService := secret.NewService(
		secret.NewRepository(storage.NewDbBucket[string](db, "Secret")),
		encryption.NewService(),
	)
	if err := secretService.SetDefault(); err != nil {
		t.Fatalf("failed to set default secret: %v", err)
	}

	permissionService := permission.NewService(
		permission.NewRepository(storage.NewDbBucket[permission.PermissionMap](db, "Permission")),
	)

	return NewService(token.NewService(secretService), createTestBasicProvider(t), ldap.NewProvider(), oidc.NewProvider(), permissionService)
}

func TestServiceGetSupportedTypes(t *testing.T) {
	s := createTestAuthService(t)

	supported := s.GetSupportedTypes()
	if len(supported) != 1 || supported[0] != BASIC {
		t.Fatalf("expected only BASIC to be supported, got %v", supported)
	}
}

func TestServiceGenerateBasicAuthToken(t *testing.T) {
	s := createTestAuthService(t)

	t.Run("valid credentials produce a usable token", func(t *testing.T) {
		token, exp, err := s.GenerateBasicAuthToken(basic.Login{Username: "admin", Password: "password123"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if token == "" {
			t.Fatalf("expected a non-empty token")
		}
		if exp == nil || !exp.After(time.Now()) {
			t.Fatalf("expected a future expiration time, got %v", exp)
		}

		valid, username, authType, errParse := s.ParseAuthToken(token, nil)
		if !valid || errParse != nil {
			t.Fatalf("expected the generated token to parse as valid, got valid=%v err=%v", valid, errParse)
		}
		if username != "admin" {
			t.Fatalf("expected username 'admin', got %q", username)
		}
		if authType == nil || *authType != BASIC {
			t.Fatalf("expected auth type BASIC, got %v", authType)
		}
	})

	t.Run("invalid credentials are rejected", func(t *testing.T) {
		_, _, err := s.GenerateBasicAuthToken(basic.Login{Username: "admin", Password: "wrong"})
		if err == nil {
			t.Fatalf("expected an error for invalid credentials")
		}
	})
}

func TestServiceParseAuthToken(t *testing.T) {
	s := createTestAuthService(t)

	t.Run("propagates an upstream token error", func(t *testing.T) {
		upstreamErr := jwt.ErrTokenMalformed
		valid, _, _, err := s.ParseAuthToken("", upstreamErr)
		if valid {
			t.Fatalf("expected valid=false")
		}
		if err != upstreamErr {
			t.Fatalf("expected the upstream error to be returned, got %v", err)
		}
	})

	t.Run("garbage token is invalid", func(t *testing.T) {
		valid, _, _, err := s.ParseAuthToken("garbage", nil)
		if valid {
			t.Fatalf("expected valid=false for a garbage token")
		}
		if err == nil {
			t.Fatalf("expected an error for a garbage token")
		}
	})

	t.Run("token signed with a different secret is rejected", func(t *testing.T) {
		wrongToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iss": "ivory",
			"sub": "admin",
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(time.Hour).Unix(),
			"frm": BASIC,
		})
		signed, errSign := wrongToken.SignedString([]byte("0123456789012345"))
		if errSign != nil {
			t.Fatalf("failed to sign token: %v", errSign)
		}

		valid, _, _, err := s.ParseAuthToken(signed, nil)
		if valid {
			t.Fatalf("expected valid=false for a token signed with the wrong secret")
		}
		if err == nil {
			t.Fatalf("expected an error")
		}
	})

	t.Run("token missing the frm claim is valid but reports the parse error", func(t *testing.T) {
		signed, _, errSign := s.tokenService.Generate("admin", nil, time.Hour)
		if errSign != nil {
			t.Fatalf("failed to sign token: %v", errSign)
		}

		valid, username, authType, err := s.ParseAuthToken(signed, nil)
		if !valid {
			t.Fatalf("expected valid=true even without a parsable frm claim")
		}
		if username != "" {
			t.Fatalf("expected an empty username, got %q", username)
		}
		if authType != nil {
			t.Fatalf("expected a nil auth type, got %v", authType)
		}
		if err != ErrInvalidTokenCannotParseAuthType {
			t.Fatalf("expected ErrInvalidTokenCannotParseAuthType, got %v", err)
		}
	})

	t.Run("when no auth type is configured, the check is skipped entirely", func(t *testing.T) {
		unconfigured := NewService(s.tokenService, basic.NewProvider(&fakeUserVerifier{}), ldap.NewProvider(), oidc.NewProvider(), s.permissionService)
		valid, username, authType, err := unconfigured.ParseAuthToken("anything", nil)
		if !valid {
			t.Fatalf("expected valid=true when auth is disabled")
		}
		if username != "" || authType != nil {
			t.Fatalf("expected no username/authType, got %q %v", username, authType)
		}
		if err != ErrAuthDisabled {
			t.Fatalf("expected ErrAuthDisabled, got %v", err)
		}
	})
}

func TestServiceParseAuthTokenWithFallback(t *testing.T) {
	s := createTestAuthService(t)
	token, _, err := s.GenerateBasicAuthToken(basic.Login{Username: "admin", Password: "password123"})
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	t.Run("uses the primary token when it is valid", func(t *testing.T) {
		valid, username, _, errParse := s.ParseAuthTokenWithFallback(token, nil, "garbage", nil)
		if !valid || errParse != nil {
			t.Fatalf("expected the primary token to validate, got valid=%v err=%v", valid, errParse)
		}
		if username != "admin" {
			t.Fatalf("expected username 'admin', got %q", username)
		}
	})

	t.Run("falls back to the secondary token when the primary is invalid", func(t *testing.T) {
		valid, username, _, errParse := s.ParseAuthTokenWithFallback("garbage", nil, token, nil)
		if !valid || errParse != nil {
			t.Fatalf("expected the fallback token to validate, got valid=%v err=%v", valid, errParse)
		}
		if username != "admin" {
			t.Fatalf("expected username 'admin', got %q", username)
		}
	})

	t.Run("reports the fallback error when both tokens are invalid", func(t *testing.T) {
		valid, _, _, errParse := s.ParseAuthTokenWithFallback("garbage", nil, "also-garbage", nil)
		if valid {
			t.Fatalf("expected valid=false")
		}
		if errParse == nil {
			t.Fatalf("expected an error")
		}
	})
}

func TestServiceGenerateLdapAndOidcAuthToken(t *testing.T) {
	s := createTestAuthService(t)

	t.Run("ldap generation fails when the provider is not configured", func(t *testing.T) {
		if _, _, err := s.GenerateLdapAuthToken(ldap.Login{Username: "user", Password: "pw"}); err == nil {
			t.Fatalf("expected an error when the ldap provider is not configured")
		}
	})

	t.Run("oidc generation fails when the provider is not configured", func(t *testing.T) {
		if _, _, err := s.GenerateOidcAuthToken("code"); err == nil {
			t.Fatalf("expected an error when the oidc provider is not configured")
		}
	})
}
