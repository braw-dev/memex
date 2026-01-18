package types

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

// SchemaType represents the detected AI provider schema type
type SchemaType int

const (
	SchemaUnknown SchemaType = iota
	SchemaAnthropic
	SchemaOpenAI
)

func (s SchemaType) String() string {
	switch s {
	case SchemaAnthropic:
		return "Anthropic"
	case SchemaOpenAI:
		return "OpenAI"
	default:
		return "Unknown"
	}
}

// ProxyRequest encapsulates an incoming HTTP request with metadata
type ProxyRequest struct {
	Request     *http.Request
	Schema      SchemaType
	UpstreamURL *url.URL
	StartTime   time.Time
}

// ProxyResponse represents the upstream response being forwarded
type ProxyResponse struct {
	StatusCode int
	Headers    http.Header
	Body       io.ReadCloser
	Request    *ProxyRequest
	Duration   time.Duration
}
