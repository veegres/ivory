# Rust profile

- Source lives in `server/src/`. Declare `features/auth/service.rs` from `features/auth.rs`; use the existing non-`mod.rs` layout. Module paths disambiguate names: `auth::Service`, `pin::Repository`. Do not apply Go's package prefixes to Rust filenames.
- A feature uses `model.rs`, `repository.rs`, `service.rs`, and `router.rs`. Split large services into `service_<concern>.rs` containing additional `impl Service` blocks.
- Model files contain requests, responses, criteria, and mappers. Mappers use inherent methods on local types; foreign input conversions can remain free functions. Use existing conversion traits where they express the same contract without changing the public API.
- Requests expose `validate()` for shape constraints. Services enforce relationships, authorization, and uniqueness; alternate callers must not bypass required validation.
- Use typed errors and the project's central HTTP mapping. For Veembly/DTDb-style APIs, expected refusals carry a `Reason`, responses are camelCase in a `response` envelope, and absent optional fields use `skip_serializing_if`. Those wire choices are selected by the project, not properties of Rust.
- SQL belongs in repositories. For SQLx/SQLite projects, use the existing storage pool, migration mechanism, and `storage::memory()` fixture; do not introduce an independent connection lifecycle per feature.
- Tests live in `#[cfg(test)] mod tests` with table-driven cases. Keep parsers pure and use captured input fixtures for upstream markup/JSON.
- Follow `rustfmt.toml` and clippy configuration; these projects use width 120 and reject clippy warnings. Non-obvious explanatory comments use `// NOTE:`; public contracts may have concise `///` documentation.
