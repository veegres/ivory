package postgres

import (
	"ivory/src/clients/database"
	"testing"
)

func TestClient_GetApplicationName(t *testing.T) {
	client := NewClient("TestApp")

	t.Run("should truncate session ID to 7 characters", func(t *testing.T) {
		result := client.GetApplicationName("1234567890abcdef")

		expected := "TestApp [1234567]"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("should handle short session ID", func(t *testing.T) {
		result := client.GetApplicationName("abc")

		expected := "TestApp [abc]"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("should handle exactly 7 characters", func(t *testing.T) {
		result := client.GetApplicationName("abcdefg")

		expected := "TestApp [abcdefg]"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})
}

func TestClient_trimQuery(t *testing.T) {
	client := NewClient("TestApp")

	t.Run("should remove single line comment", func(t *testing.T) {
		query := "SELECT * FROM users -- this is a comment\nWHERE id = 1"
		result := client.trimQuery(query)

		expected := "SELECT * FROM usersWHERE id = 1"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("should remove multiple comments", func(t *testing.T) {
		query := "SELECT * FROM users -- comment 1\nWHERE id = 1 -- comment 2\nLIMIT 10"
		result := client.trimQuery(query)

		expected := "SELECT * FROM usersWHERE id = 1LIMIT 10"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("should keep query without comments", func(t *testing.T) {
		query := "SELECT * FROM users WHERE id = 1"
		result := client.trimQuery(query)

		if result != query {
			t.Errorf("Expected '%s', got '%s'", query, result)
		}
	})

	t.Run("should handle comment at end without newline", func(t *testing.T) {
		query := "SELECT * FROM users -- comment"
		result := client.trimQuery(query)

		expected := "SELECT * FROM users"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("should handle query with only comment", func(t *testing.T) {
		query := "-- this is just a comment"
		result := client.trimQuery(query)

		expected := ""
		if result != expected {
			t.Errorf("Expected empty string, got '%s'", result)
		}
	})
}

func TestClient_parseQuery(t *testing.T) {
	client := NewClient("TestApp")

	t.Run("should parse basic SELECT query", func(t *testing.T) {
		query := "SELECT * FROM users"
		result := client.parseQuery(query)

		if result.SELECT != 1 {
			t.Errorf("Expected SELECT count 1, got %d", result.SELECT)
		}
		if result.FROM != 1 {
			t.Errorf("Expected FROM count 1, got %d", result.FROM)
		}
		if result.LIMIT != 0 {
			t.Errorf("Expected LIMIT count 0, got %d", result.LIMIT)
		}
		if result.Semicolon {
			t.Error("Expected no semicolon")
		}
	})

	t.Run("should parse SELECT query with LIMIT", func(t *testing.T) {
		query := "SELECT * FROM users LIMIT 10"
		result := client.parseQuery(query)

		if result.SELECT != 1 {
			t.Errorf("Expected SELECT count 1, got %d", result.SELECT)
		}
		if result.FROM != 1 {
			t.Errorf("Expected FROM count 1, got %d", result.FROM)
		}
		if result.LIMIT != 1 {
			t.Errorf("Expected LIMIT count 1, got %d", result.LIMIT)
		}
	})

	t.Run("should detect semicolon at end", func(t *testing.T) {
		query := "SELECT * FROM users;"
		result := client.parseQuery(query)

		if !result.Semicolon {
			t.Error("Expected semicolon to be detected")
		}
	})

	t.Run("should parse UPDATE query", func(t *testing.T) {
		query := "UPDATE users SET name = 'John' WHERE id = 1"
		result := client.parseQuery(query)

		if result.UPDATE != 1 {
			t.Errorf("Expected UPDATE count 1, got %d", result.UPDATE)
		}
		if result.SELECT != 0 {
			t.Errorf("Expected SELECT count 0, got %d", result.SELECT)
		}
	})

	t.Run("should parse INSERT query", func(t *testing.T) {
		query := "INSERT INTO users (name) VALUES ('John')"
		result := client.parseQuery(query)

		if result.INSERT != 1 {
			t.Errorf("Expected INSERT count 1, got %d", result.INSERT)
		}
	})

	t.Run("should parse DELETE query", func(t *testing.T) {
		query := "DELETE FROM users WHERE id = 1"
		result := client.parseQuery(query)

		if result.DELETE != 1 {
			t.Errorf("Expected DELETE count 1, got %d", result.DELETE)
		}
		if result.FROM != 1 {
			t.Errorf("Expected FROM count 1, got %d", result.FROM)
		}
	})

	t.Run("should parse EXPLAIN query", func(t *testing.T) {
		query := "EXPLAIN SELECT * FROM users"
		result := client.parseQuery(query)

		if result.EXPLAIN != 1 {
			t.Errorf("Expected EXPLAIN count 1, got %d", result.EXPLAIN)
		}
		if result.SELECT != 1 {
			t.Errorf("Expected SELECT count 1, got %d", result.SELECT)
		}
	})

	t.Run("should handle case insensitivity", func(t *testing.T) {
		query := "SeLeCt * FrOm users LiMiT 10"
		result := client.parseQuery(query)

		if result.SELECT != 1 {
			t.Errorf("Expected SELECT count 1, got %d", result.SELECT)
		}
		if result.FROM != 1 {
			t.Errorf("Expected FROM count 1, got %d", result.FROM)
		}
		if result.LIMIT != 1 {
			t.Errorf("Expected LIMIT count 1, got %d", result.LIMIT)
		}
	})

	t.Run("should not count keywords inside parentheses without spaces", func(t *testing.T) {
		// Note: parseQuery uses strings.Fields which splits on whitespace
		// So "(SELECT" without space is treated as one word and not matched
		query := "SELECT * FROM (SELECT id FROM users) AS subquery"
		result := client.parseQuery(query)

		// Only counts "SELECT" and "FROM" that are standalone words
		// "(SELECT" is not matched because it includes the parenthesis
		if result.SELECT != 1 {
			t.Errorf("Expected SELECT count 1, got %d", result.SELECT)
		}
		if result.FROM != 2 {
			t.Errorf("Expected FROM count 2, got %d", result.FROM)
		}
	})

	t.Run("should ignore keywords after AS", func(t *testing.T) {
		query := "SELECT id AS select FROM users"
		result := client.parseQuery(query)

		// Should count only the first SELECT, not the aliased 'select'
		if result.SELECT != 1 {
			t.Errorf("Expected SELECT count 1 (ignoring alias), got %d", result.SELECT)
		}
	})

	t.Run("should handle semicolon with spaces", func(t *testing.T) {
		query := "SELECT * FROM users  ;  "
		result := client.parseQuery(query)

		// Note: parseQuery looks at the last word, not the actual last character
		// So ";" as a separate word should be detected
		if !result.Semicolon {
			t.Error("Expected semicolon to be detected")
		}
	})
}

func TestClient_addLimitToQuery(t *testing.T) {
	client := NewClient("TestApp")

	t.Run("should add LIMIT to SELECT query without semicolon", func(t *testing.T) {
		query := "SELECT * FROM users"
		analysis := database.QueryAnalysis{
			SELECT:    1,
			FROM:      1,
			LIMIT:     0,
			Semicolon: false,
		}

		newQuery, newLimit := client.addLimitToQuery(query, analysis, "100")

		expected := "SELECT * FROM users LIMIT 100;"
		if newQuery != expected {
			t.Errorf("Expected '%s', got '%s'", expected, newQuery)
		}
		if newLimit == nil || *newLimit != "100" {
			t.Errorf("Expected limit '100', got %v", newLimit)
		}
	})

	t.Run("should replace semicolon when adding LIMIT", func(t *testing.T) {
		query := "SELECT * FROM users;"
		analysis := database.QueryAnalysis{
			SELECT:    1,
			FROM:      1,
			LIMIT:     0,
			Semicolon: true,
		}

		newQuery, newLimit := client.addLimitToQuery(query, analysis, "50")

		expected := "SELECT * FROM users LIMIT 50;"
		if newQuery != expected {
			t.Errorf("Expected '%s', got '%s'", expected, newQuery)
		}
		if newLimit == nil || *newLimit != "50" {
			t.Errorf("Expected limit '50', got %v", newLimit)
		}
	})

	t.Run("should replace semicolon with spaces", func(t *testing.T) {
		query := "SELECT * FROM users  ;  "
		analysis := database.QueryAnalysis{
			SELECT:    1,
			FROM:      1,
			LIMIT:     0,
			Semicolon: true,
		}

		newQuery, newLimit := client.addLimitToQuery(query, analysis, "25")

		expected := "SELECT * FROM users LIMIT 25;"
		if newQuery != expected {
			t.Errorf("Expected '%s', got '%s'", expected, newQuery)
		}
		if newLimit == nil || *newLimit != "25" {
			t.Errorf("Expected limit '25', got %v", newLimit)
		}
	})

	t.Run("should not add LIMIT if already present", func(t *testing.T) {
		query := "SELECT * FROM users LIMIT 10"
		analysis := database.QueryAnalysis{
			SELECT: 1,
			FROM:   1,
			LIMIT:  1,
		}

		newQuery, newLimit := client.addLimitToQuery(query, analysis, "100")

		if newQuery != query {
			t.Errorf("Expected query unchanged, got '%s'", newQuery)
		}
		if newLimit != nil {
			t.Errorf("Expected nil limit, got %v", newLimit)
		}
	})

	t.Run("should not add LIMIT to UPDATE query", func(t *testing.T) {
		query := "UPDATE users SET name = 'John'"
		analysis := database.QueryAnalysis{
			UPDATE: 1,
			LIMIT:  0,
		}

		newQuery, newLimit := client.addLimitToQuery(query, analysis, "100")

		if newQuery != query {
			t.Errorf("Expected query unchanged, got '%s'", newQuery)
		}
		if newLimit != nil {
			t.Errorf("Expected nil limit, got %v", newLimit)
		}
	})

	t.Run("should not add LIMIT to INSERT query", func(t *testing.T) {
		query := "INSERT INTO users (name) VALUES ('John')"
		analysis := database.QueryAnalysis{
			INSERT: 1,
			LIMIT:  0,
		}

		newQuery, newLimit := client.addLimitToQuery(query, analysis, "100")

		if newQuery != query {
			t.Errorf("Expected query unchanged, got '%s'", newQuery)
		}
		if newLimit != nil {
			t.Errorf("Expected nil limit, got %v", newLimit)
		}
	})

	t.Run("should not add LIMIT to DELETE query", func(t *testing.T) {
		query := "DELETE FROM users WHERE id = 1"
		analysis := database.QueryAnalysis{
			DELETE: 1,
			FROM:   1,
			LIMIT:  0,
		}

		newQuery, newLimit := client.addLimitToQuery(query, analysis, "100")

		if newQuery != query {
			t.Errorf("Expected query unchanged, got '%s'", newQuery)
		}
		if newLimit != nil {
			t.Errorf("Expected nil limit, got %v", newLimit)
		}
	})

	t.Run("should not add LIMIT to EXPLAIN query", func(t *testing.T) {
		query := "EXPLAIN SELECT * FROM users"
		analysis := database.QueryAnalysis{
			EXPLAIN: 1,
			SELECT:  1,
			FROM:    1,
			LIMIT:   0,
		}

		newQuery, newLimit := client.addLimitToQuery(query, analysis, "100")

		if newQuery != query {
			t.Errorf("Expected query unchanged, got '%s'", newQuery)
		}
		if newLimit != nil {
			t.Errorf("Expected nil limit, got %v", newLimit)
		}
	})

	t.Run("should not add LIMIT to SELECT without FROM", func(t *testing.T) {
		query := "SELECT 1"
		analysis := database.QueryAnalysis{
			SELECT: 1,
			FROM:   0,
			LIMIT:  0,
		}

		newQuery, newLimit := client.addLimitToQuery(query, analysis, "100")

		if newQuery != query {
			t.Errorf("Expected query unchanged, got '%s'", newQuery)
		}
		if newLimit != nil {
			t.Errorf("Expected nil limit, got %v", newLimit)
		}
	})
}

func TestClient_normalizeQuery(t *testing.T) {
	client := NewClient("TestApp")

	t.Run("should return original query when trim is nil", func(t *testing.T) {
		query := "SELECT * FROM users -- comment"
		result, limit, err := client.normalizeQuery(query, nil, nil)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if result != query {
			t.Errorf("Expected query unchanged, got '%s'", result)
		}
		if limit != nil {
			t.Errorf("Expected nil limit, got %v", limit)
		}
	})

	t.Run("should return original query when trim is false", func(t *testing.T) {
		query := "SELECT * FROM users -- comment"
		trimFalse := false
		result, limit, err := client.normalizeQuery(query, &trimFalse, nil)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if result != query {
			t.Errorf("Expected query unchanged, got '%s'", result)
		}
		if limit != nil {
			t.Errorf("Expected nil limit, got %v", limit)
		}
	})

	t.Run("should return error when limit without trim", func(t *testing.T) {
		query := "SELECT * FROM users"
		limitStr := "100"
		trimFalse := false

		_, _, err := client.normalizeQuery(query, &trimFalse, &limitStr)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if err.Error() != "cannot limit query without trimming it" {
			t.Errorf("Expected 'cannot limit query without trimming it', got: %v", err)
		}
	})

	t.Run("should trim query without adding limit", func(t *testing.T) {
		query := "SELECT * FROM users -- comment\nWHERE id = 1"
		trimTrue := true
		result, limit, err := client.normalizeQuery(query, &trimTrue, nil)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		expected := "SELECT * FROM usersWHERE id = 1"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
		if limit != nil {
			t.Errorf("Expected nil limit, got %v", limit)
		}
	})

	t.Run("should trim and add limit to SELECT query", func(t *testing.T) {
		query := "SELECT * FROM users -- comment"
		trimTrue := true
		limitStr := "50"
		result, newLimit, err := client.normalizeQuery(query, &trimTrue, &limitStr)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		expected := "SELECT * FROM users LIMIT 50;"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
		if newLimit == nil || *newLimit != "50" {
			t.Errorf("Expected limit '50', got %v", newLimit)
		}
	})

	t.Run("should trim but not add limit to UPDATE query", func(t *testing.T) {
		query := "UPDATE users SET name = 'John' -- comment"
		trimTrue := true
		limitStr := "100"
		result, newLimit, err := client.normalizeQuery(query, &trimTrue, &limitStr)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		expected := "UPDATE users SET name = 'John'"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
		if newLimit != nil {
			t.Errorf("Expected nil limit (UPDATE not eligible), got %v", newLimit)
		}
	})

	t.Run("should trim but not add limit when LIMIT already exists", func(t *testing.T) {
		query := "SELECT * FROM users LIMIT 10 -- comment"
		trimTrue := true
		limitStr := "100"
		result, newLimit, err := client.normalizeQuery(query, &trimTrue, &limitStr)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		expected := "SELECT * FROM users LIMIT 10"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
		if newLimit != nil {
			t.Errorf("Expected nil limit (already has LIMIT), got %v", newLimit)
		}
	})

	t.Run("should handle complex query with trim and limit", func(t *testing.T) {
		query := `SELECT u.id, u.name -- get user data
		FROM users u -- user table
		WHERE u.active = true; -- only active`
		trimTrue := true
		limitStr := "25"
		result, newLimit, err := client.normalizeQuery(query, &trimTrue, &limitStr)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		// Note: trimQuery removes comments and their trailing newlines
		// Regex \\s*--.*\\n* removes spaces before comment, comment text, and newlines
		// Result has tabs from original indentation but no newlines
		expected := "SELECT u.id, u.name\t\tFROM users u\t\tWHERE u.active = true LIMIT 25;"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
		if newLimit == nil || *newLimit != "25" {
			t.Errorf("Expected limit '25', got %v", newLimit)
		}
	})
}
