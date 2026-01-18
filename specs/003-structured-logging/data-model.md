# Data Model: Structured Logging

**Feature**: Structured Logging with slog
**Created**: 2026-01-18
**Status**: Draft

## Configuration Entities

### LogConfig

Configuration for the logging subsystem.

| Field | Type | Description | Default | Validation |
|-------|------|-------------|---------|------------|
| `Level` | `string` | Minimum log level severity | `"info"` | One of: `debug`, `info`, `warn`, `error` |
| `Format` | `string` | Output format for logs | `"text"` | One of: `text`, `json` |
| `Path` | `string` | Destination for log output | `"stderr"` | `stderr`, `stdout`, or valid file path |

### ProxyConfig (Updated)

Updates to the main configuration struct.

| Field | Type | Description | Change |
|-------|------|-------------|--------|
| `Debug` | `bool` | Legacy debug flag | **Removed** (Replaced by `Log.Level`) |
| `Log` | `LogConfig` | Nested logging configuration | **Added** |

## Type Definitions

### SensitiveData (Interface Implementation)

Types that implement `slog.LogValuer` for PII redaction.

```go
// Example of a sensitive type implementing LogValuer
type SensitiveString string

func (s SensitiveString) LogValue() slog.Value {
    return slog.StringValue("***REDACTED***")
}
```
