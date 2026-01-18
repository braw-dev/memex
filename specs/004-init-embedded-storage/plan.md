# Implementation Plan: Initialize Embedded Storage Engines

**Branch**: `004-init-embedded-storage` | **Date**: 2026-01-18 | **Spec**: [specs/004-init-embedded-storage/spec.md](spec.md)
**Input**: Feature specification from `/specs/004-init-embedded-storage/spec.md`

## Summary

This feature initializes the persistent storage layer for Memex, setting up BadgerDB (Key-Value) and DuckDB (OLAP) within a hidden `.memex` directory. It establishes the foundational data access patterns without external dependencies.

## Technical Context

**Language/Version**: Go 1.25.6
**Primary Dependencies**: `dgraph-io/badger/v4`, `marcboeker/go-duckdb`
**Storage**: Embedded BadgerDB v4, Embedded DuckDB
**Testing**: Go `testing` package
**Target Platform**: Cross-platform (Linux, macOS, Windows) - Single Binary
**Project Type**: Single CLI/Server binary
**Performance Goals**: Initialization < 500ms
**Constraints**: Must manage file locks and permissions gracefully
**Scale/Scope**: Local single-user storage

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [x] **Single Binary Distribution**: Using embedded Go libraries only.
- [x] **Go-Only Implementation**: All code in Go.
- [x] **Embedded Storage Architecture**: Using BadgerDB and DuckDB as mandated.
- [x] **Tri-Partite Cache Key System**: N/A (Infrastructure only).
- [x] **Middleware Chain Order**: N/A.
- [x] **Protocol Compatibility**: N/A.

**Technical Invariants Verification**:

- Language: Go 1.25.6+ ✅
- Distribution: Single static binary ✅
- OLAP Storage: Embedded DuckDB ✅
- KV Storage: Embedded BadgerDB v4 ✅

## Project Structure

### Documentation (this feature)

```text
specs/004-init-embedded-storage/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
└── tasks.md             # Phase 2 output
```

### Source Code (repository root)

```text
internal/store/
├── store.go           # Main Store struct and factory
├── badger.go          # BadgerDB implementation
├── duckdb.go          # DuckDB implementation
├── schema.go          # Database schema definitions
└── store_test.go      # Integration tests
```

**Structure Decision**: Placing all storage logic in `internal/store` encapsulates the database details from the rest of the application, exposing only high-level interfaces defined in `pkg/types` (or `internal/store` if they are private to internal, but Spec suggests `pkg/types` for shared structs, though interfaces are best close to usage or definition. The Constitution mentions `pkg/types/ -> Shared structs and interfaces`. I will define interfaces there if needed, or keep them in `internal/store` if they are implementation details of the Store).
For now, `internal/store` is the implementation package.

## Complexity Tracking

No constitution violations.
