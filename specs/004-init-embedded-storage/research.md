# Research: Initialize Embedded Storage Engines

**Feature**: Initialize Embedded Storage Engines (004-init-embedded-storage)
**Date**: 2026-01-18

## Decisions

### 1. Storage Libraries

- **KV Store**: `dgraph-io/badger/v4`
  - **Rationale**: Mandated by Constitution. Pure Go (mostly), high performance, designed for SSDs.
  - **Configuration**:
    - ValueLogFileSize: 256MB (Optimize for lower disk usage overhead while maintaining performance)
    - Compression: Snappy (Default)
    - Logger: Wrapped interface to adapter to project's structured logger (or nil for now to keep it clean)
    - SyncWrites: True (Safety first)

- **OLAP Store**: `marcboeker/go-duckdb`
  - **Rationale**: Mandated by Constitution. Best-in-class embedded OLAP.
  - **Configuration**:
    - DSN: `.memex/memex.db`
    - InitSQL: `INSTALL json; LOAD json;` (if needed for JSON logs)

### 2. Directory Structure

- **Location**: `internal/store/`
- **Files**:
  - `store.go`: Main `Store` struct and initialization logic.
  - `badger.go`: BadgerDB specific wrapper.
  - `duckdb.go`: DuckDB specific wrapper.
  - `schema.go`: Initial schema definitions (migrations).

### 3. Data Storage Location

- **Path**: `.memex/` in the current working directory (project root).
- **Rationale**: "Zero configuration" implies running where the user is. `git` uses `.git` in project root. Memex should follow similar pattern for project-local context.
- **Gitignore**: The `memex init` command (or this feature) should probably advise on adding `.memex` to `.gitignore`, but for now we just create the directory.

### 4. Dependency Management

- Need to add `dgraph-io/badger/v4` and `marcboeker/go-duckdb` to `go.mod`.
- **Note**: `go-duckdb` requires CGO. We must ensure the build process supports this while producing a static binary (using `CGO_ENABLED=1` and appropriate linker flags).

## Alternatives Considered

- **SQLite instead of DuckDB**: SQLite is easier for simple storage, but DuckDB is mandated for OLAP/Analytical queries which will be needed for "Shadow Billing" and cost analysis.
- **BoltDB instead of BadgerDB**: BoltDB is read-optimized and simpler, but BadgerDB is mandated by Constitution and offers better write performance for cache usage.

## Unknowns & Risks

- **CGO Static Linking**: ensuring `go-duckdb` links statically on all platforms (especially macOS/Linux cross-compilation) can be tricky. We will assume standard local build for now.
- **Database Locks**: If multiple shells open memex, BadgerDB/DuckDB might lock the files.
  - **Mitigation**: Fail fast with a clear error message "Memex is already running in another process".

## Reference Material

- BadgerDB Docs: <https://dgraph.io/docs/badger/>
- DuckDB Go Docs: <https://github.com/marcboeker/go-duckdb>
