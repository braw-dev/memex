# Research: HTTP Reverse Proxy Implementation

**Feature**: HTTP Reverse Proxy  
**Date**: 2026-01-17  
**Phase**: 0 - Outline & Research

## Research Questions

### 1. Go HTTP Proxy Implementation Patterns

**Question**: What's the simplest, most performant way to implement an HTTP reverse proxy in Go using stdlib?

**Decision**: Use `net/http/httputil.ReverseProxy` as the foundation with custom `Director` and `ModifyResponse` functions.

**Rationale**:

- `httputil.ReverseProxy` is part of Go stdlib - zero external dependencies
- Handles connection pooling, request/response forwarding automatically
- Supports streaming responses out of the box
- `Director` function allows request modification before forwarding
- `ModifyResponse` allows response modification after receiving from upstream
- Well-tested, production-ready code maintained by Go team

**Alternatives Considered**:

- Custom `net/http` server with manual request forwarding: More code, more error-prone, no connection pooling
- Third-party proxy libraries (e.g., `elazarl/goproxy`): Adds external dependency, violates minimal dependency principle
- `goproxy` or similar: Overkill for passthrough mode, adds complexity

**References**:

- Go stdlib `httputil.ReverseProxy` documentation
- Go blog: "Reverse Proxy Pattern" (<https://go.dev/blog/reverse-proxy>)

---

### 2. HTTPS/CONNECT Tunneling Support

**Question**: How to handle HTTPS traffic through proxy (CONNECT method)?

**Decision**: Implement CONNECT method handler separately using `net/http` server with custom handler that establishes TCP tunnel.

**Rationale**:

- CONNECT requires raw TCP connection handling, not HTTP request/response
- `httputil.ReverseProxy` doesn't handle CONNECT - need separate handler
- Standard pattern: Accept CONNECT, establish TCP connection to upstream, copy bytes bidirectionally
- Use `net.Dial` for upstream connection, `io.Copy` for bidirectional streaming

**Implementation Pattern**:

```go
// Pseudo-code structure
func handleCONNECT(w http.ResponseWriter, r *http.Request) {
    // Extract target host:port from request
    // Establish TCP connection to upstream
    // Hijack client connection
    // Copy bytes bidirectionally
}
```

**Alternatives Considered**:

- TLS termination at proxy: Requires certificate management, adds complexity, violates transparency
- Third-party CONNECT handlers: Adds dependency, violates minimal dependency principle

---

### 3. Schema Detection Strategy

**Question**: Where and how to detect provider schema (Anthropic vs OpenAI)?

**Decision**: Detect in `Director` function of `ReverseProxy` by parsing request URL path. Store schema type in request context for downstream use.

**Rationale**:

- URL path is available before request forwarding
- Path matching is O(1) operation (string comparison)
- No body parsing required (faster, simpler)
- Context allows passing metadata without modifying request headers
- Matches spec requirement: detect by path (`/v1/messages` vs `/v1/chat/completions`)

**Implementation**:

```go
func detectSchema(path string) SchemaType {
    if strings.HasSuffix(path, "/v1/messages") {
        return SchemaAnthropic
    }
    if strings.HasSuffix(path, "/v1/chat/completions") {
        return SchemaOpenAI
    }
    return SchemaUnknown
}
```

**Alternatives Considered**:

- Body parsing: Slower, requires buffering, violates passthrough transparency
- Header inspection: Less reliable, providers may not set unique headers
- Hostname matching: Too broad, doesn't distinguish endpoints

---

### 4. Configuration File Loading

**Question**: How to load configuration from multiple file locations with priority order?

**Decision**: Use existing `knadh/koanf/v2` dependency (already in go.mod) with filesystem provider. Check locations in priority order, use first found file.

**Rationale**:

- `koanf` already in project dependencies
- Supports YAML, TOML, JSON formats
- Filesystem provider allows checking multiple paths
- Can merge multiple config sources (if needed later)
- Minimal code - just configure provider paths

**File Priority Order** (per user requirement):

1. `memex.yml`
2. `.memex.yml`
3. `.config/memex.yml`
4. `memex.yaml`
5. `.memex.yaml`
6. `.config/memex.yaml`
7. `memex.toml`
8. `.memex.toml`
9. `.config/memex.toml`

**Implementation Pattern**:

```go
k := koanf.New(".")
for _, path := range configPaths {
    if exists(path) {
        k.Load(filesystem.Provider(path), yaml.Parser())
        break
    }
}
```

**Alternatives Considered**:

- Custom file reading: More code, need to handle YAML/TOML parsing
- Viper: Larger dependency, overkill for simple config
- Environment variables only: Doesn't meet requirement for config files

---

### 5. Performance Optimization for <5ms Overhead

**Question**: How to minimize proxy overhead to meet <5ms latency budget?

**Decision**:

- Use `httputil.ReverseProxy` (optimized C code in stdlib)
- Avoid request/response body buffering (stream directly)
- Minimize allocations in hot path (reuse buffers)
- Schema detection: simple string comparison (no regex, no parsing)
- No logging in hot path (async logging if needed)

**Rationale**:

- `ReverseProxy` uses efficient `io.Copy` with buffer pooling
- Streaming avoids memory allocation for large bodies
- String comparison is O(1) - fastest possible detection
- Logging can add 1-2ms overhead - defer to async or debug mode only

**Optimization Techniques**:

- Pre-compile path patterns (avoid repeated string operations)
- Use `strings.HasSuffix` instead of regex (faster)
- Reuse HTTP client with connection pooling (automatic in ReverseProxy)
- Avoid JSON unmarshaling in passthrough mode

**Alternatives Considered**:

- Custom proxy implementation: More code, likely slower than stdlib
- Request buffering for inspection: Adds latency, violates passthrough transparency
- Synchronous logging: Adds latency, not needed for passthrough

---

### 6. Streaming Response Handling

**Question**: How to forward Server-Sent Events (SSE) and streaming responses without buffering?

**Decision**: Use `ReverseProxy` default behavior - it streams responses automatically. Ensure `FlushInterval` is set appropriately for SSE.

**Rationale**:

- `ReverseProxy` streams by default (doesn't buffer entire response)
- `FlushInterval` controls how often to flush to client
- For SSE, set `FlushInterval: 0` to flush immediately on each chunk
- `io.Copy` handles streaming efficiently

**Configuration**:

```go
proxy := &httputil.ReverseProxy{
    Director: director,
    FlushInterval: 0, // Immediate flush for SSE
}
```

**Alternatives Considered**:

- Manual streaming: More complex, error-prone, ReverseProxy already handles it
- Buffering then forwarding: Violates streaming requirement, adds latency

---

### 7. Error Handling and Timeouts

**Question**: How to handle upstream failures and timeouts gracefully?

**Decision**: Use `http.Client` with configurable timeouts. Return appropriate HTTP error responses (502 Bad Gateway, 504 Gateway Timeout).

**Rationale**:

- `ReverseProxy` uses `http.Client` internally - can configure timeouts
- Standard HTTP error codes communicate failure type to client
- Timeout prevents hanging connections
- Error responses should be returned quickly (<5ms even on failure)

**Configuration**:

```go
transport := &http.Transport{
    ResponseHeaderTimeout: 30 * time.Second,
    IdleConnTimeout: 90 * time.Second,
}
client := &http.Client{
    Transport: transport,
    Timeout: 60 * time.Second,
}
```

**Error Response Pattern**:

- Upstream unreachable: 502 Bad Gateway
- Upstream timeout: 504 Gateway Timeout
- Client disconnect: Close connection, cancel upstream request

---

### 8. Standard Proxy Environment Variables

**Question**: How to support `HTTP_PROXY`, `HTTPS_PROXY`, `NO_PROXY`?

**Decision**: These are CLIENT-side environment variables (used by clients to FIND the proxy). The proxy itself doesn't need to read them. However, we should document that clients set these to point TO our proxy.

**Rationale**:

- `HTTP_PROXY=http://localhost:8080` tells client to use our proxy
- Proxy server itself doesn't read these vars
- We DO need to handle `Proxy-Authorization` header if clients send it
- `NO_PROXY` is client-side exclusion list

**Implementation**:

- Proxy reads `Proxy-Authorization` header if present
- Proxy handles `Proxy-Connection` header per HTTP spec
- Documentation explains how clients configure proxy via env vars

**Alternatives Considered**:

- Reading env vars in proxy: Misunderstanding - these are for clients, not proxy server

---

## Summary of Technical Decisions

1. **HTTP Proxy**: `net/http/httputil.ReverseProxy` (stdlib)
2. **HTTPS/CONNECT**: Custom CONNECT handler with TCP tunneling
3. **Schema Detection**: Path string matching in `Director` function
4. **Configuration**: `koanf/v2` filesystem provider (already in dependencies)
5. **Performance**: Stream responses, minimize allocations, no logging in hot path
6. **Streaming**: Use `FlushInterval: 0` for SSE
7. **Errors**: Standard HTTP error codes with configurable timeouts
8. **Proxy Headers**: Handle `Proxy-Authorization`, `Proxy-Connection` per HTTP spec

**Dependencies Added**: None (using stdlib + existing koanf)

**Dependencies Avoided**:

- No third-party proxy libraries
- No HTTP frameworks (using stdlib)
- No logging libraries (use stdlib `log` or defer to future)
