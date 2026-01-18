# Research: Structured Logging

**Feature**: Structured Logging with slog
**Created**: 2026-01-18
**Status**: Complete

## Topic: PII Masking Mechanism

**Decision**: Use `slog.LogValuer` interface for types containing sensitive data.

**Rationale**:

- `slog` provides the `LogValuer` interface specifically for customizing how a type is marshaled to a log attribute.
- This allows us to implement the `Value()` method on sensitive types (e.g., `UserRequest`) to return a redacted `slog.Value`.
- It ensures that PII masking is tied to the type definition, preventing accidental leakage if a developer logs the struct directly.
- It is zero-allocation if the log level is disabled.

**Alternatives Considered**:

- **Middleware-based scrubbing**: Parsing log strings to replace patterns. *Rejected*: Slow, error-prone, and works on finalized output rather than structured data.
- **Custom Wrapper Functions**: `log.Sensitive("key", val)`. *Rejected*: Relies on developer discipline to use the wrapper every time.

## Topic: Log Configuration Structure

**Decision**: Introduce a `LogConfig` struct nested within `ProxyConfig`.

**Rationale**:

- Keeps logging configuration logically grouped.
- Allows for easy expansion (e.g., adding sampling rate later).
- Maps cleanly to JSON/YAML configuration files.

**Alternatives Considered**:

- **Top-level config fields**: `LogLevel`, `LogFormat`. *Rejected*: Clutters the root config object.

## Topic: Standard Library vs Third Party

**Decision**: Use Go standard library `log/slog` exclusively.

**Rationale**:

- **Constitution Compliance**: The project constitution emphasizes "Single Binary" and minimizing dependencies. `slog` is built-in (Go 1.21+).
- **Performance**: `slog` is highly optimized for low allocation.
- **Interoperability**: It is the standard interface that other libraries are now adapting to.

**Alternatives Considered**:

- **Zap/Zerolog**: *Rejected*: Adds external dependencies, contrary to project principles.
