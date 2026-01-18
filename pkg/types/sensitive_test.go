package types

import (
	"log/slog"
	"testing"
)

func TestSensitiveString_LogValue(t *testing.T) {
	tests := []struct {
		name     string
		input    SensitiveString
		expected string
	}{
		{
			name:     "Non-empty string",
			input:    SensitiveString("secret"),
			expected: RedactedValue,
		},
		{
			name:     "Empty string",
			input:    SensitiveString(""),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := tt.input.LogValue()
			if val.Kind() != slog.KindString {
				t.Errorf("expected string kind, got %v", val.Kind())
			}
			if val.String() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, val.String())
			}
		})
	}
}
