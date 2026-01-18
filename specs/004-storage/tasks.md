# Tasks: Unified DuckDB Storage

**Feature**: Unified DuckDB Storage
**Feature Branch**: `004-unified-duckdb-storage`
**Status**: Ready to Implement

## Phase 1: Setup

- [ ] T001 Add `github.com/marcboeker/go-duckdb` and `github.com/jmoiron/sqlx` to `go.mod`
- [ ] T002 Create `internal/store` directory

## Phase 2: Foundational - Database Initialization

- [ ] T003 Implement `NewStore` in `internal/store/store.go` to initialize `.memex/brain.duckdb`
- [ ] T004 Implement table creation logic for `audit_logs` and `cache_entries` in `internal/store/store.go`
- [ ] T005 Implement connection pooling using `sqlx.Connect` and `db.SetMaxOpenConns(1)` (DuckDB recommendation for single-process)

## Phase 3: Audit Logging Implementation

- [ ] T006 Define `AuditLog` struct in `internal/store/audit.go`
- [ ] T007 Implement `WriteLog(log *AuditLog)` in `internal/store/audit.go`
- [ ] T008 Update `internal/proxy/middleware.go` or equivalent to use the new audit logger (async)

## Phase 4: Cache Implementation

- [ ] T009 Define `CacheEntry` struct in `internal/store/cache.go`
- [ ] T010 Implement `GetCache(hashKey string)` in `internal/store/cache.go`
- [ ] T011 Implement `SetCache(entry *CacheEntry)` in `internal/store/cache.go`

## Phase 5: Verification & Cleanup

- [ ] T012 Create integration test for unified storage in `tests/integration/storage_test.go`
- [ ] T013 Verify removal of all BadgerDB references and dependencies
- [ ] T014 Run performance benchmark to ensure <10ms lookup latency
