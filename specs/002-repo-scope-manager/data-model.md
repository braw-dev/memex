# Data Model: Repo-Aware Scope Manager

## Entities

### ScopeContext

Represents the project boundaries for the current execution session.

| Field | Type | Description |
|-------|------|-------------|
| `ID` | `string` | The unique identifier for the scope (Remote URL or Path Hash) |
| `Type` | `enum` | `ScopeTypeGitRemote` or `ScopeTypePathHash` |
| `Salt` | `[]byte` | Cryptographic salt derived from ID (used in cache key generation) |

## Constants

```go
type ScopeType int

const (
    ScopeTypeGitRemote ScopeType = iota
    ScopeTypePathHash
)
```

## State Transitions

- **Initialization**: Created at startup based on environment scan.
- **Immutable**: Once set for the process, it does not change (unless we support dynamic re-scoping later, but for now it's static per run).
