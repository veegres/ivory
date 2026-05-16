package env

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseUrlPath(t *testing.T) {
	t.Run("should return / for empty string", func(t *testing.T) {
		result := parseUrlPath("")

		if result != "/" {
			t.Errorf("Expected '/', got '%s'", result)
		}
	})

	t.Run("should keep single slash", func(t *testing.T) {
		result := parseUrlPath("/")

		if result != "/" {
			t.Errorf("Expected '/', got '%s'", result)
		}
	})

	t.Run("should remove trailing slash", func(t *testing.T) {
		result := parseUrlPath("/api/")

		if result != "/api" {
			t.Errorf("Expected '/api', got '%s'", result)
		}
	})

	t.Run("should keep path without trailing slash", func(t *testing.T) {
		result := parseUrlPath("/api")

		if result != "/api" {
			t.Errorf("Expected '/api', got '%s'", result)
		}
	})

	t.Run("should handle nested paths", func(t *testing.T) {
		result := parseUrlPath("/api/v1/")

		if result != "/api/v1" {
			t.Errorf("Expected '/api/v1', got '%s'", result)
		}
	})

	t.Run("should panic for path without leading slash", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for path without leading slash")
			}
		}()

		parseUrlPath("api")
	})

	t.Run("should panic for path starting with non-slash", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for invalid path")
			}
		}()

		parseUrlPath("api/v1")
	})
}

func TestUpdateBaseUrlTag(t *testing.T) {
	t.Run("should update base href in index.html", func(t *testing.T) {
		// Create temporary directory
		tmpDir, err := os.MkdirTemp("", "env-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		// Create index.html with base href
		indexPath := filepath.Join(tmpDir, "index.html")
		originalContent := `<!DOCTYPE html>
<html>
<head>
  <base href="/" />
  <title>Test</title>
</head>
<body></body>
</html>`
		os.WriteFile(indexPath, []byte(originalContent), 0666)

		// Update base href
		updateBaseUrlTag(tmpDir, "/api")

		// Read and verify
		content, _ := os.ReadFile(indexPath)
		contentStr := string(content)

		if !contains(contentStr, `<base href="/api/" />`) {
			t.Errorf("Expected base href to be updated to '/api/', got:\n%s", contentStr)
		}
		if contains(contentStr, `<base href="/" />`) {
			t.Error("Expected old base href to be replaced")
		}
	})

	t.Run("should handle root path by making it empty", func(t *testing.T) {
		tmpDir, _ := os.MkdirTemp("", "env-test-*")
		defer os.RemoveAll(tmpDir)

		indexPath := filepath.Join(tmpDir, "index.html")
		originalContent := `<base href="/old/" />`
		os.WriteFile(indexPath, []byte(originalContent), 0666)

		// Root path should become empty string which results in "/"
		updateBaseUrlTag(tmpDir, "/")

		content, _ := os.ReadFile(indexPath)
		contentStr := string(content)

		if !contains(contentStr, `<base href="/" />`) {
			t.Errorf("Expected base href to be '/', got:\n%s", contentStr)
		}
	})

	t.Run("should handle complex paths", func(t *testing.T) {
		tmpDir, _ := os.MkdirTemp("", "env-test-*")
		defer os.RemoveAll(tmpDir)

		indexPath := filepath.Join(tmpDir, "index.html")
		originalContent := `<base href="/" />`
		os.WriteFile(indexPath, []byte(originalContent), 0666)

		updateBaseUrlTag(tmpDir, "/app/v1")

		content, _ := os.ReadFile(indexPath)
		contentStr := string(content)

		if !contains(contentStr, `<base href="/app/v1/" />`) {
			t.Errorf("Expected base href to be '/app/v1/', got:\n%s", contentStr)
		}
	})

	t.Run("should panic if index.html does not exist", func(t *testing.T) {
		tmpDir, _ := os.MkdirTemp("", "env-test-*")
		defer os.RemoveAll(tmpDir)

		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic when index.html doesn't exist")
			}
		}()

		updateBaseUrlTag(tmpDir, "/api")
	})
}

func TestNewAppEnv(t *testing.T) {
	// Save original env vars to restore after tests
	originalVars := map[string]string{
		"GIN_MODE":                 os.Getenv("GIN_MODE"),
		"IVORY_URL_ADDRESS":        os.Getenv("IVORY_URL_ADDRESS"),
		"IVORY_URL_PATH":           os.Getenv("IVORY_URL_PATH"),
		"IVORY_STATIC_FILES_PATH":  os.Getenv("IVORY_STATIC_FILES_PATH"),
		"IVORY_CERT_FILE_PATH":     os.Getenv("IVORY_CERT_FILE_PATH"),
		"IVORY_CERT_KEY_FILE_PATH": os.Getenv("IVORY_CERT_KEY_FILE_PATH"),
		"IVORY_VERSION_TAG":        os.Getenv("IVORY_VERSION_TAG"),
		"IVORY_VERSION_COMMIT":     os.Getenv("IVORY_VERSION_COMMIT"),
	}

	// Cleanup function to restore env vars
	cleanup := func() {
		for key, val := range originalVars {
			if val == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, val)
			}
		}
	}
	defer cleanup()

	t.Run("should use default values when no env vars set", func(t *testing.T) {
		// Clear all env vars
		for key := range originalVars {
			os.Unsetenv(key)
		}

		appEnv := NewAppEnv()

		if appEnv.Config.UrlAddress != ":8080" {
			t.Errorf("Expected default URL address ':8080', got '%s'", appEnv.Config.UrlAddress)
		}
		if appEnv.Config.UrlPath != "/" {
			t.Errorf("Expected default URL path '/', got '%s'", appEnv.Config.UrlPath)
		}
		if appEnv.Config.StaticFilesPath != "" {
			t.Errorf("Expected empty static files path, got '%s'", appEnv.Config.StaticFilesPath)
		}
		if appEnv.Config.TlsEnabled {
			t.Error("Expected TLS to be disabled by default")
		}
		if appEnv.Version.Tag != "v0.0.0" {
			t.Errorf("Expected default version tag 'v0.0.0', got '%s'", appEnv.Version.Tag)
		}
		if appEnv.Version.Commit != "000000000000" {
			t.Errorf("Expected default commit '000000000000', got '%s'", appEnv.Version.Commit)
		}
	})

	t.Run("should use custom URL address from env", func(t *testing.T) {
		os.Unsetenv("GIN_MODE")
		os.Setenv("IVORY_URL_ADDRESS", ":9090")

		appEnv := NewAppEnv()

		if appEnv.Config.UrlAddress != ":9090" {
			t.Errorf("Expected URL address ':9090', got '%s'", appEnv.Config.UrlAddress)
		}
	})

	t.Run("should use custom URL path from env", func(t *testing.T) {
		os.Setenv("IVORY_URL_PATH", "/api/v1/")

		appEnv := NewAppEnv()

		if appEnv.Config.UrlPath != "/api/v1" {
			t.Errorf("Expected URL path '/api/v1' (trailing slash removed), got '%s'", appEnv.Config.UrlPath)
		}
	})

	t.Run("should enable TLS when both cert paths are provided", func(t *testing.T) {
		os.Setenv("IVORY_CERT_FILE_PATH", "/path/to/cert.pem")
		os.Setenv("IVORY_CERT_KEY_FILE_PATH", "/path/to/key.pem")

		appEnv := NewAppEnv()

		if !appEnv.Config.TlsEnabled {
			t.Error("Expected TLS to be enabled when cert paths are provided")
		}
		if appEnv.Config.CertFilePath != "/path/to/cert.pem" {
			t.Errorf("Expected cert file path, got '%s'", appEnv.Config.CertFilePath)
		}
		if appEnv.Config.CertKeyFilePath != "/path/to/key.pem" {
			t.Errorf("Expected cert key file path, got '%s'", appEnv.Config.CertKeyFilePath)
		}
	})

	t.Run("should not enable TLS when only cert file path is provided", func(t *testing.T) {
		os.Setenv("IVORY_CERT_FILE_PATH", "/path/to/cert.pem")
		os.Unsetenv("IVORY_CERT_KEY_FILE_PATH")

		appEnv := NewAppEnv()

		if appEnv.Config.TlsEnabled {
			t.Error("Expected TLS to be disabled when only cert file path is provided")
		}
	})

	t.Run("should not enable TLS when only cert key path is provided", func(t *testing.T) {
		os.Unsetenv("IVORY_CERT_FILE_PATH")
		os.Setenv("IVORY_CERT_KEY_FILE_PATH", "/path/to/key.pem")

		appEnv := NewAppEnv()

		if appEnv.Config.TlsEnabled {
			t.Error("Expected TLS to be disabled when only cert key path is provided")
		}
	})

	t.Run("should use port 80 in release mode without TLS", func(t *testing.T) {
		os.Setenv("GIN_MODE", "release")
		os.Unsetenv("IVORY_URL_ADDRESS")
		os.Unsetenv("IVORY_CERT_FILE_PATH")
		os.Unsetenv("IVORY_CERT_KEY_FILE_PATH")

		appEnv := NewAppEnv()

		if appEnv.Config.UrlAddress != ":80" {
			t.Errorf("Expected URL address ':80' in release mode, got '%s'", appEnv.Config.UrlAddress)
		}
	})

	t.Run("should use port 443 in release mode with TLS", func(t *testing.T) {
		os.Setenv("GIN_MODE", "release")
		os.Unsetenv("IVORY_URL_ADDRESS")
		os.Setenv("IVORY_CERT_FILE_PATH", "/path/to/cert.pem")
		os.Setenv("IVORY_CERT_KEY_FILE_PATH", "/path/to/key.pem")

		appEnv := NewAppEnv()

		if appEnv.Config.UrlAddress != ":443" {
			t.Errorf("Expected URL address ':443' in release mode with TLS, got '%s'", appEnv.Config.UrlAddress)
		}
	})

	t.Run("should use custom version from env", func(t *testing.T) {
		os.Setenv("IVORY_VERSION_TAG", "v1.2.3")
		os.Setenv("IVORY_VERSION_COMMIT", "abc123def456")

		appEnv := NewAppEnv()

		if appEnv.Version.Tag != "v1.2.3" {
			t.Errorf("Expected version tag 'v1.2.3', got '%s'", appEnv.Version.Tag)
		}
		if appEnv.Version.Commit != "abc123def456" {
			t.Errorf("Expected commit 'abc123def456', got '%s'", appEnv.Version.Commit)
		}
	})

	t.Run("should generate version label correctly", func(t *testing.T) {
		os.Setenv("IVORY_VERSION_TAG", "v1.2.3")
		os.Setenv("IVORY_VERSION_COMMIT", "abc123def456789")

		appEnv := NewAppEnv()

		expectedLabel := "Ivory v1.2.3 (abc123d)"
		if appEnv.Version.Label != expectedLabel {
			t.Errorf("Expected version label '%s', got '%s'", expectedLabel, appEnv.Version.Label)
		}
	})

	t.Run("should truncate commit to 7 characters in label", func(t *testing.T) {
		os.Setenv("IVORY_VERSION_TAG", "v1.0.0")
		os.Setenv("IVORY_VERSION_COMMIT", "1234567890abcdef")

		appEnv := NewAppEnv()

		if !contains(appEnv.Version.Label, "(1234567)") {
			t.Errorf("Expected commit to be truncated to 7 chars in label, got '%s'", appEnv.Version.Label)
		}
	})
}

func TestVersion(t *testing.T) {
	t.Run("should have all required fields", func(t *testing.T) {
		version := Version{
			Tag:    "v1.0.0",
			Commit: "abc123",
			Label:  "Ivory v1.0.0 (abc123)",
		}

		if version.Tag != "v1.0.0" {
			t.Errorf("Expected Tag to be 'v1.0.0', got '%s'", version.Tag)
		}
		if version.Commit != "abc123" {
			t.Errorf("Expected Commit to be 'abc123', got '%s'", version.Commit)
		}
		if version.Label != "Ivory v1.0.0 (abc123)" {
			t.Errorf("Expected Label to be 'Ivory v1.0.0 (abc123)', got '%s'", version.Label)
		}
	})
}

func TestConfig(t *testing.T) {
	t.Run("should have all required fields", func(t *testing.T) {
		config := Config{
			UrlAddress:      ":8080",
			UrlPath:         "/api",
			StaticFilesPath: "/var/www",
			CertFilePath:    "/certs/cert.pem",
			CertKeyFilePath: "/certs/key.pem",
			TlsEnabled:      true,
		}

		if config.UrlAddress != ":8080" {
			t.Errorf("Expected UrlAddress ':8080', got '%s'", config.UrlAddress)
		}
		if config.UrlPath != "/api" {
			t.Errorf("Expected UrlPath '/api', got '%s'", config.UrlPath)
		}
		if !config.TlsEnabled {
			t.Error("Expected TlsEnabled to be true")
		}
	})
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
