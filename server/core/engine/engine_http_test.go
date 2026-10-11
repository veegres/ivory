package engine

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// serveSpa builds just the static/SPA half of the http engine, which is the
// only part of engine_http.go with branching of its own - the rest is a
// declarative route table.
func serveSpa(t *testing.T, urlPath string) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	static := t.TempDir()
	if err := os.WriteFile(filepath.Join(static, "index.html"), []byte("<!doctype html>spa"), 0600); err != nil {
		t.Fatalf("failed to write index.html: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(static, "assets"), 0700); err != nil {
		t.Fatalf("failed to create assets dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(static, "assets", "app.js"), []byte("console.log(1)"), 0600); err != nil {
		t.Fatalf("failed to write app.js: %v", err)
	}

	engine := gin.New()
	spaServe(engine, urlPath, "/api", static)
	return engine
}

func get(t *testing.T, engine *gin.Engine, path string) *httptest.ResponseRecorder {
	t.Helper()
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
	return recorder
}

// TestSpaServeAtRoot covers the cases spaServe's own comments enumerate for
// urlPath "/". The rule that matters is the first one: an unmatched api path
// must 404 rather than be answered with the SPA, or a client gets index.html
// where it expected json.
func TestSpaServeAtRoot(t *testing.T) {
	engine := serveSpa(t, "/")

	tests := []struct {
		name   string
		path   string
		status int
		body   string
	}{
		{name: "a deep link falls back to the spa", path: "/clusters", status: http.StatusOK, body: "spa"},
		{name: "a registration token survives a full page load", path: "/user/123", status: http.StatusOK, body: "spa"},
		{name: "an unmatched api path is not answered by the spa", path: "/api/nope", status: http.StatusNotFound},
		{name: "the bare api path is not either", path: "/api", status: http.StatusNotFound},
		{name: "a static asset is served from disk", path: "/assets/app.js", status: http.StatusOK, body: "console.log(1)"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := get(t, engine, test.path)
			if recorder.Code != test.status {
				t.Errorf("status = %d, want %d", recorder.Code, test.status)
			}
			if test.body != "" && !strings.Contains(recorder.Body.String(), test.body) {
				t.Errorf("body = %q, want it to contain %q", recorder.Body.String(), test.body)
			}
		})
	}
}

// TestSpaServeUnderSubPath covers the reverse-proxy case: only urls belonging
// to the SPA fall back to index.html, so a request outside the sub-path is a
// 404 rather than somebody else's app being handed Ivory's page.
func TestSpaServeUnderSubPath(t *testing.T) {
	engine := serveSpa(t, "/ivory")

	tests := []struct {
		name   string
		path   string
		status int
		body   string
	}{
		{name: "a deep link inside the sub-path falls back to the spa", path: "/ivory/clusters", status: http.StatusOK, body: "spa"},
		{name: "and so does a registration link", path: "/ivory/user/123", status: http.StatusOK, body: "spa"},
		{name: "the sub-path itself does too", path: "/ivory", status: http.StatusOK, body: "spa"},
		{name: "a path outside the sub-path is not ours", path: "/clusters", status: http.StatusNotFound},
		{name: "nor is a bare registration link", path: "/user/123", status: http.StatusNotFound},
		{name: "the api under the sub-path still 404s rather than serving the spa", path: "/ivory/api/nope", status: http.StatusNotFound},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := get(t, engine, test.path)
			if recorder.Code != test.status {
				t.Errorf("status = %d, want %d", recorder.Code, test.status)
			}
			if test.body != "" && !strings.Contains(recorder.Body.String(), test.body) {
				t.Errorf("body = %q, want it to contain %q", recorder.Body.String(), test.body)
			}
		})
	}
}

func TestSpaServeRejectsPathsWithoutLeadingSlash(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("expected a panic for a url path with no leading slash")
		}
	}()
	spaServe(gin.New(), "ivory", "/api", t.TempDir())
}

func TestPong(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.GET("/ping", pong)

	recorder := get(t, engine, "/ping")
	if recorder.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), "pong") {
		t.Errorf("body = %q, want it to contain pong", recorder.Body.String())
	}
}
