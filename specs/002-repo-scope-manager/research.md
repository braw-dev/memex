# Research: Repo-Aware Scope Manager

**Branch**: `002-repo-scope-manager`

## Dependencies

### 1. Git Repository Detection

**Decision**: Use `go-git/v5` library.
**Rationale**: Pure Go implementation, compliant with Constitution Principle II (Go Only) and I (Single Binary). Avoids `os/exec` to `git` binary which might not be installed or in PATH.
**Alternatives Considered**:

- `os/exec` calling `git`: Violates pure Go preference and external dependency risk.
- Parsing `.git/config` manually: Error-prone, re-implementing wheels. `go-git` is robust.

### 2. Path Hashing

**Decision**: `crypto/sha256` from Go standard library.
**Rationale**: Standard, secure, collision-resistant enough for this purpose.
**Alternatives Considered**:

- `fnv` or `crc32`: Faster but higher collision risk. SHA-256 is fast enough (<1ms for short strings).

## Implementation Details

### Scope Detection Logic

1. Attempt to open repository at CWD using `go-git`.
2. If successful, get `origin` remote URL.
3. If no `origin`, try `upstream`.
4. If failure at any step, fall back to Path Hash.

### Cache Key Structure

New Key Format: `SHA256(ScopeID + SystemPromptHash + ASTHash)`
Instead of just `SHA256(SystemPromptHash + ASTHash)`.
This effectively namespaces the cache.
