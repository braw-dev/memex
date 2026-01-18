package proxy

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"

	"github.com/braw-dev/memex/pkg/types"
	"github.com/go-git/go-git/v5"
)

// DetectScope determines the scope context for the given directory path.
// It tries to find a git remote origin, otherwise falls back to hashing the absolute path.
func DetectScope(path string) (*types.ScopeContext, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Try to open git repo
	// We use PlainOpen which searches for .git in the path
	r, err := git.PlainOpen(absPath)
	if err == nil {
		// Repo found, check remotes
		remotes, err := r.Remotes()
		if err == nil {
			// Find 'origin'
			for _, remote := range remotes {
				if remote.Config().Name == "origin" {
					urls := remote.Config().URLs
					if len(urls) > 0 {
						return &types.ScopeContext{
							ID:   urls[0],
							Type: types.ScopeTypeGitRemote,
							Salt: generateSalt(urls[0]),
						}, nil
					}
				}
			}
			// If no origin, try upstream
			for _, remote := range remotes {
				if remote.Config().Name == "upstream" {
					urls := remote.Config().URLs
					if len(urls) > 0 {
						return &types.ScopeContext{
							ID:   urls[0],
							Type: types.ScopeTypeGitRemote,
							Salt: generateSalt(urls[0]),
						}, nil
					}
				}
			}
		}
	}

	// Fallback to path hash
	hash := sha256.Sum256([]byte(absPath))
	hashStr := hex.EncodeToString(hash[:])

	return &types.ScopeContext{
		ID:   hashStr,
		Type: types.ScopeTypePathHash,
		Salt: hash[:], // Use raw hash bytes as salt
	}, nil
}

func generateSalt(input string) []byte {
	hash := sha256.Sum256([]byte(input))
	return hash[:]
}

// contextKey is a type for context keys to avoid collisions
type scopeContextKey struct{}

// WithScope adds the scope context to the context
func WithScope(ctx context.Context, scope *types.ScopeContext) context.Context {
	return context.WithValue(ctx, scopeContextKey{}, scope)
}

// FromContext retrieves the scope context from the context
func FromContext(ctx context.Context) *types.ScopeContext {
	if scope, ok := ctx.Value(scopeContextKey{}).(*types.ScopeContext); ok {
		return scope
	}
	return nil
}
