# Implementation Plan: HTTP Reverse Proxy

**Branch**: `001-http-proxy` | **Date**: 2026-01-17 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-http-proxy/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Build a high-performance HTTP reverse proxy in Go that intercepts traffic between clients (editors/CLIs) and AI providers (Anthropic/OpenAI). The proxy must detect provider schemas via URL path matching (`/v1/messages` for Anthropic, `/v1/chat/completions` for OpenAI) and forward all traffic transparently with <5ms overhead. Initial implementation focuses on passthrough mode with schema detection, enabling future caching and modification features.

## Technical Context

**Language/Version**: Go 1.25.6+ (per constitution)  
**Primary Dependencies**:

- Standard library `net/http` for HTTP proxy functionality
- `knadh/koanf/v2` (already in go.mod) for configuration file reading (YAML/TOML)
- Minimal external dependencies - prefer stdlib where possible

**Storage**: N/A for this feature (passthrough only; future features will use DuckDB)  
**Testing**: Standard Go testing (`testing` package) with `net/http/httptest` for HTTP handlers  
**Target Platform**: Cross-platform (Linux, macOS, Windows) - single static binary  
**Project Type**: Single binary CLI/service  
**Performance Goals**:

- <5ms proxy overhead for passthrough requests (p95)
- Support 100+ concurrent connections without degradation
- Streaming response forwarding with <50ms latency to first byte

**Constraints**:

- Must be single static binary (no external runtime deps)
- Must support standard HTTP proxy environment variables (`HTTP_PROXY`, `HTTPS_PROXY`)
- Must handle CONNECT tunneling for HTTPS
- Must support streaming responses (SSE) without buffering
- Configuration must be read from: `memex.yml`, `.memex.yml`, `.config/memex.yml`, `memex.yaml`, `.memex.yaml`, `.config/memex.yaml`, `memex.toml`, `.memex.toml`, `.config/memex.toml` (in that order)

**Scale/Scope**:

- Single binary proxy service
- Handles HTTP/HTTPS traffic forwarding
- Schema detection for 2 provider types (Anthropic, OpenAI)
- No caching or modification logic in this feature

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### ✅ Single Binary Distribution

- **Status**: PASS
- **Rationale**: Using Go stdlib `net/http` for proxy functionality. No external runtime dependencies required. Binary will be statically linked.

### ✅ Go-Only Implementation

- **Status**: PASS
- **Rationale**: Entire implementation in Go 1.25.6+. Using stdlib `net/http` for proxy server. Configuration via existing `koanf` dependency (already in go.mod).

### ✅ Embedded Storage Architecture

- **Status**: N/A (not applicable for passthrough-only feature)
- **Rationale**: This feature does not use storage. Future caching features will use DuckDB per constitution.

### ✅ Protocol Compatibility

- **Status**: PASS
- **Rationale**: Detecting both `/v1/messages` (Anthropic) and `/v1/chat/completions` (OpenAI) schemas. All requests forwarded unchanged (passthrough mode).

### ✅ Performance Requirements

- **Status**: PASS
- **Rationale**: <5ms overhead target aligns with constitution's <10ms requirement for cache lookups. Passthrough should be faster than cache operations.

### ⚠️ Zero Configuration

- **Status**: PARTIAL
- **Rationale**: Basic operation works with environment variables only (no config file required). However, configuration file support is specified for future extensibility. This is acceptable as config files are optional - defaults enable zero-config operation.

**Overall**: ✅ PASS - All critical gates pass. Configuration file support is additive and doesn't violate zero-config principle (config is optional).

## Project Structure

### Documentation (this feature)

```text
specs/001-http-proxy/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
cmd/proxy/
└── main.go              # Entrypoint - HTTP proxy server initialization

internal/proxy/
├── server.go            # HTTP proxy server setup and request handling
├── handler.go           # Request/response forwarding logic
├── schema.go            # Schema detection (Anthropic/OpenAI/Unknown)
├── config.go            # Configuration loading (koanf integration)
└── errors.go            # Error handling and response formatting

pkg/types/
└── proxy.go              # Shared types (SchemaType enum, ProxyRequest, ProxyResponse)

tests/
├── integration/
│   └── proxy_test.go    # Integration tests for proxy behavior
└── unit/
    ├── schema_test.go   # Schema detection tests
    └── handler_test.go # Request forwarding tests
```

**Structure Decision**: Following constitution's directory structure. Proxy logic in `internal/proxy/`, shared types in `pkg/types/`, entrypoint in `cmd/proxy/`. Tests mirror source structure.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No violations - all gates pass.
