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
	"strings"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
)

func createTestSecretService(t *testing.T) *secret.Service {
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
	return secretService
}

func createTestRouter(t *testing.T) *Router {
	t.Helper()
	users := newFakeUsers()
	authService := NewService(token.NewService(createTestSecretService(t)), createTestBasicProvider(t, users), ldap.NewProvider(), oidc.NewProvider(), users)
	return NewRouter(authService, "/", false)
}

// createDisabledTestRouter builds a router with no provider configured, which is
// how a deployment with authentication turned off looks to auth.Service:
// GetSupportedTypes() is empty and every resolveAuth call reports ErrAuthDisabled.
func createDisabledTestRouter(t *testing.T) *Router {
	t.Helper()
	users := newFakeUsers()
	authService := NewService(token.NewService(createTestSecretService(t)), basic.NewProvider(users), ldap.NewProvider(), oidc.NewProvider(), users)
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

// runMiddleware runs handlers in sequence over one request/context, stopping
// early if one of them aborts - exactly how gin runs a route's middleware chain.
func runMiddleware(authHeader string, tokenCookie string, handlers ...gin.HandlerFunc) (*httptest.ResponseRecorder, *gin.Context) {
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

	for _, handler := range handlers {
		handler(context)
		if context.IsAborted() {
			break
		}
	}

	return w, context
}

// runRejectMiddleware and runAllowMiddleware chain ValidateWithContextMiddleware
// in front of the middleware under test, exactly as engine_http.go wires
// yesSecret (ValidateWithContextMiddleware) in front of yesAuth (RejectMiddleware)
// and noAuth (AllowMiddleware). RejectMiddleware and AllowMiddleware only read
// context keys they never populate themselves, so testing either in isolation
// from ValidateWithContextMiddleware does not exercise the real request path.
func runRejectMiddleware(r *Router, authHeader string, tokenCookie string) (*httptest.ResponseRecorder, *gin.Context) {
	return runMiddleware(authHeader, tokenCookie, r.ValidateWithContextMiddleware(), r.RejectMiddleware())
}

func runAllowMiddleware(r *Router, authHeader string, tokenCookie string) (*httptest.ResponseRecorder, *gin.Context) {
	return runMiddleware(authHeader, tokenCookie, r.ValidateWithContextMiddleware(), r.AllowMiddleware())
}

func runAuthContextMiddleware(r *Router, authHeader string, tokenCookie string) (*httptest.ResponseRecorder, *gin.Context) {
	return runMiddleware(authHeader, tokenCookie, r.ValidateWithContextMiddleware())
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

// TestRejectMiddleware covers the gate in front of every route that requires a
// signed-in request (yesAuth in engine_http.go). It always runs
// ValidateWithContextMiddleware first, matching how the router actually chains
// the two - RejectMiddleware only reads context keys that middleware sets.
func TestRejectMiddleware(t *testing.T) {
	t.Run("valid header token is accepted", func(t *testing.T) {
		r := createTestRouter(t)
		token := createValidToken(t, r)

		w, context := runRejectMiddleware(r, "Bearer "+token, "")

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

		w, context := runRejectMiddleware(r, "", token)

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

		w, context := runRejectMiddleware(r, "Bearer garbage-token", token)

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

		w, context := runRejectMiddleware(r, "Bearer garbage-token", "also-garbage")

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected status 401, got %d", w.Code)
		}
		if !context.IsAborted() {
			t.Fatalf("expected middleware to abort")
		}
		if w.Header().Get("WWW-Authenticate") == "" {
			t.Fatalf("expected a WWW-Authenticate header on a rejected request")
		}
	})

	t.Run("rejects request with no header and no cookie", func(t *testing.T) {
		r := createTestRouter(t)

		w, _ := runRejectMiddleware(r, "", "")

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected status 401, got %d", w.Code)
		}
	})

	// Regression test: RejectMiddleware used to skip straight to reading
	// Username/Type without checking whether auth is enabled at all, so a
	// deployment with authentication switched off had every protected route
	// (cluster, node, user, ...) answer 401 "username cannot be empty" - see
	// ValidateWithContextMiddleware, which never sets Username when it reports
	// ErrAuthDisabled.
	t.Run("allows every request through when auth is disabled", func(t *testing.T) {
		r := createDisabledTestRouter(t)

		w, context := runRejectMiddleware(r, "", "")

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200 when auth is disabled, got %d, body: %s", w.Code, w.Body.String())
		}
		if context.IsAborted() {
			t.Fatalf("expected middleware to not abort when auth is disabled")
		}
	})

	t.Run("authorised but empty username is rejected", func(t *testing.T) {
		w, context := runMiddleware("", "", func(context *gin.Context) {
			context.Set(config.AuthContextKey.Enabled, true)
			context.Set(config.AuthContextKey.Authorised, true)
		})
		createTestRouter(t).RejectMiddleware()(context)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected status 401, got %d", w.Code)
		}
		if !context.IsAborted() {
			t.Fatalf("expected middleware to abort")
		}
	})

	t.Run("authorised but unrecognised auth type is rejected", func(t *testing.T) {
		w, context := runMiddleware("", "", func(context *gin.Context) {
			context.Set(config.AuthContextKey.Enabled, true)
			context.Set(config.AuthContextKey.Authorised, true)
			context.Set(config.AuthContextKey.Username, "admin")
			context.Set(config.AuthContextKey.Type, "smoke-signals")
		})
		createTestRouter(t).RejectMiddleware()(context)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected status 401, got %d", w.Code)
		}
		if !context.IsAborted() {
			t.Fatalf("expected middleware to abort")
		}
	})

	// Regression test: an unauthorised request with no recorded Error used to
	// dereference a nil error and panic instead of answering 401. This can
	// only happen if RejectMiddleware runs without ValidateWithContextMiddleware
	// ahead of it - which is exactly what the pre-fix test suite did - but the
	// guard is cheap and turns a server crash into a clean rejection.
	t.Run("does not panic when unauthorised with no recorded error", func(t *testing.T) {
		w, context := runMiddleware("", "", func(context *gin.Context) {
			context.Set(config.AuthContextKey.Enabled, true)
			context.Set(config.AuthContextKey.Authorised, false)
		})
		createTestRouter(t).RejectMiddleware()(context)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected status 401, got %d", w.Code)
		}
		if !context.IsAborted() {
			t.Fatalf("expected middleware to abort")
		}
	})
}

// TestAllowMiddleware covers the gate in front of every route a signed-in
// request must not reach (noAuth in engine_http.go: login, config and
// registration). It always runs ValidateWithContextMiddleware first for the
// same reason TestRejectMiddleware does.
func TestAllowMiddleware(t *testing.T) {
	t.Run("lets an unauthenticated request through", func(t *testing.T) {
		r := createTestRouter(t)

		w, context := runAllowMiddleware(r, "", "")

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}
		if context.IsAborted() {
			t.Fatalf("expected middleware to not abort")
		}
	})

	t.Run("aborts a request that already carries a valid session", func(t *testing.T) {
		r := createTestRouter(t)
		token := createValidToken(t, r)

		w, context := runAllowMiddleware(r, "Bearer "+token, "")

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", w.Code)
		}
		if !context.IsAborted() {
			t.Fatalf("expected middleware to abort")
		}
		if !strings.Contains(w.Body.String(), ErrAlreadyAuthenticated.Error()) {
			t.Fatalf("expected the error body to mention %q, got %q", ErrAlreadyAuthenticated.Error(), w.Body.String())
		}
	})

	t.Run("lets a request through when auth is disabled, even with a stale token", func(t *testing.T) {
		r := createDisabledTestRouter(t)

		w, context := runAllowMiddleware(r, "Bearer garbage-token", "")

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200 when auth is disabled, got %d", w.Code)
		}
		if context.IsAborted() {
			t.Fatalf("expected middleware to not abort when auth is disabled")
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

	t.Run("reports auth disabled and never sets username when no provider is configured", func(t *testing.T) {
		r := createDisabledTestRouter(t)

		w, context := runAuthContextMiddleware(r, "", "")

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}
		if enabled, _ := context.Get(config.AuthContextKey.Enabled); enabled != false {
			t.Fatalf("expected auth enabled to be false, got %v", enabled)
		}
		if _, exists := context.Get(config.AuthContextKey.Username); exists {
			t.Fatalf("expected no username to be set when auth is disabled")
		}
	})
}
