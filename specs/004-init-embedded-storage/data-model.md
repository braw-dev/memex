# Data Model: Initialize Embedded Storage Engines

## Key-Value Store (BadgerDB)

### Key Namespaces

Keys are byte arrays. To namespace data, we use prefixes.

| Prefix | Description | Key Format | Value Format |
|--------|-------------|------------|--------------|
| `sys:` | System metadata | `sys:version` | JSON `{"version": "1.0.0", "initialized_at": "..."}` |
| `cache:` | HTTP Response Cache | `cache:<hash>` | Protobuf/Gob encoded Response |

*Note: The Tri-Partite Cache Key structure will be implemented in a future feature, but the `cache:` prefix is reserved now.*

## OLAP Store (DuckDB)

### Schema: `main`

#### Table: `audit_logs`

Stores request/response metrics for "Shadow Billing".

| Column | Type | Description |
|--------|------|-------------|
| `id` | `UUID` | Unique event ID |
| `timestamp` | `TIMESTAMPTZ` | When the event occurred |
| `action` | `VARCHAR` | Event type (e.g., "proxy_request", "cache_hit") |
| `model` | `VARCHAR` | Model ID (e.g., "gpt-4") |
| `input_tokens` | `INTEGER` | Token count input |
| `output_tokens` | `INTEGER` | Token count output |
| `cost_usd` | `DECIMAL(10, 6)` | Estimated cost in USD |
| `metadata` | `JSON` | Additional context (latency, status code, etc.) |

### Initialization SQL

```sql
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY,
    timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    action VARCHAR NOT NULL,
    model VARCHAR,
    input_tokens INTEGER,
    output_tokens INTEGER,
    cost_usd DECIMAL(10, 6),
    metadata JSON
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_timestamp ON audit_logs(timestamp);
```
