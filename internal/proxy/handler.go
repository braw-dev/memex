package proxy

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"time"
)

type contextKey string

const schemaContextKey contextKey = "schema"

// proxyHandler handles HTTP requests and forwards them to upstream servers
type proxyHandler struct {
	config   *ProxyConfig
	proxy    *httputil.ReverseProxy
	detector *SchemaDetector
}

// NewServer creates a new proxy server handler
// Returns http.Handler that can be used with http.Server
func NewServer(config *ProxyConfig) http.Handler {
	mux := http.NewServeMux()

	// Initialize proxy handler components
	detector := NewSchemaDetector()

	// Initialize ReverseProxy transport
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment, // Support upstream proxies
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       config.IdleTimeout,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	// Create reverse proxy
	reverseProxy := &httputil.ReverseProxy{
		Director:      makeDirector(detector, config),
		Transport:     transport,
		FlushInterval: config.FlushInterval, // 0 for immediate flushing (SSE)
		ErrorLog:      slog.NewLogLogger(slog.Default().Handler(), slog.LevelError),
		ErrorHandler:  makeErrorHandler(config),
	}

	// Create proxy handler instance
	handler := &proxyHandler{
		config:   config,
		proxy:    reverseProxy,
		detector: detector,
	}

	// Register routes
	mux.HandleFunc("GET /healthz", handleHealthz())
	mux.HandleFunc("/", handler.handleProxy)

	// Apply middleware
	var h http.Handler = mux
	h = debugMiddleware(h, config)

	// ScopeMiddleware (Principle V: Auth/Scope is first)
	h = ScopeMiddleware(h, config)

	return h
}

// handleHealthz returns a health check handler
func handleHealthz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

// handleProxy handles HTTP proxy requests
// TLS termination is handled by an upstream reverse proxy
func (h *proxyHandler) handleProxy(w http.ResponseWriter, r *http.Request) {
	// Track request start
	startTime := time.Now()

	// Detect schema
	schema := h.detector.Detect(r)

	// Store schema in context
	ctx := context.WithValue(r.Context(), schemaContextKey, schema)

	// Forward request
	h.proxy.ServeHTTP(w, r.WithContext(ctx))

	duration := time.Since(startTime)
	slog.Debug("Completed request", "method", r.Method, "path", r.URL.Path, "duration", duration, "schema", schema)
}

// makeDirector creates a director function for ReverseProxy
func makeDirector(detector *SchemaDetector, config *ProxyConfig) func(*http.Request) {
	return func(req *http.Request) {
		// ReverseProxy requires Scheme and Host to be set if they are empty
		if req.URL.Scheme == "" {
			req.URL.Scheme = "http"
		}
		if req.URL.Host == "" {
			req.URL.Host = req.Host
		}

		// Remove Proxy- headers
		req.Header.Del("Proxy-Connection")
		req.Header.Del("Proxy-Authenticate")
		req.Header.Del("Proxy-Authorization")
	}
}

// makeErrorHandler creates an error handler for ReverseProxy
func makeErrorHandler(config *ProxyConfig) func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		if err != nil {
			slog.Error("Proxy error", "err", err, "path", r.URL.Path)
			w.WriteHeader(http.StatusBadGateway)
			// Minimal error response
			w.Write([]byte("Bad Gateway"))
		}
	}
}

// debugMiddleware adds debug logging to requests
func debugMiddleware(next http.Handler, config *ProxyConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Request started", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
