package store

import (
	"context"
)

const initSchemaSQL = `
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY,
    timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    action VARCHAR NOT NULL,
    model VARCHAR,
    input_tokens INTEGER,
    output_tokens INTEGER,
    cost_usd DECIMAL(10, 6),
    metadata JSON
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_timestamp ON audit_logs(timestamp);
`

func (d *duckDBOLAP) initSchema(ctx context.Context) error {
	_, err := d.db.ExecContext(ctx, initSchemaSQL)
	return err
}
