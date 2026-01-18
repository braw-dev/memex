package unit

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/braw-dev/memex/internal/proxy"
	"github.com/braw-dev/memex/pkg/types"
)

func TestSchemaDetector(t *testing.T) {
	detector := proxy.NewSchemaDetector()

	tests := []struct {
		name     string
		path     string
		expected types.SchemaType
	}{
		{
			name:     "Anthropic Messages",
			path:     "/v1/messages",
			expected: types.SchemaAnthropic,
		},
		{
			name:     "Anthropic Messages with prefix",
			path:     "/api/v1/messages",
			expected: types.SchemaAnthropic,
		},
		{
			name:     "OpenAI Chat Completions",
			path:     "/v1/chat/completions",
			expected: types.SchemaOpenAI,
		},
		{
			name:     "OpenAI with prefix",
			path:     "/openai/v1/chat/completions",
			expected: types.SchemaOpenAI,
		},
		{
			name:     "Unknown Path",
			path:     "/v1/other",
			expected: types.SchemaUnknown,
		},
		{
			name:     "Root",
			path:     "/",
			expected: types.SchemaUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, _ := url.Parse("http://example.com" + tt.path)
			req := &http.Request{URL: u}
			got := detector.Detect(req)
			if got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}
