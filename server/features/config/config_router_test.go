package config

import (
	"ivory/clients/auth/basic"
	"ivory/clients/auth/ldap"
	"ivory/clients/auth/oidc"
	"ivory/clients/storage"
	"ivory/core/service/encryption"
	"ivory/core/service/secret"
	"ivory/core/service/token"
	"ivory/features/permission"
	"ivory/features/user"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
)

func createTestConfigRouter(t *testing.T) *Router {
	t.Helper()

	tmpDir, errDir := os.MkdirTemp("", "config-router-test-*")
	if errDir != nil {
		t.Fatalf("failed to create temp dir: %v", errDir)
	}
	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	db, errOpen := bolt.Open(filepath.Join(tmpDir, "test.db"), 0600, nil)
	if errOpen != nil {
		t.Fatalf("failed to open test database: %v", errOpen)
	}
	t.Cleanup(func() {
		db.Close()
	})

	oldWd, errWd := os.Getwd()
	if errWd != nil {
		t.Fatalf("failed to get working dir: %v", errWd)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	t.Cleanup(func() { os.Chdir(oldWd) })

	secretService := secret.NewService(
		secret.NewRepository(storage.NewDbBucket[string](db, "Secret")),
		encryption.NewService(),
	)
	if err := secretService.SetDefault(); err != nil {
		t.Fatalf("failed to set default secret: %v", err)
	}

	userService := user.NewService(
		user.NewRepository(storage.NewDbBucket[user.User](db, "User")),
		encryption.NewService(),
		secretService,
		permission.NewService(permission.NewRepository(storage.NewDbBucket[permission.PermissionMap](db, "Permission"))),
		token.NewService(secretService),
	)

	service := NewService(
		storage.NewFileStorage("config", ""),
		encryption.NewService(),
		secretService,
		nil,
		userService,
		basic.NewProvider(nopVerifier{}),
		ldap.NewProvider(),
		oidc.NewProvider(),
	)
	return NewRouter(service)
}

func runInitialiseMiddleware(r *Router) (*httptest.ResponseRecorder, *gin.Context) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(w)
	context.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	r.InitialiseMiddleware()(context)

	return w, context
}

func TestInitialiseMiddleware(t *testing.T) {
	t.Run("aborts when the app is not configured", func(t *testing.T) {
		r := createTestConfigRouter(t)

		w, context := runInitialiseMiddleware(r)

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected status 403, got %d", w.Code)
		}
		if !context.IsAborted() {
			t.Fatalf("expected middleware to abort")
		}
	})

	t.Run("passes through once the app is configured", func(t *testing.T) {
		r := createTestConfigRouter(t)
		if err := r.service.SetAppConfig(NewAppConfig{AppConfig: AppConfig{Company: "Acme"}}); err != nil {
			t.Fatalf("failed to set app config: %v", err)
		}

		w, context := runInitialiseMiddleware(r)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}
		if context.IsAborted() {
			t.Fatalf("expected middleware to not abort")
		}
	})
}
