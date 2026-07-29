package secret

import (
	"ivory/clients/storage"
	"ivory/core/service/encryption"
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

	service := NewService(
		NewRepository(storage.NewDbBucket[string](db, "Secret")),
		encryption.NewService(),
	)
	return NewRouter(service)
}

func runMiddleware(handler gin.HandlerFunc) (*httptest.ResponseRecorder, *gin.Context) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(w)
	context.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	handler(context)

	return w, context
}

func TestExistMiddleware(t *testing.T) {
	t.Run("aborts when the secret is empty", func(t *testing.T) {
		r := createTestRouter(t)

		w, context := runMiddleware(r.ExistMiddleware())

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected status 403, got %d", w.Code)
		}
		if !context.IsAborted() {
			t.Fatalf("expected middleware to abort")
		}
	})

	t.Run("passes through once the secret is set", func(t *testing.T) {
		r := createTestRouter(t)
		if err := r.secretService.SetDefault(); err != nil {
			t.Fatalf("failed to set default secret: %v", err)
		}

		w, context := runMiddleware(r.ExistMiddleware())

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}
		if context.IsAborted() {
			t.Fatalf("expected middleware to not abort")
		}
	})
}

func TestEmptyMiddleware(t *testing.T) {
	t.Run("passes through while the secret is empty", func(t *testing.T) {
		r := createTestRouter(t)

		w, context := runMiddleware(r.EmptyMiddleware())

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}
		if context.IsAborted() {
			t.Fatalf("expected middleware to not abort")
		}
	})

	t.Run("aborts once the secret is set", func(t *testing.T) {
		r := createTestRouter(t)
		if err := r.secretService.SetDefault(); err != nil {
			t.Fatalf("failed to set default secret: %v", err)
		}

		w, context := runMiddleware(r.EmptyMiddleware())

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected status 403, got %d", w.Code)
		}
		if !context.IsAborted() {
			t.Fatalf("expected middleware to abort")
		}
	})
}
