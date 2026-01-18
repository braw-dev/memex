package store

import (
	"time"
)

// CacheEntry represents a cached AI response
type CacheEntry struct {
	HashKey      string    `db:"hash_key"`
	ScopeID      string    `db:"scope_id"`
	SystemHash   string    `db:"system_hash"`
	PromptVector []float32 `db:"prompt_vector"`
	ResponseBlob []byte    `db:"response_blob"`
	CreatedAt    time.Time `db:"created_at"`
}

// GetCache retrieves a cache entry by its hash key
func (s *Store) GetCache(hashKey string) (*CacheEntry, error) {
	entry := &CacheEntry{}
	query := `SELECT * FROM cache_entries WHERE hash_key = ?`
	err := s.db.Get(entry, query, hashKey)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

// SetCache inserts or updates a cache entry
func (s *Store) SetCache(entry *CacheEntry) error {
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now()
	}

	query := `
	INSERT OR REPLACE INTO cache_entries (hash_key, scope_id, system_hash, prompt_vector, response_blob, created_at)
	VALUES (:hash_key, :scope_id, :system_hash, :prompt_vector, :response_blob, :created_at)
	`
	_, err := s.db.NamedExec(query, entry)
	return err
}
