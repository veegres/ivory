# Persistence profiles

Choose persistence from access patterns, consistency, deployment, and backup needs, independently of Go/Rust or frontend choice. SQLite is SQL, not NoSQL. A small application does not imply a document database, and a NoSQL choice does not imply a remote service.

## Embedded SQL

Suitable for relational constraints, joins, reporting, and transactions in a small deployment. Repositories own SQL; migrations are ordered and released migrations are immutable. Exercise the real schema in isolated repository tests. Document concurrent writers, locking, restore, and any multi-instance limitations relevant to the selected engine.

## Embedded key/value or document storage

Define record keys, document shape/version, indexes or scan costs, transaction boundaries, atomic uniqueness checks, and schema evolution explicitly. A read/check/write sequence must share the required transaction or conditional-write operation. Avoid recreating a query engine in services merely to force a NoSQL choice.

## Remote document/key/value storage

Use when its deployment or access-pattern advantages are needed. Record consistency guarantees, conditional writes, pagination, indexes, and connection ownership. Keep network retries bounded by the operation's idempotency; do not assume a timeout means a write failed.

## Compatibility

Separate internal storage evolution from an exported backup format. Follow the project's release and migration policy; don't copy Ivory's internal-data migration policy into another application. Released export formats need explicit versions and fixtures proving old imports still work. Document what is excluded (credentials, external references), and validate restore behavior rather than just export serialization.
