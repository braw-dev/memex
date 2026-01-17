# Memex Agent Guidance

This file provides runtime development guidance for AI agents working on the Memex project.

## Constitution Reference

All development MUST comply with the project constitution at `.specify/memory/constitution.md`.

Key constraints to verify before any implementation:

1. **Single Binary**: No external dependencies, Docker, or runtime downloads
2. **Go Only**: All code in Go 1.25.6+, no other languages
3. **Embedded Storage**: DuckDB for OLAP, BadgerDB for KV - no external databases
4. **Performance**: Proxy overhead <10ms
5. **Safety**: PII blocking is P0 - never persist sensitive data

## Directory Structure

```
cmd/proxy/          → Main entrypoint only
internal/brain/     → Tree-Sitter parsing, ONNX embeddings
internal/proxy/     → HTTP handlers, middleware chain
internal/store/     → DuckDB and BadgerDB wrappers
pkg/types/          → Shared structs and interfaces
assets/             → Embedded model weights (go:embed)
```

## Development Workflow

1. Check constitution compliance BEFORE starting work
2. Follow the middleware chain order (Auth → Magic → PII → Cache → Audit)
3. Use tri-partite cache keys (Context + System + Intent)
4. Test with both Anthropic and OpenAI schemas

## Specification System

Use the `.specify/` templates for feature planning:

- `/speckit.constitution` - Update project principles
- `/speckit.specify` - Create feature specifications
- `/speckit.plan` - Generate implementation plans
- `/speckit.tasks` - Break down into actionable tasks
