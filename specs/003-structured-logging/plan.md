# Implementation Plan - Structured Logging with slog

**Feature**: Structured Logging with slog
**Feature Branch**: `003-structured-logging`
**Specification**: [spec.md](spec.md)

## Technical Context

**Language/Version**: Go 1.25.6+
**Primary Dependencies**: `log/slog` (Standard Library)
**Storage**: N/A
**Project Type**: Backend Service
**Architectural Pattern**: Middleware + Global Logger

## Constitution Check

| Invariant | Status | Notes |
|-----------|--------|-------|
| Go 1.25.6+ | PASS | Using standard library |
| Single Binary | PASS | No external deps |
| Embedded Storage | PASS | N/A |
| Performance | PASS | slog is low-overhead |
| Safety P0 | PASS | PII masking included |

## Phase 0: Research & Design (Completed)

- [x] Research PII masking (`research.md`)
- [x] Design Config Structure (`data-model.md`)
- [x] Write Quickstart Guide (`quickstart.md`)

## Phase 1: Implementation

### Step 1.1: Configuration

- Update `ProxyConfig` in `internal/proxy/config.go`
- Add `LogConfig` struct
- Update config loader to handle new fields

### Step 1.2: PII Masking

- Define `LogValuer` wrapper or interface for sensitive data
- Implement masking logic

### Step 1.3: Global Logger Setup

- In `cmd/proxy/main.go`, initialize `slog.Logger` based on config
- Call `slog.SetDefault()`
- Handle file opening for `Path`

### Step 1.4: Refactor Proxy Module

- Replace `fmt.Printf` in `internal/proxy/handler.go`
- Replace `fmt.Printf` in `internal/proxy/middleware.go`
- Remove legacy `Debug` flag usage

### Step 1.5: Verification

- Verify JSON output
- Verify File output
- Verify PII masking
