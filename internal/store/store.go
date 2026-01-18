package store

import (
	"context"
	"database/sql"
	"fmt"
	"os"
)

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

	// QueryRow executes a query that is expected to return at most one row
	QueryRow(ctx context.Context, query string, args ...any) *sql.Row

	// Close closes the OLAP store
	Close() error
}

type storeImpl struct {
	path string
	kv   KVStore
	olap OLAPStore
}

// New creates a new Store instance in the specified directory.
func New(path string) (Store, error) {
	// Ensure directory exists
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, fmt.Errorf("failed to create store directory: %w", err)
	}

	kv, err := newBadgerKV(path)
	if err != nil {
		return nil, err
	}

	olap, err := newDuckDBOLAP(path)
	if err != nil {
		kv.Close() // Cleanup KV if OLAP fails
		return nil, err
	}

	return &storeImpl{
		path: path,
		kv:   kv,
		olap: olap,
	}, nil
}

func (s *storeImpl) KV() KVStore {
	return s.kv
}

func (s *storeImpl) OLAP() OLAPStore {
	return s.olap
}

func (s *storeImpl) Close() error {
	var errs []error
	if s.kv != nil {
		if err := s.kv.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if s.olap != nil {
		if err := s.olap.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("errors closing store: %v", errs)
	}
	return nil
}
