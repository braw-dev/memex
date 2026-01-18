package types

import "fmt"

// ScopeType represents the source of the scope identifier
type ScopeType int

const (
	ScopeTypeGitRemote ScopeType = iota
	ScopeTypePathHash
)

// String returns the string representation of ScopeType
func (s ScopeType) String() string {
	switch s {
	case ScopeTypeGitRemote:
		return "GitRemote"
	case ScopeTypePathHash:
		return "PathHash"
	default:
		return "Unknown"
	}
}

// ScopeContext represents the project boundaries for the current execution session
type ScopeContext struct {
	// ID is the unique identifier for the scope (Remote URL or Path Hash)
	ID string
	// Type indicates how the ID was derived
	Type ScopeType
	// Salt is a cryptographic salt derived from ID (used in cache key generation)
	Salt []byte
}

// String returns a string representation of the ScopeContext
func (s *ScopeContext) String() string {
	return fmt.Sprintf("%s (%s)", s.ID, s.Type)
}
