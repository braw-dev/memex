# Feature Specification: HTTP Reverse Proxy

**Feature Branch**: `001-http-proxy`  
**Created**: 2026-01-17  
**Status**: Draft  
**Input**: User description: "Implement a high-performance HTTP reverse proxy capable of intercepting and modifying traffic between the client (local editor/CLI) and the upstream provider (Anthropic/OpenAI)."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Transparent Passthrough (Priority: P1)

A developer using Cursor or another AI-enabled editor wants to route their AI requests through Memex without any noticeable change in behavior. The proxy intercepts traffic, determines it doesn't match a known AI provider schema, and forwards it untouched to the upstream server.

**Why this priority**: This is the foundational capability. Without reliable passthrough, the proxy cannot be trusted as infrastructure. Users must have confidence that non-AI traffic and unrecognized requests flow through without modification or delay.

**Independent Test**: Can be fully tested by configuring any HTTP client to use the proxy and verifying requests reach their destination unchanged. Delivers value by proving the proxy is safe to deploy.

**Acceptance Scenarios**:

1. **Given** the proxy is running and configured as the HTTP proxy, **When** a client sends a request to a non-AI endpoint (e.g., `https://api.github.com/repos`), **Then** the request is forwarded unchanged and the response is returned to the client without modification.
2. **Given** the proxy is running, **When** a client sends a request to an AI provider but the path does not match `/v1/messages` or `/v1/chat/completions`, **Then** the request is forwarded unchanged (passthrough mode).
3. **Given** the proxy is running, **When** measuring round-trip time for passthrough requests, **Then** the proxy adds less than 5ms of overhead compared to direct connection.

---

### User Story 2 - Anthropic Schema Detection (Priority: P2)

A developer using Claude via the Anthropic API wants Memex to recognize their requests so that future caching and modification features can be applied. The proxy detects requests to `/v1/messages` and marks them as Anthropic-schema traffic.

**Why this priority**: Anthropic is a primary target provider per the constitution. Detection is required before any caching or modification logic can be applied.

**Independent Test**: Can be tested by sending requests to `/v1/messages` endpoint and verifying the proxy correctly identifies the schema type (via logs or response headers in debug mode).

**Acceptance Scenarios**:

1. **Given** the proxy is running, **When** a client sends a POST request to `*/v1/messages`, **Then** the proxy identifies this as Anthropic-schema traffic.
2. **Given** the proxy detects Anthropic-schema traffic, **When** no cache or modification logic applies, **Then** the request is forwarded to the upstream provider unchanged.
3. **Given** the proxy is running in debug mode, **When** an Anthropic request is processed, **Then** the detected schema type is logged or indicated in response metadata.

---

### User Story 3 - OpenAI Schema Detection (Priority: P2)

A developer using GPT models via the OpenAI API wants Memex to recognize their requests so that future caching and modification features can be applied. The proxy detects requests to `/v1/chat/completions` and marks them as OpenAI-schema traffic.

**Why this priority**: OpenAI is a primary target provider per the constitution. Detection enables future middleware chain features.

**Independent Test**: Can be tested by sending requests to `/v1/chat/completions` endpoint and verifying the proxy correctly identifies the schema type.

**Acceptance Scenarios**:

1. **Given** the proxy is running, **When** a client sends a POST request to `*/v1/chat/completions`, **Then** the proxy identifies this as OpenAI-schema traffic.
2. **Given** the proxy detects OpenAI-schema traffic, **When** no cache or modification logic applies, **Then** the request is forwarded to the upstream provider unchanged.
3. **Given** the proxy is running in debug mode, **When** an OpenAI request is processed, **Then** the detected schema type is logged or indicated in response metadata.

---

### User Story 4 - Standard Proxy Configuration (Priority: P3)

A developer wants to configure their tools to use Memex as an HTTP proxy using standard environment variables and proxy headers. They should not need custom configuration beyond what any HTTP proxy requires.

**Why this priority**: Standard proxy support enables broad compatibility with existing tools (editors, CLIs, SDKs) without requiring tool-specific integrations.

**Independent Test**: Can be tested by setting `HTTP_PROXY`/`HTTPS_PROXY` environment variables and verifying traffic routes through Memex.

**Acceptance Scenarios**:

1. **Given** a client application respects `HTTP_PROXY` and `HTTPS_PROXY` environment variables, **When** these variables point to the Memex proxy address, **Then** HTTP and HTTPS traffic is routed through Memex.
2. **Given** the proxy is running, **When** a client sends a request with standard proxy headers (`Proxy-Authorization`, `Proxy-Connection`), **Then** the proxy handles these headers appropriately.
3. **Given** the proxy is configured to listen on a specific port, **When** the user sets `HTTPS_PROXY=http://localhost:<port>`, **Then** HTTPS traffic is correctly proxied via CONNECT tunneling or TLS termination.

---

### Edge Cases

- What happens when the upstream provider is unreachable? The proxy MUST return an appropriate error response to the client without hanging indefinitely.
- What happens when a request body is malformed JSON? The proxy MUST forward the request unchanged (passthrough) and let the upstream provider handle validation.
- What happens when the client disconnects mid-request? The proxy MUST clean up resources and cancel the upstream request if possible.
- What happens when the upstream response is a streaming response (SSE)? The proxy MUST forward the stream to the client in real-time without buffering the entire response.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST accept incoming HTTP/HTTPS connections and forward them to upstream servers.
- **FR-002**: System MUST detect Anthropic-schema requests by matching URL path `/v1/messages`.
- **FR-003**: System MUST detect OpenAI-schema requests by matching URL path `/v1/chat/completions`.
- **FR-004**: System MUST forward all non-matching requests unchanged (passthrough mode).
- **FR-005**: System MUST support standard HTTP proxy environment variables (`HTTP_PROXY`, `HTTPS_PROXY`, `NO_PROXY`).
- **FR-006**: System MUST handle standard HTTP proxy headers (`Proxy-Authorization`, `Proxy-Connection`, etc.).
- **FR-007**: System MUST support HTTPS traffic via CONNECT tunneling.
- **FR-008**: System MUST support streaming responses (Server-Sent Events) without buffering.
- **FR-009**: System MUST handle connection timeouts gracefully and return appropriate error responses.
- **FR-010**: System MUST clean up resources when client connections are terminated unexpectedly.
- **FR-011**: System MUST add less than 5ms of latency overhead for passthrough requests.

### Key Entities

- **ProxyRequest**: Represents an incoming request with metadata about detected schema type (Anthropic, OpenAI, or Unknown), original headers, body reference, and upstream target.
- **ProxyResponse**: Represents the upstream response being forwarded to the client, including status code, headers, and body stream.
- **SchemaType**: Enumeration of recognized AI provider schemas (Anthropic, OpenAI, Unknown) used for routing decisions.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Passthrough requests complete with less than 5ms additional latency compared to direct connection (measured at p95).
- **SC-002**: 100% of non-AI requests pass through unchanged (byte-for-byte identical to direct requests).
- **SC-003**: Proxy correctly identifies Anthropic schema for all requests to `/v1/messages` path.
- **SC-004**: Proxy correctly identifies OpenAI schema for all requests to `/v1/chat/completions` path.
- **SC-005**: Streaming responses begin delivery to client within 50ms of first upstream byte received.
- **SC-006**: Proxy handles 100 concurrent connections without degradation.
- **SC-007**: Users can configure the proxy using only standard environment variables (no custom config files required for basic operation).

## Assumptions

- The proxy will initially operate in passthrough mode only; caching and modification features will be added in subsequent features.
- HTTPS interception (for schema detection on encrypted traffic) requires the proxy to perform TLS termination, which may require certificate configuration. For initial implementation, schema detection works on unencrypted traffic or traffic where the proxy terminates TLS.
- The 5ms latency budget applies to the proxy's processing overhead, not network latency to upstream providers.
- Standard HTTP proxy semantics apply (RFC 7230, RFC 7231).
