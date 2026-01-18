# Implementation Plan: Unified DuckDB Storage

**Branch**: `004-unified-duckdb-storage` | **Date**: 2026-01-18 | **Spec**: [specs/004-storage/spec.md](specs/004-storage/spec.md)
**Input**: Feature specification from `/specs/004-storage/spec.md`

## Summary

Consolidate all persistent storage into a single embedded DuckDB instance, replacing the dual-database (BadgerDB + DuckDB) approach. This simplifies the codebase and aligns with the single-binary principle while maintaining analytical capabilities for audit logs and efficient key-value lookups for caching.

## Technical Context

**Language/Version**: Go 1.25.6+  
**Primary Dependencies**: `github.com/marcboeker/go-duckdb`, `github.com/jmoiron/sqlx`  
**Storage**: Embedded DuckDB (`.memex/brain.duckdb`)  
**Testing**: `go test` (Unit and Integration)  
**Target Platform**: Darwin/Linux/Windows (Single Static Binary)
**Project Type**: Single project (Go CLI/Proxy)  
**Performance Goals**: <10ms for cache lookups, <1ms overhead for audit logging (async)  
**Constraints**: Single database file, no CGO where possible (though go-duckdb requires CGO for the DuckDB library itself), thread-safe access.

## Constitution Check

| Principle | Check | Status |
| ---------- | ------- | -------- |
| I. Single Binary | DuckDB is embedded, no external server. | ✅ |
| II. Go-Only | Using Go bindings for DuckDB and sqlx. | ✅ |
| III. Embedded Storage | Consolidation on DuckDB satisfies "Embedded Storage". | ✅ |
| V. Middleware Chain | Storage layers will support the existing middleware order. | ✅ |

## Project Structure

### Documentation (this feature)

```text
specs/004-storage/
├── plan.md              # This file
├── spec.md              # Feature specification
└── tasks.md             # Task breakdown
```

### Source Code

```text
internal/
└── store/               # Unified DuckDB storage wrappers
    ├── store.go         # SQLX connection and table initialization
    ├── audit.go         # Audit log operations
    └── cache.go         # Cache entry operations
pkg/types/
└── storage.go           # Storage-specific types (if not in proxy types)
```

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
| ----------- | ------------ | ------------------------------------- |
| CGO Dependency | `go-duckdb` requires CGO to link with the DuckDB C++ library. | Pure Go DuckDB driver does not exist yet; alternatives like SQLite lack the ARRAY/BLOB performance and analytical features planned for Memex. |
