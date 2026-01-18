package store

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/marcboeker/go-duckdb"
)

// Store represents the DuckDB storage engine
type Store struct {
	db *sqlx.DB
}

// NewStore initializes a new DuckDB store in the .memex directory
func NewStore(dbPath string) (*Store, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Connect to DuckDB
	// Using a connection pool with MaxOpenConns(1) as recommended for DuckDB single-process
	db, err := sqlx.Connect("duckdb", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to duckdb: %w", err)
	}

	db.SetMaxOpenConns(1)

	s := &Store{db: db}

	if err := s.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return s, nil
}

// initSchema creates the necessary tables if they don't exist
func (s *Store) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS audit_logs (
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		scope_id TEXT,
		tokens_in INTEGER,
		tokens_out INTEGER,
		cost DOUBLE,
		latency INTEGER
	);

	CREATE TABLE IF NOT EXISTS cache_entries (
		hash_key TEXT PRIMARY KEY,
		scope_id TEXT,
		system_hash TEXT,
		prompt_vector FLOAT[],
		response_blob BLOB,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := s.db.Exec(schema)
	return err
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.db.Close()
}

// DB returns the underlying sqlx.DB instance
func (s *Store) DB() *sqlx.DB {
	return s.db
}
