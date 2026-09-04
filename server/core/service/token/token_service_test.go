package token

import (
	"errors"
	"ivory/clients/storage"
	"ivory/core/service/encryption"
	"ivory/core/service/secret"
	"path/filepath"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/golang-jwt/jwt/v5"
)

func newTestService(t *testing.T) *Service {
	t.Helper()

	db, err := bolt.Open(filepath.Join(t.TempDir(), "test.db"), 0600, nil)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	secretService := secret.NewService(
		secret.NewRepository(storage.NewDbBucket[string](db, "Secret")),
		encryption.NewService(),
	)
	if errSecret := secretService.SetDefault(); errSecret != nil {
		t.Fatalf("failed to set default secret: %v", errSecret)
	}
	return NewService(secretService)
}

func TestServiceIssuer(t *testing.T) {
	if issuer := newTestService(t).Issuer(); issuer != "ivory" {
		t.Fatalf("expected issuer 'ivory', got %q", issuer)
	}
}

func TestServiceGenerateAndParse(t *testing.T) {
	s := newTestService(t)

	t.Run("round trips the subject and the extra claims", func(t *testing.T) {
		token, exp, err := s.Generate("alice", jwt.MapClaims{"jti": "id-1"}, time.Hour)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if token == "" {
			t.Fatalf("expected a non-empty token")
		}
		if exp == nil || !exp.After(time.Now()) {
			t.Fatalf("expected a future expiration, got %v", exp)
		}

		claims, errParse := s.Parse(token)
		if errParse != nil {
			t.Fatalf("expected no error, got %v", errParse)
		}
		subject, errSubject := claims.GetSubject()
		if errSubject != nil || subject != "alice" {
			t.Fatalf("expected subject 'alice', got %q (%v)", subject, errSubject)
		}
		if claims["jti"] != "id-1" {
			t.Fatalf("expected the extra claim to survive, got %v", claims["jti"])
		}
		if claims["iss"] != "ivory" {
			t.Fatalf("expected issuer 'ivory', got %v", claims["iss"])
		}
	})

	t.Run("extra claims cannot overwrite the standard ones", func(t *testing.T) {
		token, _, err := s.Generate("alice", jwt.MapClaims{"iss": "somebody-else", "sub": "mallory"}, time.Hour)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		claims, errParse := s.Parse(token)
		if errParse != nil {
			t.Fatalf("expected no error, got %v", errParse)
		}
		if claims["iss"] != "ivory" {
			t.Fatalf("expected issuer 'ivory', got %v", claims["iss"])
		}
		if claims["sub"] != "alice" {
			t.Fatalf("expected subject 'alice', got %v", claims["sub"])
		}
	})

	t.Run("nil extra claims are fine", func(t *testing.T) {
		token, _, err := s.Generate("alice", nil, time.Hour)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if _, errParse := s.Parse(token); errParse != nil {
			t.Fatalf("expected no error, got %v", errParse)
		}
	})
}

func TestServiceParseRejections(t *testing.T) {
	s := newTestService(t)

	t.Run("garbage is rejected", func(t *testing.T) {
		if _, err := s.Parse("garbage"); err == nil {
			t.Fatalf("expected an error")
		}
	})

	t.Run("empty token is rejected", func(t *testing.T) {
		if _, err := s.Parse(""); err == nil {
			t.Fatalf("expected an error")
		}
	})

	t.Run("an expired token reports jwt.ErrTokenExpired", func(t *testing.T) {
		token, _, err := s.Generate("alice", nil, -time.Minute)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		_, errParse := s.Parse(token)
		if !errors.Is(errParse, jwt.ErrTokenExpired) {
			t.Fatalf("expected jwt.ErrTokenExpired, got %v", errParse)
		}
	})

	t.Run("a token signed with another key is rejected", func(t *testing.T) {
		other := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iss": "ivory",
			"sub": "alice",
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(time.Hour).Unix(),
		})
		signed, errSign := other.SignedString([]byte("0123456789012345"))
		if errSign != nil {
			t.Fatalf("failed to sign: %v", errSign)
		}
		if _, err := s.Parse(signed); err == nil {
			t.Fatalf("expected an error")
		}
	})

	t.Run("a token from another issuer is rejected", func(t *testing.T) {
		other := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iss": "somebody-else",
			"sub": "alice",
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(time.Hour).Unix(),
		})
		signed, errSign := other.SignedString(s.secretService.GetByte())
		if errSign != nil {
			t.Fatalf("failed to sign: %v", errSign)
		}
		if _, err := s.Parse(signed); err == nil {
			t.Fatalf("expected an error")
		}
	})
}
