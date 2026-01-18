# Implementation Tasks: Initialize Embedded Storage Engines

**Branch**: `004-init-embedded-storage` | **Spec**: [spec.md](spec.md) | **Plan**: [plan.md](plan.md)

## Phase 1: Setup

- [x] T001 Add `dgraph-io/badger/v4` and `marcboeker/go-duckdb` to `go.mod`
- [x] T002 Create directory structure `internal/store` with placeholders

## Phase 2: Foundation (Shared Interfaces)

- [x] T003 Define `Store`, `KVStore`, `OLAPStore` interfaces in `internal/store/store.go`
- [x] T004 Create `Store` struct and `New(path string)` factory in `internal/store/store.go`

## Phase 3: User Story 1 - Initialization (P1)

*Goal: Initialize storage engines in .memex directory*

- [x] T005 [P] [US1] Implement BadgerDB wrapper in `internal/store/badger.go` (Open, Close, Get, Set)
- [x] T006 [P] [US1] Implement DuckDB wrapper in `internal/store/duckdb.go` (Open, Close, Exec)
- [x] T007 [US1] Update `New()` in `internal/store/store.go` to initialize both engines
- [x] T008 [US1] Create integration test `internal/store/store_test.go` verifying `.memex` creation and file existence

## Phase 4: User Story 2 - Persistence & Restart (P2)

*Goal: Ensure data persists across restarts*

- [x] T009 [US2] Add CRUD methods to `Store` interface wrappers to support testing
- [x] T010 [US2] Update `store_test.go` to test writing data, closing, reopening, and reading data

## Phase 5: User Story 3 - Graceful Failure (P3)

*Goal: Handle errors gracefully*

- [x] T011 [US3] Add specific error types for initialization failures in `internal/store/store.go`
- [x] T012 [US3] Add test case for read-only filesystem/permission denied in `store_test.go`

## Phase 6: Polish & Schema

- [x] T013 Implement schema initialization (audit_logs table) in `internal/store/schema.go` and call from `duckdb.go`
- [x] T014 Configure BadgerDB for large binary payloads (ValueLogFileSize) in `internal/store/badger.go`

## Dependencies

- Phase 2 depends on Phase 1
- Phase 3 depends on Phase 2
- Phase 4 depends on Phase 3
- Phase 5 depends on Phase 3
- Phase 6 depends on Phase 3

## Implementation Strategy

We will implement the interfaces first, then the concrete wrappers for Badger and DuckDB. We will start with a simple `New` that does nothing, then wire them up.
Tests will be added incrementally.
