# Quickstart: HTTP Reverse Proxy

**Feature**: HTTP Reverse Proxy  
**Date**: 2026-01-17  
**Phase**: 1 - Design & Contracts

## Prerequisites

- Go 1.25.6+ installed
- Memex binary built (or `go run cmd/proxy/main.go`)

## Basic Usage

### 1. Start the Proxy

```bash
# Default: listens on :8080
./memex proxy

# Custom port
./memex proxy --listen :9090

# With config file
./memex proxy --config .config/memex.yml
```

### 2. Configure Your Client

Set environment variables to route traffic through the proxy:

```bash
# HTTP traffic
export HTTP_PROXY=http://localhost:8080

# HTTPS traffic
export HTTPS_PROXY=http://localhost:8080

# Exclude specific hosts (optional)
export NO_PROXY=localhost,127.0.0.1,.local
```

### 3. Test Passthrough

```bash
# Test with curl
curl -x http://localhost:8080 https://api.github.com/repos/octocat/Hello-World

# Test with HTTP client that respects HTTP_PROXY
HTTP_PROXY=http://localhost:8080 curl https://api.github.com/repos/octocat/Hello-World
```

## Configuration

### Configuration File Locations

The proxy checks for configuration files in this order (first found wins):

1. `memex.yml`
2. `.memex.yml`
3. `.config/memex.yml`
4. `memex.yaml`
5. `.memex.yaml`
6. `.config/memex.yaml`
7. `memex.toml`
8. `.memex.toml`
9. `.config/memex.toml`

### Configuration Format (YAML)

```yaml
# .config/memex.yml
proxy:
  listen: ":8080"              # Address to listen on
  upstream_timeout: "60s"      # Timeout for upstream requests
  idle_timeout: "90s"          # Idle connection timeout
  flush_interval: "0s"          # Response flush interval (0 = immediate)
  debug: false                  # Enable debug logging
```

### Configuration Format (TOML)

```toml
# .config/memex.toml
[proxy]
listen = ":8080"
upstream_timeout = "60s"
idle_timeout = "90s"
flush_interval = "0s"
debug = false
```

### Default Configuration

If no config file is found, defaults are used:

- `listen`: `:8080`
- `upstream_timeout`: `60s`
- `idle_timeout`: `90s`
- `flush_interval`: `0s` (immediate flush for streaming)
- `debug`: `false`

## Schema Detection

The proxy automatically detects AI provider schemas:

- **Anthropic**: Requests to `*/v1/messages` → Detected as `SchemaAnthropic`
- **OpenAI**: Requests to `*/v1/chat/completions` → Detected as `SchemaOpenAI`
- **Other**: All other requests → `SchemaUnknown` (passthrough)

### Testing Schema Detection

```bash
# Anthropic request (will be detected)
curl -x http://localhost:8080 \
  -X POST https://api.anthropic.com/v1/messages \
  -H "Content-Type: application/json" \
  -d '{"model":"claude-3-opus","messages":[]}'

# OpenAI request (will be detected)
curl -x http://localhost:8080 \
  -X POST https://api.openai.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{"model":"gpt-4","messages":[]}'

# Non-AI request (passthrough)
curl -x http://localhost:8080 https://api.github.com/user
```

## HTTPS/CONNECT Tunneling

The proxy supports HTTPS via CONNECT tunneling:

```bash
# HTTPS request through proxy
HTTPS_PROXY=http://localhost:8080 curl https://api.anthropic.com/v1/messages
```

The proxy establishes a TCP tunnel for HTTPS traffic, forwarding encrypted bytes without inspection.

## Debug Mode

Enable debug logging to see schema detection and request details:

```bash
# Via config file
echo 'proxy: { debug: true }' > .memex.yml
./memex proxy

# Via command line (if supported)
./memex proxy --debug
```

Debug output shows:

- Detected schema type for each request
- Upstream target URL
- Request/response timing
- Error details

## Integration with Editors/CLIs

### Cursor Editor

1. Set environment variables:

   ```bash
   export HTTP_PROXY=http://localhost:8080
   export HTTPS_PROXY=http://localhost:8080
   ```

2. Restart Cursor

3. Cursor's AI requests will route through Memex proxy

### OpenAI CLI

```bash
export HTTPS_PROXY=http://localhost:8080
openai api chat_completions.create --model gpt-4 --messages '[{"role":"user","content":"Hello"}]'
```

### Anthropic CLI

```bash
export HTTPS_PROXY=http://localhost:8080
anthropic messages create --model claude-3-opus --messages '[{"role":"user","content":"Hello"}]'
```

## Performance

Expected performance characteristics:

- **Passthrough overhead**: <5ms (p95)
- **Concurrent connections**: 100+ without degradation
- **Streaming latency**: <50ms to first byte
- **Memory**: Bounded (no buffering)

## Troubleshooting

### Proxy Not Receiving Requests

1. Verify proxy is running: `curl http://localhost:8080` (should return error or proxy response)
2. Check environment variables: `echo $HTTP_PROXY $HTTPS_PROXY`
3. Verify client respects proxy env vars (some tools require explicit `--proxy` flag)

### High Latency

1. Check upstream provider latency (proxy adds <5ms)
2. Verify network connectivity to upstream
3. Enable debug mode to see timing details

### CONNECT Method Not Working

1. Verify `HTTPS_PROXY` is set (not just `HTTP_PROXY`)
2. Check proxy logs for CONNECT errors
3. Some clients may require explicit proxy configuration

## Next Steps

- **Caching**: Future features will add caching for detected AI requests
- **PII Scrubbing**: Future features will block sensitive data
- **Audit Logging**: Future features will log requests to DuckDB

## Examples

### Complete Workflow

```bash
# 1. Start proxy
./memex proxy --listen :8080

# 2. In another terminal, configure client
export HTTPS_PROXY=http://localhost:8080

# 3. Make AI request (will be detected and forwarded)
curl -X POST https://api.anthropic.com/v1/messages \
  -H "Authorization: Bearer $ANTHROPIC_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "claude-3-opus",
    "messages": [{"role": "user", "content": "Hello"}]
  }'

# 4. Check proxy logs (if debug enabled)
# Should show: SchemaAnthropic detected for /v1/messages
```

### Testing Passthrough

```bash
# Non-AI request (should pass through unchanged)
curl -x http://localhost:8080 https://httpbin.org/get

# Response should be identical to direct request
curl https://httpbin.org/get
```
