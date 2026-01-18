package types

import "log/slog"

const RedactedValue = "==REDACTED=="

// SensitiveString represents a string that contains PII and should be redacted in logs.
type SensitiveString string

// LogValue implements slog.LogValuer to mask the sensitive value.
func (s SensitiveString) LogValue() slog.Value {
	if s == "" {
		return slog.StringValue("")
	}
	return slog.StringValue(RedactedValue)
}
