# Quickstart: Using the Storage Engine

## Initialization

The storage engine is initialized by pointing to a directory.

```go
import "github.com/braw-dev/memex/internal/store"

func main() {
    // Initialize store in .memex directory
    s, err := store.New(".memex")
    if err != nil {
        panic(err)
    }
    defer s.Close()

    // Use KV Store
    err = s.KV().Set(ctx, []byte("my-key"), []byte("my-value"))

    // Use OLAP Store
    err = s.OLAP().Exec(ctx, "INSERT INTO audit_logs (id, action) VALUES (?, ?)", uuid.New(), "test")
}
```

## Running locally

When running the application:

1. Ensure you have write permissions to the current directory.
2. The `.memex` directory will be created automatically.
3. If you encounter lock errors, ensure no other instance is running.
