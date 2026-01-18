package store

import (
	"context"
	"fmt"
	"path/filepath"

	badger "github.com/dgraph-io/badger/v4"
)

type badgerKV struct {
	db *badger.DB
}

func newBadgerKV(basePath string) (*badgerKV, error) {
	path := filepath.Join(basePath, "badger")
	opts := badger.DefaultOptions(path)
	// Optimization for large binary payloads (as per requirements)
	// Default ValueLogFileSize is 1GB. Lowering it helps with cleanup on smaller disks,
	// but keeping it large is better for performance.
	// Spec/Research said 256MB.
	opts.ValueLogFileSize = 256 * 1024 * 1024

	// Reduce log noise
	opts.Logger = nil

	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open badger db: %w", err)
	}

	return &badgerKV{db: db}, nil
}

func (b *badgerKV) Set(ctx context.Context, key, value []byte) error {
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

func (b *badgerKV) Get(ctx context.Context, key []byte) ([]byte, error) {
	var val []byte
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		val, err = item.ValueCopy(nil)
		return err
	})
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (b *badgerKV) Close() error {
	return b.db.Close()
}
