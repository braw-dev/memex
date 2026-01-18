# Interface Contracts

Since this feature is internal infrastructure, "Contracts" refer to the Go Interfaces that other components will use.

## Store Interface

Located in `internal/store/types.go` (or similar).

```go
package store

import "context"

// Store is the main entrypoint for storage operations
type Store interface {
    // KV returns the Key-Value store handle
    KV() KVStore

    // OLAP returns the OLAP store handle
    OLAP() OLAPStore

    // Close cleanly shuts down both engines
    Close() error
}

// KVStore defines key-value operations (subset of BadgerDB)
type KVStore interface {
    // Set saves a value
    Set(ctx context.Context, key, value []byte) error

    // Get retrieves a value
    Get(ctx context.Context, key []byte) ([]byte, error)

    // Close closes the KV store
    Close() error
}

// OLAPStore defines analytical operations (subset of DuckDB)
type OLAPStore interface {
    // Exec executes a query without returning rows
    Exec(ctx context.Context, query string, args ...any) error

    // Close closes the OLAP store
    Close() error
}
```
