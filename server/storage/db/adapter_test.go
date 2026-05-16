package db

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/boltdb/bolt"
)

// Test model for testing generic Bucket operations
type TestModel struct {
	ID   string
	Name string
	Age  int
}

// Helper to create an in-memory BoltDB for testing
func createTestDB(t *testing.T) *bolt.DB {
	t.Helper()

	// Create temporary directory for test database
	tmpDir, err := os.MkdirTemp("", "bolt-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Cleanup after test
	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	// Open BoltDB in temporary location
	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func TestBucket_Create(t *testing.T) {
	db := createTestDB(t)
	bucket := NewBucket[TestModel](db, "test-bucket")

	t.Run("should create new element", func(t *testing.T) {
		model := TestModel{ID: "1", Name: "Alice", Age: 30}

		result, err := bucket.Create("test-key", model)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if result.ID != model.ID {
			t.Errorf("Expected ID %s, got %s", model.ID, result.ID)
		}
	})

	t.Run("should return error for empty key", func(t *testing.T) {
		model := TestModel{ID: "2", Name: "Bob", Age: 25}

		_, err := bucket.Create("", model)

		if err == nil {
			t.Fatal("Expected error for empty key, got nil")
		}
		if err.Error() != "element identifier cannot be empty" {
			t.Errorf("Expected 'element identifier cannot be empty', got: %v", err)
		}
	})

	t.Run("should return error for duplicate key", func(t *testing.T) {
		model1 := TestModel{ID: "3", Name: "Charlie", Age: 35}
		model2 := TestModel{ID: "4", Name: "David", Age: 40}

		_, err := bucket.Create("duplicate-key", model1)
		if err != nil {
			t.Fatalf("First create failed: %v", err)
		}

		_, err = bucket.Create("duplicate-key", model2)

		if err == nil {
			t.Fatal("Expected error for duplicate key, got nil")
		}
		if err.Error() != "such an element already exists" {
			t.Errorf("Expected 'such an element already exists', got: %v", err)
		}
	})
}

func TestBucket_Get(t *testing.T) {
	db := createTestDB(t)
	bucket := NewBucket[TestModel](db, "test-bucket")

	t.Run("should get existing element", func(t *testing.T) {
		model := TestModel{ID: "1", Name: "Alice", Age: 30}
		bucket.Create("key1", model)

		result, err := bucket.Get("key1")

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if result.ID != model.ID || result.Name != model.Name || result.Age != model.Age {
			t.Errorf("Expected %+v, got %+v", model, result)
		}
	})

	t.Run("should return error for non-existent element", func(t *testing.T) {
		_, err := bucket.Get("non-existent-key")

		if err == nil {
			t.Fatal("Expected error for non-existent key, got nil")
		}
		if err.Error() != "element doesn't exist" {
			t.Errorf("Expected 'element doesn't exist', got: %v", err)
		}
	})
}

func TestBucket_Update(t *testing.T) {
	db := createTestDB(t)
	bucket := NewBucket[TestModel](db, "test-bucket")

	t.Run("should update existing element", func(t *testing.T) {
		model := TestModel{ID: "1", Name: "Alice", Age: 30}
		bucket.Create("key1", model)

		updatedModel := TestModel{ID: "1", Name: "Alice Updated", Age: 31}
		err := bucket.Update("key1", updatedModel)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		result, _ := bucket.Get("key1")
		if result.Name != "Alice Updated" || result.Age != 31 {
			t.Errorf("Expected updated values, got %+v", result)
		}
	})

	t.Run("should create element if not exists", func(t *testing.T) {
		model := TestModel{ID: "2", Name: "Bob", Age: 25}

		err := bucket.Update("new-key", model)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		result, _ := bucket.Get("new-key")
		if result.Name != "Bob" {
			t.Errorf("Expected Bob, got %s", result.Name)
		}
	})

	t.Run("should return error for empty key", func(t *testing.T) {
		model := TestModel{ID: "3", Name: "Charlie", Age: 35}

		err := bucket.Update("", model)

		if err == nil {
			t.Fatal("Expected error for empty key, got nil")
		}
		if err.Error() != "element identifier cannot be empty" {
			t.Errorf("Expected 'element identifier cannot be empty', got: %v", err)
		}
	})
}

func TestBucket_Delete(t *testing.T) {
	db := createTestDB(t)
	bucket := NewBucket[TestModel](db, "test-bucket")

	t.Run("should delete existing element", func(t *testing.T) {
		model := TestModel{ID: "1", Name: "Alice", Age: 30}
		bucket.Create("key1", model)

		err := bucket.Delete("key1")

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		_, err = bucket.Get("key1")
		if err == nil {
			t.Fatal("Expected error after delete, got nil")
		}
	})

	t.Run("should not return error for non-existent element", func(t *testing.T) {
		err := bucket.Delete("non-existent-key")

		if err != nil {
			t.Errorf("Expected no error for deleting non-existent key, got: %v", err)
		}
	})
}

func TestBucket_GetList(t *testing.T) {
	db := createTestDB(t)
	bucket := NewBucket[TestModel](db, "test-bucket")

	// Setup test data
	bucket.Create("key1", TestModel{ID: "1", Name: "Alice", Age: 30})
	bucket.Create("key2", TestModel{ID: "2", Name: "Bob", Age: 25})
	bucket.Create("key3", TestModel{ID: "3", Name: "Charlie", Age: 35})

	t.Run("should get all elements without filter", func(t *testing.T) {
		result, err := bucket.GetList(nil, nil)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if len(result) != 3 {
			t.Errorf("Expected 3 elements, got %d", len(result))
		}
	})

	t.Run("should filter elements by condition", func(t *testing.T) {
		// Filter: Age > 28
		result, err := bucket.GetList(func(el TestModel) bool {
			return el.Age > 28
		}, nil)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("Expected 2 elements (Age > 28), got %d", len(result))
		}
		for _, el := range result {
			if el.Age <= 28 {
				t.Errorf("Expected Age > 28, got %d", el.Age)
			}
		}
	})

	t.Run("should sort elements", func(t *testing.T) {
		// Sort by Age ascending
		result, err := bucket.GetList(nil, func(list []TestModel, i, j int) bool {
			return list[i].Age < list[j].Age
		})

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if len(result) != 3 {
			t.Errorf("Expected 3 elements, got %d", len(result))
		}
		// Check ascending order
		if result[0].Age != 25 || result[1].Age != 30 || result[2].Age != 35 {
			t.Errorf("Expected sorted by Age [25, 30, 35], got [%d, %d, %d]",
				result[0].Age, result[1].Age, result[2].Age)
		}
	})

	t.Run("should filter and sort elements", func(t *testing.T) {
		// Filter: Age > 26, Sort by Age descending
		result, err := bucket.GetList(
			func(el TestModel) bool {
				return el.Age > 26
			},
			func(list []TestModel, i, j int) bool {
				return list[i].Age > list[j].Age
			},
		)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("Expected 2 elements (Age > 26), got %d", len(result))
		}
		// Check descending order
		if result[0].Age < result[1].Age {
			t.Errorf("Expected descending order, got [%d, %d]", result[0].Age, result[1].Age)
		}
	})

	t.Run("should return empty list when no matches", func(t *testing.T) {
		result, err := bucket.GetList(func(el TestModel) bool {
			return el.Age > 100
		}, nil)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("Expected empty list, got %d elements", len(result))
		}
	})
}

func TestBucket_GetKeyList(t *testing.T) {
	db := createTestDB(t)
	bucket := NewBucket[TestModel](db, "test-bucket")

	t.Run("should return empty list for empty bucket", func(t *testing.T) {
		result, err := bucket.GetKeyList()

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("Expected empty list, got %d keys", len(result))
		}
	})

	t.Run("should return all keys", func(t *testing.T) {
		bucket.Create("key1", TestModel{ID: "1", Name: "Alice", Age: 30})
		bucket.Create("key2", TestModel{ID: "2", Name: "Bob", Age: 25})
		bucket.Create("key3", TestModel{ID: "3", Name: "Charlie", Age: 35})

		result, err := bucket.GetKeyList()

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if len(result) != 3 {
			t.Errorf("Expected 3 keys, got %d", len(result))
		}

		keyMap := make(map[string]bool)
		for _, key := range result {
			keyMap[key] = true
		}
		if !keyMap["key1"] || !keyMap["key2"] || !keyMap["key3"] {
			t.Errorf("Expected keys [key1, key2, key3], got %v", result)
		}
	})
}

func TestBucket_GetMap(t *testing.T) {
	db := createTestDB(t)
	bucket := NewBucket[TestModel](db, "test-bucket")

	// Setup test data
	bucket.Create("key1", TestModel{ID: "1", Name: "Alice", Age: 30})
	bucket.Create("key2", TestModel{ID: "2", Name: "Bob", Age: 25})
	bucket.Create("key3", TestModel{ID: "3", Name: "Charlie", Age: 35})

	t.Run("should get all elements as map without filter", func(t *testing.T) {
		result, err := bucket.GetMap(nil)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if len(result) != 3 {
			t.Errorf("Expected 3 elements, got %d", len(result))
		}
		if result["key1"].Name != "Alice" {
			t.Errorf("Expected Alice, got %s", result["key1"].Name)
		}
	})

	t.Run("should filter elements by condition", func(t *testing.T) {
		// Filter: Age >= 30
		result, err := bucket.GetMap(func(el TestModel) bool {
			return el.Age >= 30
		})

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("Expected 2 elements (Age >= 30), got %d", len(result))
		}
		for key, el := range result {
			if el.Age < 30 {
				t.Errorf("Key %s: Expected Age >= 30, got %d", key, el.Age)
			}
		}
	})

	t.Run("should return empty map when no matches", func(t *testing.T) {
		result, err := bucket.GetMap(func(el TestModel) bool {
			return el.Age > 100
		})

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("Expected empty map, got %d elements", len(result))
		}
	})
}

func TestBucket_DeleteAll(t *testing.T) {
	db := createTestDB(t)
	bucket := NewBucket[TestModel](db, "test-bucket")

	t.Run("should delete all elements", func(t *testing.T) {
		// Add some data
		bucket.Create("key1", TestModel{ID: "1", Name: "Alice", Age: 30})
		bucket.Create("key2", TestModel{ID: "2", Name: "Bob", Age: 25})
		bucket.Create("key3", TestModel{ID: "3", Name: "Charlie", Age: 35})

		// Verify data exists
		list, _ := bucket.GetList(nil, nil)
		if len(list) != 3 {
			t.Fatalf("Expected 3 elements before delete, got %d", len(list))
		}

		// Delete all
		err := bucket.DeleteAll()

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Verify all deleted
		list, _ = bucket.GetList(nil, nil)
		if len(list) != 0 {
			t.Errorf("Expected 0 elements after DeleteAll, got %d", len(list))
		}
	})

	t.Run("should allow operations after DeleteAll", func(t *testing.T) {
		bucket.DeleteAll()

		// Should be able to create new elements
		model := TestModel{ID: "1", Name: "Alice", Age: 30}
		_, err := bucket.Create("key1", model)

		if err != nil {
			t.Fatalf("Expected to create after DeleteAll, got error: %v", err)
		}

		result, _ := bucket.Get("key1")
		if result.Name != "Alice" {
			t.Errorf("Expected Alice, got %s", result.Name)
		}
	})
}

func TestNewBucket_CreatesIfNotExists(t *testing.T) {
	db := createTestDB(t)

	// Should not panic when creating bucket
	bucket := NewBucket[TestModel](db, "new-bucket")

	if bucket == nil {
		t.Fatal("Expected bucket to be created")
	}

	// Should be able to use the bucket
	model := TestModel{ID: "1", Name: "Test", Age: 20}
	_, err := bucket.Create("key1", model)
	if err != nil {
		t.Fatalf("Expected to use newly created bucket, got error: %v", err)
	}
}

func TestBucket_MultipleTypes(t *testing.T) {
	db := createTestDB(t)

	// Create buckets with different types
	bucket1 := NewBucket[TestModel](db, "bucket1")
	bucket2 := NewBucket[string](db, "bucket2")

	// Add data to both
	bucket1.Create("key1", TestModel{ID: "1", Name: "Alice", Age: 30})
	bucket2.Create("key1", "test-string")

	// Verify both work independently
	result1, err1 := bucket1.Get("key1")
	result2, err2 := bucket2.Get("key1")

	if err1 != nil || err2 != nil {
		t.Fatalf("Expected no errors, got: %v, %v", err1, err2)
	}

	if result1.Name != "Alice" {
		t.Errorf("Expected Alice, got %s", result1.Name)
	}
	if result2 != "test-string" {
		t.Errorf("Expected test-string, got %s", result2)
	}
}
