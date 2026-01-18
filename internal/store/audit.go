package store

import (
	"time"
)

// AuditLog represents a single audit log entry
type AuditLog struct {
	Timestamp time.Time `db:"timestamp"`
	ScopeID   string    `db:"scope_id"`
	TokensIn  int       `db:"tokens_in"`
	TokensOut int       `db:"tokens_out"`
	Cost      float64   `db:"cost"`
	Latency   int       `db:"latency"` // in milliseconds
}

// WriteLog inserts a new audit log entry into the database
func (s *Store) WriteLog(log *AuditLog) error {
	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now()
	}

	query := `
	INSERT INTO audit_logs (timestamp, scope_id, tokens_in, tokens_out, cost, latency)
	VALUES (:timestamp, :scope_id, :tokens_in, :tokens_out, :cost, :latency)
	`
	_, err := s.db.NamedExec(query, log)
	return err
}
