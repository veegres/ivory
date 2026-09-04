package permission

import (
	"ivory/clients/storage"
	"ivory/core/config"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"slices"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
)

func createTestPermissionRouter(t *testing.T) *Router {
	t.Helper()

	db, errOpen := bolt.Open(filepath.Join(t.TempDir(), "test.db"), 0600, nil)
	if errOpen != nil {
		t.Fatalf("failed to open test database: %v", errOpen)
	}
	t.Cleanup(func() {
		db.Close()
	})

	service := NewService(NewRepository(storage.NewDbBucket[PermissionMap](db, "Permission")))
	return NewRouter(service)
}

func runValidateMiddleware(r *Router, authEnabled bool, superuser bool, username string) (*httptest.ResponseRecorder, *gin.Context) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(w)
	context.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	context.Set(config.AuthContextKey.Enabled, authEnabled)
	context.Set(config.AuthContextKey.Superuser, superuser)
	context.Set(config.AuthContextKey.Username, username)

	r.ValidateMiddleware()(context)

	return w, context
}

func TestValidateMiddleware(t *testing.T) {
	t.Run("auth disabled grants every permission without a lookup", func(t *testing.T) {
		r := createTestPermissionRouter(t)

		w, context := runValidateMiddleware(r, false, false, "")

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}
		if context.IsAborted() {
			t.Fatalf("expected middleware to not abort")
		}
		perms, ok := context.Get("permissions")
		if !ok {
			t.Fatalf("expected permissions to be set in context")
		}
		permMap := perms.(PermissionMap)
		withheld := config.Withheld(false)
		for _, feature := range config.All {
			if slices.Contains(withheld, feature) {
				continue
			}
			if permMap[feature] != GRANTED {
				t.Fatalf("expected all features granted when auth is disabled, got %v for %s", permMap[feature], feature)
			}
		}
	})

	// NOTE: without authentication there is nobody to register and nobody to
	// grant anything to, so the features that only mean anything under a session
	// are not handed out at all
	t.Run("auth disabled withholds the features that need a session", func(t *testing.T) {
		r := createTestPermissionRouter(t)

		_, context := runValidateMiddleware(r, false, false, "")

		perms, _ := context.Get("permissions")
		permMap := perms.(PermissionMap)
		for _, feature := range config.Withheld(false) {
			if _, ok := permMap[feature]; ok {
				t.Fatalf("expected %s to be withheld when auth is disabled", feature)
			}
		}
	})

	t.Run("known user resolves their stored permissions", func(t *testing.T) {
		r := createTestPermissionRouter(t)
		if _, err := r.permissionService.CreateUserPermissions("alice", NOT_PERMITTED); err != nil {
			t.Fatalf("failed to seed alice: %v", err)
		}

		w, context := runValidateMiddleware(r, true, false, "alice")

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}
		if context.IsAborted() {
			t.Fatalf("expected middleware to not abort")
		}
		if _, ok := context.Get("permissions"); !ok {
			t.Fatalf("expected permissions to be set in context")
		}
	})

	// NOTE: a superuser holds every permission and cannot be talked out of it,
	// so what a record says about them is never asked - which is also why a
	// superuser with no record at all still gets through
	t.Run("superuser is granted everything whatever the record says", func(t *testing.T) {
		r := createTestPermissionRouter(t)
		if _, err := r.permissionService.CreateUserPermissions("root", NOT_PERMITTED); err != nil {
			t.Fatalf("failed to seed root: %v", err)
		}

		w, context := runValidateMiddleware(r, true, true, "root")

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}
		perms, ok := context.Get("permissions")
		if !ok {
			t.Fatalf("expected permissions to be set in context")
		}
		permMap := perms.(PermissionMap)
		for _, feature := range config.All {
			if permMap[feature] != GRANTED {
				t.Fatalf("expected all features granted for a superuser, got %v for %s", permMap[feature], feature)
			}
		}
	})

	t.Run("unknown user aborts with forbidden", func(t *testing.T) {
		r := createTestPermissionRouter(t)

		w, context := runValidateMiddleware(r, true, false, "unknown")

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected status 403, got %d", w.Code)
		}
		if !context.IsAborted() {
			t.Fatalf("expected middleware to abort")
		}
	})
}

func runValidateMethodMiddleware(feature config.Feature, seedPermissions *PermissionMap) (*httptest.ResponseRecorder, *gin.Context) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(w)
	context.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	if seedPermissions != nil {
		context.Set("permissions", *seedPermissions)
	}

	r := &Router{}
	r.ValidateMethodMiddleware(feature)(context)

	return w, context
}

func TestValidateMethodMiddleware(t *testing.T) {
	t.Run("no permissions in context aborts with forbidden", func(t *testing.T) {
		w, context := runValidateMethodMiddleware(config.ViewClusterList, nil)

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected status 403, got %d", w.Code)
		}
		if !context.IsAborted() {
			t.Fatalf("expected middleware to abort")
		}
	})

	t.Run("granted feature passes through", func(t *testing.T) {
		perms := PermissionMap{config.ViewClusterList: GRANTED}
		w, context := runValidateMethodMiddleware(config.ViewClusterList, &perms)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}
		if context.IsAborted() {
			t.Fatalf("expected middleware to not abort")
		}
	})

	t.Run("not permitted feature aborts with forbidden", func(t *testing.T) {
		perms := PermissionMap{config.ViewClusterList: NOT_PERMITTED}
		w, context := runValidateMethodMiddleware(config.ViewClusterList, &perms)

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected status 403, got %d", w.Code)
		}
		if !context.IsAborted() {
			t.Fatalf("expected middleware to abort")
		}
	})

	t.Run("pending feature aborts with forbidden", func(t *testing.T) {
		perms := PermissionMap{config.ViewClusterList: PENDING}
		w, context := runValidateMethodMiddleware(config.ViewClusterList, &perms)

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected status 403, got %d", w.Code)
		}
		if !context.IsAborted() {
			t.Fatalf("expected middleware to abort")
		}
	})
}
