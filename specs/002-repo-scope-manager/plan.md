# Implementation Plan: Repo-Aware Scope Manager

**Branch**: `002-repo-scope-manager` | **Date**: 2026-01-18 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-repo-scope-manager/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Implement a mechanism to detect the "Project Context" of the active session to isolate cache data (Multi-Tenancy on Localhost) by identifying the current working directory's Git Remote Origin URL on initialization or falling back to a hashed absolute path.

## Technical Context

**Language/Version**: Go 1.25.6+
**Primary Dependencies**: `go-git/v5` (for git remote detection)
**Storage**: N/A (modifies cache key generation logic)
**Testing**: `go test` (Unit tests for scope detection)
**Target Platform**: Cross-platform (macOS, Linux, Windows)
**Project Type**: CLI/Service
**Performance Goals**: <5ms for scope detection on startup
**Constraints**: Must run first in the middleware chain.
**Scale/Scope**: Local execution, per-process scope.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [x] **I. Single Binary Distribution**: Uses `go-git` (pure Go), no external git binary required.
- [x] **II. Go-Only Implementation**: All logic in Go.
- [x] **III. Embedded Storage Architecture**: N/A (no new storage).
- [x] **IV. Tri-Partite Cache Key System**: Enhances this by scoping the keys.
- [x] **V. Middleware Chain Order**: Implements the *first* step (Auth/Scope).
- [x] **VI. Protocol Compatibility**: Transparent to the client.

## Project Structure

### Documentation (this feature)

```text
specs/002-repo-scope-manager/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
src/
├── internal/
│   ├── proxy/         # Scope detection middleware
│   └── store/         # Cache key generation updates
├── pkg/
│   └── types/         # Scope context definition
└── tests/
    └── unit/          # Scope detection tests
```

**Structure Decision**: Standard Go service structure, adding middleware to `internal/proxy` and types to `pkg/types`.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| N/A | | |
