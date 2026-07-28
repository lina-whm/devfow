# ADR-005: PostgreSQL as Primary Database

**Status:** Accepted

**Context:** DevFlow needs a relational database with strong consistency, JSON support, full-text search, and mature tooling. The database must handle complex queries (task listing with filters, sorting, and cursor pagination) and transactional operations.

**Decision:** Use PostgreSQL 16+ as the primary database.

**Rationale:** PostgreSQL provides the best combination of features for this project:
- **ACID compliance:** Transactional integrity for task moves, sprint operations, and membership changes.
- **JSONB:** Flexible metadata storage for activity log changes (`changes JSONB`) and notification payloads.
- **Full-text search:** Native `tsvector`/`tsquery` for task search, eliminating the need for a separate search index in Phase 1.
- **UUID support:** Native `gen_random_uuid()` with UUIDv7 capability.
- **Extensibility:** Extensions like `pgcrypto`, `pg_stat_statements`.
- **Maturity:** Battle-tested in production at every major tech company.

Alternatives considered:
- **MySQL 8:** Similar feature set but weaker JSON support and no native UUID type. Less robust full-text search.
- **CockroachDB:** PostgreSQL-compatible with horizontal scaling, but adds operational complexity and latency overhead for a single-node deployment.
- **SQLite:** Excellent for small projects but lacks concurrency, JSONB, and full-text search performance needed for a multi-tenant SaaS.

**Consequences:** All team members must be familiar with PostgreSQL-specific features (JSONB operators, `ON CONFLICT`, `RETURNING`). The project leverages PostgreSQL-specific types (`UUID`, `TIMESTAMPTZ`, `JSONB`), making migration to another database difficult — an acceptable trade-off given the benefits.
