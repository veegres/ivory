# Backend development

## Build

Use `go mod tidy` to download all libraries.

In the project directory, you can run:

- `go build ivory` - Runs the app in the development mode

## Testing

Run tests with:
- `go test ./...` - Run all tests
- `go test ./... -v` - Run with verbose output
- `go test ./... -cover` - Run with coverage report
- `go test ./src/storage/...` - Run tests for specific package

### Testing Guidelines

#### ✅ What to Test (Business Logic & Data Transformations)

1. **Storage Adapters** (`src/storage/*`)
   - CRUD operations and validation
   - Path/key validation and security (e.g., path traversal prevention)
   - Configuration parsing and defaults
   - Use in-memory databases/temp directories with `t.Cleanup()`

2. **Client Mappers** (`src/clients/*`)
   - API response → internal model transformations
   - Data type handling (e.g., `json.RawMessage` for flexible fields)
   - Query parsing and normalization logic
   - Configuration validation (required fields, format checks)
   - **Do not test** actual HTTP/network calls

3. **Business Logic Services** (`src/features/*`)
   - Core algorithms (encryption, calculations, etc.)
   - Input validation and error handling
   - Edge cases (empty values, special characters, unicode)
   - Integration scenarios (encrypt → decrypt cycles)

4. **Test Patterns to Follow:**
   - Use subtests with `t.Run()` for organization
   - Table-driven tests for multiple scenarios
   - Test error cases explicitly
   - Use `t.Helper()` for test utilities
   - Clean up resources with `t.Cleanup()`

#### ❌ What NOT to Test (Infrastructure & HTTP Layers)

1. **HTTP Routers** (`src/routers/*`)
   - These are thin HTTP handlers that delegate to services
   - Testing would require complex HTTP mocking
   - Focus on testing the underlying services instead

2. **Repositories** (`src/repositories/*`)
   - Thin wrappers over storage adapters
   - No business logic to test
   - Storage adapters are already tested

3. **External Network Calls**
   - LDAP connections
   - OIDC provider communication
   - Actual database queries (test query parsing, not execution)
   - These require external services or extensive mocking

4. **Main/Initialization Code**
   - Application bootstrap
   - Dependency injection setup
   - Router registration

### Example Test Structure

```go
func TestService_Method(t *testing.T) {
    service := NewService()

    t.Run("should handle valid input", func(t *testing.T) {
        result, err := service.Method(validInput)

        if err != nil {
            t.Fatalf("Expected no error, got: %v", err)
        }
        if result != expected {
            t.Errorf("Expected %v, got %v", expected, result)
        }
    })

    t.Run("should return error for invalid input", func(t *testing.T) {
        _, err := service.Method(invalidInput)

        if err == nil {
            t.Fatal("Expected error, got nil")
        }
        if err.Error() != "expected message" {
            t.Errorf("Expected specific error, got: %v", err)
        }
    })
}
```

## Environment variables

- `IVORY_URL_ADDRESS` - change tcp network address, `default: :8080` (if `GIN_MODE=release` it either `:80` or `:443` when cert files are set)
- `IVORY_URL_PATH` - change sub path under which Ivory serve its routes, `default: /`
- `IVORY_STATIC_FILES_PATH` - add static files to host, `default: empty`
- `IVORY_CERT_FILE_PATH` - provide your sever cert to enable tls, `default: empty`
- `IVORY_CERT_KEY_FILE_PATH` - provide your server key to enable tls, `default: empty`
- `IVORY_VERSION_TAG` - change version tag, `default: v0.0.0`
- `IVORY_VERSION_COMMIT` - change version commit, `default: 000000000000`

## Learn More

This app is build by:

- [Bolt](https://github.com/boltdb/bolt) - Embedded key/value database
- [Gin](https://github.com/gin-gonic/gin) - Web framework

