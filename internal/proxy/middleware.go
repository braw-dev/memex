package proxy

import (
	"log/slog"
	"net/http"
	"os"
)

// ScopeMiddleware detects and injects the project scope into the request context.
// It must run early in the chain (before caching).
func ScopeMiddleware(next http.Handler, config *ProxyConfig) http.Handler {
	// Detect scope once at initialization as per requirements.
	cwd, err := os.Getwd()
	if err != nil {
		slog.Error("Error getting CWD for scope detection", "err", err)
		cwd = "."
	}

	scope, err := DetectScope(cwd)
	if err != nil {
		slog.Debug("Error detecting scope", "err", err)
	} else {
		slog.Debug("Initialized Scope", "scope", scope)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if scope != nil {
			ctx := WithScope(r.Context(), scope)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}
