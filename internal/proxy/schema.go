package proxy

import (
	"net/http"
	"strings"

	"github.com/braw-dev/memex/pkg/types"
)

// SchemaDetector detects AI provider schema from request
type SchemaDetector struct{}

// NewSchemaDetector creates a new SchemaDetector
func NewSchemaDetector() *SchemaDetector {
	return &SchemaDetector{}
}

// Detect identifies the schema type based on the request URL path
func (d *SchemaDetector) Detect(r *http.Request) types.SchemaType {
	path := r.URL.Path

	// Detect Anthropic
	if strings.HasSuffix(path, "/v1/messages") {
		return types.SchemaAnthropic
	}

	// Detect OpenAI
	if strings.HasSuffix(path, "/v1/chat/completions") {
		return types.SchemaOpenAI
	}

	return types.SchemaUnknown
}
