package auth

import (
	"ivory/clients/auth/basic"
	"ivory/clients/auth/ldap"
	"ivory/clients/auth/oidc"
	"ivory/clients/storage"
	"ivory/core/config"
	"ivory/core/service/encryption"
	"ivory/core/service/secret"
	"ivory/features/permission"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
)

func createTestRouter(t *testing.T) *Router {
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

	basicProvider := basic.NewProvider()
	if err := basicProvider.SetConfig(basic.Config{Username: "admin", Password: "password123"}); err != nil {
		t.Fatalf("failed to configure basic provider: %v", err)
	}

	permissionService := permission.NewService(
		permission.NewRepository(storage.NewDbBucket[permission.PermissionMap](db, "Permission")),
	)

	authService := NewService(secretService, basicProvider, ldap.NewProvider(), oidc.NewProvider(), permissionService)
	return NewRouter(authService, "/", false)
}

func createValidToken(t *testing.T, r *Router) string {
	t.Helper()
	token, _, err := r.authService.GenerateBasicAuthToken(basic.Login{Username: "admin", Password: "password123"})
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	return token
}

func runMiddleware(handler gin.HandlerFunc, authHeader string, tokenCookie string) (*httptest.ResponseRecorder, *gin.Context) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(w)
	context.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	if authHeader != "" {
		context.Request.Header.Set("Authorization", authHeader)
	}
	if tokenCookie != "" {
		context.Request.AddCookie(&http.Cookie{Name: "token", Value: tokenCookie})
	}

	handler(context)

	return w, context
}

func runValidateMiddleware(r *Router, authHeader string, tokenCookie string) (*httptest.ResponseRecorder, *gin.Context) {
	return runMiddleware(r.ValidateMiddleware(), authHeader, tokenCookie)
}

func runAuthContextMiddleware(r *Router, authHeader string, tokenCookie string) (*httptest.ResponseRecorder, *gin.Context) {
	return runMiddleware(r.AuthContextMiddleware(), authHeader, tokenCookie)
}

func TestValidateMiddleware(t *testing.T) {
	t.Run("valid header token is accepted", func(t *testing.T) {
		r := createTestRouter(t)
		token := createValidToken(t, r)

		w, context := runValidateMiddleware(r, "Bearer "+token, "")

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}
		if context.IsAborted() {
			t.Fatalf("expected middleware to not abort")
		}
		if username, _ := context.Get(env.AuthContextKey.Username); username != "admin" {
			t.Fatalf("expected username 'admin', got %v", username)
		}
	})

	t.Run("valid cookie token is accepted when no header is present", func(t *testing.T) {
		r := createTestRouter(t)
		token := createValidToken(t, r)

		w, context := runValidateMiddleware(r, "", token)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}
		if username, _ := context.Get(env.AuthContextKey.Username); username != "admin" {
			t.Fatalf("expected username 'admin', got %v", username)
		}
	})

	t.Run("falls back to cookie token when header token is invalid", func(t *testing.T) {
		r := createTestRouter(t)
		token := createValidToken(t, r)

		w, context := runValidateMiddleware(r, "Bearer garbage-token", token)

		if w.Code != http.StatusOK {
			t.Fatalf("expected fallback to cookie to succeed with status 200, got %d", w.Code)
		}
		if context.IsAborted() {
			t.Fatalf("expected middleware to not abort when a valid cookie token is available as fallback")
		}
		if username, _ := context.Get(env.AuthContextKey.Username); username != "admin" {
			t.Fatalf("expected username 'admin', got %v", username)
		}
	})

	t.Run("rejects request when both header and cookie tokens are invalid", func(t *testing.T) {
		r := createTestRouter(t)

		w, context := runValidateMiddleware(r, "Bearer garbage-token", "also-garbage")

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected status 401, got %d", w.Code)
		}
		if !context.IsAborted() {
			t.Fatalf("expected middleware to abort")
		}
	})

	t.Run("rejects request with no header and no cookie", func(t *testing.T) {
		r := createTestRouter(t)

		w, _ := runValidateMiddleware(r, "", "")

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected status 401, got %d", w.Code)
		}
	})
}

func TestAuthContextMiddleware(t *testing.T) {
	t.Run("valid token sets context and never aborts", func(t *testing.T) {
		r := createTestRouter(t)
		token := createValidToken(t, r)

		w, context := runAuthContextMiddleware(r, "Bearer "+token, "")

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}
		if context.IsAborted() {
			t.Fatalf("expected middleware to never abort")
		}
		if authorised, _ := context.Get(env.AuthContextKey.Authorised); authorised != true {
			t.Fatalf("expected authorised true, got %v", authorised)
		}
		if username, _ := context.Get(env.AuthContextKey.Username); username != "admin" {
			t.Fatalf("expected username 'admin', got %v", username)
		}
		if _, exists := context.Get(env.AuthContextKey.Error); exists {
			t.Fatalf("expected no authError to be set")
		}
	})

	t.Run("falls back to cookie token when header token is invalid", func(t *testing.T) {
		r := createTestRouter(t)
		token := createValidToken(t, r)

		w, context := runAuthContextMiddleware(r, "Bearer garbage-token", token)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}
		if authorised, _ := context.Get(env.AuthContextKey.Authorised); authorised != true {
			t.Fatalf("expected authorised true, got %v", authorised)
		}
		if username, _ := context.Get(env.AuthContextKey.Username); username != "admin" {
			t.Fatalf("expected username 'admin', got %v", username)
		}
	})

	t.Run("invalid tokens do not abort but report unauthorised with an error", func(t *testing.T) {
		r := createTestRouter(t)

		w, context := runAuthContextMiddleware(r, "Bearer garbage-token", "also-garbage")

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200 (never aborts), got %d", w.Code)
		}
		if context.IsAborted() {
			t.Fatalf("expected middleware to never abort")
		}
		if authorised, _ := context.Get(env.AuthContextKey.Authorised); authorised != false {
			t.Fatalf("expected authorised false, got %v", authorised)
		}
		authError, exists := context.Get(env.AuthContextKey.Error)
		if !exists || authError == "" {
			t.Fatalf("expected a non-empty authError to be set")
		}
	})

	t.Run("no header and no cookie reports unauthorised without aborting", func(t *testing.T) {
		r := createTestRouter(t)

		w, context := runAuthContextMiddleware(r, "", "")

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200 (never aborts), got %d", w.Code)
		}
		if authorised, _ := context.Get(env.AuthContextKey.Authorised); authorised != false {
			t.Fatalf("expected authorised false, got %v", authorised)
		}
	})
}
