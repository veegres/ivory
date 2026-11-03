# Go profile

These naming and organization choices describe this package's Go style, derived from Ivory; they are not requirements of Go itself.

- Source packages live directly under `server/` in the existing layout. Use `<package>_<concern>.go`, including `<feature>_model.go`, `<feature>_repository.go`, `<feature>_service.go`, and `<feature>_service_<concern>.go`. `main.go` is the entrypoint exception. Test files mirror the source name exactly with `_test.go`.
- Where a feature has HTTP and MCP adapters, use `<feature>_router_http.go` and `<feature>_router_mcp.go`, with `HttpRouter` and `McpRouter`. Keep transport-specific binding in those adapters, not in services.
- Put all domain/transport structs and mappers in the model file. Use a receiver method when the mapper's input type is owned by the package; foreign inputs require free functions. Do not move a mapper merely because only one service uses it.
- Feature service/repository behavior is discoverable through methods on `Service`/`Repository`. A function introducing its own type parameters remains a free function because Go methods cannot introduce additional type parameters.
- Use import aliases only for an actual collision or misleading name. An import sharing the enclosing package's name is not itself a collision. Preserve a dependency's declared package name when it differs from the last path segment.
- Use standard `testing`, table-driven cases, and `t.Run`. Test storage with the project's isolated test store. Keep pure parsing and mapping independently testable.
- Use the project's formatter, vet/lint, and test commands; do not invent a Makefile target.
- Comments explain a non-obvious reason or invariant. Exported functions, interfaces, and models may have a concise one-line behavioral contract; do not add multi-line documentation merely to restate the identifier.

## Gin binding, when Gin is used

Declare field-shape validation with binding tags and ensure every transport invokes the same validator. Cross-field rules, permissions, registry lookups, and uniqueness belong in services. Other entry points that bypass binding must still validate new input.

In Ivory's configured validator, `notblank` and `valid` are registered extensions; `Valid() bool` backs rich enum validation. Integer ports use `min=1,max=65535`, because the built-in `port` rule expects an unsigned kind. `required` on a struct marks its MCP schema requirement but does not validate the struct's presence with that validator configuration; nested field rules do the checking. Keep MCP schema enum values aligned with accepted values. Confirm these configuration details before adopting them in a new Gin project; do not assume custom rules are built in.
