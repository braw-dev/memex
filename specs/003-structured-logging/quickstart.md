# Quickstart: Structured Logging Configuration

This guide explains how to configure the new structured logging system.

## Configuration

You can configure logging via the configuration file (YAML/TOML) or environment variables.

### Config File (memex.yaml)

```yaml
proxy:
  listen: ":8080"
  log:
    level: "debug"  # Options: debug, info, warn, error
    format: "json"  # Options: text, json
    path: "stderr"  # Options: stderr, stdout, or /path/to/file.log
```

### Environment Variables

Environment variables override configuration file settings. Use `.` as a separator for nested keys.

```bash
# Set log level to debug
export MEMEX_PROXY_LOG_LEVEL="debug"

# Set output format to JSON
export MEMEX_PROXY_LOG_FORMAT="json"

# Write logs to a file
export MEMEX_PROXY_LOG_PATH="./app.log"
```

## Usage in Code

Do not use `fmt.Printf` or `log.Println`. Use the global `slog` functions.

```go
import "log/slog"

func someFunction() {
    // Info log
    slog.Info("server starting", "port", 8080)

    // Error log with error object
    if err := doSomething(); err != nil {
        slog.Error("operation failed", "err", err, "attempt", 1)
    }

    // Debug log (only shows if level is debug)
    slog.Debug("processing request", "request_id", "123")
}
```

## PII Masking

Sensitive data should be wrapped in types that implement `LogValuer` or manually redacted before logging.

```go
// Good - PII is masked by the type implementation
slog.Info("user login", "user", userObject)
```
