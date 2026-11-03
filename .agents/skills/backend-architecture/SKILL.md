---
name: backend-architecture
description: Design or change a feature-oriented backend, including models, services, repositories, transports, and integration boundaries. Use the project's language profile for implementation details.
---

# Backend architecture

Read the project's AGENTS.md for its selected language, persistence, contracts, and exceptions. These are this package's conventions, not universal language rules. Apply the shared responsibilities below, then read only the applicable profile: [Go](references/go.md) or [Rust](references/rust.md). For another language, retain the responsibilities and use its native module/type conventions; propose a profile when actual work establishes useful differences.

## Responsibilities

- Organize by feature: one business entity or cohesive capability owns its model, repository, service, and transport adapters. Do not introduce empty layers just to satisfy a diagram.
- Models own domain values, requests, responses, criteria, and boundary mappers. Service/repository/router implementation types stay beside their implementations. A transport-only request or response is named `Router<What>Request/Response/Query`; a type consumed by a service keeps a domain name.
- Repositories own persistence operations, queries, indexes, and storage mappings. Services own business decisions, permissions, cross-field rules, and orchestration. Routers bind inputs, authenticate, invoke services, and serialize responses.
- Request types state input shape. Validate at every external entry point, including CLI, background jobs, and alternate transports. Reuse validation rather than depending on an HTTP handler having run. Loading old persisted data is a separate compatibility path, not a new-request validation path.
- Keep related behavior in one feature; split a large service by concern without creating a second feature merely because a file is large.
- `core/` composes configuration, lifecycle, errors, routing, and dependencies. `clients/` owns reusable transport mechanics. Optional `plugins/` owns interchangeable integrations; features own product-specific interpretation.
- Features expose the types their consumers need, aliasing or mapping integration types at the boundary. Dependencies follow the consuming direction; avoid cycles and a generic common package of unrelated domain types.
- Prefer a narrow interface at the consuming boundary when substitution is useful. Do not add interfaces or an adapter registry for hypothetical alternatives.
- Error codes, response envelopes, serialization casing, and authentication rules are explicit project contracts. Preserve them across transports; never copy another application's wire format implicitly.

## Persistence and verification

Storage choice is independent of backend language. Read [persistence](references/persistence.md) when designing or changing stored data, transactions, or migrations.

Test model decisions, repository behavior, services, and integration parsing at their owning boundaries. Prefer table-driven cases and deterministic fixtures. Thin transport glue need not duplicate service tests; test transport-specific behavior when that is the behavior at risk. Project test requirements still apply. Preserve existing tests during refactors.

Use the actual commands in AGENTS.md and build configuration. Explain non-obvious invariants in comments; leave syntax and formatting to the configured tools.

## Evolving conventions

When a task exposes a missing or conflicting rule, follow [convention-maintenance](../convention-maintenance/SKILL.md). Complete the authorized task using the nearest established pattern; suggest the smallest concrete rule and whether it belongs in this core, a language profile, or the project's own instructions. Do not generalize a product-specific exception into the shared architecture.
