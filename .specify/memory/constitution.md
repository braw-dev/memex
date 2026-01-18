<!--
SYNC IMPACT REPORT
==================
Version Change: 0.0.0 → 1.0.0 (MAJOR - Initial constitution ratification)

Modified Principles: N/A (new constitution)

Added Sections:
  - Core Principles (6 principles)
  - Technical Invariants
  - Architectural Patterns
  - User Experience Laws
  - Directory Structure
  - Strategic Alignment
  - Governance

Removed Sections: N/A (new constitution)

Templates Requiring Updates:
  - .specify/templates/plan-template.md ✅ Compatible (Constitution Check section exists)
  - .specify/templates/spec-template.md ✅ Compatible (technology-agnostic)
  - .specify/templates/tasks-template.md ✅ Compatible (phase structure aligns)

Follow-up TODOs: None
==================
-->

# Memex Constitution

## Mission Directive

Build "Memex": A single-binary, local-first "Universal Brain" for AI coding tools.

**Goal**: Create an indispensable infrastructure layer for Enterprise AI adoption.

## Core Principles

### I. Single Binary Distribution

All functionality MUST be delivered as a single static binary with no external runtime dependencies.

- **No Docker** containers required for deployment
- **No Python** or other interpreter dependencies
- **No model downloads** at runtime - all assets embedded via `//go:embed`
- Users MUST be able to run `memex init` and have a working system immediately

**Rationale**: Enterprise adoption requires minimal operational overhead. A single binary eliminates dependency management, version conflicts, and deployment complexity.

### II. Go-Only Implementation

The entire codebase MUST be written in Go (Golang) 1.25.6+.

- All components, libraries, and tools MUST use Go
- Third-party dependencies MUST have Go bindings (no CGO shelling out to external processes)
- Build process MUST produce a statically-linked binary

**Rationale**: Go's compilation model, cross-platform support, and static linking capabilities directly enable the Single Binary Distribution principle.

### III. Embedded Storage Architecture

All persistent storage MUST use embedded DuckDB for all storage needs (OLAP, KV, and Cache).

- **Unified Storage**: Embedded DuckDB via `go-duckdb` for analytical queries, audit logging, and key-value cache operations.
- **No BadgerDB**: BadgerDB is deprecated and MUST NOT be used.
- **No external database servers** (PostgreSQL, Redis, etc.) permitted.

**Rationale**: Consolidating on DuckDB reduces codebase complexity, simplifies the single-binary distribution, and provides superior analytical capabilities for audit logs while still supporting efficient key-value lookups via primary keys.

### IV. Tri-Partite Cache Key System

Every cache lookup MUST use three distinct signals to prevent collisions and ensure semantic accuracy:

1. **Context Hash (Strict)**: Tree-Sitter AST hash of input files (ignores whitespace/comments)
2. **System Hash (Strict)**: SHA-256 of the System Prompt
3. **Intent Vector (Fuzzy)**: Cosine similarity (>0.97 threshold, configurable) of User Prompt using embedded ONNX model (`all-MiniLM-L6-v2`)

**Rationale**: Single-dimension caching leads to false positives. The tri-partite approach ensures cache hits only occur when context, system instructions, AND user intent all match.

### V. Middleware Chain Order

The HTTP Proxy MUST process requests in this exact sequence:

1. **Auth/Scope**: Detect Git Repo URL to salt the cache key
2. **Magic Command**: Intercept `!reset` or `!bust` to clear cache
3. **PII Scrubber**: Block requests containing AWS/Stripe keys (Regex-based, extendable)
4. **Cache Lookaside**: Check DuckDB `cache_entries` using Tri-Partite Key
5. **Audit Logger**: Async write to DuckDB `audit_logs` (Shadow Billing)

**Rationale**: Order matters for security and correctness. PII blocking MUST occur before any caching or logging to prevent sensitive data persistence.

### VI. Protocol Compatibility

The proxy MUST support both major AI provider schemas transparently:

- **Anthropic**: `/v1/messages` endpoint schema
- **OpenAI**: `/v1/chat/completions` endpoint schema

Clients MUST NOT need to modify their existing integrations to use Memex.

**Rationale**: Enterprise adoption requires drop-in compatibility with existing tooling. Supporting both schemas maximizes addressable market.

## Technical Invariants

These constraints are NON-NEGOTIABLE and MUST NOT be violated under any circumstances:

| Invariant | Requirement |
|-----------|-------------|
| Language | Go 1.25.6+ |
| Distribution | Single static binary |
| OLAP Storage | Embedded DuckDB (`go-duckdb`) |
| KV Storage | Embedded DuckDB (`go-duckdb`) |
| Embeddings | `all-MiniLM-L6-v2` (ONNX) via `//go:embed` |
| Parsing | Tree-Sitter (Go bindings) |
| Protocol | HTTP Proxy (Anthropic + OpenAI schemas) |
| Testability | No `os.Setenv` in tests (inject `getenv`) |

## User Experience Laws

### Zero Configuration

- The binary MUST work immediately after `memex init`
- No model downloads, no config file editing, no environment setup required
- Sensible defaults MUST cover 90% of use cases

### Cache Transparency

- Cache hits MUST inject a Markdown footer: `> ⚡ Cached. [Force Regenerate](link)`
- Users MUST always know when they're receiving cached responses
- The "bust" mechanism MUST be discoverable and easy to use

### Cost Visibility

- CLI output MUST show "Money Saved" per session
- Users MUST have visibility into cache hit rates and estimated savings
- Transparency builds trust and demonstrates value

## Directory Structure

```plaintext
/memex-ai-caching-memory
  ├── cmd/proxy          # Main entrypoint
  ├── internal/
  │   ├── brain/         # Tree-Sitter & ONNX logic
  │   ├── proxy/         # HTTP Handlers & Middleware
  │   └── store/         # DuckDB wrappers
  ├── pkg/types/         # Shared Structs
  └── assets/            # Embedded Model Weights
```

All new code MUST follow this structure. Deviations require explicit justification and constitution amendment.

## Strategic Alignment

### Reliability First

- Prioritize exact matches over fuzzy matches to prevent code hallucination
- When in doubt, cache miss is safer than false positive
- Fuzzy matching threshold (0.97) MUST be conservative

### Performance Requirements

- Proxy overhead MUST be <10ms for cache lookups
- Cache operations MUST NOT block request processing beyond this threshold
- Async logging MUST NOT impact response latency

### Safety as P0

- PII blocking is a P0 (highest priority) feature
- AWS keys, Stripe keys, and other secrets MUST be detected and blocked
- Regex patterns MUST be extendable for organization-specific requirements
- No sensitive data MUST ever be persisted to cache or logs

## Governance

### Amendment Process

1. Proposed changes MUST be documented with rationale
2. Changes to Technical Invariants require MAJOR version bump
3. All amendments MUST update the `LAST_AMENDED_DATE`
4. Breaking changes MUST include migration guidance

### Compliance Verification

- All PRs MUST verify compliance with this constitution
- Constitution Check in plan templates MUST pass before implementation
- Violations MUST be justified in Complexity Tracking section

### Version Policy

- **MAJOR**: Backward-incompatible changes to principles or invariants
- **MINOR**: New principles, sections, or materially expanded guidance
- **PATCH**: Clarifications, wording improvements, typo fixes

**Version**: 1.0.1 | **Ratified**: 2026-01-17 | **Last Amended**: 2026-01-18
