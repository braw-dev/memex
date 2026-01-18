# Contracts: HTTP Reverse Proxy Interfaces

**Feature**: HTTP Reverse Proxy  
**Date**: 2026-01-17  
**Phase**: 1 - Design & Contracts

## Overview

This document defines the internal interfaces and contracts for the HTTP reverse proxy. Since this is a proxy service (not an API), these contracts define internal component interfaces and behavioral guarantees.

## Internal Interfaces

### SchemaDetector

**Purpose**: Detects AI provider schema type from HTTP request.

**Interface**:

```go
type SchemaDetector interface {
    Detect(request *http.Request) SchemaType
}
```

**Contract**:

- **Input**: `*http.Request` - HTTP request from client
- **Output**: `SchemaType` - Detected schema (Anthropic, OpenAI, or Unknown)
- **Behavior**:
  - Examines request URL path
  - Returns `SchemaAnthropic` if path ends with `/v1/messages`
  - Returns `SchemaOpenAI` if path ends with `/v1/chat/completions`
  - Returns `SchemaUnknown` otherwise
- **Performance**: Must complete in <1ms (simple string comparison)
- **Thread Safety**: Must be safe for concurrent use

**Implementation**: `internal/proxy/schema.go`

---

### ProxyHandler

**Purpose**: Handles HTTP requests and forwards them to upstream servers.

**Interface**:

```go
type ProxyHandler interface {
    ServeHTTP(w http.ResponseWriter, r *http.Request)
    HandleCONNECT(w http.ResponseWriter, r *http.Request)
}
```

**Contract**:

- **ServeHTTP**: Handles standard HTTP requests (GET, POST, etc.)
  - Detects schema type
  - Forwards request to upstream unchanged
  - Streams response back to client
  - Returns error responses (502, 504) on upstream failures
- **HandleCONNECT**: Handles CONNECT method for HTTPS tunneling
  - Establishes TCP tunnel to upstream
  - Copies bytes bidirectionally
  - Handles client disconnection gracefully
- **Performance**:
  - Passthrough overhead <5ms (p95)
  - Streaming responses forwarded with <50ms latency to first byte
- **Error Handling**:
  - Upstream unreachable → 502 Bad Gateway
  - Upstream timeout → 504 Gateway Timeout
  - Client disconnect → Close connection, cancel upstream

**Implementation**: `internal/proxy/handler.go`

---

### ConfigLoader

**Purpose**: Loads proxy configuration from filesystem or environment.

**Interface**:

```go
type ConfigLoader interface {
    Load() (*ProxyConfig, error)
}
```

**Contract**:

- **Input**: None (reads from filesystem/environment)
- **Output**: `*ProxyConfig` or error
- **Behavior**:
  - Checks config file locations in priority order
  - Uses first found file (YAML or TOML)
  - Applies defaults for missing values
  - Returns error only if config file exists but is invalid
- **File Priority** (first found wins):
  1. `memex.yml`
  2. `.memex.yml`
  3. `.config/memex.yml`
  4. `memex.yaml`
  5. `.memex.yaml`
  6. `.config/memex.yaml`
  7. `memex.toml`
  8. `.memex.toml`
  9. `.config/memex.toml`
- **Environment Overrides**: Environment variables can override config file values (future enhancement)

**Implementation**: `internal/proxy/config.go`

---

## Behavioral Contracts

### Request Forwarding Contract

**Guarantee**: Requests forwarded to upstream are byte-for-byte identical to original request (except for `Host` header which is set to upstream host).

**Exceptions**:

- `Host` header modified to upstream host (required by HTTP spec)
- `Proxy-*` headers removed (standard proxy behavior)
- Connection headers may be modified (connection pooling)

**Verification**: Integration tests compare request bytes before/after forwarding.

---

### Response Forwarding Contract

**Guarantee**: Responses forwarded to client are byte-for-byte identical to upstream response.

**Exceptions**:

- Connection headers may be modified (connection pooling)
- `Transfer-Encoding` preserved for streaming

**Verification**: Integration tests compare response bytes before/after forwarding.

---

### Schema Detection Contract

**Guarantee**: Schema detection is deterministic and based solely on URL path.

**Rules**:

- Path ending with `/v1/messages` → `SchemaAnthropic`
- Path ending with `/v1/chat/completions` → `SchemaOpenAI`
- All other paths → `SchemaUnknown`
- Detection is case-sensitive (matches provider behavior)
- Query parameters ignored (only path matters)

**Verification**: Unit tests cover all path variations.

---

### Performance Contract

**Guarantees**:

- Passthrough overhead: <5ms (p95 latency)
- Concurrent connections: 100+ without degradation
- Streaming latency: <50ms to first byte
- Memory: Bounded (no request/response buffering)

**Measurement**: Benchmark tests measure latency under load.

---

### Error Handling Contract

**Guarantees**:

- Upstream errors return appropriate HTTP status codes
- Client disconnection cancels upstream request (when possible)
- Timeouts prevent hanging connections
- Error responses returned quickly (<5ms even on failure)

**Error Mappings**:

- Upstream unreachable → 502 Bad Gateway
- Upstream timeout → 504 Gateway Timeout
- Invalid request → 400 Bad Request (if detectable before forwarding)
- Client disconnect → Connection closed (no response)

---

## External Contracts

### HTTP Proxy Protocol Contract

**Compliance**: RFC 7230 (HTTP/1.1), RFC 7231 (HTTP Methods)

**Supported Methods**:

- `GET`, `POST`, `PUT`, `DELETE`, `PATCH`, `HEAD`, `OPTIONS` → Forwarded via `ReverseProxy`
- `CONNECT` → Handled via TCP tunneling

**Supported Headers**:

- `Proxy-Authorization`: Handled (if present)
- `Proxy-Connection`: Handled per HTTP spec
- `Host`: Modified to upstream host
- All other headers: Forwarded unchanged

**Environment Variables** (client-side):

- `HTTP_PROXY`: Client uses to find proxy (proxy doesn't read this)
- `HTTPS_PROXY`: Client uses for HTTPS (proxy doesn't read this)
- `NO_PROXY`: Client-side exclusion list (proxy doesn't read this)

---

## Testing Contracts

### Unit Test Contract

**Coverage Requirements**:

- Schema detection: 100% path coverage
- Config loading: All file locations tested
- Error handling: All error paths tested

### Integration Test Contract

**Scenarios**:

- Passthrough requests (non-AI endpoints)
- Anthropic schema detection
- OpenAI schema detection
- CONNECT tunneling (HTTPS)
- Streaming responses (SSE)
- Upstream failures (502, 504)
- Client disconnection

**Performance Tests**:

- Latency benchmarks (<5ms overhead)
- Concurrent connection tests (100+ connections)
- Streaming latency tests (<50ms to first byte)
