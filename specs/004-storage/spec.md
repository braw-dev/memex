# Feature Specification: Unified DuckDB Storage

**Feature Branch**: `004-unified-duckdb-storage`  
**Created**: 2026-01-18  
**Status**: Draft  
**Input**: Update issue #14 to only use DuckDB and not badger. Initialize a single embedded DuckDB instance for all storage needs.

## User Scenarios & Testing

### User Story 1 - Initialize Unified Storage (Priority: P1)

As a developer, I want Memex to initialize a single DuckDB database file so that all application data (logs and cache) is stored in one place.

**Why this priority**: Foundational requirement for consolidating storage and removing BadgerDB dependency.

**Independent Test**: Run Memex; verify that `.memex/brain.duckdb` is created and contains the required schema.

**Acceptance Scenarios**:

1. **Given** no `.memex` directory exists, **When** the application starts, **Then** a `.memex` directory and `brain.duckdb` file are created.
2. **Given** an existing `brain.duckdb` file, **When** the application starts, **Then** it connects to the existing database without data loss.

---

### User Story 2 - Persistent Audit Logging (Priority: P2)

As an operator, I want all proxy requests to be logged to a DuckDB table so that I can analyze usage and costs.

**Why this priority**: Required for "Savings Analytics" and "Shadow Billing" features.

**Independent Test**: Send a request through the proxy; verify a new row appears in the `audit_logs` table.

**Acceptance Scenarios**:

1. **Given** the proxy is running, **When** an AI request is processed, **Then** a record is inserted into `audit_logs` with correct token counts and latency.

---

### User Story 3 - Persistent Cache Storage (Priority: P2)

As a user, I want AI responses to be cached in DuckDB so that I can save money and time on repeated prompts.

**Why this priority**: Replaces the BadgerDB key-value store with DuckDB for caching.

**Independent Test**: Send the same prompt twice; verify the second response is served from the `cache_entries` table.

**Acceptance Scenarios**:

1. **Given** a prompt has been processed once, **When** the same prompt is sent again, **Then** the `cache_entries` table is queried and the cached response is returned.

## Requirements

### Functional Requirements

- **FR-001**: System MUST initialize a single embedded DuckDB instance at `.memex/brain.duckdb`.
- **FR-002**: System MUST create an `audit_logs` table with columns: `timestamp`, `scope_id`, `tokens_in`, `tokens_out`, `cost`, `latency`.
- **FR-003**: System MUST create a `cache_entries` table with columns: `hash_key` (PRIMARY KEY), `scope_id`, `system_hash`, `prompt_vector` (ARRAY<FLOAT>), `response_blob` (BLOB), `created_at`.
- **FR-004**: System MUST use a connection pool compliant with Go's `database/sql` interface.
- **FR-005**: System MUST use `sqlx` for database operations.
- **FR-006**: System MUST ensure thread-safe access to the database.

### Key Entities

- **AuditLog**: Represents a single request's metadata for analytics.
- **CacheEntry**: Represents a cached AI response indexed by a tri-partite hash key.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Single database file `.memex/brain.duckdb` manages all persistent state.
- **SC-002**: Zero BadgerDB dependencies remain in the codebase.
- **SC-003**: Database initialization completes in under 100ms.
- **SC-004**: Audit log writes do not block the main request path (async).
