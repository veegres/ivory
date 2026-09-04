package auth

import (
	"ivory/clients/auth/basic"
	"ivory/clients/auth/ldap"
	"ivory/clients/auth/oidc"
	"ivory/clients/storage"
	"ivory/core/config"
	"ivory/core/service/encryption"
	"ivory/core/service/secret"
	"ivory/core/service/token"
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

	users := newFakeUsers()
	authService := NewService(token.NewService(secretService), createTestBasicProvider(t, users), ldap.NewProvider(), oidc.NewProvider(), users)
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

func runSessionMiddleware(r *Router, sessionCookie string) (*httptest.ResponseRecorder, *gin.Context) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(w)
	context.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	if sessionCookie != "" {
		context.Request.AddCookie(&http.Cookie{Name: "session", Value: sessionCookie})
	}

	r.SessionMiddleware()(context)

	return w, context
}

func TestSessionMiddleware(t *testing.T) {
	t.Run("generates and cookies a new session when none is present", func(t *testing.T) {
		r := createTestRouter(t)

		w, context := runSessionMiddleware(r, "")

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}
		if context.IsAborted() {
			t.Fatalf("expected middleware to not abort")
		}
		session, _ := context.Get(config.AuthContextKey.Session)
		sessionStr, ok := session.(string)
		if !ok || sessionStr == "" {
			t.Fatalf("expected a non-empty session to be set in context, got %v", session)
		}

		var cookieSet bool
		for _, c := range w.Result().Cookies() {
			if c.Name == "session" && c.Value == sessionStr {
				cookieSet = true
			}
		}
		if !cookieSet {
			t.Fatalf("expected a session cookie matching %q to be set, got %v", sessionStr, w.Result().Cookies())
		}
	})

	t.Run("reuses an existing session cookie without setting a new one", func(t *testing.T) {
		r := createTestRouter(t)

		w, context := runSessionMiddleware(r, "existing-session-id")

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}
		session, _ := context.Get(config.AuthContextKey.Session)
		if session != "existing-session-id" {
			t.Fatalf("expected the existing session to be reused, got %v", session)
		}
		if len(w.Result().Cookies()) != 0 {
			t.Fatalf("expected no new cookie to be set, got %v", w.Result().Cookies())
		}
	})
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
		if username, _ := context.Get(config.AuthContextKey.Username); username != "admin" {
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
		if username, _ := context.Get(config.AuthContextKey.Username); username != "admin" {
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
		if username, _ := context.Get(config.AuthContextKey.Username); username != "admin" {
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
		if authorised, _ := context.Get(config.AuthContextKey.Authorised); authorised != true {
			t.Fatalf("expected authorised true, got %v", authorised)
		}
		if username, _ := context.Get(config.AuthContextKey.Username); username != "admin" {
			t.Fatalf("expected username 'admin', got %v", username)
		}
		if _, exists := context.Get(config.AuthContextKey.Error); exists {
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
		if authorised, _ := context.Get(config.AuthContextKey.Authorised); authorised != true {
			t.Fatalf("expected authorised true, got %v", authorised)
		}
		if username, _ := context.Get(config.AuthContextKey.Username); username != "admin" {
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
		if authorised, _ := context.Get(config.AuthContextKey.Authorised); authorised != false {
			t.Fatalf("expected authorised false, got %v", authorised)
		}
		authError, exists := context.Get(config.AuthContextKey.Error)
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
		if authorised, _ := context.Get(config.AuthContextKey.Authorised); authorised != false {
			t.Fatalf("expected authorised false, got %v", authorised)
		}
	})
}
