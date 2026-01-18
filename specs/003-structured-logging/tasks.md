# Tasks: Structured Logging with slog

**Feature**: Structured Logging with slog
**Feature Branch**: `003-structured-logging`
**Status**: Ready to Implement

## Implementation Strategy

We will implement this feature in priority order (P0 -> P1).

1. **Setup**: Define configuration structures.
2. **P0 (PII)**: Define sensitive types and masking logic first to ensure they are available for logging.
3. **P1 (Config)**: Initialize the logger with the new configuration.
4. **P1 (Refactor)**: Replace existing print calls with the new logger.

## Phase 1: Setup

*Goal: Initialize configuration structures.*

- [x] T001 Update `ProxyConfig` and add `LogConfig` struct in `internal/proxy/config.go`
- [x] T002 Update `DefaultConfigLoader` to load logging configuration in `internal/proxy/config.go`

## Phase 2: Foundational

*Goal: Blocking prerequisites.*

(None identified - PII moved to P0 User Story phase)

## Phase 3: User Story 3 - PII Protection (P0)

*Goal: Ensure sensitive data is masked.*
*Independent Test: Unit tests for LogValuer implementations.*

- [x] T003 [US3] Create `pkg/types/sensitive.go` with `LogValuer` implementations for sensitive strings
- [x] T004 [US3] Add unit tests for PII masking in `pkg/types/sensitive_test.go`

## Phase 4: User Story 1 - Configure Logging (P1)

*Goal: Control logging behavior via configuration.*
*Independent Test: Verify application starts with different log levels and formats.*

- [x] T005 [US1] Add unit tests for `LogConfig` parsing in `internal/proxy/config_test.go`
- [x] T006 [US1] Initialize global `slog.Logger` in `cmd/proxy/main.go` based on configuration

## Phase 5: User Story 2 - Structured Output (P1)

*Goal: Replace unstructured printing with structured logs.*
*Independent Test: Run proxy and verify JSON/Text output format.*

- [x] T007 [P] [US2] Replace `fmt.Printf` and `log.Print` with `slog` in `internal/proxy/handler.go`
- [x] T008 [P] [US2] Replace `fmt.Printf` and `log.Print` with `slog` in `internal/proxy/middleware.go`
- [x] T009 [US2] Remove `Debug` field from `ProxyConfig` and cleanup usages in `internal/proxy/config.go`

## Final Phase: Polish & Verification

*Goal: Final consistency checks.*

- [x] T010 Verify no `fmt.Printf` calls remain in the proxy path (Manual Grep)
- [x] T011 Verify PII is masked in actual logs during a request (Manual Test)

## Dependencies

- US3 (PII) depends on Setup
- US1 (Config) depends on Setup
- US2 (Output) depends on US1 (Logger Init) and US3 (PII Types)
