# Data Model: HTTP Reverse Proxy

**Feature**: HTTP Reverse Proxy  
**Date**: 2026-01-17  
**Phase**: 1 - Design & Contracts

## Entities

### SchemaType

**Type**: Enumeration  
**Purpose**: Represents the detected AI provider schema type for a request.

**Values**:

- `SchemaUnknown` (0): Request doesn't match any known AI provider schema
- `SchemaAnthropic` (1): Request matches Anthropic API schema (`/v1/messages`)
- `SchemaOpenAI` (2): Request matches OpenAI API schema (`/v1/chat/completions`)

**Usage**: Used for routing decisions and future middleware chain integration. Stored in request context during proxy handling.

**Validation**: No validation needed - enum type ensures only valid values.

---

### ProxyRequest

**Type**: Struct (internal representation)  
**Purpose**: Encapsulates incoming HTTP request with proxy-specific metadata.

**Fields**:

- `Request *http.Request`: Original HTTP request from client
- `Schema SchemaType`: Detected provider schema (from path matching)
- `UpstreamURL *url.URL`: Target upstream URL (parsed from request)
- `StartTime time.Time`: Request start timestamp (for latency measurement)

**Relationships**:

- Contains `*http.Request` (standard library type)
- Contains `SchemaType` enum
- Used by `ProxyResponse` for correlation

**State Transitions**: None - this is a value object representing a single request.

**Validation Rules**:

- `Request` must not be nil
- `UpstreamURL` must be valid URL (parsed from request)
- `Schema` must be valid enum value

---

### ProxyResponse

**Type**: Struct (internal representation)  
**Purpose**: Represents upstream response being forwarded to client.

**Fields**:

- `StatusCode int`: HTTP status code from upstream
- `Headers http.Header`: Response headers from upstream
- `Body io.ReadCloser`: Response body stream (for forwarding)
- `Request *ProxyRequest`: Reference to original request (for correlation)
- `Duration time.Duration`: Time taken to receive response from upstream

**Relationships**:

- References `ProxyRequest` (one-to-one)
- Contains `http.Header` (standard library type)

**State Transitions**: None - this is a value object representing a single response.

**Validation Rules**:

- `StatusCode` must be valid HTTP status code (100-599)
- `Headers` must not be nil
- `Body` may be nil for HEAD requests or empty responses

---

### ProxyConfig

**Type**: Struct (configuration)  
**Purpose**: Configuration for proxy server behavior.

**Fields**:

- `ListenAddr string`: Address to listen on (e.g., `:8080`, `127.0.0.1:8080`)
- `UpstreamTimeout time.Duration`: Timeout for upstream requests (default: 60s)
- `IdleTimeout time.Duration`: Idle connection timeout (default: 90s)
- `FlushInterval time.Duration`: Response flush interval for streaming (default: 0 for immediate)
- `Debug bool`: Enable debug logging (default: false)

**Relationships**: None - standalone configuration object.

**State Transitions**: None - configuration is immutable after loading.

**Validation Rules**:

- `ListenAddr` must be valid network address format
- `UpstreamTimeout` must be > 0
- `IdleTimeout` must be > 0
- `FlushInterval` must be >= 0

**Default Values**:

- `ListenAddr`: `:8080` (if not specified)
- `UpstreamTimeout`: `60s`
- `IdleTimeout`: `90s`
- `FlushInterval`: `0` (immediate flush for streaming)
- `Debug`: `false`

---

## Relationships

```text
ProxyRequest
  ├── contains SchemaType (enum)
  ├── contains *http.Request (stdlib)
  └── referenced by ProxyResponse

ProxyResponse
  ├── references ProxyRequest
  └── contains http.Header, io.ReadCloser (stdlib)

ProxyConfig
  └── standalone (no relationships)
```

## Notes

- **No Persistent Storage**: This feature is passthrough-only. No data is persisted. Future caching features will introduce storage entities (BadgerDB keys, DuckDB records).

- **Request Context**: `SchemaType` is stored in `request.Context()` using `context.WithValue()` to pass metadata through proxy chain without modifying request headers.

- **Streaming**: `ProxyResponse.Body` is an `io.ReadCloser` to support streaming responses. The proxy forwards bytes as they arrive from upstream without buffering.

- **Configuration Loading**: `ProxyConfig` is loaded from filesystem using `koanf` (priority order specified in research.md). Environment variables can override config file values.
