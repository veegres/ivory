package files

import (
	"os"
	"path/filepath"
	"testing"
)

// Helper to create a test storage with cleanup
func createTestStorage(t *testing.T) *Storage {
	t.Helper()

	// Create temporary directory for test storage
	tmpDir, err := os.MkdirTemp("", "files-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Cleanup after test
	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	// Create storage in temporary location
	storage := &Storage{
		path: tmpDir,
		ext:  ".txt",
	}

	return storage
}

func TestNewStorage(t *testing.T) {
	t.Run("should create storage with directory", func(t *testing.T) {
		tmpDir, _ := os.MkdirTemp("", "storage-test-*")
		defer os.RemoveAll(tmpDir)

		// Use the actual NewStorage which creates data/name structure
		oldWd, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(oldWd)

		storage := NewStorage("test-storage", ".json")

		if storage == nil {
			t.Fatal("Expected storage to be created")
		}
		if storage.ext != ".json" {
			t.Errorf("Expected extension .json, got %s", storage.ext)
		}

		// Check directory was created
		if _, err := os.Stat("data/test-storage"); os.IsNotExist(err) {
			t.Error("Expected directory to be created")
		}
	})
}

func TestStorage_getPath(t *testing.T) {
	storage := createTestStorage(t)

	t.Run("should generate valid path", func(t *testing.T) {
		path, err := storage.getPath("testfile")

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		expectedPath := storage.path + "/testfile.txt"
		if path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, path)
		}
	})

	t.Run("should return error for empty name", func(t *testing.T) {
		_, err := storage.getPath("")

		if err == nil {
			t.Fatal("Expected error for empty name, got nil")
		}
		if err.Error() != "file name cannot be empty" {
			t.Errorf("Expected 'file name cannot be empty', got: %v", err)
		}
	})

	t.Run("should return error for name with slash", func(t *testing.T) {
		_, err := storage.getPath("test/file")

		if err == nil {
			t.Fatal("Expected error for name with slash, got nil")
		}
		if err.Error() != "file name contains invalid characters '/', '.'" {
			t.Errorf("Expected invalid characters error, got: %v", err)
		}
	})

	t.Run("should return error for name with dot", func(t *testing.T) {
		_, err := storage.getPath("test.file")

		if err == nil {
			t.Fatal("Expected error for name with dot, got nil")
		}
		if err.Error() != "file name contains invalid characters '/', '.'" {
			t.Errorf("Expected invalid characters error, got: %v", err)
		}
	})

	t.Run("should return error for path traversal attempt", func(t *testing.T) {
		_, err := storage.getPath("../etc/passwd")

		if err == nil {
			t.Fatal("Expected error for path traversal, got nil")
		}
	})
}

func TestStorage_CreateByName(t *testing.T) {
	storage := createTestStorage(t)

	t.Run("should create new file", func(t *testing.T) {
		path, err := storage.CreateByName("testfile")

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if path == "" {
			t.Error("Expected path to be returned")
		}

		// Check file exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Error("Expected file to be created")
		}
	})

	t.Run("should return error for invalid name", func(t *testing.T) {
		_, err := storage.CreateByName("")

		if err == nil {
			t.Fatal("Expected error for empty name, got nil")
		}
	})

	t.Run("should return error for name with invalid characters", func(t *testing.T) {
		_, err := storage.CreateByName("test/file")

		if err == nil {
			t.Fatal("Expected error for invalid name, got nil")
		}
	})
}

func TestStorage_ExistByName(t *testing.T) {
	storage := createTestStorage(t)

	t.Run("should return false for non-existent file", func(t *testing.T) {
		exists := storage.ExistByName("nonexistent")

		if exists {
			t.Error("Expected false for non-existent file")
		}
	})

	t.Run("should return true for existing file", func(t *testing.T) {
		storage.CreateByName("existingfile")

		exists := storage.ExistByName("existingfile")

		if !exists {
			t.Error("Expected true for existing file")
		}
	})

	t.Run("should return false for invalid name", func(t *testing.T) {
		exists := storage.ExistByName("")

		if exists {
			t.Error("Expected false for invalid name")
		}
	})

	t.Run("should return false for name with invalid characters", func(t *testing.T) {
		exists := storage.ExistByName("test/file")

		if exists {
			t.Error("Expected false for name with invalid characters")
		}
	})
}

func TestStorage_ReadByName(t *testing.T) {
	storage := createTestStorage(t)

	t.Run("should read file contents", func(t *testing.T) {
		// Create file with content
		path, _ := storage.CreateByName("readtest")
		testContent := []byte("test content")
		os.WriteFile(path, testContent, 0666)

		content, err := storage.ReadByName("readtest")

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if string(content) != string(testContent) {
			t.Errorf("Expected content '%s', got '%s'", testContent, content)
		}
	})

	t.Run("should return error for non-existent file", func(t *testing.T) {
		_, err := storage.ReadByName("nonexistent")

		if err == nil {
			t.Fatal("Expected error for non-existent file, got nil")
		}
	})

	t.Run("should return error for invalid name", func(t *testing.T) {
		_, err := storage.ReadByName("")

		if err == nil {
			t.Fatal("Expected error for empty name, got nil")
		}
	})

	t.Run("should read empty file", func(t *testing.T) {
		storage.CreateByName("emptyfile")

		content, err := storage.ReadByName("emptyfile")

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if len(content) != 0 {
			t.Errorf("Expected empty content, got %d bytes", len(content))
		}
	})
}

func TestStorage_OpenByName(t *testing.T) {
	storage := createTestStorage(t)

	t.Run("should open existing file", func(t *testing.T) {
		storage.CreateByName("opentest")

		file, err := storage.OpenByName("opentest")

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if file == nil {
			t.Fatal("Expected file handle, got nil")
		}
		defer file.Close()

		// Try to write to verify it's opened with RDWR
		_, writeErr := file.WriteString("test")
		if writeErr != nil {
			t.Errorf("Expected to write to opened file, got error: %v", writeErr)
		}
	})

	t.Run("should return error for non-existent file", func(t *testing.T) {
		_, err := storage.OpenByName("nonexistent")

		if err == nil {
			t.Fatal("Expected error for non-existent file, got nil")
		}
	})

	t.Run("should return error for invalid name", func(t *testing.T) {
		_, err := storage.OpenByName("")

		if err == nil {
			t.Fatal("Expected error for empty name, got nil")
		}
	})
}

func TestStorage_DeleteByName(t *testing.T) {
	storage := createTestStorage(t)

	t.Run("should delete existing file", func(t *testing.T) {
		storage.CreateByName("deletetest")

		err := storage.DeleteByName("deletetest")

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Verify file is deleted
		if storage.ExistByName("deletetest") {
			t.Error("Expected file to be deleted")
		}
	})

	t.Run("should return error for non-existent file", func(t *testing.T) {
		err := storage.DeleteByName("nonexistent")

		if err == nil {
			t.Fatal("Expected error for non-existent file, got nil")
		}
	})

	t.Run("should return error for invalid name", func(t *testing.T) {
		err := storage.DeleteByName("")

		if err == nil {
			t.Fatal("Expected error for empty name, got nil")
		}
	})

	t.Run("should return error for name with invalid characters", func(t *testing.T) {
		err := storage.DeleteByName("test/file")

		if err == nil {
			t.Fatal("Expected error for invalid name, got nil")
		}
	})
}

func TestStorage_ExistByPath(t *testing.T) {
	storage := createTestStorage(t)

	t.Run("should return false for non-existent path", func(t *testing.T) {
		exists := storage.ExistByPath("/nonexistent/path")

		if exists {
			t.Error("Expected false for non-existent path")
		}
	})

	t.Run("should return true for existing path", func(t *testing.T) {
		path, _ := storage.CreateByName("pathtest")

		exists := storage.ExistByPath(path)

		if !exists {
			t.Error("Expected true for existing path")
		}
	})

	t.Run("should work with directory path", func(t *testing.T) {
		exists := storage.ExistByPath(storage.path)

		if !exists {
			t.Error("Expected true for existing directory")
		}
	})
}

func TestStorage_DeleteAll(t *testing.T) {
	storage := createTestStorage(t)

	t.Run("should delete all files and recreate directory", func(t *testing.T) {
		// Create multiple files
		storage.CreateByName("file1")
		storage.CreateByName("file2")
		storage.CreateByName("file3")

		err := storage.DeleteAll()

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Verify directory still exists
		if _, err := os.Stat(storage.path); os.IsNotExist(err) {
			t.Error("Expected directory to exist after DeleteAll")
		}

		// Verify files are deleted
		if storage.ExistByName("file1") || storage.ExistByName("file2") || storage.ExistByName("file3") {
			t.Error("Expected all files to be deleted")
		}
	})

	t.Run("should allow creating files after DeleteAll", func(t *testing.T) {
		storage.DeleteAll()

		_, err := storage.CreateByName("newfile")

		if err != nil {
			t.Fatalf("Expected to create file after DeleteAll, got error: %v", err)
		}

		if !storage.ExistByName("newfile") {
			t.Error("Expected file to exist after creation")
		}
	})
}

func TestStorage_MultipleExtensions(t *testing.T) {
	t.Run("should work with different extensions", func(t *testing.T) {
		tmpDir, _ := os.MkdirTemp("", "files-ext-test-*")
		defer os.RemoveAll(tmpDir)

		storageJson := &Storage{path: tmpDir + "/json", ext: ".json"}
		storageTxt := &Storage{path: tmpDir + "/txt", ext: ".txt"}

		os.MkdirAll(storageJson.path, os.ModePerm)
		os.MkdirAll(storageTxt.path, os.ModePerm)

		// Create files with different extensions
		pathJson, _ := storageJson.CreateByName("data")
		pathTxt, _ := storageTxt.CreateByName("data")

		if pathJson == pathTxt {
			t.Error("Expected different paths for different extensions")
		}

		if filepath.Ext(pathJson) != ".json" {
			t.Errorf("Expected .json extension, got %s", filepath.Ext(pathJson))
		}
		if filepath.Ext(pathTxt) != ".txt" {
			t.Errorf("Expected .txt extension, got %s", filepath.Ext(pathTxt))
		}
	})
}
