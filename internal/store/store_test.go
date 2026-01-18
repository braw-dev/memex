package store_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/braw-dev/memex/internal/store"
)

func TestStoreInitialization(t *testing.T) {
	// Setup temporary directory
	tmpDir, err := os.MkdirTemp("", "memex-store-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	memexDir := filepath.Join(tmpDir, ".memex")

	// Initialize Store
	s, err := store.New(memexDir)
	if err != nil {
		t.Fatalf("Failed to initialize store: %v", err)
	}
	defer s.Close()

	// Verify directory creation
	if _, err := os.Stat(memexDir); os.IsNotExist(err) {
		t.Errorf("Expected .memex directory to be created at %s", memexDir)
	}

	// Verify BadgerDB directory
	badgerDir := filepath.Join(memexDir, "badger")
	if _, err := os.Stat(badgerDir); os.IsNotExist(err) {
		t.Errorf("Expected badger directory to be created at %s", badgerDir)
	}

	// Verify DuckDB file
	duckDBFile := filepath.Join(memexDir, "memex.db")
	if _, err := os.Stat(duckDBFile); os.IsNotExist(err) {
		t.Errorf("Expected memex.db file to be created at %s", duckDBFile)
	}
}

func TestStorePersistence(t *testing.T) {
	// Setup temporary directory
	tmpDir, err := os.MkdirTemp("", "memex-store-persistence-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	memexDir := filepath.Join(tmpDir, ".memex")

	// 1. Open Store, Write Data, Close
	func() {
		s, err := store.New(memexDir)
		if err != nil {
			t.Fatalf("Failed to initialize store: %v", err)
		}
		defer s.Close()

		// Write to KV
		err = s.KV().Set(context.Background(), []byte("test-key"), []byte("test-value"))
		if err != nil {
			t.Fatalf("Failed to write to KV: %v", err)
		}

		// Write to OLAP (Create table and insert)
		err = s.OLAP().Exec(context.Background(), "CREATE TABLE test (id INTEGER); INSERT INTO test VALUES (42);")
		if err != nil {
			t.Fatalf("Failed to write to OLAP: %v", err)
		}
	}()

	// 2. Reopen Store, Read Data, Verify
	func() {
		s, err := store.New(memexDir)
		if err != nil {
			t.Fatalf("Failed to re-initialize store: %v", err)
		}
		defer s.Close()

		// Read from KV
		val, err := s.KV().Get(context.Background(), []byte("test-key"))
		if err != nil {
			t.Fatalf("Failed to read from KV: %v", err)
		}
		if string(val) != "test-value" {
			t.Errorf("Expected 'test-value', got '%s'", string(val))
		}

		// Read from OLAP
		var id int
		err = s.OLAP().QueryRow(context.Background(), "SELECT id FROM test LIMIT 1").Scan(&id)
		if err != nil {
			t.Fatalf("Failed to read from OLAP: %v", err)
		}
		if id != 42 {
			t.Errorf("Expected 42, got %d", id)
		}
	}()
}

func TestStoreErrors(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "memex-store-error-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Make directory read-only (chmod 0555 = read-execute, no write)
	if err := os.Chmod(tmpDir, 0555); err != nil {
		t.Fatal(err)
	}

	// Attempt to create store in read-only dir (should fail to create .memex)
	// New attempts to MkdirAll(filepath.Join(tmpDir, ".memex"), ...)
	_, err = store.New(filepath.Join(tmpDir, ".memex"))
	if err == nil {
		t.Error("Expected error when creating store in read-only directory, got nil")
	}
}
