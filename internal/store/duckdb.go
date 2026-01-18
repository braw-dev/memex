package store

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"

	_ "github.com/marcboeker/go-duckdb"
)

type duckDBOLAP struct {
	db *sql.DB
}

func newDuckDBOLAP(basePath string) (*duckDBOLAP, error) {
	path := filepath.Join(basePath, "memex.db")
	// Open DuckDB
	// Note: ?access_mode=READ_WRITE is default, but explicit is fine.
	db, err := sql.Open("duckdb", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open duckdb: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping duckdb: %w", err)
	}

	// Initial configuration
	// Note: We skip "INSTALL json" as it requires download/CGO linkage of extensions not present.
	// We assume standard DuckDB or that we rely on what's available.
	// However, we MUST initialize our schema.
	d := &duckDBOLAP{db: db}
	if err := d.initSchema(context.Background()); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return d, nil
}

func (d *duckDBOLAP) Exec(ctx context.Context, query string, args ...any) error {
	_, err := d.db.ExecContext(ctx, query, args...)
	return err
}

func (d *duckDBOLAP) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return d.db.QueryRowContext(ctx, query, args...)
}

func (d *duckDBOLAP) Close() error {
	return d.db.Close()
}
